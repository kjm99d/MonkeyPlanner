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
	StatusQA         Status = "QA"
	StatusDone       Status = "Done"
	StatusRejected   Status = "Rejected"
)

// Valid 는 Status 가 허용된 값인지 검증합니다.
func (s Status) Valid() bool {
	switch s {
	case StatusPending, StatusApproved, StatusInProgress, StatusQA, StatusDone, StatusRejected:
		return true
	}
	return false
}

// Criterion 은 이슈의 성공 기준 항목입니다.
type Criterion struct {
	Text string `json:"text"`
	Done bool   `json:"done"`
}

// Issue 는 에이전트 작업 기억의 기본 단위입니다.
type Issue struct {
	ID          string            `json:"id"`
	BoardID     string            `json:"boardId"`
	ParentID    *string           `json:"parentId,omitempty"`
	Title        string            `json:"title"`
	Body         string            `json:"body"`
	Instructions string            `json:"instructions"`
	Status       Status            `json:"status"`
	Properties  map[string]any    `json:"properties"`
	Criteria    []Criterion       `json:"criteria"`
	Position    int               `json:"position"`
	CreatedAt   time.Time         `json:"createdAt"`
	UpdatedAt   time.Time         `json:"updatedAt"`
	ApprovedAt  *time.Time        `json:"approvedAt,omitempty"`
	CompletedAt *time.Time        `json:"completedAt,omitempty"`
	BlockedBy   []string          `json:"blockedBy"`
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
//   - Pending→Rejected 전이는 허용 (거절 처리)
//   - Approved → InProgress (claim)
//   - InProgress → QA (작업 완료 후 QA 제출)
//   - QA → Done (검증 통과)
//   - QA → InProgress (검증 실패, 재작업)
//   - Rejected 는 터미널 상태 (다른 상태로 전이 불가)
func ValidateTransition(from, to Status) error {
	if !from.Valid() || !to.Valid() {
		return ErrInvalidStatus
	}
	if from == to {
		return ErrSelfSameTransition
	}
	// Rejected는 터미널 상태
	if from == StatusRejected {
		return ErrUnknownTransition
	}
	// Pending에서는 Approve 버튼 또는 Reject만 사용 가능
	if from == StatusPending {
		if to == StatusApproved {
			return ErrDirectApproval
		}
		if to == StatusRejected {
			return nil
		}
		return ErrUnknownTransition
	}
	// Pending/Rejected로 되돌리기 불가
	if to == StatusPending || to == StatusRejected {
		return ErrUnknownTransition
	}
	if to == StatusApproved {
		return ErrDirectApproval
	}
	// 허용되는 전이 정의
	allowed := map[Status][]Status{
		StatusApproved:   {StatusInProgress},
		StatusInProgress: {StatusQA},
		StatusQA:         {StatusDone, StatusInProgress},
		StatusDone:       {StatusQA},
	}
	for _, s := range allowed[from] {
		if s == to {
			return nil
		}
	}
	return ErrUnknownTransition
}
