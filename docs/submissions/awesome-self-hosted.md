# awesome-selfhosted submission

**Target**: https://github.com/awesome-selfhosted/awesome-selfhosted

## Location

Under `## Software` → `### Task Management & To-do Lists`.

## Entry to add

Awesome-Selfhosted enforces a rigid format: `[name](repo) - Description. (Demo link, Source Code, Technology)` and a strict license line.

```markdown
- [MonkeyPlanner](https://github.com/kjm99d/MonkeyPlanner) - Kanban / approval-gate issue tracker with a built-in MCP server for AI coding agents (Claude Code, Cursor, Continue). Single Go binary + SQLite, Docker image provided. `MIT` `Go/TypeScript`
```

## PR title

```
Add MonkeyPlanner to Task Management & To-do Lists
```

## PR body

```markdown
Adding MonkeyPlanner under "Task Management & To-do Lists".

**Project**: https://github.com/kjm99d/MonkeyPlanner

### Why it fits awesome-selfhosted

- 100% self-hosted: runs as a single Go binary, stores everything in
  SQLite on local disk. No cloud, no telemetry, no outbound connections
  except user-configured webhooks.
- Single-container Docker image (`ghcr.io/kjm99d/monkeyplanner`),
  distroless/static base, ~20 MB.
- MIT licensed — listed under the free license requirements.
- Actively maintained, weekly cadence at minimum.

### Format checklist

- [x] Alphabetical placement within the category
- [x] Description ≤ 250 chars, starts with a capital letter, ends with a period
- [x] License listed as `MIT`
- [x] Language listed (`Go/TypeScript`)
- [x] Public source repository linked

Thanks!
```
