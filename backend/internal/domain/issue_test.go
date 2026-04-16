package domain

import (
	"errors"
	"testing"
)

func TestStatusValid(t *testing.T) {
	valids := []Status{StatusPending, StatusApproved, StatusInProgress, StatusDone}
	for _, s := range valids {
		if !s.Valid() {
			t.Errorf("expected %q to be valid", s)
		}
	}
	if Status("Nope").Valid() {
		t.Error("expected Nope to be invalid")
	}
}

func TestValidateTransition(t *testing.T) {
	cases := []struct {
		name    string
		from    Status
		to      Status
		wantErr error
	}{
		// 전진 허용
		{"approved→inProgress", StatusApproved, StatusInProgress, nil},
		{"inProgress→done", StatusInProgress, StatusDone, nil},

		// 역행 허용 (Approved ⇄ InProgress ⇄ Done)
		{"done→inProgress", StatusDone, StatusInProgress, nil},
		{"inProgress→done again", StatusInProgress, StatusDone, nil},

		// Pending→Approved 직접 차단 (Approve 엔드포인트만)
		{"pending→approved direct PATCH", StatusPending, StatusApproved, ErrDirectApproval},

		// Pending에서 다른 곳 직접 이동 불가
		{"pending→inProgress", StatusPending, StatusInProgress, ErrUnknownTransition},
		{"pending→done", StatusPending, StatusDone, ErrUnknownTransition},

		// Approved로 PATCH 전이 차단 (Approve 버튼 전용)
		{"inProgress→approved", StatusInProgress, StatusApproved, ErrDirectApproval},
		{"done→approved", StatusDone, StatusApproved, ErrDirectApproval},

		// Pending으로 되돌리기 불가
		{"approved→pending", StatusApproved, StatusPending, ErrUnknownTransition},
		{"inProgress→pending", StatusInProgress, StatusPending, ErrUnknownTransition},
		{"done→pending", StatusDone, StatusPending, ErrUnknownTransition},

		// 유효하지 않은 상태
		{"invalid from", Status("Nope"), StatusApproved, ErrInvalidStatus},
		{"invalid to", StatusPending, Status("Nope"), ErrInvalidStatus},

		// 같은 상태
		{"same status pending", StatusPending, StatusPending, ErrSelfSameTransition},
		{"same status done", StatusDone, StatusDone, ErrSelfSameTransition},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got := ValidateTransition(tc.from, tc.to)
			if !errors.Is(got, tc.wantErr) {
				t.Fatalf("from=%s to=%s: got %v, want %v", tc.from, tc.to, got, tc.wantErr)
			}
		})
	}
}

func TestViewTypeValid(t *testing.T) {
	if !ViewKanban.Valid() || !ViewList.Valid() {
		t.Fatal("kanban/list must be valid view types")
	}
	if ViewType("grid").Valid() {
		t.Fatal("unknown view type must be invalid")
	}
}
