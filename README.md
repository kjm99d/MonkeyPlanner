**English** | [한국어](./README.ko.md) | [日本語](./README.ja.md) | [中文](./README.zh.md)

# Monkey Planner

> AI agent task memory store — Notion/JIRA-style issue tracker + MCP server

A collaborative tool where humans create and approve issues, and AI agents consume tasks via MCP (Model Context Protocol) clients.

![Monkey Planner Demo](./docs/demo/monkey-planner-demo.gif)

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
  - Events: `issue.created`, `issue.approved`, `issue.status_changed`, `issue.deleted`
- **JSON Export** — Export all issue data
- **Right-click Context Menu** — Quick actions
- **Issue Templates** — Per-board localStorage persistence

### MCP Server (AI Agent Integration)
Ten tools for AI agent automation:
1. `list_boards` — List all boards
2. `list_issues` — Query issues (filter by boardId, status)
3. `get_issue` — Issue detail including instructions, criteria, and comments
4. `create_issue` — Create a new issue
5. `approve_issue` — Approve: Pending → Approved
6. `claim_issue` — Claim: Approved → InProgress
7. `complete_issue` — Complete: InProgress → Done (optional comment)
8. `add_comment` — Add a comment to an issue
9. `update_criteria` — Check or uncheck a success criterion
10. `search_issues` — Search issues by title

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

### Using with Claude Code

A `.mcp.json` file is included in the project root:

```json
{
  "mcpServers": {
    "monkey-planner": {
      "command": "./bin/monkey-planner.exe",
      "args": ["mcp"],
      "cwd": "D:/Projects/MonkeyPlanner"
    }
  }
}
```

**Windows users**: Update the path to match your environment.

### Using with Claude Desktop

Add the following to Claude Desktop's config file (`~/.claude/claude_desktop_config.json`):

```json
{
  "mcpServers": {
    "monkey-planner": {
      "command": "/path/to/monkey-planner",
      "args": ["mcp"]
    }
  }
}
```

Restart Claude Desktop and the Monkey Planner tools will load automatically.

### MCP Tool Usage Examples

```
AI: List all boards
→ list_boards()

AI: Find issues related to "authentication"
→ search_issues(query="authentication")

AI: Approve the first pending issue, claim it, and complete it
→ approve_issue() → claim_issue() → complete_issue()
```

## Agent Workflow

```
┌────────────────┐
│  Human creates │  Enter title, body, and instructions
│  an issue      │
└────────┬───────┘
         │
         ↓
┌────────────────┐
│  Approve       │  Pending → Approved
└────────┬───────┘
         │
         ↓
┌────────────────────────────┐
│  AI Agent (MCP client)     │  list_issues or search_issues
└────────┬───────────────────┘
         │
         ↓
┌────────────────────┐
│ claim_issue()      │  Approved → InProgress
└────────┬───────────┘
         │
         ↓
┌────────────────────┐
│ Working...         │  add_comment(), update_criteria()
│                    │  (progress report & criteria checks)
└────────┬───────────┘
         │
         ↓
┌────────────────────┐
│ complete_issue()   │  InProgress → Done
│ + final comment    │
└────────┬───────────┘
         │
         ↓
┌────────────────┐
│  Human reviews │  Review results and provide feedback
└────────────────┘
```

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
