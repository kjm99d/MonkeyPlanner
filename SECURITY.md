# Security Policy

MonkeyPlanner is a local-first task manager that — by design — runs on the
user's own machine or a private network. That said, the HTTP API and MCP
bridge can both be misused if a vulnerability slips in. We take reports
seriously.

## Supported versions

Only the **latest tagged release** receives security fixes. Please upgrade
before filing a report unless you are reproducing on `main`.

| Version | Supported |
|---------|-----------|
| Latest tagged `vX.Y.Z` | ✅ |
| Older releases | ❌ (upgrade to latest) |
| `main` branch | ✅ (best effort) |

## Reporting a vulnerability

**Do NOT open a public GitHub issue for security problems.**

Instead, email **kjm99d@naver.com** with:

1. A short summary of the class of issue (e.g. "SSRF in webhook delivery",
   "unauthenticated write to /api/issues").
2. Steps to reproduce — ideally a single `curl` command or minimal script.
3. Affected version (`monkey-planner --version` or commit SHA).
4. Your assessment of impact if any.
5. Whether you want to be credited in the fix commit / release notes.

Optional but appreciated: a suggested patch or mitigation.

## What to expect

| Step | Target time |
|------|-------------|
| Acknowledgement of receipt | **72 hours** |
| First triage + severity decision | 7 days |
| Fix (critical / high) | 14 days from triage |
| Coordinated disclosure | after fix ships |

If you do not get an acknowledgement within 72 hours, open an unrelated
GitHub issue pinging the maintainer (do **not** mention the vulnerability
there) — email delivery sometimes fails silently.

## Scope

In scope:

- Code in this repository (`backend/`, `frontend/`, CI configs)
- The published Docker images on GHCR
- The `mcp install` path-handling and config-merge logic

Out of scope:

- Social engineering, physical access
- Third-party services you chose to connect (Discord, Slack, Telegram webhooks)
- Vulnerabilities in upstream dependencies that already have CVEs — please
  report those to the upstream project (we will update our lock files once
  a patched version exists)

Thanks for keeping MonkeyPlanner safe for everyone.
