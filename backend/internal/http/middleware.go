package http

import (
	"bytes"
	"io"
	"net/http"
	"unicode/utf8"
)

// ValidateUTF8 rejects non-UTF-8 JSON request bodies with 400 Bad Request.
// This prevents mojibake from shells with non-UTF-8 encodings (e.g. Windows
// cp949) when users curl multibyte text through the API.
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
