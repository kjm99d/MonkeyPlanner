package domain

import (
	"errors"
	"testing"
)

func TestStatusValid(t *testing.T) {
	valids := []Status{StatusPending, StatusApproved, StatusInProgress, StatusQA, StatusDone}
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
		// Forward transitions allowed.
		{"approved→inProgress", StatusApproved, StatusInProgress, nil},
		{"inProgress→qa", StatusInProgress, StatusQA, nil},
		{"qa→done", StatusQA, StatusDone, nil},

		// QA → InProgress (rework after reviewer rejects QA).
		{"qa→inProgress", StatusQA, StatusInProgress, nil},

		// Done → QA (re-review).
		{"done→qa", StatusDone, StatusQA, nil},

		// InProgress → Done is blocked; must go through QA.
		{"inProgress→done blocked", StatusInProgress, StatusDone, ErrUnknownTransition},

		// Pending→Approved is blocked via PATCH (use the dedicated Approve endpoint).
		{"pending→approved direct PATCH", StatusPending, StatusApproved, ErrDirectApproval},

		// From Pending no direct move to post-approval states is allowed.
		{"pending→inProgress", StatusPending, StatusInProgress, ErrUnknownTransition},
		{"pending→done", StatusPending, StatusDone, ErrUnknownTransition},

		// Approved via PATCH is always blocked (Approve-endpoint only).
		{"inProgress→approved", StatusInProgress, StatusApproved, ErrDirectApproval},
		{"done→approved", StatusDone, StatusApproved, ErrDirectApproval},

		// Moving back to Pending is forbidden from any post-approval state.
		{"approved→pending", StatusApproved, StatusPending, ErrUnknownTransition},
		{"inProgress→pending", StatusInProgress, StatusPending, ErrUnknownTransition},
		{"done→pending", StatusDone, StatusPending, ErrUnknownTransition},

		// Unknown status values.
		{"invalid from", Status("Nope"), StatusApproved, ErrInvalidStatus},
		{"invalid to", StatusPending, Status("Nope"), ErrInvalidStatus},

		// Same-status transitions are a no-op.
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
