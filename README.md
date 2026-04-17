**English** | [한국어](./README.ko.md) | [日本語](./README.ja.md) | [中文](./README.zh.md)

<p align="center">
  <a href="https://github.com/kjm99d/monkey-planner/stargazers"><img alt="GitHub stars" src="https://img.shields.io/github/stars/kjm99d/monkey-planner?style=flat&logo=github"></a>
  <a href="https://github.com/kjm99d/monkey-planner/releases/latest"><img alt="Latest release" src="https://img.shields.io/github/v/release/kjm99d/monkey-planner?include_prereleases&sort=semver"></a>
  <a href="./LICENSE"><img alt="License: MIT" src="https://img.shields.io/badge/license-MIT-blue.svg"></a>
  <a href="https://github.com/kjm99d/monkey-planner/pkgs/container/monkey-planner"><img alt="Docker image" src="https://img.shields.io/badge/docker-ghcr.io-2496ed?logo=docker&logoColor=white"></a>
  <a href="https://github.com/kjm99d/monkey-planner/actions"><img alt="CI" src="https://img.shields.io/github/actions/workflow/status/kjm99d/monkey-planner/docker.yml?branch=main&label=build"></a>
</p>

# MonkeyPlanner

> **Local-first task memory for your AI coding agents.**
> Approve with a click; your agents do the rest. No cloud. No telemetry. Forever free, forever MIT.

Works with **Claude Code** · **Claude Desktop** · **Cursor** · **Continue** · any MCP-compatible client.

![MonkeyPlanner Demo](./docs/demo/monkey-planner-demo.gif)

## Quickstart

```bash
# Docker (recommended)
docker run -p 8080:8080 -v $(pwd)/data:/data ghcr.io/kjm99d/monkey-planner:latest

# then wire up your agent
monkey-planner mcp install --for claude-code     # or --for cursor / --for claude-desktop
```

Open http://localhost:8080 — the built-in Welcome board walks you through the rest.

## Features

### Issue & Board Management
- **Kanban Board** — Drag and drop, horizontal scroll, filtering, sorting, and table view toggle
- **Issue Creation** — Title, markdown body, and custom properties
- **Custom Properties** — Six supported types:
  - Text
  - Number
  - Select
  - Multi-select
  - Date
  - Checkbox

### Approval Gate
- **Pending → Approved** via a dedicated approval endpoint (cannot be done via generic PATCH)
- **Approval Queue** — Bulk-approve all Pending issues across boards
- **Approved → InProgress → Done** — Flexible status transitions
- **Rejected status** — Record a rejection reason

### Agent Features
- **Agent Instructions field** — Provide detailed instructions for MCP agents to follow
- **Success Criteria** — Manage completion conditions as a checklist
- **Comments** — Log progress and communicate per issue
- **Dependencies** — Express blocking relationships between issues

### Data Visualization
- **Calendar** — Monthly grid + daily activity (created, approved, completed counts)
- **Dashboard** — Stats cards + weekly activity chart
- **Sidebar** — Board list, issue counts, and recent items

### User Experience
- **Global Search** — Quick search with Cmd+K
- **Keyboard Shortcuts**
  - `h` — Go to dashboard
  - `a` — Go to approval queue
  - `?` — Show shortcut help
  - `Cmd+S` — Save
  - `Escape` — Close modal/dialog
- **Collapsible Sidebar** — Maximize screen space
- **Dark Mode** — Theme toggle
- **Internationalization** — Korean, English, Japanese, and Chinese

### Automation & Integrations
- **Webhooks** — Discord, Slack, and Telegram support
  - Events: `issue.created`, `issue.approved`, `issue.status_changed`, `issue.updated`, `issue.deleted`, `comment.created`
- **Real-time UI sync (SSE)** — Changes via MCP/CLI automatically reflect in open browser tabs, no refresh needed
- **JSON Export** — Export all issue data
- **Right-click Context Menu** — Quick actions
- **Issue Templates** — Per-board localStorage persistence

### MCP Server (AI Agent Integration)
Thirteen tools for AI agent automation:
1. `list_boards` — List all boards
2. `list_issues` — Query issues (filter by boardId, status)
3. `get_issue` — Issue detail including instructions, criteria, and comments
4. `create_issue` — Create a new issue
5. `approve_issue` — Approve: Pending → Approved
6. `claim_issue` — Claim: Approved → InProgress
7. `submit_qa` — Submit for QA: InProgress → QA
8. `complete_issue` — Complete: QA → Done (optional comment)
9. `reject_issue` — Reject: QA → InProgress with required reason
10. `add_comment` — Add a comment to an issue
11. `update_criteria` — Check or uncheck a success criterion
12. `search_issues` — Search issues by title
13. `get_version` — Get the MCP server version (for diagnostics)

## Tech Stack

### Backend
- **Language**: Go 1.26
- **Router**: chi/v5
- **Database**: SQLite / PostgreSQL (configurable)
- **Migrations**: goose/v3
- **Embedded files**: embed.FS (single-binary deployment)

### Frontend
- **Framework**: React 18
- **Language**: TypeScript
- **Bundler**: Vite 6
- **CSS**: Tailwind CSS
- **State management**: React Query (TanStack)
- **Drag and drop**: @dnd-kit/core, @dnd-kit/sortable
- **Icons**: lucide-react
- **Charts**: recharts
- **i18n**: react-i18next
- **Markdown**: react-markdown + rehype-sanitize

### MCP
- Protocol: JSON-RPC 2.0 over stdio
- Targets: Claude Code, Claude Desktop

## Getting Started

### Requirements
- Go 1.26 or later
- Node.js 18 or later
- npm or yarn

### Installation & Running

#### 1. Clone and initialize
```bash
git clone https://github.com/kjm99d/monkey-planner.git
cd monkey-planner
make init
```

#### 2. Production build (single binary)
```bash
make build
./bin/monkey-planner
```

The server runs at `http://localhost:8080` with the frontend embedded.

#### 3. Development mode (separate processes)

Terminal 1 — backend:
```bash
make run-backend
```

Terminal 2 — frontend (Vite dev server, :5173):
```bash
make run-frontend
```

The frontend automatically proxies `/api` requests to `:8080`.

### Environment Variables

```bash
# Server address (default: :8080)
export MP_ADDR=":8080"

# Database connection string
# SQLite (default: sqlite://./data/monkey.db)
export MP_DSN="sqlite://./data/monkey.db"

# PostgreSQL example
export MP_DSN="postgres://user:password@localhost:5432/monkey_planner"
```

## MCP Server Setup

### Quick Setup (Auto-Update)

1. Download the latest binary from [Releases](https://github.com/kjm99d/MonkeyPlanner/releases) to a directory (e.g. `D:/mp/`)
2. Download `update-and-run.sh` from the repo to the same directory
3. Add to `.mcp.json` or Claude Code settings:

```json
{
  "mcpServers": {
    "monkey-planner": {
      "command": "bash",
      "args": ["/path/to/update-and-run.sh", "mcp"],
      "env": {
        "MP_DSN": "sqlite:///path/to/data/monkey.db",
        "MP_BASE_URL": "http://localhost:8080"
      }
    }
  }
}
```

**Windows (without bash)**: Use `update-and-run.bat` instead:

```json
{
  "mcpServers": {
    "monkey-planner": {
      "command": "D:/mp/update-and-run.bat",
      "args": ["mcp"],
      "env": {
        "MP_DSN": "sqlite://D:/mp/data/monkey.db",
        "MP_BASE_URL": "http://localhost:8080"
      }
    }
  }
}
```

The wrapper script automatically checks for new releases on GitHub and updates the binary before each launch.

### Manual Setup

If you prefer not to use auto-update:

```json
{
  "mcpServers": {
    "monkey-planner": {
      "command": "/path/to/monkey-planner",
      "args": ["mcp"],
      "env": {
        "MP_DSN": "sqlite:///path/to/data/monkey.db",
        "MP_BASE_URL": "http://localhost:8080"
      }
    }
  }
}
```

This config works for both Claude Code (`.mcp.json`) and Claude Desktop (`~/.claude/claude_desktop_config.json`). Restart after changes.

### MCP Tool Usage Examples

```
AI: List all boards
→ list_boards()

AI: Find issues related to "authentication"
→ search_issues(query="authentication")

AI: Approve the first pending issue, claim it, work on it, and submit for QA
→ approve_issue() → claim_issue() → submit_qa()
```

## Workflow — Real Usage Scenario

Below is a real workflow from fixing a language switcher bug, showing how a human and AI agent collaborate through MonkeyPlanner.

### Status Flow

```
Pending → Approved → InProgress → QA → Done
                         ↑              │ (reject with reason)
                         └──────────────┘
```

### Step-by-Step

**1. Create Issue** — Human finds a bug, asks AI to register it
```
Human: "The language selector dropdown doesn't appear when clicking the button. Create an issue."
AI:    create_issue(boardId, title, body, instructions)  →  status: Pending
```

**2. Approve** — Human reviews and approves
```
Human: (clicks Approve on the board or tells AI)
AI:    approve_issue(issueId)  →  status: Approved
```

**3. Start Work** — AI claims the issue and begins coding
```
AI:    claim_issue(issueId)  →  status: InProgress
       - Reads code, identifies root cause
       - Implements fix, runs tests
       - Commits changes
```

**4. Submit for QA** — AI finishes and submits for review
```
AI:    submit_qa(issueId, comment: "commit abc1234 — fixed click handler")
       →  status: QA
       add_comment(issueId, "Commit info: ...")
```

**5. Review** — Human tests the fix
```
Human: Tests in browser, finds the dropdown is clipped by sidebar
       →  reject_issue(issueId, reason: "Dropdown is hidden behind sidebar")
       →  status: InProgress  (back to step 3)

Human: Tests again after fix, everything works
       →  complete_issue(issueId)  →  status: Done
```

**6. Feedback Loop** — Communication via comments throughout
```
Human: add_comment("Dropdown is clipped on the left side, fix it")
AI:    get_issue() → reads comment → fixes → commit → submit_qa()
Human: Tests → complete_issue()  →  Done ✓
```

### Key Takeaways

- **Human controls the gates**: Approve, QA pass/reject, Complete
- **AI does the work**: Code analysis, implementation, testing, commits
- **Comments are the communication channel**: Both sides use `add_comment` to exchange feedback
- **QA loop prevents premature completion**: Issues must pass human review before Done

## API Reference

OpenAPI 3.0 spec: [backend/docs/swagger.yaml](./backend/docs/swagger.yaml)

### Key Endpoints

#### Boards
```
GET    /api/boards                  # List boards
POST   /api/boards                  # Create board
PATCH  /api/boards/{id}             # Update board
DELETE /api/boards/{id}             # Delete board
```

#### Issues
```
GET    /api/issues                  # List issues (filter: boardId, status, parentId)
POST   /api/issues                  # Create issue
GET    /api/issues/{id}             # Issue detail + child issues
PATCH  /api/issues/{id}             # Update issue (status, properties, title, etc.)
DELETE /api/issues/{id}             # Delete issue
POST   /api/issues/{id}/approve     # Approve issue (Pending → Approved)
```

#### Comments
```
GET    /api/issues/{issueId}/comments    # List comments
POST   /api/issues/{issueId}/comments    # Add comment
DELETE /api/comments/{commentId}         # Delete comment
```

#### Properties (Custom Attributes)
```
GET    /api/boards/{boardId}/properties      # List property definitions
POST   /api/boards/{boardId}/properties      # Create property
PATCH  /api/boards/{boardId}/properties/{propId}  # Update property
DELETE /api/boards/{boardId}/properties/{propId}  # Delete property
```

#### Webhooks
```
GET    /api/boards/{boardId}/webhooks           # List webhooks
POST   /api/boards/{boardId}/webhooks           # Create webhook
PATCH  /api/boards/{boardId}/webhooks/{whId}    # Update webhook
DELETE /api/boards/{boardId}/webhooks/{whId}    # Delete webhook
```

#### Calendar
```
GET /api/calendar           # Monthly stats (year, month required)
GET /api/calendar/day       # Daily issue list (date required)
```

For full schema details, see [backend/docs/swagger.yaml](./backend/docs/swagger.yaml).

## Project Structure

```
monkey-planner/
├── backend/
│   ├── cmd/monkey-planner/
│   │   ├── main.go              # Entry point (HTTP server)
│   │   └── mcp.go               # MCP server (JSON-RPC stdio)
│   ├── internal/
│   │   ├── domain/              # Domain models (Issue, Board, etc.)
│   │   ├── service/             # Business logic
│   │   ├── storage/             # Database layer (SQLite/PostgreSQL)
│   │   ├── http/                # HTTP handlers & router
│   │   └── migrations/          # goose migration files
│   ├── web/                     # Embedded frontend (embed.FS)
│   ├── docs/
│   │   └── swagger.yaml         # OpenAPI 3.0 spec
│   ├── go.mod
│   └── go.sum
│
├── frontend/
│   ├── src/
│   │   ├── components/          # Reusable components
│   │   ├── features/            # Page & feature components
│   │   │   ├── home/           # Dashboard
│   │   │   ├── board/          # Board & Kanban
│   │   │   ├── issue/          # Issue detail
│   │   │   ├── calendar/       # Calendar
│   │   │   └── approval/       # Approval queue
│   │   ├── api/                 # API hooks & client
│   │   ├── design/              # Tailwind tokens
│   │   ├── i18n/                # Translations (en.json, ko.json, ja.json, zh.json)
│   │   ├── App.tsx              # Router
│   │   ├── index.css            # Global styles
│   │   └── main.tsx
│   ├── package.json
│   ├── vite.config.ts
│   ├── tsconfig.json
│   └── tailwind.config.js
│
├── .mcp.json                    # Claude Code MCP config
├── Makefile                     # Build & dev commands
├── .githooks/                   # Git hooks
└── data/                        # SQLite database (default)
```

## Testing

### Backend tests
```bash
make test-backend
```

### Frontend tests
```bash
make test-frontend
```

### Accessibility tests
```bash
make test-a11y
```

### All tests
```bash
make test
```

## Common Commands

```bash
# Initial setup after cloning
make init

# Production build
make build

# Run production server
./bin/monkey-planner

# Development mode
make run-backend        # Terminal 1
make run-frontend       # Terminal 2

# Clean build artifacts
make clean
```

## Status Transition Rules

```
Pending
  ↓ (approve endpoint)
Approved
  ↓ (PATCH status)
InProgress
  ↓ (PATCH status)
Done

Pending → Approved: POST /api/issues/{id}/approve (dedicated endpoint only)
Approved ↔ InProgress ↔ Done: Free transitions via PATCH
Pending: Cannot be re-entered from other statuses
Rejected: Separate rejection state with reason tracking
```

## License

MIT

## Contributing

Issues and pull requests are welcome.

## Contact

For questions or feedback about the project, please open a GitHub Issue.
