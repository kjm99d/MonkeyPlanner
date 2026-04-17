package http_test

import (
	"bytes"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	mphttp "github.com/kjm99d/MonkeyPlanner/backend/internal/http"
)

func TestValidateUTF8_RejectsNonUTF8(t *testing.T) {
	handler := mphttp.ValidateUTF8(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(200)
	}))

	// "테스트" encoded in cp949 (0xC5 0xD7 0xBD 0xBA 0xC6 0xAE) — intentionally invalid UTF-8.
	cp949Body := []byte{0x7b, 0x22, 0x74, 0x22, 0x3a, 0x22, 0xC5, 0xD7, 0xBD, 0xBA, 0xC6, 0xAE, 0x22, 0x7d}
	req := httptest.NewRequest(http.MethodPost, "/api/issues", bytes.NewReader(cp949Body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, req)

	if w.Code != 400 {
		t.Fatalf("expected 400 for non-UTF-8, got %d", w.Code)
	}
	body, _ := io.ReadAll(w.Result().Body)
	if !bytes.Contains(body, []byte("invalid_encoding")) {
		t.Fatalf("expected invalid_encoding error, got %s", body)
	}
}

func TestValidateUTF8_AllowsValidUTF8(t *testing.T) {
	var called bool
	handler := mphttp.ValidateUTF8(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		called = true
		body, _ := io.ReadAll(r.Body)
		if string(body) != `{"t":"한글"}` {
			t.Fatalf("body mangled: %s", body)
		}
		w.WriteHeader(200)
	}))

	req := httptest.NewRequest(http.MethodPost, "/api/issues", bytes.NewReader([]byte(`{"t":"한글"}`)))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, req)

	if !called {
		t.Fatal("next handler was not called")
	}
	if w.Code != 200 {
		t.Fatalf("expected 200 for valid UTF-8, got %d", w.Code)
	}
}

func TestValidateUTF8_SkipsNonJSON(t *testing.T) {
	var called bool
	handler := mphttp.ValidateUTF8(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		called = true
		w.WriteHeader(200)
	}))

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, req)

	if !called || w.Code != 200 {
		t.Fatal("GET without body should pass through")
	}
}
