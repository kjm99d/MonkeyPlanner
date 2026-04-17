import { useEffect } from 'react';
import { useQueryClient } from '@tanstack/react-query';

type ServerEvent = {
  type: string;
  boardId: string;
  issueId?: string;
  status?: string;
  timestamp: string;
};

// How often to full-refresh as a safety net if SSE events are missed.
// Covers the gap between an SSE drop and the next event-driven invalidate.
const FALLBACK_INTERVAL_MS = 5 * 60 * 1000; // 5 minutes

/**
 * Subscribes to server events for a board and invalidates scoped React Query caches.
 * EventSource handles auto-reconnect automatically.
 *
 * Two recovery paths ensure UI correctness even when events are dropped:
 *   1. onerror — invalidates immediately so a reconnect fetches current state.
 *   2. 5-minute interval — catches any events missed during a connection gap.
 */
export function useEventStream(boardId: string | undefined) {
  const qc = useQueryClient();

  useEffect(() => {
    if (!boardId) return;

    const invalidateBoard = () => {
      qc.invalidateQueries({ queryKey: ['issues', { boardId }] });
    };

    const es = new EventSource(`/api/events?boardId=${encodeURIComponent(boardId)}`);

    es.onmessage = (e) => {
      if (!e.data) return;
      let ev: ServerEvent;
      try {
        ev = JSON.parse(e.data);
      } catch {
        return;
      }
      switch (ev.type) {
        case 'issue.created':
        case 'issue.updated':
        case 'issue.status_changed':
        case 'issue.approved':
        case 'issue.deleted':
          // Partial match: invalidates any ['issues', { boardId, ... }] query for this board.
          qc.invalidateQueries({ queryKey: ['issues', { boardId: ev.boardId }] });
          qc.invalidateQueries({ queryKey: ['calendar'] });
          if (ev.issueId) {
            qc.invalidateQueries({ queryKey: ['issue', ev.issueId] });
          }
          break;
        case 'comment.created':
          if (ev.issueId) {
            qc.invalidateQueries({ queryKey: ['comments', ev.issueId] });
            qc.invalidateQueries({ queryKey: ['issue', ev.issueId] });
          }
          break;
      }
    };

    // On error (proxy timeout, brief disconnect) force a refetch immediately
    // so the UI does not stay stale while EventSource auto-reconnects.
    es.onerror = () => invalidateBoard();

    // Periodic safety-net: full refresh every 5 minutes regardless of SSE state.
    const timer = setInterval(invalidateBoard, FALLBACK_INTERVAL_MS);

    return () => {
      es.close();
      clearInterval(timer);
    };
  }, [boardId, qc]);
}
