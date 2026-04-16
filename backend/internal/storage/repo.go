package storage

import (
	"context"
	"errors"
	"time"

	"github.com/kjm99d/monkey-planner/backend/internal/domain"
)

// 공통 에러 (어댑터가 래핑해서 반환)
var (
	ErrNotFound  = errors.New("not found")
	ErrCycle     = errors.New("parent_id would create a cycle")
	ErrConflict  = errors.New("conflict")
)

// IssueFilter 는 List 쿼리의 선택적 조건입니다.
type IssueFilter struct {
	BoardID  *string
	ParentID *string // 빈 문자열이면 "루트(parent=NULL)"만
	Status   *domain.Status
}

// DayCount 는 캘린더 월 집계 응답의 일별 카운트입니다.
type DayCount struct {
	Date      time.Time `json:"date"` // UTC 자정
	Created   int       `json:"created"`
	Approved  int       `json:"approved"`
	Completed int       `json:"completed"`
}

// DayStats 는 특정 일자의 3분할 이슈 목록입니다.
type DayStats struct {
	Created   []domain.Issue `json:"created"`
	Approved  []domain.Issue `json:"approved"`
	Completed []domain.Issue `json:"completed"`
}

// IssuePatch 는 PATCH 요청에서 변경할 필드만 담습니다 (nil = 미변경).
type IssuePatch struct {
	Title    *string
	Body     *string
	ParentID **string        // 이중 포인터: nil(미변경) vs *nil(=NULL로 설정)
	Status   *domain.Status  // Approved로 전이 시도 시 서비스 계층이 409 응답
}

// IssueRepo 는 이슈 저장소 인터페이스입니다.
type IssueRepo interface {
	Create(ctx context.Context, issue domain.Issue) (domain.Issue, error)
	GetByID(ctx context.Context, id string) (domain.Issue, error)
	ListChildren(ctx context.Context, parentID string) ([]domain.Issue, error)
	List(ctx context.Context, f IssueFilter) ([]domain.Issue, error)
	Update(ctx context.Context, id string, patch IssuePatch) (domain.Issue, error)
	Delete(ctx context.Context, id string) error
	// Approve 는 멱등: 이미 Approved인 이슈에 호출해도 approved_at 유지.
	Approve(ctx context.Context, id string, now time.Time) (domain.Issue, error)
	// Complete 는 InProgress→Done 전이 + completed_at 기록.
	Complete(ctx context.Context, id string, now time.Time) (domain.Issue, error)
	// GetMonthStats 는 지정 연/월의 일별 created/approved/completed 카운트.
	GetMonthStats(ctx context.Context, year int, month time.Month) ([]DayCount, error)
	// GetDayStats 는 특정 날짜의 3분할 이슈 목록.
	GetDayStats(ctx context.Context, day time.Time) (DayStats, error)
}

// BoardRepo 는 보드 저장소 인터페이스입니다.
type BoardRepo interface {
	Create(ctx context.Context, board domain.Board) (domain.Board, error)
	GetByID(ctx context.Context, id string) (domain.Board, error)
	List(ctx context.Context) ([]domain.Board, error)
	Update(ctx context.Context, id string, name *string, viewType *domain.ViewType) (domain.Board, error)
	Delete(ctx context.Context, id string) error
}

// Repo 는 이슈/보드 레포를 함께 제공하는 상위 인터페이스입니다.
type Repo interface {
	Issues() IssueRepo
	Boards() BoardRepo
	Close() error
}
