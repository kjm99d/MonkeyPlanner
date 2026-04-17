# awesome-mcp-servers submission

**Target**: https://github.com/punkpeye/awesome-mcp-servers
(also applicable to https://github.com/appcypher/awesome-mcp-servers)

## Location

Under `## Server Implementations` → `### 🗂️ <appropriate category>` — most
likely **Project Management** or a new "Task Trackers" subcategory. If a
`Project Management` section does not exist, **Developer Tools** is the
safe fallback.

## Entry to add

```markdown
- [MonkeyPlanner](https://github.com/kjm99d/MonkeyPlanner) 🏠 📇 🐍 - Local-first kanban / approval-gate issue tracker with a native MCP server (13 tools). Humans approve; agents (Claude Code, Cursor, Continue, …) consume tasks, report progress, and submit for QA. Single Go binary + SQLite.
```

Legend (match whatever the upstream README uses):
- 🏠 — self-hosted / local
- 📇 — written in Go (use the Go flag the list uses; some use 🐹)
- 🐍 — substitute for your language conventions. Drop flags that do not apply.

## PR title

```
Add MonkeyPlanner — local-first task tracker with native MCP
```

## PR body

```markdown
Hi — adding MonkeyPlanner to the list.

**Link**: https://github.com/kjm99d/MonkeyPlanner
**Description**: Local-first task manager for AI coding agents. Exposes
13 MCP tools for creating, approving, claiming, and completing issues —
the human owns the approval gate, the agent owns the execution. Runs as
a single Go binary, stores everything in SQLite, MIT licensed.

### What makes it different from the other task trackers already listed

- The approval gate is a first-class concept, not a convention. Pending →
  Approved can only happen through a dedicated endpoint — agents literally
  cannot self-approve.
- The MCP bridge and the HTTP API share the same service layer, so every
  agent action shows up in the web UI in real time via SSE.
- Ships with `mcp install --for <claude-code|claude-desktop|cursor>` so
  users do not hand-edit JSON paths.

### Checklist

- [x] Entry placed in the correct category and sorted
- [x] Description ≤ 2 lines
- [x] Emoji legend follows the README's convention
- [x] Link is public and works

Thanks for maintaining this list!
```
