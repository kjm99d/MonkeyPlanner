package domain

import (
	"errors"
	"time"
)

// Status 는 이슈의 라이프사이클 상태입니다.
type Status string

const (
	StatusPending    Status = "Pending"
	StatusApproved   Status = "Approved"
	StatusInProgress Status = "InProgress"
	StatusDone       Status = "Done"
)

// Valid 는 Status 가 허용된 값인지 검증합니다.
func (s Status) Valid() bool {
	switch s {
	case StatusPending, StatusApproved, StatusInProgress, StatusDone:
		return true
	}
	return false
}

// Issue 는 에이전트 작업 기억의 기본 단위입니다.
type Issue struct {
	ID          string            `json:"id"`
	BoardID     string            `json:"boardId"`
	ParentID    *string           `json:"parentId,omitempty"`
	Title       string            `json:"title"`
	Body        string            `json:"body"`
	Status      Status            `json:"status"`
	Properties  map[string]any    `json:"properties"`
	CreatedAt   time.Time         `json:"createdAt"`
	UpdatedAt   time.Time         `json:"updatedAt"`
	ApprovedAt  *time.Time        `json:"approvedAt,omitempty"`
	CompletedAt *time.Time        `json:"completedAt,omitempty"`
}

// 전이 규칙 에러
var (
	ErrInvalidStatus       = errors.New("invalid status value")
	ErrBackwardTransition  = errors.New("backward status transition is forbidden")
	ErrDirectApproval      = errors.New("direct Pending→Approved via PATCH is forbidden; use POST /api/issues/:id/approve")
	ErrSelfSameTransition  = errors.New("transition to the same status is a no-op")
	ErrUnknownTransition   = errors.New("unknown status transition")
)

// ValidateTransition 은 PATCH 경로로 허용되는 상태 전이를 검증합니다.
// 규칙:
//   - Pending→Approved 직접 전이는 차단 (Approve 전용 엔드포인트 사용 강제)
//   - Pending→(InProgress|Done) 직접 전이는 차단 (반드시 Approve를 거쳐야 함)
//   - Approved ⇄ InProgress ⇄ Done 사이는 자유 이동 허용 (단일 사용자 유연성)
func ValidateTransition(from, to Status) error {
	if !from.Valid() || !to.Valid() {
		return ErrInvalidStatus
	}
	if from == to {
		return ErrSelfSameTransition
	}
	// Pending에서는 Approve 버튼만 사용 가능
	if from == StatusPending {
		if to == StatusApproved {
			return ErrDirectApproval
		}
		return ErrUnknownTransition
	}
	// Approved/InProgress/Done 사이는 자유 이동
	if to == StatusPending {
		return ErrUnknownTransition // Pending으로 되돌리기는 불가
	}
	if to == StatusApproved {
		return ErrDirectApproval // Approved로 가려면 Approve 버튼 사용
	}
	return nil
}
