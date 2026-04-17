package contract_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/kjm99d/monkey-planner/backend/internal/storage"
	"github.com/kjm99d/monkey-planner/backend/internal/storage/contract"
	_ "github.com/kjm99d/monkey-planner/backend/internal/storage/postgres"
	_ "github.com/kjm99d/monkey-planner/backend/internal/storage/sqlite"
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
		t.Log("MP_PG_DSN unset — skipping PostgreSQL contract suite (SQLite-only default)")
		t.Skip("skip: set MP_PG_DSN to run PG contract tests")
	}
	contract.RunAll(t, func(t *testing.T) storage.Repo {
		t.Helper()
		r, err := storage.NewRepo(dsn)
		if err != nil {
			t.Fatalf("postgres open: %v", err)
		}
		t.Cleanup(func() {
			// Raw DB access for TRUNCATE is awkward here, so just Close — CI
			// should use a fresh schema per run anyway.
			_ = r.Close()
		})
		return r
	})
}
