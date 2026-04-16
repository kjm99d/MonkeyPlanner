// Package service는 도메인·storage 계층을 얇게 감싼 유스케이스 계층입니다.
// HTTP 핸들러가 의존하는 지점이며, 상태 전이 규칙과 ID/시간 생성 책임을 가집니다.
package service

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"

	"github.com/ckmdevb/monkey-planner/backend/internal/domain"
	"github.com/ckmdevb/monkey-planner/backend/internal/storage"
)

// Service는 이슈·보드·캘린더 유스케이스를 묶은 퍼사드입니다.
type Service struct {
	repo storage.Repo
	now  func() time.Time
}

// New는 storage.Repo 를 감싸 Service 를 만듭니다. now 가 nil이면 time.Now().UTC() 사용.
func New(repo storage.Repo, now func() time.Time) *Service {
	if now == nil {
		now = func() time.Time { return time.Now().UTC() }
	}
	return &Service{repo: repo, now: now}
}

// ---- Issue 유스케이스 ----

// CreateIssueInput은 이슈 생성 요청 본문입니다.
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
	return created, nil
}

// UpdateIssueInput은 PATCH 본문입니다. nil 은 미변경.
type UpdateIssueInput struct {
	Title    *string
	Body     *string
	ParentID **string        // 이중 포인터로 "미변경" vs "NULL로" 구분
	Status   *domain.Status
}

// UpdateIssue는 PATCH 전이 규칙을 강제합니다. Approved 전이는 차단(409).
func (s *Service) UpdateIssue(ctx context.Context, id string, in UpdateIssueInput) (domain.Issue, error) {
	// 상태 전이 검증 먼저 (DB 왕복 전)
	var patchStatus *domain.Status
	if in.Status != nil {
		cur, err := s.repo.Issues().GetByID(ctx, id)
		if err != nil {
			return domain.Issue{}, err
		}
		if *in.Status == cur.Status {
			// 같은 상태는 no-op로 허용
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
		// InProgress → Done 은 Complete 메서드로 completed_at 기록
		if cur.Status == domain.StatusInProgress && *in.Status == domain.StatusDone {
			done, err := s.repo.Issues().Complete(ctx, id, s.now())
			if err != nil {
				return domain.Issue{}, err
			}
			// Title/Body/ParentID 도 추가로 변경이 필요할 수 있음
			if in.Title == nil && in.Body == nil && in.ParentID == nil {
				return done, nil
			}
			// status 는 이미 적용됐으므로 추가 Update 에서 빼기
			in.Status = nil
		} else {
			patchStatus = in.Status
		}
	}

	patch := storage.IssuePatch{
		Title:    in.Title,
		Body:     in.Body,
		ParentID: in.ParentID,
		Status:   patchStatus,
	}
	return s.repo.Issues().Update(ctx, id, patch)
}

func (s *Service) ApproveIssue(ctx context.Context, id string) (domain.Issue, error) {
	approved, err := s.repo.Issues().Approve(ctx, id, s.now())
	if err != nil {
		return domain.Issue{}, err
	}
	s.DispatchWebhook(approved.BoardID, domain.EventIssueApproved, &approved)
	return approved, nil
}

func (s *Service) CompleteIssue(ctx context.Context, id string) (domain.Issue, error) {
	return s.repo.Issues().Complete(ctx, id, s.now())
}

func (s *Service) DeleteIssue(ctx context.Context, id string) error {
	iss, _ := s.repo.Issues().GetByID(ctx, id)
	err := s.repo.Issues().Delete(ctx, id)
	if err != nil {
		return err
	}
	s.DispatchWebhook(iss.BoardID, domain.EventIssueDeleted, &iss)
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

// ---- Board 유스케이스 ----

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

// ---- Issue Properties 유스케이스 ----

// UpdateIssueProperties 는 이슈의 커스텀 속성을 merge 업데이트합니다.
func (s *Service) UpdateIssueProperties(ctx context.Context, id string, props map[string]any) (domain.Issue, error) {
	cur, err := s.repo.Issues().GetByID(ctx, id)
	if err != nil {
		return domain.Issue{}, err
	}
	if cur.Properties == nil {
		cur.Properties = map[string]any{}
	}
	for k, v := range props {
		if v == nil {
			delete(cur.Properties, k)
		} else {
			cur.Properties[k] = v
		}
	}
	merged := cur.Properties
	return s.repo.Issues().Update(ctx, id, storage.IssuePatch{Properties: &merged})
}

// ---- Board Properties 유스케이스 ----

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

// ---- Calendar 유스케이스 ----

func (s *Service) GetMonthStats(ctx context.Context, year int, month time.Month) ([]storage.DayCount, error) {
	return s.repo.Issues().GetMonthStats(ctx, year, month)
}

func (s *Service) GetDayStats(ctx context.Context, day time.Time) (storage.DayStats, error) {
	return s.repo.Issues().GetDayStats(ctx, day)
}
