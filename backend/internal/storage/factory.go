package storage

import (
	"fmt"
	"strings"
)

// NewRepo 는 DSN 접두사로 어댑터를 선택하는 팩토리입니다.
// 지원 형식:
//   - "sqlite://<path>" 또는 "file:<path>" → SQLite 어댑터 (modernc.org/sqlite)
//   - "postgres://user:pass@host:port/db?opts" → PostgreSQL 어댑터 (pgx/v5)
//
// 실제 어댑터 구현은 import 순환을 피하기 위해 외부에서 NewFromDSN 을
// 등록(Register)하는 방식을 사용합니다.
type Opener func(dsn string) (Repo, error)

var openers = map[string]Opener{}

// Register 는 스킴별 어댑터 Opener 를 등록합니다. 각 어댑터 패키지가 init() 에서 호출합니다.
func Register(scheme string, o Opener) {
	openers[scheme] = o
}

// NewRepo 는 DSN 을 읽어 등록된 어댑터로 위임합니다.
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
