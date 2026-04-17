package sqlite

import (
	"context"
	"errors"
	"path/filepath"
	"testing"
	"time"

	"github.com/kjm99d/MonkeyPlanner/backend/internal/domain"
	"github.com/kjm99d/MonkeyPlanner/backend/internal/storage"
	"github.com/google/uuid"
)

func newTestRepo(t *testing.T) *Repo {
	t.Helper()
	dir := t.TempDir()
	repo, err := Open(filepath.Join(dir, "test.db"))
	if err != nil {
		t.Fatalf("open: %v", err)
	}
	t.Cleanup(func() { _ = repo.Close() })
	return repo
}

func seedBoard(t *testing.T, r *Repo) domain.Board {
	t.Helper()
	b := domain.Board{
		ID:        uuid.NewString(),
		Name:      "test",
		ViewType:  domain.ViewKanban,
		CreatedAt: time.Now().UTC(),
	}
	_, err := r.Boards().Create(context.Background(), b)
	if err != nil {
		t.Fatalf("seed board: %v", err)
	}
	return b
}

func TestCRUDAndApproveIdempotent(t *testing.T) {
	ctx := context.Background()
	r := newTestRepo(t)
	b := seedBoard(t, r)

	iss := domain.Issue{
		ID: uuid.NewString(), BoardID: b.ID, Title: "t", Body: "",
		Status: domain.StatusPending, CreatedAt: time.Now().UTC(), UpdatedAt: time.Now().UTC(),
	}
	if _, err := r.Issues().Create(ctx, iss); err != nil {
		t.Fatalf("create: %v", err)
	}

	got, err := r.Issues().GetByID(ctx, iss.ID)
	if err != nil || got.Status != domain.StatusPending {
		t.Fatalf("getByID: err=%v status=%s", err, got.Status)
	}

	now := time.Now().UTC()
	app, err := r.Issues().Approve(ctx, iss.ID, now)
	if err != nil || app.Status != domain.StatusApproved || app.ApprovedAt == nil {
		t.Fatalf("approve 1: err=%v status=%s approved_at=%v", err, app.Status, app.ApprovedAt)
	}
	firstApprovedAt := *app.ApprovedAt

	time.Sleep(20 * time.Millisecond)
	app2, err := r.Issues().Approve(ctx, iss.ID, time.Now().UTC())
	if err != nil {
		t.Fatalf("approve 2: %v", err)
	}
	if !app2.ApprovedAt.Equal(firstApprovedAt) {
		t.Fatalf("approve not idempotent: first=%v second=%v", firstApprovedAt, *app2.ApprovedAt)
	}
}

func TestCycleDetection(t *testing.T) {
	ctx := context.Background()
	r := newTestRepo(t)
	b := seedBoard(t, r)
	ir := r.Issues()

	// A ← B ← C
	mk := func(title string) domain.Issue {
		i := domain.Issue{
			ID: uuid.NewString(), BoardID: b.ID, Title: title,
			Status: domain.StatusPending, CreatedAt: time.Now().UTC(), UpdatedAt: time.Now().UTC(),
		}
		if _, err := ir.Create(ctx, i); err != nil {
			t.Fatalf("create %s: %v", title, err)
		}
		return i
	}
	a, bi, c := mk("A"), mk("B"), mk("C")
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
	// Try to set A's parent to C, which would form a cycle A→C→B→A.
	err := setParent(a.ID, c.ID)
	if !errors.Is(err, storage.ErrCycle) {
		t.Fatalf("expected ErrCycle, got %v", err)
	}
	// An issue cannot be its own parent.
	err = setParent(a.ID, a.ID)
	if !errors.Is(err, storage.ErrCycle) {
		t.Fatalf("self-parent expected ErrCycle, got %v", err)
	}
}

func TestCascadeDelete(t *testing.T) {
	ctx := context.Background()
	r := newTestRepo(t)
	b := seedBoard(t, r)
	ir := r.Issues()

	parent := domain.Issue{
		ID: uuid.NewString(), BoardID: b.ID, Title: "P",
		Status: domain.StatusPending, CreatedAt: time.Now().UTC(), UpdatedAt: time.Now().UTC(),
	}
	if _, err := ir.Create(ctx, parent); err != nil {
		t.Fatal(err)
	}
	pid := parent.ID
	child := domain.Issue{
		ID: uuid.NewString(), BoardID: b.ID, ParentID: &pid, Title: "C",
		Status: domain.StatusPending, CreatedAt: time.Now().UTC(), UpdatedAt: time.Now().UTC(),
	}
	if _, err := ir.Create(ctx, child); err != nil {
		t.Fatal(err)
	}
	if err := ir.Delete(ctx, parent.ID); err != nil {
		t.Fatalf("delete parent: %v", err)
	}
	if _, err := ir.GetByID(ctx, child.ID); !errors.Is(err, storage.ErrNotFound) {
		t.Fatalf("cascade failed: child still exists, err=%v", err)
	}
}

func TestMonthStats(t *testing.T) {
	ctx := context.Background()
	r := newTestRepo(t)
	b := seedBoard(t, r)
	ir := r.Issues()

	day := time.Date(2026, 4, 16, 12, 0, 0, 0, time.UTC)
	iss := domain.Issue{
		ID: uuid.NewString(), BoardID: b.ID, Title: "x",
		Status: domain.StatusPending, CreatedAt: day, UpdatedAt: day,
	}
	if _, err := ir.Create(ctx, iss); err != nil {
		t.Fatal(err)
	}
	if _, err := ir.Approve(ctx, iss.ID, day); err != nil {
		t.Fatal(err)
	}

	stats, err := ir.GetMonthStats(ctx, 2026, 4)
	if err != nil {
		t.Fatalf("month stats: %v", err)
	}
	var found bool
	for _, d := range stats {
		if d.Date.Year() == 2026 && d.Date.Month() == 4 && d.Date.Day() == 16 {
			found = true
			if d.Created != 1 || d.Approved != 1 {
				t.Fatalf("expected created=1 approved=1, got %+v", d)
			}
		}
	}
	if !found {
		t.Fatalf("day 2026-04-16 not found in stats: %+v", stats)
	}
}
