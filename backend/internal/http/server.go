// Package http는 chi 기반 HTTP 라우터와 핸들러를 제공합니다.
package http

import (
	"io/fs"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	"github.com/ckmdevb/monkey-planner/backend/internal/service"
)

// NewRouter는 /api/* 경로에 핸들러를 바인딩한 라우터를 반환합니다.
// static 이 nil이 아니면 /api 이외 경로는 SPA fallback 으로 서빙됩니다.
func NewRouter(svc *service.Service, static fs.FS) http.Handler {
	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Logger)
	r.Use(ValidateUTF8)

	ih := &issueHandler{svc: svc}
	bh := &boardHandler{svc: svc}
	ch := &calendarHandler{svc: svc}
	ph := &propertyHandler{svc: svc}
	wh := &webhookHandler{svc: svc}
	cmh := &commentHandler{svc: svc}

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

		// 보드 속성(커스텀 프로퍼티)
		api.Route("/boards/{boardId}/properties", func(p chi.Router) {
			p.Get("/", ph.list)
			p.Post("/", ph.create)
			p.Patch("/{propId}", ph.update)
			p.Delete("/{propId}", ph.delete)
		})

		// 웹훅
		api.Route("/boards/{boardId}/webhooks", func(w chi.Router) {
			w.Get("/", wh.list)
			w.Post("/", wh.create)
			w.Patch("/{whId}", wh.update)
			w.Delete("/{whId}", wh.delete)
		})

		// 보드별 이슈 순서 재정렬
		api.Post("/boards/{boardId}/issues/reorder", ih.reorder)

		api.Route("/issues", func(i chi.Router) {
			i.Get("/", ih.list)
			i.Post("/", ih.create)
			i.Get("/{id}", ih.get)
			i.Patch("/{id}", ih.patch)
			i.Delete("/{id}", ih.delete)
			i.Post("/{id}/approve", ih.approve)
			// 이슈 댓글
			i.Get("/{issueId}/comments", cmh.list)
			i.Post("/{issueId}/comments", cmh.create)
		})

		// 댓글 삭제 (issueId 불필요)
		api.Delete("/comments/{commentId}", cmh.delete)

		api.Route("/calendar", func(c chi.Router) {
			c.Get("/", ch.month)
			c.Get("/day", ch.day)
		})
	})

	if static != nil {
		spa := SPAHandler(static)
		r.NotFound(spa.ServeHTTP)
		r.MethodNotAllowed(spa.ServeHTTP)
	}

	return r
}
