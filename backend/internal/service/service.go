// Package service is a thin use-case layer over domain and storage. HTTP
// handlers depend on it; it owns ID/time generation and enforces status-
// transition rules that the storage layer cannot express alone.
package service

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"

	"github.com/kjm99d/monkey-planner/backend/internal/domain"
	"github.com/kjm99d/monkey-planner/backend/internal/events"
	"github.com/kjm99d/monkey-planner/backend/internal/storage"
)

// Service is the facade that bundles issue, board, and calendar use cases.
type Service struct {
	repo   storage.Repo
	now    func() time.Time
	broker *events.Broker
}

// New wraps a storage.Repo into a Service. If now is nil, time.Now().UTC()
// is used — pass a fake clock in tests for deterministic timestamps.
func New(repo storage.Repo, now func() time.Time) *Service {
	if now == nil {
		now = func() time.Time { return time.Now().UTC() }
	}
	return &Service{repo: repo, now: now, broker: events.New()}
}

// Broker returns the event broker that SSE handlers subscribe to.
func (s *Service) Broker() *events.Broker { return s.broker }

// publishEvent fans an Event out to all SSE subscribers of the board.
func (s *Service) publishEvent(boardID, eventType, issueID, status string) {
	if s.broker == nil {
		return
	}
	s.broker.Publish(events.Event{
		Type:    eventType,
		BoardID: boardID,
		IssueID: issueID,
		Status:  status,
	})
}

// ---- Issue use cases ----

// CreateIssueInput is the request body for creating an issue.
type CreateIssueInput struct {
	BoardID  string
	ParentID *string
	Title    string
	Body     string
}

var (
	ErrEmptyTitle          = errors.New("title must not be empty")
	ErrMissingBoard        = errors.New("boardId is required")
	ErrApproveViaPatch     = errors.New("use POST /api/issues/:id/approve to set Approved")
	ErrBackwardTransition  = errors.New("backward transition is forbidden")
	ErrInvalidTransition   = errors.New("invalid status transition")
)

func (s *Service) CreateIssue(ctx context.Context, in CreateIssueInput) (domain.Issue, error) {
	if in.Title == "" {
		return domain.Issue{}, ErrEmptyTitle
	}
	if in.BoardID == "" {
		return domain.Issue{}, ErrMissingBoard
	}
	now := s.now()
	iss := domain.Issue{
		ID:         uuid.NewString(),
		BoardID:    in.BoardID,
		ParentID:   in.ParentID,
		Title:      in.Title,
		Body:       in.Body,
		Status:     domain.StatusPending,
		Properties: map[string]any{},
		CreatedAt:  now,
		UpdatedAt:  now,
	}
	created, err := s.repo.Issues().Create(ctx, iss)
	if err != nil {
		return domain.Issue{}, err
	}
	s.DispatchWebhook(created.BoardID, domain.EventIssueCreated, &created)
	s.publishEvent(created.BoardID, "issue.created", created.ID, string(created.Status))
	return created, nil
}

// UpdateIssueInput is the PATCH body; any nil field leaves the value unchanged.
type UpdateIssueInput struct {
	Title        *string
	Body         *string
	Instructions *string
	ParentID     **string              // double pointer: nil = unchanged, *nil = set NULL
	Status       *domain.Status
	Criteria     *[]domain.Criterion   // nil = unchanged
}

// UpdateIssue enforces PATCH transition rules; direct Approved transitions
// are blocked (409) and must use the dedicated approve endpoint.
func (s *Service) UpdateIssue(ctx context.Context, id string, in UpdateIssueInput) (domain.Issue, error) {
	// Validate the transition first to avoid a pointless DB round-trip.
	var patchStatus *domain.Status
	if in.Status != nil {
		cur, err := s.repo.Issues().GetByID(ctx, id)
		if err != nil {
			return domain.Issue{}, err
		}
		if *in.Status == cur.Status {
			// Same-status transition is a silent no-op at the service layer.
		} else if err := domain.ValidateTransition(cur.Status, *in.Status); err != nil {
			switch {
			case errors.Is(err, domain.ErrDirectApproval):
				return domain.Issue{}, ErrApproveViaPatch
			case errors.Is(err, domain.ErrBackwardTransition):
				return domain.Issue{}, ErrBackwardTransition
			default:
				return domain.Issue{}, ErrInvalidTransition
			}
		}
		// QA → Done goes through Complete so completed_at gets recorded.
		if cur.Status == domain.StatusQA && *in.Status == domain.StatusDone {
			done, err := s.repo.Issues().Complete(ctx, id, s.now())
			if err != nil {
				return domain.Issue{}, err
			}
			// Title/Body/ParentID may still need to change alongside the status.
			if in.Title == nil && in.Body == nil && in.ParentID == nil {
				return done, nil
			}
			// status was already applied by Complete; strip it from the follow-up patch.
			in.Status = nil
		} else {
			patchStatus = in.Status
		}
	}

	patch := storage.IssuePatch{
		Title:        in.Title,
		Body:         in.Body,
		Instructions: in.Instructions,
		ParentID:     in.ParentID,
		Status:       patchStatus,
		Criteria:     in.Criteria,
	}
	updated, err := s.repo.Issues().Update(ctx, id, patch)
	if err != nil {
		return updated, err
	}
	if patchStatus != nil {
		s.publishEvent(updated.BoardID, "issue.status_changed", updated.ID, string(updated.Status))
	} else {
		s.publishEvent(updated.BoardID, "issue.updated", updated.ID, string(updated.Status))
	}
	return updated, nil
}

func (s *Service) ApproveIssue(ctx context.Context, id string) (domain.Issue, error) {
	approved, err := s.repo.Issues().Approve(ctx, id, s.now())
	if err != nil {
		return domain.Issue{}, err
	}
	s.DispatchWebhook(approved.BoardID, domain.EventIssueApproved, &approved)
	s.publishEvent(approved.BoardID, "issue.approved", approved.ID, string(approved.Status))
	return approved, nil
}

func (s *Service) CompleteIssue(ctx context.Context, id string) (domain.Issue, error) {
	done, err := s.repo.Issues().Complete(ctx, id, s.now())
	if err != nil {
		return done, err
	}
	s.publishEvent(done.BoardID, "issue.status_changed", done.ID, string(done.Status))
	return done, nil
}

func (s *Service) DeleteIssue(ctx context.Context, id string) error {
	iss, _ := s.repo.Issues().GetByID(ctx, id)
	err := s.repo.Issues().Delete(ctx, id)
	if err != nil {
		return err
	}
	s.DispatchWebhook(iss.BoardID, domain.EventIssueDeleted, &iss)
	s.publishEvent(iss.BoardID, "issue.deleted", iss.ID, "")
	return nil
}

func (s *Service) GetIssue(ctx context.Context, id string) (issue domain.Issue, children []domain.Issue, err error) {
	issue, err = s.repo.Issues().GetByID(ctx, id)
	if err != nil {
		return
	}
	children, err = s.repo.Issues().ListChildren(ctx, id)
	return
}

func (s *Service) ListIssues(ctx context.Context, f storage.IssueFilter) ([]domain.Issue, error) {
	return s.repo.Issues().List(ctx, f)
}

func (s *Service) ReorderIssues(ctx context.Context, issueIDs []string) error {
	return s.repo.Issues().ReorderIssues(ctx, issueIDs)
}

func (s *Service) AddDependency(ctx context.Context, blockerID, blockedID string) error {
	return s.repo.Issues().AddDependency(ctx, blockerID, blockedID)
}

func (s *Service) RemoveDependency(ctx context.Context, blockerID, blockedID string) error {
	return s.repo.Issues().RemoveDependency(ctx, blockerID, blockedID)
}

func (s *Service) GetBlockedBy(ctx context.Context, issueID string) ([]string, error) {
	return s.repo.Issues().GetBlockedBy(ctx, issueID)
}

// ---- Comment use cases ----

func (s *Service) CreateComment(ctx context.Context, issueID, body string) (*domain.Comment, error) {
	c, err := s.repo.Comments().Create(ctx, issueID, body)
	if err != nil {
		return c, err
	}
	if iss, gerr := s.repo.Issues().GetByID(ctx, issueID); gerr == nil {
		s.publishEvent(iss.BoardID, "comment.created", issueID, "")
	}
	return c, nil
}

func (s *Service) ListComments(ctx context.Context, issueID string) ([]domain.Comment, error) {
	return s.repo.Comments().List(ctx, issueID)
}

func (s *Service) DeleteComment(ctx context.Context, commentID string) error {
	return s.repo.Comments().Delete(ctx, commentID)
}

// ---- Board use cases ----

func (s *Service) CreateBoard(ctx context.Context, name string, viewType domain.ViewType) (domain.Board, error) {
	if name == "" {
		return domain.Board{}, errors.New("board name must not be empty")
	}
	if viewType == "" {
		viewType = domain.ViewKanban
	}
	if !viewType.Valid() {
		return domain.Board{}, errors.New("invalid viewType")
	}
	b := domain.Board{
		ID:        uuid.NewString(),
		Name:      name,
		ViewType:  viewType,
		CreatedAt: s.now(),
	}
	return s.repo.Boards().Create(ctx, b)
}

func (s *Service) ListBoards(ctx context.Context) ([]domain.Board, error) {
	return s.repo.Boards().List(ctx)
}

func (s *Service) GetBoard(ctx context.Context, id string) (domain.Board, error) {
	return s.repo.Boards().GetByID(ctx, id)
}

func (s *Service) UpdateBoard(ctx context.Context, id string, name *string, viewType *domain.ViewType) (domain.Board, error) {
	if viewType != nil && !viewType.Valid() {
		return domain.Board{}, errors.New("invalid viewType")
	}
	return s.repo.Boards().Update(ctx, id, name, viewType)
}

func (s *Service) DeleteBoard(ctx context.Context, id string) error {
	return s.repo.Boards().Delete(ctx, id)
}

// ---- Issue property use cases ----

// UpdateIssueProperties atomically merges props into the issue's properties.
// Keys with nil values are removed (RFC 7396 merge-patch semantics). The merge
// happens at the SQL level so concurrent calls from multiple MCP clients do not
// lose each other's writes.
func (s *Service) UpdateIssueProperties(ctx context.Context, id string, props map[string]any) (domain.Issue, error) {
	return s.repo.Issues().MergeProperties(ctx, id, props)
}

// ---- Board property use cases ----

func (s *Service) CreateBoardProperty(ctx context.Context, boardID, name string, propType domain.PropertyType, options []string) (domain.BoardProperty, error) {
	if name == "" {
		return domain.BoardProperty{}, errors.New("property name must not be empty")
	}
	if !propType.Valid() {
		return domain.BoardProperty{}, errors.New("invalid property type")
	}
	if options == nil {
		options = []string{}
	}
	existing, _ := s.repo.BoardProperties().List(ctx, boardID)
	p := domain.BoardProperty{
		BoardID:  boardID,
		Name:     name,
		Type:     propType,
		Options:  options,
		Position: len(existing),
	}
	return s.repo.BoardProperties().Create(ctx, p)
}

func (s *Service) ListBoardProperties(ctx context.Context, boardID string) ([]domain.BoardProperty, error) {
	return s.repo.BoardProperties().List(ctx, boardID)
}

func (s *Service) UpdateBoardProperty(ctx context.Context, id string, name *string, options *[]string, position *int) (domain.BoardProperty, error) {
	return s.repo.BoardProperties().Update(ctx, id, name, options, position)
}

func (s *Service) DeleteBoardProperty(ctx context.Context, id string) error {
	return s.repo.BoardProperties().Delete(ctx, id)
}

// ---- Calendar use cases ----

func (s *Service) GetMonthStats(ctx context.Context, year int, month time.Month) ([]storage.DayCount, error) {
	return s.repo.Issues().GetMonthStats(ctx, year, month)
}

func (s *Service) GetDayStats(ctx context.Context, day time.Time) (storage.DayStats, error) {
	return s.repo.Issues().GetDayStats(ctx, day)
}
