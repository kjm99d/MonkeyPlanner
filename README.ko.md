[English](./README.md) | **한국어** | [日本語](./README.ja.md) | [中文](./README.zh.md)

# 몽키 플래너 (Monkey Planner)

> AI 에이전트 작업 기억 저장소 — Notion/JIRA 스타일 이슈 트래커 + MCP 서버

사람이 이슈를 생성하고 승인하면, AI 에이전트가 MCP(Model Context Protocol) 클라이언트를 통해 작업을 소비하는 협업 도구입니다.

![Monkey Planner](./docs/screenshots/d-home-l.png)

## 주요 기능

### 이슈 & 보드 관리
- **칸반 보드** — 드래그 앤 드롭, 수평 스크롤, 필터링, 정렬, 테이블 뷰 전환
- **이슈 생성** — 제목, 마크다운 본문, 커스텀 속성 지원
- **커스텀 속성** — 6가지 타입 지원
  - 텍스트 (text)
  - 숫자 (number)
  - 선택 (select)
  - 다중 선택 (multi_select)
  - 날짜 (date)
  - 체크박스 (checkbox)

### 승인 게이트 (Approval Flow)
- **Pending → Approved** 전용 승인 엔드포인트 (일반 PATCH로는 불가)
- **Approval Queue** — 전체 보드의 Pending 이슈 일괄 승인
- **Approved → InProgress → Done** — 자유로운 상태 전환
- **Rejected 상태** — 거절 이유 기록 가능

### 에이전트 기능
- **에이전트 지시사항 필드** — MCP 에이전트가 참조할 구체적 지시사항 입력
- **성공 기준** — 완료 조건을 체크리스트로 관리
- **댓글** — 이슈별 진행 상황 기록 및 소통
- **의존성** — 이슈 간 차단 관계 표현

### 데이터 시각화
- **캘린더** — 월간 그리드 + 일별 실적 (생성, 승인, 완료 카운트)
- **대시보드** — 통계 카드 + 주간 활동 차트
- **사이드바** — 보드 목록, 이슈 카운트, 최근 항목

### 사용자 경험
- **글로벌 검색** — Cmd+K로 빠른 검색
- **키보드 단축키**
  - `h` — 대시보드로 이동
  - `a` — 승인 큐로 이동
  - `?` — 단축키 도움말
  - `Cmd+S` — 저장
  - `Escape` — 모달/다이얼로그 닫기
- **사이드바 접기/펼치기** — 화면 공간 최적화
- **다크 모드** — 테마 전환
- **다국어** — 한국어, 영어, 일본어, 중국어 지원

### 자동화 & 연동
- **Webhook** — Discord, Slack, Telegram 지원
  - 이벤트: `issue.created`, `issue.approved`, `issue.status_changed`, `issue.deleted`
- **JSON 내보내기** — 모든 이슈 데이터 내보내기
- **우클릭 컨텍스트 메뉴** — 빠른 작업 메뉴
- **이슈 템플릿** — 보드별 localStorage 저장

### MCP 서버 (AI 에이전트 연동)
10가지 도구로 AI 에이전트 자동화:
1. `list_boards` — 모든 보드 조회
2. `list_issues` — 이슈 조회 (boardId, status 필터)
3. `get_issue` — 이슈 상세 (지시사항, 기준, 댓글 포함)
4. `create_issue` — 새 이슈 생성
5. `approve_issue` — Pending → Approved 승인
6. `claim_issue` — Approved → InProgress 전환
7. `complete_issue` — InProgress → Done 완료 (선택 댓글)
8. `add_comment` — 이슈에 댓글 추가
9. `update_criteria` — 성공 기준 체크/언체크
10. `search_issues` — 제목 기반 이슈 검색

## 기술 스택

### Backend
- **언어**: Go 1.26
- **라우터**: chi/v5
- **데이터베이스**: SQLite / PostgreSQL (선택 가능)
- **마이그레이션**: goose/v3
- **내장 파일**: embed.FS (단일 바이너리 배포)

### Frontend
- **프레임워크**: React 18
- **언어**: TypeScript
- **번들러**: Vite 6
- **CSS**: Tailwind CSS
- **상태 관리**: React Query (TanStack)
- **드래그**: @dnd-kit/core, @dnd-kit/sortable
- **아이콘**: lucide-react
- **차트**: recharts
- **i18n**: react-i18next
- **마크다운**: react-markdown + rehype-sanitize

### MCP
- 프로토콜: JSON-RPC 2.0 over stdio
- 대상: Claude Code, Claude Desktop

## 시작하기

### 요구사항
- Go 1.26 이상
- Node.js 18 이상
- npm 또는 yarn

### 설치 & 실행

#### 1. 저장소 클론 및 초기화
```bash
git clone https://github.com/ckmdevb/monkey-planner.git
cd monkey-planner
make init
```

#### 2. 프로덕션 빌드 (단일 바이너리)
```bash
make build
./bin/monkey-planner
```

서버는 `http://localhost:8080`에서 실행되며, 프론트엔드는 내장됩니다.

#### 3. 개발 모드 (분리 실행)

터미널 1 — 백엔드:
```bash
make run-backend
```

터미널 2 — 프론트엔드 (Vite dev server, :5173):
```bash
make run-frontend
```

프론트엔드는 자동으로 `/api` 요청을 `:8080`으로 프록시합니다.

### 환경 변수

```bash
# 서버 주소 (기본값: :8080)
export MP_ADDR=":8080"

# 데이터베이스 연결 문자열
# SQLite (기본값: sqlite://./data/monkey.db)
export MP_DSN="sqlite://./data/monkey.db"

# PostgreSQL 예시
export MP_DSN="postgres://user:password@localhost:5432/monkey_planner"
```

## MCP 서버 설정

### Claude Code에서 사용

`.mcp.json` 파일이 프로젝트 루트에 있습니다:

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

**Windows 사용자**: 경로를 자신의 환경에 맞게 수정하세요.

### Claude Desktop에서 사용

Claude Desktop의 설정 파일(`~/.claude/claude_desktop_config.json`)에 추가:

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

그후 Claude Desktop을 재시작하면 Monkey Planner 도구들이 자동으로 로드됩니다.

### MCP 도구 사용 예시

```
AI: 모든 보드를 나열해주세요
→ list_boards() 호출

AI: "인증" 관련 이슈를 찾아주세요
→ search_issues(query="인증") 호출

AI: 첫 번째 Pending 이슈를 승인하고 진행 중으로 전환한 후 완료하겠습니다
→ approve_issue() → claim_issue() → complete_issue() 순차 호출
```

## 에이전트 작업 플로우

```
┌────────────────┐
│  사람이 이슈   │  제목, 본문, 지시사항 입력
└────────┬───────┘
         │
         ↓
┌────────────────┐
│  Approve 버튼  │  Pending → Approved
└────────┬───────┘
         │
         ↓
┌────────────────────────────┐
│  AI 에이전트 (MCP 클라이언트) │  list_issues 또는 search_issues
└────────┬───────────────────┘
         │
         ↓
┌────────────────────┐
│ claim_issue()      │  Approved → InProgress
└────────┬───────────┘
         │
         ↓
┌────────────────────┐
│ 작업 진행 중...    │  add_comment(), update_criteria()
│                    │  (상황 보고 & 기준 체크)
└────────┬───────────┘
         │
         ↓
┌────────────────────┐
│ complete_issue()   │  InProgress → Done
│ + 최종 댓글        │
└────────┬───────────┘
         │
         ↓
┌────────────────┐
│  사람이 확인   │  결과 검토 및 피드백
└────────────────┘
```

## API 문서

OpenAPI 3.0 스펙: [backend/docs/swagger.yaml](./backend/docs/swagger.yaml)

### 주요 엔드포인트

#### Boards
```
GET    /api/boards                  # 보드 목록
POST   /api/boards                  # 보드 생성
PATCH  /api/boards/{id}             # 보드 수정
DELETE /api/boards/{id}             # 보드 삭제
```

#### Issues
```
GET    /api/issues                  # 이슈 목록 (필터: boardId, status, parentId)
POST   /api/issues                  # 이슈 생성
GET    /api/issues/{id}             # 이슈 상세 + 자식 이슈
PATCH  /api/issues/{id}             # 이슈 수정 (상태, 속성, 제목 등)
DELETE /api/issues/{id}             # 이슈 삭제
POST   /api/issues/{id}/approve     # 이슈 승인 (Pending → Approved)
```

#### Comments
```
GET    /api/issues/{issueId}/comments    # 댓글 목록
POST   /api/issues/{issueId}/comments    # 댓글 추가
DELETE /api/comments/{commentId}         # 댓글 삭제
```

#### Properties (커스텀 속성)
```
GET    /api/boards/{boardId}/properties      # 속성 정의 목록
POST   /api/boards/{boardId}/properties      # 속성 생성
PATCH  /api/boards/{boardId}/properties/{propId}  # 속성 수정
DELETE /api/boards/{boardId}/properties/{propId}  # 속성 삭제
```

#### Webhooks
```
GET    /api/boards/{boardId}/webhooks       # Webhook 목록
POST   /api/boards/{boardId}/webhooks       # Webhook 생성
PATCH  /api/boards/{boardId}/webhooks/{whId}    # Webhook 수정
DELETE /api/boards/{boardId}/webhooks/{whId}    # Webhook 삭제
```

#### Calendar
```
GET /api/calendar           # 월간 통계 (year, month 필수)
GET /api/calendar/day       # 일일 이슈 목록 (date 필수)
```

자세한 스키마는 [backend/docs/swagger.yaml](./backend/docs/swagger.yaml)을 참조하세요.

## 프로젝트 구조

```
monkey-planner/
├── backend/
│   ├── cmd/monkey-planner/
│   │   ├── main.go              # 엔트리 포인트 (HTTP 서버)
│   │   └── mcp.go               # MCP 서버 (JSON-RPC stdio)
│   ├── internal/
│   │   ├── domain/              # 도메인 모델 (Issue, Board, etc.)
│   │   ├── service/             # 비즈니스 로직
│   │   ├── storage/             # 데이터베이스 (SQLite/PostgreSQL)
│   │   ├── http/                # HTTP 핸들러 & 라우터
│   │   └── migrations/          # goose 마이그레이션 파일
│   ├── web/                     # 프론트엔드 임베드 (embed.FS)
│   ├── docs/
│   │   └── swagger.yaml         # OpenAPI 3.0 스펙
│   ├── go.mod
│   └── go.sum
│
├── frontend/
│   ├── src/
│   │   ├── components/          # 재사용 가능한 컴포넌트
│   │   ├── features/            # 페이지 & 기능별 컴포넌트
│   │   │   ├── home/           # 대시보드
│   │   │   ├── board/          # 보드 & 칸반
│   │   │   ├── issue/          # 이슈 상세
│   │   │   ├── calendar/       # 캘린더
│   │   │   └── approval/       # 승인 큐
│   │   ├── api/                 # API 훅 & 클라이언트
│   │   ├── design/              # Tailwind 토큰
│   │   ├── i18n/                # 다국어 (en.json, ko.json, ja.json, zh.json)
│   │   ├── App.tsx              # 라우터
│   │   ├── index.css            # 글로벌 스타일
│   │   └── main.tsx
│   ├── package.json
│   ├── vite.config.ts
│   ├── tsconfig.json
│   └── tailwind.config.js
│
├── .mcp.json                    # Claude Code MCP 설정
├── Makefile                     # 빌드 & 개발 명령어
├── .githooks/                   # Git hooks
└── data/                        # SQLite 데이터베이스 (기본값)
```

## 테스트

### 백엔드 테스트
```bash
make test-backend
```

### 프론트엔드 테스트
```bash
make test-frontend
```

### 접근성 테스트
```bash
make test-a11y
```

### 전체 테스트
```bash
make test
```

## 주요 명령어

```bash
# 새로 클론한 후 초기 설정
make init

# 프로덕션 빌드
make build

# 프로덕션 서버 실행
./bin/monkey-planner

# 개발 모드
make run-backend        # 터미널 1
make run-frontend       # 터미널 2

# 정리
make clean
```

## 상태 전환 규칙

```
Pending
  ↓ (approve endpoint)
Approved
  ↓ (PATCH status)
InProgress
  ↓ (PATCH status)
Done

Pending → Approved: POST /api/issues/{id}/approve (전용)
Approved ↔ InProgress ↔ Done: PATCH로 자유로운 전환
Pending: 다른 상태에서 돌아올 수 없음
Rejected: 거절 상태 (별도 추적)
```

## 라이선스

MIT

## 기여

이슈와 풀 리퀘스트는 환영합니다.

## 문의

프로젝트에 대한 질문이나 피드백은 GitHub Issues를 통해 제출해 주세요.
