# Contributing to MonkeyPlanner

Thanks for taking the time — small or large, every PR helps.

MonkeyPlanner is an **AI-agent task manager**. It is maintained mostly on
weekends by one person, so response times are 24–72 hours. That is fine.
Please do not panic-ping.

## Ways to help

| What | Good for |
|------|----------|
| Try it and open an issue with your setup if anything breaks | Everyone |
| Improve docs, examples, or screenshots | First-time contributors |
| Translate UI strings to a new locale | Bilingual users |
| Add a new MCP client integration (`mcp install --for ...`) | Go devs |
| Fix a "good first issue" | New Go/TypeScript contributors |
| Propose a new feature via Discussions | Everyone |

Before starting non-trivial work, open an issue or a Discussion to align on
approach. Nothing hurts more than a great PR that has to be rejected.

## Development setup

Prereqs:

- Go **1.26+**
- Node **20+**
- Make (optional but recommended)

```bash
git clone https://github.com/<you>/monkey-planner
cd monkey-planner
make init               # installs git hooks + frontend deps

# In two terminals:
make run-backend        # http://localhost:8080
make run-frontend       # http://localhost:5173 (Vite dev server, proxies /api)

# Tests:
make test               # backend + frontend + a11y
make test-backend       # Go tests only
make test-frontend      # Vitest
```

If you only touch the backend, you do not need Node. If you only touch
the frontend, keep the backend running in one terminal.

## Project layout

```
backend/
  cmd/monkey-planner/    entry point + MCP stdio bridge
  internal/
    domain/              pure types + state transition rules
    service/             use cases (the layer HTTP + MCP sit on)
    storage/             SQLite and PostgreSQL adapters
    events/              SSE pub/sub broker
    http/                chi router, handlers, middleware
frontend/
  src/
    features/            one folder per page
    components/          shared UI primitives
    i18n/                en / ko / ja / zh
    api/                 React Query hooks + fetch client
```

## Pull request checklist

- [ ] Ran `make test` locally and it passes
- [ ] New public Go symbols have a godoc comment (English)
- [ ] New i18n keys are added to **all four** locales (en/ko/ja/zh)
- [ ] Commit messages follow [Conventional Commits](https://www.conventionalcommits.org/)
      (`fix(...)`, `feat(...)`, `docs(...)`, `refactor(...)`, `test(...)`, `ci(...)`, `chore(...)`)
- [ ] The PR description explains **why**, not just **what**

## Code style

- **Go**: `gofmt` and `go vet` are enforced. Prefer small, named types
  over `map[string]any` on public APIs.
- **TypeScript**: strict mode is on. Avoid `any` in new code.
- **Comments**: English only. If you need to explain a subtle invariant,
  do — but if the code is self-explanatory, skip the comment.
- **SQL**: every query must be parameterized. Use the existing migration
  pattern (`0010_xxx.sql`) when changing the schema.

## Writing good issues

Include, in order of importance:

1. The smallest steps to reproduce
2. What you expected to happen
3. What actually happened
4. OS + `monkey-planner --version` output
5. Which MCP client(s) are connected

## Community

- **Questions**: open a [Discussion](https://github.com/kjm99d/MonkeyPlanner/discussions/categories/q-a)
- **Bugs**: open an [Issue](https://github.com/kjm99d/MonkeyPlanner/issues/new?template=bug_report.yml)
- **Security vulnerabilities**: see [SECURITY.md](./SECURITY.md) — do not file them publicly

## License

By contributing you agree that your contributions will be licensed under
the MIT License (see [LICENSE](./LICENSE)). MonkeyPlanner will remain
**forever free, forever MIT**.
