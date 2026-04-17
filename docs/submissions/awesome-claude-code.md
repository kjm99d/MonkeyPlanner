# awesome-claude-code submission

**Target**: https://github.com/hesreallyhim/awesome-claude-code
(or any similarly-named community list — verify the active one before submitting)

## Location

Under a section like `## MCP Servers` or `## Integrations`. If neither
exists, propose a new `## MCP Servers` section in the PR.

## Entry to add

```markdown
- [MonkeyPlanner](https://github.com/kjm99d/MonkeyPlanner) - Local-first kanban / approval-gate issue tracker that hooks Claude Code into a 13-tool MCP server. Wire it up with `monkey-planner mcp install --for claude-code`, then ask Claude to `list_issues`, `claim_issue`, `submit_qa`. Forever MIT.
```

## PR body

```markdown
Adds MonkeyPlanner — a task tracker that was built specifically so Claude
Code sessions stop losing context between restarts.

- Installs the MCP server with one command: `monkey-planner mcp install
  --for claude-code` writes the right `.mcp.json` entry.
- 13 MCP tools covering the full lifecycle: create → approve → claim →
  submit for QA → complete, plus comments, criteria, dependencies, search.
- Human owns the approval gate, agent owns the execution. The approve
  endpoint is the only path from Pending → Approved, so agents cannot
  accidentally self-approve work.
- Single Go binary + SQLite, Docker image at ghcr.io/kjm99d/monkeyplanner.

Happy to rework the entry if a different shape fits better.
```
