package domain

import (
	"errors"
	"time"
)

// Status is an issue's lifecycle state.
type Status string

const (
	StatusPending    Status = "Pending"
	StatusApproved   Status = "Approved"
	StatusInProgress Status = "InProgress"
	StatusQA         Status = "QA"
	StatusDone       Status = "Done"
	StatusRejected   Status = "Rejected"
)

// Valid reports whether the status is one of the allowed constants.
func (s Status) Valid() bool {
	switch s {
	case StatusPending, StatusApproved, StatusInProgress, StatusQA, StatusDone, StatusRejected:
		return true
	}
	return false
}

// Criterion is a single acceptance-criterion checklist item on an issue.
type Criterion struct {
	Text string `json:"text"`
	Done bool   `json:"done"`
}

// Issue is the atomic unit of agent task memory.
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

// Status-transition sentinels returned by ValidateTransition.
var (
	ErrInvalidStatus       = errors.New("invalid status value")
	ErrBackwardTransition  = errors.New("backward status transition is forbidden")
	ErrDirectApproval      = errors.New("direct Pending→Approved via PATCH is forbidden; use POST /api/issues/:id/approve")
	ErrSelfSameTransition  = errors.New("transition to the same status is a no-op")
	ErrUnknownTransition   = errors.New("unknown status transition")
)

// ValidateTransition checks whether a status change is allowed via the PATCH path.
// Rules:
//   - Pending→Approved is blocked via PATCH (force use of the dedicated Approve endpoint)
//   - Pending→Rejected is allowed (rejection path)
//   - Approved → InProgress (claim)
//   - InProgress → QA (work finished, submit for review)
//   - QA → Done (review passed)
//   - QA → InProgress (review failed, rework)
//   - Rejected is terminal (no further transitions allowed)
func ValidateTransition(from, to Status) error {
	if !from.Valid() || !to.Valid() {
		return ErrInvalidStatus
	}
	if from == to {
		return ErrSelfSameTransition
	}
	// Rejected is a terminal state.
	if from == StatusRejected {
		return ErrUnknownTransition
	}
	// From Pending only Approve (via dedicated endpoint) or Reject are allowed.
	if from == StatusPending {
		if to == StatusApproved {
			return ErrDirectApproval
		}
		if to == StatusRejected {
			return nil
		}
		return ErrUnknownTransition
	}
	// Moving back to Pending or Rejected is never allowed from a post-approval state.
	if to == StatusPending || to == StatusRejected {
		return ErrUnknownTransition
	}
	if to == StatusApproved {
		return ErrDirectApproval
	}
	// Allowed transitions per source state.
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
