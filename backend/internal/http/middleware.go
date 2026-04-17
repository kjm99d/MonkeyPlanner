package http

import (
	"bytes"
	"crypto/subtle"
	"io"
	"net/http"
	"strings"
	"unicode/utf8"
)

// RequireBearerToken returns a middleware that enforces a shared-secret
// Authorization header on every /api/* request. When token is empty the
// middleware returns the next handler unchanged (off-by-default stays the
// default — local single-user installs keep working without auth).
//
// When enabled, requests must carry `Authorization: Bearer <token>`. The
// comparison uses subtle.ConstantTimeCompare so rejection does not leak
// byte-by-byte timing.
//
// Health and SSE endpoints still require the token — they run under the
// same /api/* tree. Use a dedicated proxy if a subset needs to stay open.
func RequireBearerToken(token string) func(http.Handler) http.Handler {
	if token == "" {
		return func(next http.Handler) http.Handler { return next }
	}
	expected := []byte("Bearer " + token)
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			got := r.Header.Get("Authorization")
			if len(got) != len(expected) ||
				subtle.ConstantTimeCompare([]byte(got), expected) != 1 {
				w.Header().Set("WWW-Authenticate", `Bearer realm="monkey-planner"`)
				writeErr(w, http.StatusUnauthorized, "unauthorized",
					"this MonkeyPlanner requires an Authorization: Bearer token")
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}

// SecurityHeaders sets a conservative set of response headers on every
// response. The values are deliberately strict but compatible with a
// React SPA served from the same origin:
//
//   - X-Content-Type-Options: nosniff blocks MIME sniffing
//   - X-Frame-Options: DENY blocks clickjacking (also covered by CSP)
//   - Referrer-Policy: strict-origin-when-cross-origin limits leakage
//   - Permissions-Policy strips camera/microphone/geolocation defaults
//   - Content-Security-Policy allows self + inline styles (Tailwind) +
//     two Google Fonts CDNs the frontend index.html pulls from
//
// The CSP is skipped for /api/* responses because API clients never
// render HTML and the header adds noise to responses.
func SecurityHeaders(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		h := w.Header()
		h.Set("X-Content-Type-Options", "nosniff")
		h.Set("X-Frame-Options", "DENY")
		h.Set("Referrer-Policy", "strict-origin-when-cross-origin")
		h.Set("Permissions-Policy", "camera=(), microphone=(), geolocation=(), payment=()")

		if !strings.HasPrefix(r.URL.Path, "/api/") {
			// Matches what the current index.html actually loads. Tighten once
			// we drop the jsdelivr font dependency.
			h.Set("Content-Security-Policy",
				"default-src 'self'; "+
					"script-src 'self'; "+
					"style-src 'self' 'unsafe-inline' https://fonts.googleapis.com https://cdn.jsdelivr.net; "+
					"font-src 'self' https://fonts.gstatic.com https://cdn.jsdelivr.net; "+
					"img-src 'self' data: blob:; "+
					"connect-src 'self'; "+
					"frame-ancestors 'none'; "+
					"base-uri 'self'; "+
					"form-action 'self'")
		}

		next.ServeHTTP(w, r)
	})
}

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
