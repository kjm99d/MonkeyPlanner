package postgres

import "github.com/ckmdevb/monkey-planner/backend/internal/storage"

func init() {
	storage.Register("postgres", func(dsn string) (storage.Repo, error) {
		return Open(dsn)
	})
}
