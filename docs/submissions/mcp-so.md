# mcp.so submission

**Target**: https://mcp.so (and the source repo behind it)

## Listing content

Copy-paste into the submission form.

**Name**: MonkeyPlanner

**One-line**: Local-first task memory for AI coding agents. Approve with a click; your agents do the rest.

**Category**: Task Management / Project Management / Productivity

**Description**:
```
MonkeyPlanner is an open-source kanban board with an MCP server built in.
Humans create issues and approve them via a dedicated approval gate; AI
coding agents (Claude Code, Claude Desktop, Cursor, Continue) claim work,
submit for QA, and add comments — all through 13 MCP tools.

Key traits:
- Local-first: single Go binary, SQLite, no cloud, no telemetry
- Forever MIT licensed
- Real-time web UI updates via SSE so agent-driven changes appear instantly
- Multi-language UI (English, Korean, Japanese, Chinese)
- Easy wiring: `monkey-planner mcp install --for claude-code` writes the
  right .mcp.json entry for you

Docker quickstart:
  docker run -p 8080:8080 -v $(pwd)/data:/data ghcr.io/kjm99d/monkeyplanner:latest
```

**Tags**: mcp, claude-code, claude-desktop, cursor, kanban, task-management, self-hosted, local-first, go, react

**Repository**: https://github.com/kjm99d/MonkeyPlanner
**Docker**: ghcr.io/kjm99d/monkeyplanner:latest
**License**: MIT
**Language**: Go (backend), TypeScript (frontend)
**Author**: kjm99d

## Screenshots to attach

- `docs/screenshots/d-home-l.png` — dashboard, light mode
- `docs/screenshots/d-board-l.png` — kanban board, light mode
- `docs/screenshots/d-issue-l.png` — issue detail, light mode
- `docs/demo/monkey-planner-demo.gif` — 30s end-to-end flow
