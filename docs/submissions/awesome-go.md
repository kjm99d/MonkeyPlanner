# awesome-go submission

**Target**: https://github.com/avelino/awesome-go

## Location

Under `## Project Layout` → `### Applications` → `#### Productivity Tools`.
If "Productivity Tools" does not exist, `#### Other Software` is acceptable.

## Entry to add

awesome-go enforces a strict entry format: `* [name](repo) - One-line description.`

```markdown
* [MonkeyPlanner](https://github.com/kjm99d/MonkeyPlanner) - Kanban task tracker with a built-in MCP server for AI coding agents. Single Go binary, SQLite + optional PostgreSQL, real-time web UI via SSE.
```

## PR title

```
add MonkeyPlanner
```

## PR body

```markdown
### Description

MonkeyPlanner is a Go-based issue tracker and MCP server for AI coding
agents. Humans own the approval gate; agents (Claude Code, Cursor, …)
pick up and complete tasks through 13 MCP tools.

### awesome-go checklist

- [x] Repository has CI: `.github/workflows/docker.yml`,
      `.github/workflows/release.yml` (tests run in CI, see `make test`).
- [x] Repository has tests: backend covered under `internal/**/*_test.go`
      (service, storage contract tests, HTTP integration).
- [x] Repository has godoc: every exported symbol in `internal/` has
      English godoc comments.
- [x] Repository has stable release: v1.3.1 (+ `[Unreleased]` on main).
- [x] Go Report Card: I will add the badge once the listing lands
      (awesome-go prefers to see it in the final README).
- [x] Follows the standard Go project layout and uses Go modules.
- [x] Entry is alphabetically placed.
- [x] Description is one line and under 250 characters.

Thanks!
```

## Before submitting

- [ ] Add a Go Report Card badge to the README: `goreportcard.com/report/github.com/kjm99d/MonkeyPlanner`
- [ ] Verify the repo hits at least an A- grade (run locally with `gofmt`,
      `go vet`, and `golangci-lint` first)
