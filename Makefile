.PHONY: init run run-backend run-frontend build build-frontend build-backend test test-backend test-frontend test-a11y clean

# OS 감지 — 바이너리 확장자 결정
ifeq ($(OS),Windows_NT)
    BIN_EXT := .exe
else
    BIN_EXT :=
endif

BIN := bin/monkey-planner$(BIN_EXT)

## 신규 클론 후 한 번 실행: git hooks 경로 설정 + 프론트 의존성 설치
init:
	git config core.hooksPath .githooks
	@echo "✓ git hooks path set to .githooks"
	cd frontend && npm install
	@echo "✓ frontend npm install 완료. 이제 'make run'으로 개발 서버를 띄우세요."

## 개발 모드: 백엔드(8080) + 프론트(5173, /api 프록시)
run:
	@echo "백엔드·프론트를 별도 터미널에서 띄우려면 make run-backend / make run-frontend"
	@echo "혹은 단일 바이너리를 쓰려면 make build && ./$(BIN)"

run-backend:
	cd backend && go run ./cmd/monkey-planner

run-frontend:
	cd frontend && npm run dev

## 프로덕션 빌드: 프론트 → dist, 백엔드가 embed.FS로 내장
build: build-frontend build-backend

build-frontend:
	cd frontend && npm run build

build-backend:
	@mkdir -p bin
	cd backend && go build -tags prod -o ../$(BIN) ./cmd/monkey-planner

## 테스트
test: test-backend test-frontend

test-backend:
	cd backend && go test ./...

test-frontend:
	cd frontend && npm test -- --run

test-a11y:
	cd frontend && npm run test:a11y

## 정리
clean:
	rm -rf bin frontend/dist frontend/node_modules backend/bin
	@echo "✓ 빌드 산출물 제거"
