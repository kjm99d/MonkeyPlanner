package http

import (
	"io/fs"
	"net/http"
	"strings"
)

// SPAHandler는 static fs에서 파일을 찾고, 없으면 index.html(SPA 진입)로 돌려 보냅니다.
// /api 경로는 이 핸들러 바깥의 chi 라우터가 먼저 처리합니다.
func SPAHandler(staticFS fs.FS) http.Handler {
	fileServer := http.FileServer(http.FS(staticFS))
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		upath := strings.TrimPrefix(r.URL.Path, "/")
		if upath == "" {
			upath = "index.html"
		}
		if _, err := fs.Stat(staticFS, upath); err != nil {
			// 존재하지 않는 경로 → index.html 로 재작성 (SPA 라우팅)
			r2 := r.Clone(r.Context())
			r2.URL.Path = "/"
			fileServer.ServeHTTP(w, r2)
			return
		}
		fileServer.ServeHTTP(w, r)
	})
}
