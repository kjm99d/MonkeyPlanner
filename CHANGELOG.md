# Changelog

All notable changes to MonkeyPlanner are documented here.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

## [1.4.1] — 2026-04-17

### Added
- **Homebrew tap** — every tagged release now auto-pushes a formula to
  [`kjm99d/homebrew-tap`](https://github.com/kjm99d/homebrew-tap). Install
  with `brew tap kjm99d/tap && brew install monkey-planner`. See
  [`docs/HOMEBREW_SETUP.md`](./docs/HOMEBREW_SETUP.md) for the one-time
  maintainer PAT setup.
- **Go Report Card badge** in README; link surface to the tool's analysis page.
- `LICENSE` file (MIT) — previously referenced everywhere but never
  committed; goreleaser now includes it in every archive.
- **One-shot `fix-lockfile` workflow** to regenerate `package-lock.json`
  on demand whenever Dependabot leaves it out of sync.

### Changed
- **All four READMEs** (en/ko/ja/zh) lead with `monkey-planner mcp install
  --for <client>` instead of the old `update-and-run.sh/bat` wrappers —
  downloading extra scripts is no longer necessary.
- `golangci-lint` + `govulncheck` jobs in CI are non-blocking until the
  upstream tools catch up to Go 1.26.
- goreleaser pipeline moved the frontend build to explicit GitHub Actions
  steps (visible npm errors) and fixed an invalid `changelog.use` value.

### Fixed
- `@vitejs/plugin-react` downgraded from 6→4.7 — v6 requires `vite@^8` and
  would break `npm ci` on a fresh clone. Merged Dependabot PR #11 was too
  aggressive; now pinned to the v4 track compatible with our Vite 6.

## [1.4.0] — 2026-04-17

### Added
- **Docker** — `Dockerfile` (multi-stage, distroless/static, ~20 MB),
  `docker-compose.yml`, and a multi-arch (amd64+arm64) GHCR publish workflow.
  Run the whole stack with one command: `docker run -p 8080:8080 ghcr.io/kjm99d/monkeyplanner`.
- **`monkey-planner mcp install`** — writes the correct `.mcp.json` entry for
  Claude Code, Claude Desktop, or Cursor automatically. Supports `--dry-run`,
  `--force`, `--scope user`, `--base-url`, and `--name`.
- **First-run Welcome board** — new databases are seeded with a Welcome board
  and a demo issue that walks through the approve-then-agent flow. Idempotent.
- **Agent Presence Bar** — the board page now shows a real-time strip that
  lists which issues are currently InProgress or in QA review, updated via SSE.
- **⌘↵ / Ctrl+↵ approve shortcut** — pressing ⌘↵ on a focused Pending issue
  card triggers approval. Success shows a brief "Approved ✓" green flash.
- **Security headers** — every response now carries `X-Content-Type-Options`,
  `X-Frame-Options: DENY`, `Referrer-Policy`, `Permissions-Policy`, and a
  Content-Security-Policy on non-API paths.
- **Optional API token auth** — set `MP_API_TOKEN` to require
  `Authorization: Bearer <token>` on all `/api/*` requests. Unset = no auth
  (default, backward compatible).
- **CORS** — strict localhost-only policy by default; override via
  `MP_CORS_ORIGINS`.
- **Webhook SSRF defense** — `validateWebhookURL` blocks cloud metadata
  endpoints (169.254.169.254, metadata.google.internal, etc.) and
  RFC1918/loopback addresses. Opt out with `MP_WEBHOOK_ALLOW_PRIVATE=1`.
- **SSE heartbeat** — 30-second comment line keeps connections alive through
  reverse proxies. `onerror` now triggers an immediate refetch; a 5-minute
  fallback interval handles silent drops.
- **Graceful shutdown** — SIGINT/SIGTERM drains in-flight requests with a
  15-second window.
- **Server timeouts** — `ReadHeaderTimeout: 10s`, `ReadTimeout: 30s`,
  `IdleTimeout: 120s` (Slowloris mitigation).
- **Atomic properties merge** — `IssueRepo.MergeProperties` uses SQLite
  `json_patch()` / Postgres `jsonb || jsonb_strip_nulls` so concurrent MCP
  clients no longer race on property updates.
- **SQLite concurrency** — `busy_timeout=5000ms` + `_txlock=immediate` +
  `MaxOpenConns=1`. Multi-client MCP setups (Claude Code + Cursor + Desktop)
  no longer hit `SQLITE_BUSY`.
- **CI workflow** — `.github/workflows/ci.yml` runs `go test -race`,
  `govulncheck`, `golangci-lint`, and TypeScript check on every push and PR.
- **`CONTRIBUTING.md`**, **`CODE_OF_CONDUCT.md`**, **`SECURITY.md`**,
  **`CHANGELOG.md`** (this file), goreleaser config.

### Changed
- `mcp.go` (606 lines) split into 7 focused files: `mcp_jsonrpc.go`,
  `mcp_client.go`, `mcp_tools_registry.go`, `mcp_tools_issues.go`,
  `mcp_tools_comments.go`, `mcp_helpers.go` + slimmed `mcp.go` (38 lines).
- English display name unified to `MonkeyPlanner` across all READMEs, OpenAPI
  spec, HTML title, and i18n. Kebab form `monkey-planner` kept for binary,
  repo slug, and Docker image.
- All Go source comments translated from Korean to English.
- `useEventStream` issue invalidation scoped to `{ boardId }` (was global).
- recharts lazy-loaded with `React.lazy()` to reduce initial bundle by ~400 KB.
- `.mcp.json` removed from version control (personal paths); `.mcp.json.example`
  added as a reference template.
- OpenAPI `Status` enum updated to include `QA` and `Rejected` (added in
  migration 0009 but missing from the spec).

### Fixed
- **ON DELETE CASCADE not working** — migration 0009 accidentally dropped the
  `REFERENCES issues(id) ON DELETE CASCADE` constraint from `parent_id`.
  Restored; `TestCascadeDelete` now passes.
- **SSE invalidation scope** — events for board A no longer invalidated board B.
- Missing i18n keys (`board.showDone`, `board.hideDone`, `board.fullscreen`,
  `board.exitFullscreen`, `chart.noData`) added to all four locales.
- Status-transition tests updated to route through QA (migration 0009 inserted
  QA between InProgress and Done).
- Server test `NewRouter` call updated to match the 3-arg signature.

## [1.3.1] — 2026-04-17

### Fixed
- MCP tool calls now route through the HTTP API instead of hitting the DB
  directly, so mutations made via MCP appear in real time on open web UI
  tabs via SSE.

## [1.3.0] — 2026-04-17

### Added
- Server-Sent Events (`/api/events?boardId=...`) for real-time web UI
  updates when issues change.
- `get_version` MCP tool + build-time version injection (`-ldflags -X`).

### Fixed
- DB `CHECK` constraint extended to include `QA` and `Rejected` statuses
  (migration 0009) so the new workflow no longer trips the old constraint.
- Health-check endpoint now reports the version the binary was built with.

## [1.2.0] — 2026-04-17

### Added
- Release binaries embed the frontend via `-tags prod`, so the downloaded
  binary runs the UI with no extra setup.
- `update-and-run.bat` wrapper for Windows.
- Auto-update wrapper script + documentation for MCP installs.

## [1.1.1] — 2026-04-17

### Fixed
- Removed unused import in `IssueCard`.
- Added missing `recharts` dependency (regression from dashboard work).

## [1.1.0] — 2026-04-16

### Added
- **QA** status + Kanban column, turning the pipeline into
  `Pending → Approved → InProgress → QA → Done`.
- GitHub Release automation: tag push builds multi-platform binaries and
  attaches them to the release.

### Fixed
- Language switcher dropdown failing to open on click.

## [1.0.0] — 2026-04-16

Initial public release.

### Added
- Kanban + list board with drag-and-drop.
- Issue approval gate (`POST /api/issues/:id/approve`) — the only path from
  Pending to Approved.
- MCP server with 10 tools for Claude Code / Claude Desktop agents.
- Multi-language UI (en/ko/ja/zh) with separated README files.
- Agent Activity feed and metrics dashboard.
- SQLite primary storage with Goose migrations, PostgreSQL adapter skeleton.
- Issue templates, success-criteria checklists, agent Instructions field,
  comments, dependencies.
- Webhooks (Discord / Slack / Telegram compatible).
- Full-text issue search (Cmd+K), keyboard shortcuts, dark mode.

[Unreleased]: https://github.com/kjm99d/MonkeyPlanner/compare/v1.4.1...HEAD
[1.4.1]: https://github.com/kjm99d/MonkeyPlanner/compare/v1.4.0...v1.4.1
[1.4.0]: https://github.com/kjm99d/MonkeyPlanner/compare/v1.3.1...v1.4.0
[1.3.1]: https://github.com/kjm99d/MonkeyPlanner/compare/v1.3.0...v1.3.1
[1.3.0]: https://github.com/kjm99d/MonkeyPlanner/compare/v1.2.0...v1.3.0
[1.2.0]: https://github.com/kjm99d/MonkeyPlanner/compare/v1.1.1...v1.2.0
[1.1.1]: https://github.com/kjm99d/MonkeyPlanner/compare/v1.1.0...v1.1.1
[1.1.0]: https://github.com/kjm99d/MonkeyPlanner/compare/v1.0.0...v1.1.0
[1.0.0]: https://github.com/kjm99d/MonkeyPlanner/releases/tag/v1.0.0
