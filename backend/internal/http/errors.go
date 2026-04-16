package http

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"

	"github.com/kjm99d/monkey-planner/backend/internal/service"
	"github.com/kjm99d/monkey-planner/backend/internal/storage"
)

// errBody는 공통 JSON 에러 응답 본문입니다.
type errBody struct {
	Error struct {
		Code    string `json:"code"`
		Message string `json:"message"`
	} `json:"error"`
}

func writeJSON(w http.ResponseWriter, status int, body any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	if body == nil {
		return
	}
	if err := json.NewEncoder(w).Encode(body); err != nil {
		log.Printf("writeJSON: %v", err)
	}
}

func writeErr(w http.ResponseWriter, status int, code, msg string) {
	b := errBody{}
	b.Error.Code = code
	b.Error.Message = msg
	writeJSON(w, status, b)
}

// mapError 는 service/storage 에러를 HTTP 상태로 매핑합니다.
func mapError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, storage.ErrNotFound):
		writeErr(w, http.StatusNotFound, "not_found", err.Error())
	case errors.Is(err, storage.ErrCycle):
		writeErr(w, http.StatusBadRequest, "cycle", "parent_id would create a cycle")
	case errors.Is(err, storage.ErrConflict):
		writeErr(w, http.StatusConflict, "conflict", err.Error())
	case errors.Is(err, service.ErrApproveViaPatch):
		writeErr(w, http.StatusConflict, "use_approve_endpoint",
			"use POST /api/issues/:id/approve")
	case errors.Is(err, service.ErrBackwardTransition):
		writeErr(w, http.StatusBadRequest, "backward_transition", err.Error())
	case errors.Is(err, service.ErrInvalidTransition),
		errors.Is(err, service.ErrEmptyTitle),
		errors.Is(err, service.ErrMissingBoard):
		writeErr(w, http.StatusBadRequest, "invalid_input", err.Error())
	default:
		log.Printf("http: unhandled err: %v", err)
		writeErr(w, http.StatusInternalServerError, "internal", "internal server error")
	}
}
