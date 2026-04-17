package http

import (
	"io/fs"
	"net/http"
	"strings"
)

// SPAHandler serves files from the static FS, falling back to index.html for
// unknown paths so the SPA router can take over. The chi /api routes run
// ahead of this handler and are never hit here.
func SPAHandler(staticFS fs.FS) http.Handler {
	fileServer := http.FileServer(http.FS(staticFS))
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		upath := strings.TrimPrefix(r.URL.Path, "/")
		if upath == "" {
			upath = "index.html"
		}
		if _, err := fs.Stat(staticFS, upath); err != nil {
			// Unknown path → rewrite to index.html so the SPA router runs.
			r2 := r.Clone(r.Context())
			r2.URL.Path = "/"
			fileServer.ServeHTTP(w, r2)
			return
		}
		fileServer.ServeHTTP(w, r)
	})
}
