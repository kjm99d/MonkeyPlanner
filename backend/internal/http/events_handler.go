package http

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/kjm99d/monkey-planner/backend/internal/events"
)

type eventsHandler struct {
	broker *events.Broker
}

// stream 은 SSE 엔드포인트입니다. GET /api/events?boardId=xxx
func (h *eventsHandler) stream(w http.ResponseWriter, r *http.Request) {
	boardID := r.URL.Query().Get("boardId")
	if boardID == "" {
		writeErr(w, http.StatusBadRequest, "missing_board_id", "boardId is required")
		return
	}

	flusher, ok := w.(http.Flusher)
	if !ok {
		writeErr(w, http.StatusInternalServerError, "streaming_unsupported", "server does not support streaming")
		return
	}

	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
	w.Header().Set("X-Accel-Buffering", "no") // nginx 버퍼링 방지

	ch := h.broker.Subscribe(boardID)
	defer h.broker.Unsubscribe(boardID, ch)

	// 연결 즉시 핑 이벤트로 헤더 flush
	fmt.Fprintf(w, ": connected\n\n")
	flusher.Flush()

	ctx := r.Context()
	for {
		select {
		case <-ctx.Done():
			return
		case ev, ok := <-ch:
			if !ok {
				return
			}
			data, err := json.Marshal(ev)
			if err != nil {
				continue
			}
			fmt.Fprintf(w, "data: %s\n\n", data)
			flusher.Flush()
		}
	}
}
