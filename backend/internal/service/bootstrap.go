package service

import (
	"context"

	"github.com/kjm99d/monkey-planner/backend/internal/domain"
)

// welcomeBody is shown as the body of the first issue on a brand-new install.
// Written in Markdown so it renders nicely in the issue detail view.
const welcomeBody = `Welcome to **MonkeyPlanner** — the local-first task manager for your AI coding agents.

## What to try first

1. **Approve this issue** with the ✅ button (or press ⌘↵). Approval is the
   only way an agent can start working on a task.
2. **Connect an AI agent**: run ` + "`monkey-planner mcp install --for claude-code`" + `
   (or ` + "`--for claude-desktop`" + ` / ` + "`--for cursor`" + `) to wire up this server.
3. **Ask your agent** ` + "`list_issues`" + ` — it should see this ticket.

## Learn more

- Every issue flows through: Pending → Approved → InProgress → QA → Done.
- Agents cannot self-approve; approval is the human's job, by design.
- Everything is stored locally in SQLite. No cloud, no telemetry.

Delete this issue (or the whole board) when you are ready to make it yours.
`

// welcomeInstructions is the hidden-to-humans, agent-facing Instructions field.
// Agents are told explicitly that this is a demo ticket and not to write code.
const welcomeInstructions = `This is a demo issue auto-created on the first run of MonkeyPlanner.

Do NOT start implementation work. Only read the human's approval, then add
a short comment acknowledging you can see the board. The human will delete
this issue once they verify the connection works.`

// SeedWelcomeIfEmpty creates a Welcome board and issue, but only when no
// boards exist yet. It is safe to call on every server start — subsequent
// calls are no-ops.
//
// The goal is to collapse the first-run empty-state cliff: a user who just
// installed the binary sees a populated board, knows what to click, and has
// instructions their agent can immediately consume.
func (s *Service) SeedWelcomeIfEmpty(ctx context.Context) error {
	boards, err := s.repo.Boards().List(ctx)
	if err != nil {
		return err
	}
	if len(boards) > 0 {
		return nil
	}

	board, err := s.CreateBoard(ctx, "Welcome", domain.ViewKanban)
	if err != nil {
		return err
	}

	iss, err := s.CreateIssue(ctx, CreateIssueInput{
		BoardID: board.ID,
		Title:   "👋 Welcome — approve me to get started",
		Body:    welcomeBody,
	})
	if err != nil {
		return err
	}

	// Attach agent-facing instructions via a follow-up update so the write
	// path matches the normal PATCH flow (and thus publishes the same events).
	instr := welcomeInstructions
	if _, err := s.UpdateIssue(ctx, iss.ID, UpdateIssueInput{Instructions: &instr}); err != nil {
		return err
	}
	return nil
}
