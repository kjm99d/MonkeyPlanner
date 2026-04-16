package http

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"

	"github.com/kjm99d/monkey-planner/backend/internal/domain"
	"github.com/kjm99d/monkey-planner/backend/internal/service"
	"github.com/kjm99d/monkey-planner/backend/internal/storage"
)

type issueHandler struct{ svc *service.Service }

func (h *issueHandler) list(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	var f storage.IssueFilter
	if v := q.Get("board_id"); v != "" {
		f.BoardID = &v
	}
	if v := q.Get("status"); v != "" {
		s := domain.Status(v)
		if !s.Valid() {
			writeErr(w, http.StatusBadRequest, "invalid_status", "unknown status")
			return
		}
		f.Status = &s
	}
	if q.Has("parent_id") {
		v := q.Get("parent_id") // 빈 문자열이면 루트 필터
		f.ParentID = &v
	}
	out, err := h.svc.ListIssues(r.Context(), f)
	if err != nil {
		mapError(w, err)
		return
	}
	if out == nil {
		out = []domain.Issue{}
	}
	writeJSON(w, http.StatusOK, out)
}

func (h *issueHandler) get(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	iss, children, err := h.svc.GetIssue(r.Context(), id)
	if err != nil {
		mapError(w, err)
		return
	}
	if children == nil {
		children = []domain.Issue{}
	}
	writeJSON(w, http.StatusOK, map[string]any{
		"issue":    iss,
		"children": children,
	})
}

func (h *issueHandler) create(w http.ResponseWriter, r *http.Request) {
	var in struct {
		BoardID  string  `json:"boardId"`
		ParentID *string `json:"parentId,omitempty"`
		Title    string  `json:"title"`
		Body     string  `json:"body"`
	}
	if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
		writeErr(w, http.StatusBadRequest, "invalid_json", err.Error())
		return
	}
	out, err := h.svc.CreateIssue(r.Context(), service.CreateIssueInput{
		BoardID:  in.BoardID,
		ParentID: in.ParentID,
		Title:    in.Title,
		Body:     in.Body,
	})
	if err != nil {
		mapError(w, err)
		return
	}
	writeJSON(w, http.StatusCreated, out)
}

func (h *issueHandler) patch(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	// parentId 는 3상태: 없음(미변경) / null(루트) / 문자열(지정)
	raw := map[string]json.RawMessage{}
	if err := json.NewDecoder(r.Body).Decode(&raw); err != nil {
		writeErr(w, http.StatusBadRequest, "invalid_json", err.Error())
		return
	}
	var in service.UpdateIssueInput
	if v, ok := raw["title"]; ok {
		var s string
		if err := json.Unmarshal(v, &s); err != nil {
			writeErr(w, http.StatusBadRequest, "invalid_title", err.Error())
			return
		}
		in.Title = &s
	}
	if v, ok := raw["body"]; ok {
		var s string
		if err := json.Unmarshal(v, &s); err != nil {
			writeErr(w, http.StatusBadRequest, "invalid_body", err.Error())
			return
		}
		in.Body = &s
	}
	if v, ok := raw["status"]; ok {
		var s domain.Status
		if err := json.Unmarshal(v, &s); err != nil {
			writeErr(w, http.StatusBadRequest, "invalid_status", err.Error())
			return
		}
		if !s.Valid() {
			writeErr(w, http.StatusBadRequest, "invalid_status", "unknown status value")
			return
		}
		in.Status = &s
	}
	if v, ok := raw["properties"]; ok {
		var props map[string]any
		if err := json.Unmarshal(v, &props); err != nil {
			writeErr(w, http.StatusBadRequest, "invalid_properties", err.Error())
			return
		}
		// properties는 merge 방식으로 서비스 계층에서 처리
		updated, err := h.svc.UpdateIssueProperties(r.Context(), id, props)
		if err != nil {
			mapError(w, err)
			return
		}
		writeJSON(w, http.StatusOK, updated)
		return
	}
	if v, ok := raw["parentId"]; ok {
		if string(v) == "null" {
			var np *string // nil → root
			in.ParentID = &np
		} else {
			var s string
			if err := json.Unmarshal(v, &s); err != nil {
				writeErr(w, http.StatusBadRequest, "invalid_parent", err.Error())
				return
			}
			sp := &s
			in.ParentID = &sp
		}
	}
	out, err := h.svc.UpdateIssue(r.Context(), id, in)
	if err != nil {
		mapError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, out)
}

func (h *issueHandler) approve(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	out, err := h.svc.ApproveIssue(r.Context(), id)
	if err != nil {
		mapError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, out)
}

func (h *issueHandler) delete(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if err := h.svc.DeleteIssue(r.Context(), id); err != nil {
		mapError(w, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
