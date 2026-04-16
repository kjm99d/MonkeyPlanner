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
		// 허용 경로
		{"approved→inProgress", StatusApproved, StatusInProgress, nil},
		{"inProgress→done", StatusInProgress, StatusDone, nil},

		// Pending→Approved 직접 차단 (Approve 엔드포인트만)
		{"pending→approved direct PATCH", StatusPending, StatusApproved, ErrDirectApproval},

		// 역행 금지
		{"approved→pending", StatusApproved, StatusPending, ErrBackwardTransition},
		{"inProgress→approved", StatusInProgress, StatusApproved, ErrBackwardTransition},
		{"inProgress→pending", StatusInProgress, StatusPending, ErrBackwardTransition},
		{"done→inProgress", StatusDone, StatusInProgress, ErrBackwardTransition},
		{"done→approved", StatusDone, StatusApproved, ErrBackwardTransition},
		{"done→pending", StatusDone, StatusPending, ErrBackwardTransition},

		// 유효하지 않은 상태
		{"invalid from", Status("Nope"), StatusApproved, ErrInvalidStatus},
		{"invalid to", StatusPending, Status("Nope"), ErrInvalidStatus},

		// 같은 상태
		{"same status pending", StatusPending, StatusPending, ErrSelfSameTransition},
		{"same status done", StatusDone, StatusDone, ErrSelfSameTransition},

		// 알 수 없는 전이
		{"pending→inProgress", StatusPending, StatusInProgress, ErrUnknownTransition},
		{"pending→done", StatusPending, StatusDone, ErrUnknownTransition},
		{"approved→done", StatusApproved, StatusDone, ErrUnknownTransition},
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
