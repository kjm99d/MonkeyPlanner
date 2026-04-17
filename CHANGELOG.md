# Changelog

All notable changes to MonkeyPlanner are documented here.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Added
- **Docker** — `Dockerfile` (multi-stage, distroless/static, ~20 MB) and
  `docker-compose.yml` for one-command setup, plus a multi-arch (amd64+arm64)
  GHCR publish workflow at `.github/workflows/docker.yml`.
- **`monkey-planner mcp install`** subcommand that writes a MonkeyPlanner
  MCP server entry into Claude Code, Claude Desktop, or Cursor — no more
  hand-editing `.mcp.json` paths. `--dry-run`, `--force`, `--scope user`,
  `--base-url`, and `--name` are supported.
- **First-run bootstrap** — new DBs get a "Welcome" board and a single
  demo issue that teaches the approve-then-agent flow. Idempotent.
- **SSE heartbeat** — 30-second comment lines keep streams alive through
  reverse proxies (nginx / CloudFront / Cloudflare idle timeouts).
- **Graceful shutdown** — SIGINT/SIGTERM now triggers `http.Server.Shutdown`
  with a 15-second drain window. Safe for Docker/k8s rolling updates.
- **Server timeouts** — `ReadHeaderTimeout: 10s`, `ReadTimeout: 30s`,
  `IdleTimeout: 120s` (Slowloris mitigation, gosec G112).
- **Atomic properties merge** — new `IssueRepo.MergeProperties` uses SQLite
  `json_patch()` / Postgres `jsonb || patch + jsonb_strip_nulls` at the SQL
  level. Concurrent MCP clients writing different property keys no longer
  race.
- **SQLite concurrency** — `busy_timeout=5000ms` and `_txlock=immediate` on
  the DSN; `MaxOpenConns=1` at the pool. Multi-client MCP setups (Claude
  Code + Cursor + Claude Desktop) no longer hit spurious `SQLITE_BUSY`.
- **Custom Monkey brand** — replaced the Lucide `Squirrel` icon with an
  inline `MonkeyLogo` SVG that inherits `currentColor`.
- **`CONTRIBUTING.md`**, **`CODE_OF_CONDUCT.md`**, **`SECURITY.md`**.

### Changed
- English display name unified to **`MonkeyPlanner`** (was "Monkey Planner"
  with a space in some surfaces). Kebab form `monkey-planner` stays for the
  repo slug, Go module path, binary, and Docker image. Localized names in
  `ko/ja/zh` are kept.
- All Go-layer comments translated from Korean to English (domain → storage
  → service → http → events → web → cmd). Korean test-data strings
  (`한글`, `첫 작업`, `내용`, `테스트`, `새 제목`, `본문 수정`) are preserved
  because they exercise multibyte UTF-8 round-trips.

### Fixed
- `issue.criteria` / `issue.addCriterion` and four `board.*` / `chart.*` i18n
  keys now exist in all four locales (en/ko/ja/zh).

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

[Unreleased]: https://github.com/kjm99d/MonkeyPlanner/compare/v1.3.1...HEAD
[1.3.1]: https://github.com/kjm99d/MonkeyPlanner/compare/v1.3.0...v1.3.1
[1.3.0]: https://github.com/kjm99d/MonkeyPlanner/compare/v1.2.0...v1.3.0
[1.2.0]: https://github.com/kjm99d/MonkeyPlanner/compare/v1.1.1...v1.2.0
[1.1.1]: https://github.com/kjm99d/MonkeyPlanner/compare/v1.1.0...v1.1.1
[1.1.0]: https://github.com/kjm99d/MonkeyPlanner/compare/v1.0.0...v1.1.0
[1.0.0]: https://github.com/kjm99d/MonkeyPlanner/releases/tag/v1.0.0
