package http

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"

	"github.com/ckmdevb/monkey-planner/backend/internal/domain"
	"github.com/ckmdevb/monkey-planner/backend/internal/service"
)

type commentHandler struct{ svc *service.Service }

func (h *commentHandler) list(w http.ResponseWriter, r *http.Request) {
	issueID := chi.URLParam(r, "issueId")
	comments, err := h.svc.ListComments(r.Context(), issueID)
	if err != nil {
		mapError(w, err)
		return
	}
	if comments == nil {
		comments = []domain.Comment{}
	}
	writeJSON(w, http.StatusOK, comments)
}

func (h *commentHandler) create(w http.ResponseWriter, r *http.Request) {
	issueID := chi.URLParam(r, "issueId")
	var in struct {
		Body string `json:"body"`
	}
	if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
		writeErr(w, http.StatusBadRequest, "invalid_json", err.Error())
		return
	}
	c, err := h.svc.CreateComment(r.Context(), issueID, in.Body)
	if err != nil {
		mapError(w, err)
		return
	}
	writeJSON(w, http.StatusCreated, c)
}

func (h *commentHandler) delete(w http.ResponseWriter, r *http.Request) {
	commentID := chi.URLParam(r, "commentId")
	if err := h.svc.DeleteComment(r.Context(), commentID); err != nil {
		mapError(w, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
