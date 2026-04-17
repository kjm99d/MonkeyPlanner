package http

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/kjm99d/monkey-planner/backend/internal/events"
)

// heartbeatInterval is how often an SSE comment is flushed to keep the
// connection alive through reverse proxies (nginx idle 60s, CloudFront 60s,
// Cloudflare 100s). 30s leaves ample margin.
const heartbeatInterval = 30 * time.Second

type eventsHandler struct {
	broker *events.Broker
}

// stream is the SSE endpoint: GET /api/events?boardId=xxx
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
	w.Header().Set("X-Accel-Buffering", "no") // disable nginx response buffering

	ch := h.broker.Subscribe(boardID)
	defer h.broker.Unsubscribe(boardID, ch)

	// Send an SSE comment immediately to flush response headers to the client.
	fmt.Fprintf(w, ": connected\n\n")
	flusher.Flush()

	heartbeat := time.NewTicker(heartbeatInterval)
	defer heartbeat.Stop()

	ctx := r.Context()
	for {
		select {
		case <-ctx.Done():
			return
		case <-heartbeat.C:
			// SSE comment line — clients ignore it but proxies see traffic.
			if _, err := fmt.Fprintf(w, ": heartbeat\n\n"); err != nil {
				return
			}
			flusher.Flush()
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
