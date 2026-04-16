import {
  useQuery,
  useMutation,
  useQueryClient,
  type UseQueryOptions,
} from '@tanstack/react-query';
import { api } from './client';
import type { Board, DayCount, DayStats, Issue, IssueStatus } from './types';

// ---- Boards ----

export const boardsKey = ['boards'] as const;

export function useBoards(options?: Partial<UseQueryOptions<Board[]>>) {
  return useQuery<Board[]>({
    queryKey: boardsKey,
    queryFn: () => api.get<Board[]>('/api/boards'),
    ...options,
  });
}

export function useCreateBoard() {
  const qc = useQueryClient();
  return useMutation({
    mutationFn: (p: { name: string; viewType?: 'kanban' | 'list' }) =>
      api.post<Board>('/api/boards', p),
    onSuccess: () => {
      qc.invalidateQueries({ queryKey: boardsKey });
    },
  });
}

// ---- Issues ----

export function issuesKey(filter?: {
  boardId?: string;
  status?: IssueStatus;
  parentId?: string | null;
}) {
  return ['issues', filter ?? {}] as const;
}

export function useIssues(filter?: {
  boardId?: string;
  status?: IssueStatus;
  parentId?: string | null;
}) {
  const qs = new URLSearchParams();
  if (filter?.boardId) qs.set('board_id', filter.boardId);
  if (filter?.status) qs.set('status', filter.status);
  if (filter?.parentId !== undefined)
    qs.set('parent_id', filter.parentId === null ? '' : filter.parentId);
  return useQuery<Issue[]>({
    queryKey: issuesKey(filter),
    queryFn: () => api.get<Issue[]>(`/api/issues?${qs.toString()}`),
  });
}

export function useIssue(id: string | undefined) {
  return useQuery<{ issue: Issue; children: Issue[] }>({
    queryKey: ['issue', id],
    queryFn: () => api.get(`/api/issues/${id}`),
    enabled: !!id,
  });
}

export function useCreateIssue() {
  const qc = useQueryClient();
  return useMutation({
    mutationFn: (p: { boardId: string; title: string; body?: string; parentId?: string }) =>
      api.post<Issue>('/api/issues', p),
    onSuccess: () => {
      qc.invalidateQueries({ queryKey: ['issues'] });
      qc.invalidateQueries({ queryKey: ['calendar'] });
    },
  });
}

export function useUpdateIssue() {
  const qc = useQueryClient();
  return useMutation({
    mutationFn: ({
      id,
      patch,
    }: {
      id: string;
      patch: Partial<Pick<Issue, 'title' | 'body' | 'status'>> & {
        parentId?: string | null;
      };
    }) => api.patch<Issue>(`/api/issues/${id}`, patch),
    onSuccess: (_data, vars) => {
      qc.invalidateQueries({ queryKey: ['issues'] });
      qc.invalidateQueries({ queryKey: ['issue', vars.id] });
      qc.invalidateQueries({ queryKey: ['calendar'] });
    },
  });
}

export function useApproveIssue() {
  const qc = useQueryClient();
  return useMutation({
    mutationFn: (id: string) => api.post<Issue>(`/api/issues/${id}/approve`),
    onSuccess: (_data, id) => {
      qc.invalidateQueries({ queryKey: ['issues'] });
      qc.invalidateQueries({ queryKey: ['issue', id] });
      qc.invalidateQueries({ queryKey: ['calendar'] });
    },
  });
}

export function useDeleteIssue() {
  const qc = useQueryClient();
  return useMutation({
    mutationFn: (id: string) => api.del<void>(`/api/issues/${id}`),
    onSuccess: () => {
      qc.invalidateQueries({ queryKey: ['issues'] });
      qc.invalidateQueries({ queryKey: ['calendar'] });
    },
  });
}

// ---- Calendar ----

export function useMonthStats(year: number, month: number) {
  return useQuery<DayCount[]>({
    queryKey: ['calendar', 'month', year, month],
    queryFn: () => api.get<DayCount[]>(`/api/calendar?year=${year}&month=${month}`),
  });
}

export function useDayStats(date: string | undefined) {
  return useQuery<DayStats>({
    queryKey: ['calendar', 'day', date],
    queryFn: () => api.get<DayStats>(`/api/calendar/day?date=${date}`),
    enabled: !!date,
  });
}
