package sqlite

import (
	"strings"

	"github.com/kjm99d/monkey-planner/backend/internal/storage"
)

func init() {
	// Accept both sqlite://<path> and file:<path>.
	open := func(dsn string) (storage.Repo, error) {
		path := strings.TrimPrefix(dsn, "sqlite://")
		path = strings.TrimPrefix(path, "file:")
		if path == "" {
			path = "./data/monkey.db"
		}
		return Open(path)
	}
	storage.Register("sqlite", open)
	storage.Register("file", open)
}
