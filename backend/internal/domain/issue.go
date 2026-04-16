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
	ID          string     `json:"id"`
	BoardID     string     `json:"boardId"`
	ParentID    *string    `json:"parentId,omitempty"`
	Title       string     `json:"title"`
	Body        string     `json:"body"`
	Status      Status     `json:"status"`
	CreatedAt   time.Time  `json:"createdAt"`
	UpdatedAt   time.Time  `json:"updatedAt"`
	ApprovedAt  *time.Time `json:"approvedAt,omitempty"`
	CompletedAt *time.Time `json:"completedAt,omitempty"`
}

// 전이 규칙 에러
var (
	ErrInvalidStatus       = errors.New("invalid status value")
	ErrBackwardTransition  = errors.New("backward status transition is forbidden")
	ErrDirectApproval      = errors.New("direct Pending→Approved via PATCH is forbidden; use POST /api/issues/:id/approve")
	ErrSelfSameTransition  = errors.New("transition to the same status is a no-op")
	ErrUnknownTransition   = errors.New("unknown status transition")
)

// ValidateTransition 은 PATCH 경로로 허용되는 상태 전이만 승인합니다.
// Pending→Approved 직접 전이는 차단됩니다 (Approve 전용 엔드포인트 사용 강제).
// 역행(Done→InProgress 등)은 모두 차단됩니다.
func ValidateTransition(from, to Status) error {
	if !from.Valid() || !to.Valid() {
		return ErrInvalidStatus
	}
	if from == to {
		return ErrSelfSameTransition
	}
	switch from {
	case StatusPending:
		if to == StatusApproved {
			return ErrDirectApproval
		}
		return ErrUnknownTransition
	case StatusApproved:
		if to == StatusInProgress {
			return nil
		}
		if to == StatusPending {
			return ErrBackwardTransition
		}
		return ErrUnknownTransition
	case StatusInProgress:
		if to == StatusDone {
			return nil
		}
		if to == StatusPending || to == StatusApproved {
			return ErrBackwardTransition
		}
		return ErrUnknownTransition
	case StatusDone:
		return ErrBackwardTransition
	}
	return ErrUnknownTransition
}
