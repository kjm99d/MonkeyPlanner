package storage

import (
	"fmt"
	"strings"
)

// Opener opens a storage adapter from a DSN string.
// Supported DSN schemes (registered at init time by each adapter package):
//   - "sqlite://<path>" or "file:<path>" → SQLite (modernc.org/sqlite)
//   - "postgres://user:pass@host:port/db?opts" → PostgreSQL (pgx/v5)
//
// Adapter packages Register themselves from init() to avoid an import cycle
// between this factory and the concrete implementations.
type Opener func(dsn string) (Repo, error)

var openers = map[string]Opener{}

// Register associates a DSN scheme with an Opener. Adapter packages call this
// from their init() function.
func Register(scheme string, o Opener) {
	openers[scheme] = o
}

// NewRepo dispatches to the Opener registered for the DSN's scheme.
func NewRepo(dsn string) (Repo, error) {
	scheme := dsnScheme(dsn)
	o, ok := openers[scheme]
	if !ok {
		return nil, fmt.Errorf("storage: unsupported DSN scheme %q (registered: %v)", scheme, keys(openers))
	}
	return o(dsn)
}

func dsnScheme(dsn string) string {
	if i := strings.Index(dsn, "://"); i > 0 {
		return dsn[:i]
	}
	// file: prefix without //
	if strings.HasPrefix(dsn, "file:") {
		return "file"
	}
	return ""
}

func keys(m map[string]Opener) []string {
	out := make([]string, 0, len(m))
	for k := range m {
		out = append(out, k)
	}
	return out
}
