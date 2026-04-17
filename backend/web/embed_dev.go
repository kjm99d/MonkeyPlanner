//go:build !prod

package web

import (
	"errors"
	"io/fs"
)

// Dist always returns an error for dev builds. In dev mode the Vite dev
// server on :5173 serves the UI, so embedding the dist bundle is unnecessary.
func Dist() (fs.FS, error) {
	return nil, errors.New("web.Dist: dev build does not embed frontend (use Vite dev server on :5173)")
}
