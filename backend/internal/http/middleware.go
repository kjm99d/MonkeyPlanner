package http

import (
	"bytes"
	"io"
	"net/http"
	"unicode/utf8"
)

// ValidateUTF8 은 JSON 요청 본문이 유효한 UTF-8인지 검증합니다.
// 비-UTF-8 바이트가 포함되면 400 Bad Request를 반환합니다.
// Windows cp949 등 비-UTF-8 환경에서 curl로 한글을 보낼 때 발생하는 mojibake를 사전 차단합니다.
func ValidateUTF8(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Body == nil || r.ContentLength == 0 {
			next.ServeHTTP(w, r)
			return
		}
		ct := r.Header.Get("Content-Type")
		if ct == "" || !isJSON(ct) {
			next.ServeHTTP(w, r)
			return
		}

		body, err := io.ReadAll(r.Body)
		r.Body.Close()
		if err != nil {
			writeErr(w, http.StatusBadRequest, "read_error", "failed to read request body")
			return
		}

		if !utf8.Valid(body) {
			writeErr(w, http.StatusBadRequest, "invalid_encoding",
				"request body contains non-UTF-8 bytes; ensure your client sends UTF-8 encoded JSON")
			return
		}

		r.Body = io.NopCloser(bytes.NewReader(body))
		next.ServeHTTP(w, r)
	})
}

func isJSON(ct string) bool {
	return len(ct) >= 16 && ct[:16] == "application/json" ||
		len(ct) >= 4 && ct[:4] == "text"
}
