package contract_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/ckmdevb/monkey-planner/backend/internal/storage"
	"github.com/ckmdevb/monkey-planner/backend/internal/storage/contract"
	_ "github.com/ckmdevb/monkey-planner/backend/internal/storage/postgres"
	_ "github.com/ckmdevb/monkey-planner/backend/internal/storage/sqlite"
)

func TestContract_SQLite(t *testing.T) {
	contract.RunAll(t, func(t *testing.T) storage.Repo {
		t.Helper()
		dsn := "sqlite://" + filepath.Join(t.TempDir(), "contract.db")
		r, err := storage.NewRepo(dsn)
		if err != nil {
			t.Fatalf("sqlite open: %v", err)
		}
		t.Cleanup(func() { _ = r.Close() })
		return r
	})
}

func TestContract_Postgres(t *testing.T) {
	dsn := os.Getenv("MP_PG_DSN")
	if dsn == "" {
		t.Log("MP_PG_DSN 미설정 → PostgreSQL 계약 테스트 skip (로컬 단일 사용자 스펙 기본 동작)")
		t.Skip("skip: set MP_PG_DSN to run PG contract tests")
	}
	contract.RunAll(t, func(t *testing.T) storage.Repo {
		t.Helper()
		r, err := storage.NewRepo(dsn)
		if err != nil {
			t.Fatalf("postgres open: %v", err)
		}
		t.Cleanup(func() {
			// TRUNCATE 를 위해 원시 DB에 접근하기 어려우므로 단순 Close.
			// CI에서 각 test run마다 새 스키마 사용 권장.
			_ = r.Close()
		})
		return r
	})
}
