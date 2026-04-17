// Package contract provides the test suite every storage adapter must satisfy.
// It ensures the SQLite and PostgreSQL adapters behave identically; the
// PostgreSQL run is skipped when MP_PG_DSN is unset.
package contract

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/google/uuid"

	"github.com/kjm99d/MonkeyPlanner/backend/internal/domain"
	"github.com/kjm99d/MonkeyPlanner/backend/internal/storage"
)

// RunAll runs every registered contract case against the supplied adapter.
// The factory must return a fresh storage.Repo per subtest so cases do not
// share state.
func RunAll(t *testing.T, newRepo func(t *testing.T) storage.Repo) {
	t.Helper()
	cases := []struct {
		name string
		fn   func(t *testing.T, repo storage.Repo)
	}{
		{"CreateAndGet", testCreateAndGet},
		{"ApproveIdempotent", testApproveIdempotent},
		{"CycleDetection", testCycleDetection},
		{"CascadeDelete", testCascadeDelete},
		{"MonthStats", testMonthStats},
	}
	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			repo := newRepo(t)
			tc.fn(t, repo)
		})
	}
}

func seedBoard(t *testing.T, repo storage.Repo) domain.Board {
	t.Helper()
	b := domain.Board{
		ID:        uuid.NewString(),
		Name:      "contract-board",
		ViewType:  domain.ViewKanban,
		CreatedAt: time.Now().UTC(),
	}
	if _, err := repo.Boards().Create(context.Background(), b); err != nil {
		t.Fatalf("seed board: %v", err)
	}
	return b
}

func mkIssue(boardID string) domain.Issue {
	now := time.Now().UTC()
	return domain.Issue{
		ID: uuid.NewString(), BoardID: boardID, Title: "t",
		Status: domain.StatusPending, CreatedAt: now, UpdatedAt: now,
	}
}

func testCreateAndGet(t *testing.T, repo storage.Repo) {
	ctx := context.Background()
	b := seedBoard(t, repo)
	iss := mkIssue(b.ID)
	if _, err := repo.Issues().Create(ctx, iss); err != nil {
		t.Fatalf("create: %v", err)
	}
	got, err := repo.Issues().GetByID(ctx, iss.ID)
	if err != nil {
		t.Fatalf("get: %v", err)
	}
	if got.Status != domain.StatusPending {
		t.Errorf("status: got=%s want=Pending", got.Status)
	}
}

func testApproveIdempotent(t *testing.T, repo storage.Repo) {
	ctx := context.Background()
	b := seedBoard(t, repo)
	iss := mkIssue(b.ID)
	if _, err := repo.Issues().Create(ctx, iss); err != nil {
		t.Fatal(err)
	}
	now := time.Now().UTC()
	first, err := repo.Issues().Approve(ctx, iss.ID, now)
	if err != nil || first.ApprovedAt == nil {
		t.Fatalf("approve 1: err=%v approved=%v", err, first.ApprovedAt)
	}
	firstApproved := *first.ApprovedAt
	time.Sleep(15 * time.Millisecond)
	second, err := repo.Issues().Approve(ctx, iss.ID, time.Now().UTC())
	if err != nil {
		t.Fatalf("approve 2: %v", err)
	}
	if !second.ApprovedAt.Equal(firstApproved) {
		t.Fatalf("not idempotent: first=%v second=%v", firstApproved, *second.ApprovedAt)
	}
}

func testCycleDetection(t *testing.T, repo storage.Repo) {
	ctx := context.Background()
	b := seedBoard(t, repo)
	ir := repo.Issues()

	create := func(title string) domain.Issue {
		i := mkIssue(b.ID)
		i.Title = title
		if _, err := ir.Create(ctx, i); err != nil {
			t.Fatalf("create %s: %v", title, err)
		}
		return i
	}
	a, bi, c := create("A"), create("B"), create("C")
	setParent := func(child, parent string) error {
		p := &parent
		pp := &p
		_, err := ir.Update(ctx, child, storage.IssuePatch{ParentID: pp})
		return err
	}
	if err := setParent(bi.ID, a.ID); err != nil {
		t.Fatalf("B→A: %v", err)
	}
	if err := setParent(c.ID, bi.ID); err != nil {
		t.Fatalf("C→B: %v", err)
	}
	if err := setParent(a.ID, c.ID); !errors.Is(err, storage.ErrCycle) {
		t.Fatalf("cycle A→C expected ErrCycle, got %v", err)
	}
	if err := setParent(a.ID, a.ID); !errors.Is(err, storage.ErrCycle) {
		t.Fatalf("self parent expected ErrCycle, got %v", err)
	}
}

func testCascadeDelete(t *testing.T, repo storage.Repo) {
	ctx := context.Background()
	b := seedBoard(t, repo)
	ir := repo.Issues()
	parent := mkIssue(b.ID)
	if _, err := ir.Create(ctx, parent); err != nil {
		t.Fatal(err)
	}
	pid := parent.ID
	child := mkIssue(b.ID)
	child.ParentID = &pid
	if _, err := ir.Create(ctx, child); err != nil {
		t.Fatal(err)
	}
	if err := ir.Delete(ctx, parent.ID); err != nil {
		t.Fatal(err)
	}
	if _, err := ir.GetByID(ctx, child.ID); !errors.Is(err, storage.ErrNotFound) {
		t.Fatalf("cascade failed: %v", err)
	}
}

func testMonthStats(t *testing.T, repo storage.Repo) {
	ctx := context.Background()
	b := seedBoard(t, repo)
	ir := repo.Issues()
	day := time.Now().UTC()
	iss := mkIssue(b.ID)
	iss.CreatedAt = day
	iss.UpdatedAt = day
	if _, err := ir.Create(ctx, iss); err != nil {
		t.Fatal(err)
	}
	if _, err := ir.Approve(ctx, iss.ID, day); err != nil {
		t.Fatal(err)
	}
	stats, err := ir.GetMonthStats(ctx, day.Year(), day.Month())
	if err != nil {
		t.Fatal(err)
	}
	var seen bool
	for _, d := range stats {
		if d.Date.Day() == day.Day() && d.Created >= 1 && d.Approved >= 1 {
			seen = true
		}
	}
	if !seen {
		t.Fatalf("expected created/approved for today, got %+v", stats)
	}
}
