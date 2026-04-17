package service_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/kjm99d/MonkeyPlanner/backend/internal/domain"
	"github.com/kjm99d/MonkeyPlanner/backend/internal/service"
	"github.com/kjm99d/MonkeyPlanner/backend/internal/storage"
)

func TestBoardCRUD(t *testing.T) {
	ctx := context.Background()
	s := newService(t)

	// Empty name rejected
	if _, err := s.CreateBoard(ctx, "", ""); err == nil {
		t.Fatal("expected empty name error")
	}
	// Invalid viewType rejected
	if _, err := s.CreateBoard(ctx, "x", domain.ViewType("grid")); err == nil {
		t.Fatal("expected invalid viewType error")
	}
	b, err := s.CreateBoard(ctx, "backlog", "") // empty viewType defaults to kanban
	if err != nil || b.ViewType != domain.ViewKanban {
		t.Fatalf("create board defaults: err=%v vt=%s", err, b.ViewType)
	}
	list, err := s.ListBoards(ctx)
	if err != nil || len(list) != 1 {
		t.Fatalf("list: err=%v len=%d", err, len(list))
	}
	got, err := s.GetBoard(ctx, b.ID)
	if err != nil || got.Name != "backlog" {
		t.Fatalf("get: err=%v name=%s", err, got.Name)
	}

	newName := "todo"
	newVT := domain.ViewList
	updated, err := s.UpdateBoard(ctx, b.ID, &newName, &newVT)
	if err != nil || updated.Name != "todo" || updated.ViewType != domain.ViewList {
		t.Fatalf("update: %+v %v", updated, err)
	}

	bad := domain.ViewType("grid")
	if _, err := s.UpdateBoard(ctx, b.ID, nil, &bad); err == nil {
		t.Fatal("expected invalid viewType update error")
	}

	if err := s.DeleteBoard(ctx, b.ID); err != nil {
		t.Fatal(err)
	}
}

func TestCreateIssueValidation(t *testing.T) {
	ctx := context.Background()
	s := newService(t)
	if _, err := s.CreateIssue(ctx, service.CreateIssueInput{}); !errors.Is(err, service.ErrEmptyTitle) {
		t.Fatalf("empty title: %v", err)
	}
	if _, err := s.CreateIssue(ctx, service.CreateIssueInput{Title: "x"}); !errors.Is(err, service.ErrMissingBoard) {
		t.Fatalf("missing board: %v", err)
	}
}

func TestGetIssueWithChildren(t *testing.T) {
	ctx := context.Background()
	s := newService(t)
	b, parent := seed(t, s)

	pid := parent.ID
	_, err := s.CreateIssue(ctx, service.CreateIssueInput{BoardID: b.ID, Title: "child", ParentID: &pid})
	if err != nil {
		t.Fatal(err)
	}
	_, children, err := s.GetIssue(ctx, parent.ID)
	if err != nil || len(children) != 1 {
		t.Fatalf("get with children: err=%v count=%d", err, len(children))
	}
}

func TestListIssuesByFilter(t *testing.T) {
	ctx := context.Background()
	s := newService(t)
	b, _ := seed(t, s)

	_, err := s.CreateIssue(ctx, service.CreateIssueInput{BoardID: b.ID, Title: "another"})
	if err != nil {
		t.Fatal(err)
	}
	bid := b.ID
	pending := domain.StatusPending
	got, err := s.ListIssues(ctx, storage.IssueFilter{BoardID: &bid, Status: &pending})
	if err != nil || len(got) != 2 {
		t.Fatalf("list: err=%v len=%d", err, len(got))
	}
}

func TestUpdateIssueTitleAndBody(t *testing.T) {
	ctx := context.Background()
	s := newService(t)
	_, iss := seed(t, s)
	nt := "새 제목"
	nb := "본문 수정"
	got, err := s.UpdateIssue(ctx, iss.ID, service.UpdateIssueInput{Title: &nt, Body: &nb})
	if err != nil {
		t.Fatal(err)
	}
	if got.Title != nt || got.Body != nb {
		t.Fatalf("update title/body: %+v", got)
	}
}

func TestCalendarStats(t *testing.T) {
	ctx := context.Background()
	s := newService(t)
	_, iss := seed(t, s)
	_, err := s.ApproveIssue(ctx, iss.ID)
	if err != nil {
		t.Fatal(err)
	}

	now := time.Now().UTC()
	stats, err := s.GetMonthStats(ctx, now.Year(), now.Month())
	if err != nil || len(stats) == 0 {
		t.Fatalf("month stats: err=%v len=%d", err, len(stats))
	}
	day, err := s.GetDayStats(ctx, now)
	if err != nil {
		t.Fatal(err)
	}
	if len(day.Created) < 1 || len(day.Approved) < 1 {
		t.Fatalf("day stats created=%d approved=%d", len(day.Created), len(day.Approved))
	}
}

func TestDeleteIssue(t *testing.T) {
	ctx := context.Background()
	s := newService(t)
	_, iss := seed(t, s)
	if err := s.DeleteIssue(ctx, iss.ID); err != nil {
		t.Fatal(err)
	}
}
