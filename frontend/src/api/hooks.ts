import {
  useQuery,
  useMutation,
  useQueryClient,
  type UseQueryOptions,
} from '@tanstack/react-query';
import { api } from './client';
import type { Board, BoardProperty, DayCount, DayStats, Issue, IssueStatus, PropertyType } from './types';

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

// ---- Board Properties ----

export function useBoardProperties(boardId: string | undefined) {
  return useQuery<BoardProperty[]>({
    queryKey: ['boardProperties', boardId],
    queryFn: () => api.get<BoardProperty[]>(`/api/boards/${boardId}/properties`),
    enabled: !!boardId,
  });
}

export function useCreateBoardProperty() {
  const qc = useQueryClient();
  return useMutation({
    mutationFn: (p: { boardId: string; name: string; type: PropertyType; options?: string[] }) =>
      api.post<BoardProperty>(`/api/boards/${p.boardId}/properties`, p),
    onSuccess: (_d, v) => {
      qc.invalidateQueries({ queryKey: ['boardProperties', v.boardId] });
    },
  });
}

export function useDeleteBoardProperty() {
  const qc = useQueryClient();
  return useMutation({
    mutationFn: ({ boardId, propId }: { boardId: string; propId: string }) =>
      api.del<void>(`/api/boards/${boardId}/properties/${propId}`),
    onSuccess: (_d, v) => {
      qc.invalidateQueries({ queryKey: ['boardProperties', v.boardId] });
    },
  });
}

export function useUpdateIssueProperties() {
  const qc = useQueryClient();
  return useMutation({
    mutationFn: ({ id, properties }: { id: string; properties: Record<string, unknown> }) =>
      api.patch<Issue>(`/api/issues/${id}`, { properties }),
    onSuccess: (_d, v) => {
      qc.invalidateQueries({ queryKey: ['issue', v.id] });
      qc.invalidateQueries({ queryKey: ['issues'] });
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
