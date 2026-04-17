package http

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"

	"github.com/kjm99d/MonkeyPlanner/backend/internal/domain"
	"github.com/kjm99d/MonkeyPlanner/backend/internal/service"
)

type propertyHandler struct{ svc *service.Service }

func (h *propertyHandler) list(w http.ResponseWriter, r *http.Request) {
	boardID := chi.URLParam(r, "boardId")
	out, err := h.svc.ListBoardProperties(r.Context(), boardID)
	if err != nil {
		mapError(w, err)
		return
	}
	if out == nil {
		out = []domain.BoardProperty{}
	}
	writeJSON(w, http.StatusOK, out)
}

func (h *propertyHandler) create(w http.ResponseWriter, r *http.Request) {
	boardID := chi.URLParam(r, "boardId")
	var in struct {
		Name    string           `json:"name"`
		Type    domain.PropertyType `json:"type"`
		Options []string         `json:"options"`
	}
	if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
		writeErr(w, http.StatusBadRequest, "invalid_json", err.Error())
		return
	}
	out, err := h.svc.CreateBoardProperty(r.Context(), boardID, in.Name, in.Type, in.Options)
	if err != nil {
		mapError(w, err)
		return
	}
	writeJSON(w, http.StatusCreated, out)
}

func (h *propertyHandler) update(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "propId")
	var in struct {
		Name     *string   `json:"name,omitempty"`
		Options  *[]string `json:"options,omitempty"`
		Position *int      `json:"position,omitempty"`
	}
	if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
		writeErr(w, http.StatusBadRequest, "invalid_json", err.Error())
		return
	}
	out, err := h.svc.UpdateBoardProperty(r.Context(), id, in.Name, in.Options, in.Position)
	if err != nil {
		mapError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, out)
}

func (h *propertyHandler) delete(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "propId")
	if err := h.svc.DeleteBoardProperty(r.Context(), id); err != nil {
		mapError(w, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
