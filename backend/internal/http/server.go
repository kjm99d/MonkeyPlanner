// Package http provides the chi-based HTTP router and handlers.
package http

import (
	"io/fs"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	"github.com/kjm99d/MonkeyPlanner/backend/internal/service"
)

// NewRouter wires the /api/* handlers onto a chi router. When static is
// non-nil, any non-/api path falls through to the SPA handler.
func NewRouter(svc *service.Service, static fs.FS, version string) http.Handler {
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
	eh := &eventsHandler{broker: svc.Broker()}

	r.Route("/api", func(api chi.Router) {
		api.Get("/health", func(w http.ResponseWriter, _ *http.Request) {
			writeJSON(w, http.StatusOK, map[string]any{"ok": true, "version": version})
		})

		api.Route("/boards", func(b chi.Router) {
			b.Get("/", bh.list)
			b.Post("/", bh.create)
			b.Patch("/{id}", bh.patch)
			b.Delete("/{id}", bh.delete)
		})

		// Board-level custom property definitions.
		api.Route("/boards/{boardId}/properties", func(p chi.Router) {
			p.Get("/", ph.list)
			p.Post("/", ph.create)
			p.Patch("/{propId}", ph.update)
			p.Delete("/{propId}", ph.delete)
		})

		// Outbound webhooks.
		api.Route("/boards/{boardId}/webhooks", func(w chi.Router) {
			w.Get("/", wh.list)
			w.Post("/", wh.create)
			w.Patch("/{whId}", wh.update)
			w.Delete("/{whId}", wh.delete)
		})

		// Reorder issues within a board.
		api.Post("/boards/{boardId}/issues/reorder", ih.reorder)

		api.Route("/issues", func(i chi.Router) {
			i.Get("/", ih.list)
			i.Post("/", ih.create)
			i.Get("/{id}", ih.get)
			i.Patch("/{id}", ih.patch)
			i.Delete("/{id}", ih.delete)
			i.Post("/{id}/approve", ih.approve)
			// Issue dependencies.
			i.Post("/{issueId}/dependencies", ih.addDependency)
			i.Delete("/{issueId}/dependencies/{blockerId}", ih.removeDependency)
			// Issue comments.
			i.Get("/{issueId}/comments", cmh.list)
			i.Post("/{issueId}/comments", cmh.create)
		})

		// Comment deletion (issueId is not part of the path).
		api.Delete("/comments/{commentId}", cmh.delete)

		api.Route("/calendar", func(c chi.Router) {
			c.Get("/", ch.month)
			c.Get("/day", ch.day)
		})

		// SSE event stream for real-time UI updates.
		api.Get("/events", eh.stream)
	})

	if static != nil {
		spa := SPAHandler(static)
		r.NotFound(spa.ServeHTTP)
		r.MethodNotAllowed(spa.ServeHTTP)
	}

	return r
}
