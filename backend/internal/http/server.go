// Package http는 chi 기반 HTTP 라우터와 핸들러를 제공합니다.
package http

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	"github.com/ckmdevb/monkey-planner/backend/internal/service"
)

// NewRouter는 /api/* 경로에 핸들러를 바인딩한 라우터를 반환합니다.
// 정적 파일/SPA fallback은 호출자(main)가 바깥에서 합성합니다.
func NewRouter(svc *service.Service) http.Handler {
	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Logger)

	ih := &issueHandler{svc: svc}
	bh := &boardHandler{svc: svc}
	ch := &calendarHandler{svc: svc}

	r.Route("/api", func(api chi.Router) {
		api.Get("/health", func(w http.ResponseWriter, _ *http.Request) {
			writeJSON(w, http.StatusOK, map[string]any{"ok": true, "version": "0.0.1"})
		})

		api.Route("/boards", func(b chi.Router) {
			b.Get("/", bh.list)
			b.Post("/", bh.create)
			b.Patch("/{id}", bh.patch)
			b.Delete("/{id}", bh.delete)
		})

		api.Route("/issues", func(i chi.Router) {
			i.Get("/", ih.list)
			i.Post("/", ih.create)
			i.Get("/{id}", ih.get)
			i.Patch("/{id}", ih.patch)
			i.Delete("/{id}", ih.delete)
			i.Post("/{id}/approve", ih.approve)
		})

		api.Route("/calendar", func(c chi.Router) {
			c.Get("/", ch.month)
			c.Get("/day", ch.day)
		})
	})

	return r
}
