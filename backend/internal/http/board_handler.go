package http

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"

	"github.com/kjm99d/monkey-planner/backend/internal/domain"
	"github.com/kjm99d/monkey-planner/backend/internal/service"
)

type boardHandler struct{ svc *service.Service }

func (h *boardHandler) list(w http.ResponseWriter, r *http.Request) {
	out, err := h.svc.ListBoards(r.Context())
	if err != nil {
		mapError(w, err)
		return
	}
	if out == nil {
		out = []domain.Board{}
	}
	writeJSON(w, http.StatusOK, out)
}

func (h *boardHandler) create(w http.ResponseWriter, r *http.Request) {
	var in struct {
		Name     string          `json:"name"`
		ViewType domain.ViewType `json:"viewType"`
	}
	if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
		writeErr(w, http.StatusBadRequest, "invalid_json", err.Error())
		return
	}
	out, err := h.svc.CreateBoard(r.Context(), in.Name, in.ViewType)
	if err != nil {
		mapError(w, err)
		return
	}
	writeJSON(w, http.StatusCreated, out)
}

func (h *boardHandler) patch(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	var in struct {
		Name     *string          `json:"name,omitempty"`
		ViewType *domain.ViewType `json:"viewType,omitempty"`
	}
	if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
		writeErr(w, http.StatusBadRequest, "invalid_json", err.Error())
		return
	}
	out, err := h.svc.UpdateBoard(r.Context(), id, in.Name, in.ViewType)
	if err != nil {
		mapError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, out)
}

func (h *boardHandler) delete(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if err := h.svc.DeleteBoard(r.Context(), id); err != nil {
		mapError(w, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
