# awesome-ai-agents submission

**Target**: https://github.com/e2b-dev/awesome-ai-agents
(or another active "awesome AI agents" repo)

## Location

Under `## Tools & Infrastructure` → `### Task / Memory` or similar. If no
MCP-specific section exists, `## Tools` is the safe fallback.

## Entry to add

```markdown
- [MonkeyPlanner](https://github.com/kjm99d/MonkeyPlanner) - Task memory + approval gate for AI coding agents. Native MCP server (13 tools) that works with Claude Code, Claude Desktop, Cursor, and any MCP-compatible client. Self-hosted, local-first, MIT.
```

## PR body

```markdown
Adds MonkeyPlanner to the list.

**What it is**: An open-source issue tracker + MCP server designed around
the observation that AI coding agents need two things we never build for
them — (1) memory that survives between sessions, and (2) a human
approval checkpoint before they start writing code.

MonkeyPlanner gives them both:
- Every task has an Instructions field the agent reads, and a rejection
  comment trail the agent can learn from.
- Status pipeline is `Pending → Approved → InProgress → QA → Done`. The
  approval step is on a dedicated endpoint so agents literally cannot
  self-approve.
- 13 MCP tools cover the whole lifecycle; `monkey-planner mcp install
  --for claude-code|claude-desktop|cursor` wires it up in one command.

Stack: Go + SQLite backend, React + TypeScript web UI, distroless Docker
image. MIT forever.
```
