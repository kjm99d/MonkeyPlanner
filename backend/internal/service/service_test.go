package service_test

import (
	"context"
	"errors"
	"path/filepath"
	"testing"
	"time"

	"github.com/kjm99d/monkey-planner/backend/internal/domain"
	"github.com/kjm99d/monkey-planner/backend/internal/service"
	"github.com/kjm99d/monkey-planner/backend/internal/storage"
	_ "github.com/kjm99d/monkey-planner/backend/internal/storage/sqlite"
)

func newService(t *testing.T) *service.Service {
	t.Helper()
	dsn := "sqlite://" + filepath.Join(t.TempDir(), "svc.db")
	repo, err := storage.NewRepo(dsn)
	if err != nil {
		t.Fatalf("new repo: %v", err)
	}
	t.Cleanup(func() { _ = repo.Close() })
	return service.New(repo, func() time.Time { return time.Now().UTC() })
}

func seed(t *testing.T, s *service.Service) (domain.Board, domain.Issue) {
	t.Helper()
	ctx := context.Background()
	board, err := s.CreateBoard(ctx, "test", domain.ViewKanban)
	if err != nil {
		t.Fatalf("create board: %v", err)
	}
	iss, err := s.CreateIssue(ctx, service.CreateIssueInput{BoardID: board.ID, Title: "first"})
	if err != nil {
		t.Fatalf("create issue: %v", err)
	}
	return board, iss
}

func TestCreateIssueDefaultsToPending(t *testing.T) {
	s := newService(t)
	_, iss := seed(t, s)
	if iss.Status != domain.StatusPending {
		t.Fatalf("expected Pending, got %s", iss.Status)
	}
	if iss.ApprovedAt != nil || iss.CompletedAt != nil {
		t.Fatalf("expected approved_at/completed_at nil on create")
	}
}

func TestUpdateBlocksDirectApproval(t *testing.T) {
	ctx := context.Background()
	s := newService(t)
	_, iss := seed(t, s)

	approved := domain.StatusApproved
	_, err := s.UpdateIssue(ctx, iss.ID, service.UpdateIssueInput{Status: &approved})
	if !errors.Is(err, service.ErrApproveViaPatch) {
		t.Fatalf("expected ErrApproveViaPatch, got %v", err)
	}
}

func TestApproveIsIdempotent(t *testing.T) {
	ctx := context.Background()
	s := newService(t)
	_, iss := seed(t, s)
	first, err := s.ApproveIssue(ctx, iss.ID)
	if err != nil || first.ApprovedAt == nil {
		t.Fatalf("approve 1: %v %v", err, first.ApprovedAt)
	}
	time.Sleep(10 * time.Millisecond)
	second, err := s.ApproveIssue(ctx, iss.ID)
	if err != nil {
		t.Fatalf("approve 2: %v", err)
	}
	if !second.ApprovedAt.Equal(*first.ApprovedAt) {
		t.Fatalf("not idempotent: %v vs %v", *first.ApprovedAt, *second.ApprovedAt)
	}
}

func TestCompleteFlow(t *testing.T) {
	ctx := context.Background()
	s := newService(t)
	_, iss := seed(t, s)

	if _, err := s.ApproveIssue(ctx, iss.ID); err != nil {
		t.Fatal(err)
	}
	ip := domain.StatusInProgress
	if _, err := s.UpdateIssue(ctx, iss.ID, service.UpdateIssueInput{Status: &ip}); err != nil {
		t.Fatal(err)
	}
	done := domain.StatusDone
	got, err := s.UpdateIssue(ctx, iss.ID, service.UpdateIssueInput{Status: &done})
	if err != nil {
		t.Fatal(err)
	}
	if got.Status != domain.StatusDone || got.CompletedAt == nil {
		t.Fatalf("expected Done+completedAt set, got status=%s completed=%v", got.Status, got.CompletedAt)
	}
}

func TestBackwardTransitionBlocked(t *testing.T) {
	ctx := context.Background()
	s := newService(t)
	_, iss := seed(t, s)
	_, _ = s.ApproveIssue(ctx, iss.ID)

	pending := domain.StatusPending
	_, err := s.UpdateIssue(ctx, iss.ID, service.UpdateIssueInput{Status: &pending})
	if !errors.Is(err, service.ErrBackwardTransition) {
		t.Fatalf("expected ErrBackwardTransition, got %v", err)
	}
}
