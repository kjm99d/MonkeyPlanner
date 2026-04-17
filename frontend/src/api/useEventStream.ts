import { useEffect } from 'react';
import { useQueryClient } from '@tanstack/react-query';

type ServerEvent = {
  type: string;
  boardId: string;
  issueId?: string;
  status?: string;
  timestamp: string;
};

/**
 * Subscribes to server events for a board and invalidates scoped React Query caches.
 * EventSource handles auto-reconnect, so no explicit retry logic is needed.
 */
export function useEventStream(boardId: string | undefined) {
  const qc = useQueryClient();

  useEffect(() => {
    if (!boardId) return;
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

    es.onerror = () => {
      // EventSource auto-reconnects; no explicit handling needed.
    };

    return () => es.close();
  }, [boardId, qc]);
}
