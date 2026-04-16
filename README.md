# 몽키 플래너 (Monkey Planner)

> AI **코숭이(Code Monkey) 에이전트**들의 작업 기억을 이슈 단위로 기록·보존하는 단일 사용자 로컬 웹 도구.

여러 에이전트(Claude Code / Codex / MCP 클라이언트 등)가 끊기지 않고 컨텍스트를 인계받도록 돕는 "하드디스크" 역할을 한다. 사용자는 로컬 웹 UI에서 이슈를 생성·편집·**Approve**하고, 승인된 이슈는 (phase 2) MCP로 외부 에이전트가 소비한다.

## 문서
- 스펙: [`.omc/specs/deep-interview-monkey-planner.md`](./.omc/specs/deep-interview-monkey-planner.md)
- 플랜: [`.omc/plans/monkey-planner-plan.md`](./.omc/plans/monkey-planner-plan.md)
- PRD (이슈 목록): [`.omc/prd.json`](./.omc/prd.json)

## 스택
- 백엔드: Go + `chi` + `pressly/goose` (마이그레이션)
- 저장소: SQLite(기본) 또는 PostgreSQL (DB 어댑터로 스왑)
- 프론트엔드: React + TypeScript + Vite + Tailwind CSS + React Query + React Router
- 패키징: `embed.FS`로 프론트 번들을 Go 바이너리에 내장 → 단일 실행 파일

## 모노레포 구조
```
MonkeyPlanner/
├── .githooks/          # commit-msg 훅 (Conventional Commits + author 태그 금지)
├── backend/            # Go 백엔드
├── frontend/           # React 프론트엔드
├── .omc/               # OMC 메타(스펙·플랜·디자인)
└── Makefile            # init/run/build/test 편의 스크립트
```

## 빠른 시작
```bash
# 처음 한 번만 — git hooks 경로 설정 + npm 의존성 설치
make init

# 개발 모드: 백엔드(localhost:8080) + 프론트(localhost:5173, /api 프록시)
make run

# 프로덕션 빌드 (단일 바이너리)
make build
./bin/monkey-planner          # 또는 bin\monkey-planner.exe (Windows)
```

### DB 선택
환경 변수 `MP_DSN`으로 어댑터 지정:
```bash
# SQLite (기본값)
MP_DSN="sqlite://./data/monkey.db" ./bin/monkey-planner

# PostgreSQL (phase 2 튜닝 전까지는 스켈레톤 수준)
MP_DSN="postgres://user:pass@localhost:5432/monkey" ./bin/monkey-planner
```

## 테스트
```bash
make test           # go test + vitest
make test-a11y      # axe-core 자동 접근성 감사
```

## 커밋 규칙
- **Conventional Commits** 타입 + **한국어** 메시지 (예: `feat: 이슈 Approve 버튼 추가`)
- `Co-Authored-By` / `Signed-off-by` 등 **author 트레일러 금지**
- 하나의 논리 단위 = 하나의 커밋 (**atomic commit**)

`.githooks/commit-msg`가 위 규칙을 기계적으로 강제한다. `make init` 실행 시 자동 설치됨.

## 라이선스
내부 프로젝트 (라이선스 미정).
