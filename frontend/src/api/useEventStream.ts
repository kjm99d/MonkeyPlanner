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
 * boardId의 서버 이벤트를 구독하고 React Query 캐시를 무효화합니다.
 * EventSource는 자동 재연결을 제공하므로 별도 로직 불필요.
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
          qc.invalidateQueries({ queryKey: ['issues'] });
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
      // EventSource가 자동 재연결함. 명시적 처리 불필요.
    };

    return () => es.close();
  }, [boardId, qc]);
}
