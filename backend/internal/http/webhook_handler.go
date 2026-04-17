package http

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"

	"github.com/kjm99d/MonkeyPlanner/backend/internal/domain"
	"github.com/kjm99d/MonkeyPlanner/backend/internal/service"
)

type webhookHandler struct{ svc *service.Service }

func (h *webhookHandler) list(w http.ResponseWriter, r *http.Request) {
	boardID := chi.URLParam(r, "boardId")
	out, err := h.svc.ListWebhooks(r.Context(), boardID)
	if err != nil {
		mapError(w, err)
		return
	}
	if out == nil {
		out = []domain.Webhook{}
	}
	writeJSON(w, http.StatusOK, out)
}

func (h *webhookHandler) create(w http.ResponseWriter, r *http.Request) {
	boardID := chi.URLParam(r, "boardId")
	var in struct {
		Name   string                 `json:"name"`
		URL    string                 `json:"url"`
		Events []domain.WebhookEvent  `json:"events"`
	}
	if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
		writeErr(w, http.StatusBadRequest, "invalid_json", err.Error())
		return
	}
	out, err := h.svc.CreateWebhook(r.Context(), boardID, in.Name, in.URL, in.Events)
	if err != nil {
		mapError(w, err)
		return
	}
	writeJSON(w, http.StatusCreated, out)
}

func (h *webhookHandler) update(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "whId")
	var in struct {
		Name    *string                 `json:"name,omitempty"`
		URL     *string                 `json:"url,omitempty"`
		Events  *[]domain.WebhookEvent  `json:"events,omitempty"`
		Enabled *bool                   `json:"enabled,omitempty"`
	}
	if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
		writeErr(w, http.StatusBadRequest, "invalid_json", err.Error())
		return
	}
	out, err := h.svc.UpdateWebhook(r.Context(), id, in.Name, in.URL, in.Events, in.Enabled)
	if err != nil {
		mapError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, out)
}

func (h *webhookHandler) delete(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "whId")
	if err := h.svc.DeleteWebhook(r.Context(), id); err != nil {
		mapError(w, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
