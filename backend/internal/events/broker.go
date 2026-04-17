// Package events provides a lightweight in-memory pub/sub broker for SSE streaming.
package events

import (
	"sync"
	"time"
)

// Event is the payload delivered to SSE subscribers.
type Event struct {
	Type      string `json:"type"`
	BoardID   string `json:"boardId"`
	IssueID   string `json:"issueId,omitempty"`
	Status    string `json:"status,omitempty"`
	Timestamp string `json:"timestamp"`
}

// Broker is a per-board channel pub/sub router.
type Broker struct {
	mu      sync.RWMutex
	clients map[string]map[chan Event]struct{} // boardId → set of channels
}

// New returns an empty broker.
func New() *Broker {
	return &Broker{
		clients: make(map[string]map[chan Event]struct{}),
	}
}

// Subscribe returns a channel that receives events for the given boardId.
// The caller must call Unsubscribe to free the channel.
func (b *Broker) Subscribe(boardID string) chan Event {
	ch := make(chan Event, 16) // buffered so bursts do not block slow consumers
	b.mu.Lock()
	defer b.mu.Unlock()
	if b.clients[boardID] == nil {
		b.clients[boardID] = make(map[chan Event]struct{})
	}
	b.clients[boardID][ch] = struct{}{}
	return ch
}

// Unsubscribe detaches the channel and closes it.
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

// Publish fans ev out to every subscriber of ev.BoardID. Slow consumers with
// a full buffer drop the event rather than blocking the publisher — the
// client recovers via full refetch on reconnect.
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
			// Buffer full — drop (client re-syncs on reconnect).
		}
	}
}
