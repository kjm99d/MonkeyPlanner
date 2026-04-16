// Package events provides a lightweight in-memory pub/sub broker for SSE streaming.
package events

import (
	"sync"
	"time"
)

// Event 는 클라이언트에 전송되는 이벤트 페이로드입니다.
type Event struct {
	Type      string `json:"type"`
	BoardID   string `json:"boardId"`
	IssueID   string `json:"issueId,omitempty"`
	Status    string `json:"status,omitempty"`
	Timestamp string `json:"timestamp"`
}

// Broker 는 boardId 채널 기반 pub/sub 브로커입니다.
type Broker struct {
	mu      sync.RWMutex
	clients map[string]map[chan Event]struct{} // boardId → set of channels
}

// New 는 빈 브로커를 생성합니다.
func New() *Broker {
	return &Broker{
		clients: make(map[string]map[chan Event]struct{}),
	}
}

// Subscribe 는 해당 boardId의 이벤트를 받을 채널을 반환합니다.
// 호출자는 반드시 Unsubscribe로 해제해야 합니다.
func (b *Broker) Subscribe(boardID string) chan Event {
	ch := make(chan Event, 16) // 버퍼로 느린 소비자 대비
	b.mu.Lock()
	defer b.mu.Unlock()
	if b.clients[boardID] == nil {
		b.clients[boardID] = make(map[chan Event]struct{})
	}
	b.clients[boardID][ch] = struct{}{}
	return ch
}

// Unsubscribe 는 채널을 등록 해제하고 닫습니다.
func (b *Broker) Unsubscribe(boardID string, ch chan Event) {
	b.mu.Lock()
	defer b.mu.Unlock()
	if set, ok := b.clients[boardID]; ok {
		if _, exists := set[ch]; exists {
			delete(set, ch)
			close(ch)
		}
		if len(set) == 0 {
			delete(b.clients, boardID)
		}
	}
}

// Publish 는 해당 boardId의 모든 구독자에게 이벤트를 비동기 전송합니다.
// 버퍼가 가득 찬 느린 소비자는 해당 이벤트를 드롭합니다 (블로킹 방지).
func (b *Broker) Publish(ev Event) {
	if ev.Timestamp == "" {
		ev.Timestamp = time.Now().UTC().Format(time.RFC3339)
	}
	b.mu.RLock()
	defer b.mu.RUnlock()
	for ch := range b.clients[ev.BoardID] {
		select {
		case ch <- ev:
		default:
			// 버퍼 full → 드롭 (재연결 시 full refetch로 복구)
		}
	}
}
