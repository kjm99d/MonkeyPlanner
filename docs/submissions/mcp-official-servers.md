# modelcontextprotocol/servers submission

**Target**: https://github.com/modelcontextprotocol/servers

## Location

Add under `README.md` → `### 🤖 Community Servers` (keeps alphabetical order).

## Entry to add

```markdown
- **[MonkeyPlanner](https://github.com/kjm99d/MonkeyPlanner)** - Local-first task memory for AI coding agents. A kanban/approval-gate issue tracker with an MCP server built in (13 tools: `list_boards`, `create_issue`, `approve_issue`, `claim_issue`, `submit_qa`, `reject_issue`, `complete_issue`, `add_comment`, `update_criteria`, `search_issues`, and more). Works with Claude Code, Claude Desktop, Cursor, and any MCP-compatible client. Single binary + SQLite.
```

## PR title

```
Add MonkeyPlanner to Community Servers
```

## PR body

```markdown
## What

MonkeyPlanner is a local-first task manager for AI coding agents. It
provides a kanban/approval-gate workflow where humans review and approve
work, and agents pick up tasks, report progress, and submit for QA — all
through 13 MCP tools.

## Why it belongs here

- Ships an MCP server as a first-class, not bolted-on, feature
- Works with Claude Code, Claude Desktop, Cursor, and Continue out of the
  box (we maintain a `monkey-planner mcp install --for <client>` subcommand)
- Single Go binary + SQLite — no dependencies, no cloud, no telemetry
- MIT licensed, actively maintained

## Checklist

- [x] I have read the contribution guidelines
- [x] Entry is in alphabetical order
- [x] Entry follows the existing format: `**[Name](url)** - Description.`
- [x] Link is to the public repo, not a fork
- [x] Description is concise (≤2 lines) and describes what the server does

Happy to tweak the wording — just let me know.
```

## Before submitting

- [ ] Verify the entry is inserted alphabetically (between the last "M…"
      entry and the first "N…" entry)
- [ ] `git fetch origin main && git rebase origin/main` so the PR opens
      cleanly
- [ ] Double-check the repo URL resolves (the link is public-facing)
