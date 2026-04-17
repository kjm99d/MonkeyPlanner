package storage

import (
	"context"
	"errors"
	"time"

	"github.com/kjm99d/MonkeyPlanner/backend/internal/domain"
)

// Shared sentinel errors returned (or wrapped) by every storage adapter so
// that service-layer code can remain backend-agnostic.
var (
	ErrNotFound = errors.New("not found")
	ErrCycle    = errors.New("parent_id would create a cycle")
	ErrConflict = errors.New("conflict")
)

// IssueFilter holds optional predicates for the List query.
type IssueFilter struct {
	BoardID  *string
	ParentID *string // empty string = only top-level issues (parent IS NULL)
	Status   *domain.Status
}

// DayCount is a per-day bucket used by the calendar month aggregate.
type DayCount struct {
	Date      time.Time `json:"date"` // midnight UTC
	Created   int       `json:"created"`
	Approved  int       `json:"approved"`
	Completed int       `json:"completed"`
}

// DayStats splits a single day's issues into the three activity buckets.
type DayStats struct {
	Created   []domain.Issue `json:"created"`
	Approved  []domain.Issue `json:"approved"`
	Completed []domain.Issue `json:"completed"`
}

// IssuePatch carries only the fields that should be mutated on a PATCH
// (nil = leave unchanged).
type IssuePatch struct {
	Title        *string
	Body         *string
	Instructions *string
	ParentID     **string            // double-pointer: nil (unchanged) vs *nil (set NULL)
	Status       *domain.Status      // service layer returns 409 on direct Approved
	Properties   *map[string]any     // nil = unchanged, non-nil = whole-object replace
	Criteria     *[]domain.Criterion // nil = unchanged
}

// IssueRepo is the storage contract for issues.
type IssueRepo interface {
	Create(ctx context.Context, issue domain.Issue) (domain.Issue, error)
	GetByID(ctx context.Context, id string) (domain.Issue, error)
	ListChildren(ctx context.Context, parentID string) ([]domain.Issue, error)
	List(ctx context.Context, f IssueFilter) ([]domain.Issue, error)
	Update(ctx context.Context, id string, patch IssuePatch) (domain.Issue, error)
	// MergeProperties atomically merges props into the issue's properties JSON.
	// A key with a nil value is removed (RFC 7396 merge-patch semantics). This
	// avoids the read-modify-write race inherent in service-layer merging.
	MergeProperties(ctx context.Context, id string, props map[string]any) (domain.Issue, error)
	Delete(ctx context.Context, id string) error
	// Approve is idempotent: calling it on an already-Approved issue keeps the
	// original approved_at timestamp.
	Approve(ctx context.Context, id string, now time.Time) (domain.Issue, error)
	// Complete transitions InProgress → Done and records completed_at.
	Complete(ctx context.Context, id string, now time.Time) (domain.Issue, error)
	// GetMonthStats returns per-day created/approved/completed counts for the
	// given year+month (UTC).
	GetMonthStats(ctx context.Context, year int, month time.Month) ([]DayCount, error)
	// GetDayStats returns the three-bucket issue lists for a single day.
	GetDayStats(ctx context.Context, day time.Time) (DayStats, error)
	// ReorderIssues updates positions to match the supplied ID ordering.
	ReorderIssues(ctx context.Context, issueIDs []string) error
	// AddDependency records that blockerID must complete before blockedID.
	AddDependency(ctx context.Context, blockerID, blockedID string) error
	// RemoveDependency drops the blockerID → blockedID edge.
	RemoveDependency(ctx context.Context, blockerID, blockedID string) error
	// GetBlockedBy returns the IDs currently blocking the given issue.
	GetBlockedBy(ctx context.Context, issueID string) ([]string, error)
}

// BoardRepo is the storage contract for boards.
type BoardRepo interface {
	Create(ctx context.Context, board domain.Board) (domain.Board, error)
	GetByID(ctx context.Context, id string) (domain.Board, error)
	List(ctx context.Context) ([]domain.Board, error)
	Update(ctx context.Context, id string, name *string, viewType *domain.ViewType) (domain.Board, error)
	Delete(ctx context.Context, id string) error
}

// BoardPropertyRepo is the storage contract for per-board custom properties.
type BoardPropertyRepo interface {
	Create(ctx context.Context, prop domain.BoardProperty) (domain.BoardProperty, error)
	List(ctx context.Context, boardID string) ([]domain.BoardProperty, error)
	Update(ctx context.Context, id string, name *string, options *[]string, position *int) (domain.BoardProperty, error)
	Delete(ctx context.Context, id string) error
}

// WebhookRepo is the storage contract for per-board outbound webhooks.
type WebhookRepo interface {
	Create(ctx context.Context, wh domain.Webhook) (domain.Webhook, error)
	List(ctx context.Context, boardID string) ([]domain.Webhook, error)
	ListByEvent(ctx context.Context, boardID string, event domain.WebhookEvent) ([]domain.Webhook, error)
	Update(ctx context.Context, id string, name *string, url *string, events *[]domain.WebhookEvent, enabled *bool) (domain.Webhook, error)
	Delete(ctx context.Context, id string) error
}

// CommentRepo is the storage contract for issue comments.
type CommentRepo interface {
	Create(ctx context.Context, issueID, body string) (*domain.Comment, error)
	List(ctx context.Context, issueID string) ([]domain.Comment, error)
	Delete(ctx context.Context, commentID string) error
}

// Repo is the aggregate root that exposes all per-entity repos.
type Repo interface {
	Issues() IssueRepo
	Boards() BoardRepo
	BoardProperties() BoardPropertyRepo
	Webhooks() WebhookRepo
	Comments() CommentRepo
	Close() error
}
