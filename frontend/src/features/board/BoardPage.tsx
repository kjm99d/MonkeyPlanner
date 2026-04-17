import { useEffect, useMemo, useState } from 'react';
import { useParams, Link } from 'react-router-dom';
import { useTranslation } from 'react-i18next';
import { DndContext, DragOverlay, type DragEndEvent, type DragStartEvent, PointerSensor, useSensor, useSensors } from '@dnd-kit/core';
import { useBoards, useBoardProperties, useCreateBoardProperty, useDeleteBoardProperty, useCreateIssue, useIssues, useUpdateIssue, useDeleteBoard, useReorderIssues } from '../../api/hooks';
import { useEventStream } from '../../api/useEventStream';
import { useQueryClient } from '@tanstack/react-query';
import { api } from '../../api/client';
import { useNavigate } from 'react-router-dom';
import { Trash2, LayoutGrid, List, Filter, AlertCircle, Download, FileText, Maximize2, Minimize2 } from 'lucide-react';
import { Button } from '../../components/Button';
import { Input } from '../../components/Input';
import { Breadcrumb } from '../../components/Breadcrumb';
import { AddPropertyForm } from '../../components/PropertyEditor';
import { useToast } from '../../components/Toast';
import { WebhookSettings } from '../../components/WebhookSettings';
import { ConfirmDialog } from '../../components/ConfirmDialog';
import { TemplateDialog } from '../../components/TemplateDialog';
import { StatusBadge } from '../../components/StatusBadge';
import { KanbanColumn } from './KanbanColumn';
import { AgentPresenceBar } from './AgentPresenceBar';
import type { Issue, IssueStatus } from '../../api/types';

const COLUMN_KEYS: { status: IssueStatus; key: string }[] = [
  { status: 'Pending', key: 'kanban.pending' },
  { status: 'Approved', key: 'kanban.approved' },
  { status: 'InProgress', key: 'kanban.inProgress' },
  { status: 'QA', key: 'kanban.qa' },
  { status: 'Done', key: 'kanban.done' },
  { status: 'Rejected', key: 'kanban.rejected' },
];

export default function BoardPage() {
  const { boardId } = useParams<{ boardId: string }>();
  const boards = useBoards();
  const board = boards.data?.find((b) => b.id === boardId);
  const issues = useIssues({ boardId });
  const createIssue = useCreateIssue();
  const deleteBoard = useDeleteBoard();
  const boardPropsQuery = useBoardProperties(boardId);
  const createProp = useCreateBoardProperty();
  const deleteProp = useDeleteBoardProperty();
  const updateIssue = useUpdateIssue();
  const reorderIssues = useReorderIssues();

  const { t } = useTranslation();
  const navigate = useNavigate();
  const { toast } = useToast();
  const qc = useQueryClient();
  const [title, setTitle] = useState('');
  const [errMsg, setErrMsg] = useState<string | null>(null);
  const [activeIssue, setActiveIssue] = useState<Issue | null>(null);
  const [confirmDelete, setConfirmDelete] = useState(false);
  const [viewMode, setViewMode] = useState<'kanban' | 'table'>('kanban');
  const [filterText, setFilterText] = useState('');
  const [filterStatus, setFilterStatus] = useState<IssueStatus | 'all'>('all');
  const [editingName, setEditingName] = useState(false);
  const [boardName, setBoardName] = useState('');
  const [hideDone, setHideDone] = useState(false);
  const [templateOpen, setTemplateOpen] = useState(false);
  const [fullscreen, setFullscreen] = useState(true);

  useEffect(() => {
    localStorage.setItem('board-fullscreen', String(fullscreen));
    // 전체화면 시 Layout의 max-w 제한을 해제
    const main = document.querySelector('main > div');
    if (main) {
      if (fullscreen) {
        main.classList.remove('max-w-6xl');
        main.classList.add('max-w-full');
      } else {
        main.classList.remove('max-w-full');
        main.classList.add('max-w-6xl');
      }
    }
    return () => {
      if (main) {
        main.classList.remove('max-w-full');
        main.classList.add('max-w-6xl');
      }
    };
  }, [fullscreen]);

  useEffect(() => { if (board) setBoardName(board.name); }, [board?.name]);

  useEventStream(boardId);

  const sensors = useSensors(
    useSensor(PointerSensor, { activationConstraint: { distance: 6 } }),
  );

  const filtered = useMemo(() => {
    let items = issues.data ?? [];
    if (filterText.trim()) {
      const q = filterText.toLowerCase();
      items = items.filter(i => i.title.toLowerCase().includes(q));
    }
    if (filterStatus !== 'all') {
      items = items.filter(i => i.status === filterStatus);
    }
    return items;
  }, [issues.data, filterText, filterStatus]);

  const grouped = useMemo(() => {
    const map: Record<IssueStatus, Issue[]> = { Pending: [], Approved: [], InProgress: [], QA: [], Done: [], Rejected: [] };
    filtered.forEach((i) => map[i.status].push(i));
    return map;
  }, [filtered]);

  async function onCreate(e: React.FormEvent) {
    e.preventDefault();
    if (!title.trim() || !boardId) return;
    try {
      await createIssue.mutateAsync({ boardId, title: title.trim() });
      setTitle('');
      toast('success', t('board.issueCreated'));
    } catch (err) {
      toast('error', (err as { message?: string }).message ?? t('board.issueCreateFailed'));
    }
  }

  function onDragStart(e: DragStartEvent) {
    const issue = (issues.data ?? []).find((i) => i.id === String(e.active.id));
    setActiveIssue(issue ?? null);
  }

  async function onDragEnd(e: DragEndEvent) {
    setActiveIssue(null);
    setErrMsg(null);
    if (!e.over) return;
    const toStatus = e.over.id as IssueStatus;
    const issueId = String(e.active.id);
    const current = (issues.data ?? []).find((i) => i.id === issueId);
    if (!current || current.status === toStatus) return;

    // Optimistic update: immediately move the card in the cache
    const issuesKey = ['issues', { boardId }];
    const previousIssues = qc.getQueryData(issuesKey);
    qc.setQueryData(issuesKey, (old: any) =>
      old?.map((i: any) => i.id === issueId ? { ...i, status: toStatus } : i)
    );

    try {
      await updateIssue.mutateAsync({ id: issueId, patch: { status: toStatus } });
      toast('success', t('board.statusChanged', { status: toStatus }));
    } catch (err) {
      // Rollback on error
      qc.setQueryData(issuesKey, previousIssues);
      const msg = t('board.statusChangeFailed');
      setErrMsg(msg);
      toast('error', msg);
    }
  }

  async function handleReorder(status: IssueStatus, issueId: string, direction: 'up' | 'down') {
    if (!boardId) return;
    const columnIssues = grouped[status];
    const idx = columnIssues.findIndex((i) => i.id === issueId);
    if (idx === -1) return;
    const swapIdx = direction === 'up' ? idx - 1 : idx + 1;
    if (swapIdx < 0 || swapIdx >= columnIssues.length) return;

    // Build the new order for this column
    const newColumnOrder = [...columnIssues];
    [newColumnOrder[idx], newColumnOrder[swapIdx]] = [newColumnOrder[swapIdx], newColumnOrder[idx]];

    // Optimistic update: reorder in cache preserving issues from other columns
    const issuesKey = ['issues', { boardId }];
    const previousIssues = qc.getQueryData(issuesKey);
    const allIssues: Issue[] = (issues.data ?? []);
    const otherIssues = allIssues.filter((i) => i.status !== status);
    qc.setQueryData(issuesKey, [...otherIssues, ...newColumnOrder]);

    try {
      const allColumnIds = newColumnOrder.map((i) => i.id);
      await reorderIssues.mutateAsync({ boardId, issueIds: allColumnIds });
    } catch (err) {
      qc.setQueryData(issuesKey, previousIssues);
      toast('error', t('board.statusChangeFailed'));
    }
  }

  return (
    <section className="flex flex-col gap-6">
      <header className="flex items-start justify-between">
        <div className="flex flex-col gap-2">
          <Breadcrumb items={[
            { label: t('nav.boards'), to: '/boards' },
            { label: board?.name ?? '...' },
          ]} />
          {editingName ? (
            <input
              autoFocus
              value={boardName}
              onChange={(e) => setBoardName(e.target.value)}
              onBlur={async () => {
                setEditingName(false);
                if (boardName.trim() && boardName !== board?.name) {
                  await api.patch(`/api/boards/${boardId}`, { name: boardName.trim() });
                  qc.invalidateQueries({ queryKey: ['boards'] });
                }
              }}
              onKeyDown={(e) => {
                if (e.key === 'Enter') e.currentTarget.blur();
                if (e.key === 'Escape') { setBoardName(board?.name ?? ''); setEditingName(false); }
              }}
              className="text-2xl font-bold bg-transparent border-b-2 border-brand-500 outline-none w-full"
            />
          ) : (
            <h1
              className="text-2xl font-bold cursor-pointer hover:text-brand-500 transition-colors"
              onClick={() => setEditingName(true)}
              title="Click to rename"
            >
              {board?.name ?? t('board.title')}
            </h1>
          )}
        </div>
        {boardId && (
          <div className="flex items-center gap-2">
            <button
              type="button"
              onClick={() => {
                const data = {
                  board: { id: boardId, name: board?.name },
                  issues: issues.data ?? [],
                  exportedAt: new Date().toISOString(),
                };
                const blob = new Blob([JSON.stringify(data, null, 2)], { type: 'application/json' });
                const url = URL.createObjectURL(blob);
                const a = document.createElement('a');
                a.href = url;
                a.download = `${board?.name ?? 'board'}-export.json`;
                a.click();
                URL.revokeObjectURL(url);
              }}
              className="flex items-center gap-1.5 rounded-md px-3 py-1.5 text-xs font-medium text-ink-secondary transition-colors hover:bg-surface-muted"
            >
              <Download size={13} />
              Export
            </button>
            <button
              type="button"
              onClick={() => setConfirmDelete(true)}
              className="flex items-center gap-1.5 rounded-md px-3 py-1.5 text-xs font-medium text-red-600 transition-colors hover:bg-red-50 dark:hover:bg-red-950/30"
            >
              <Trash2 size={13} />
              {t('board.delete')}
            </button>
          </div>
        )}
      </header>

      {/* Agent activity — shows InProgress/QA issues in real time */}
      <AgentPresenceBar issues={issues.data ?? []} />

      {/* 속성 관리 */}
      <div className="flex flex-wrap items-center gap-2">
        {boardPropsQuery.data?.map((p) => (
          <span key={p.id} className="group flex items-center gap-1 rounded-full border border-edge-base bg-surface-subtle px-2.5 py-0.5 text-xs text-ink-secondary">
            {p.name} · {p.type}
            <button
              type="button"
              onClick={async () => {
                if (window.confirm(t('board.propertyDeleteConfirm', { name: p.name }))) {
                  try {
                    await deleteProp.mutateAsync({ boardId: boardId!, propId: p.id });
                    toast('success', t('board.propertyDeleted', { name: p.name }));
                  } catch {
                    toast('error', t('board.statusChangeFailed'));
                  }
                }
              }}
              className="hidden group-hover:inline-flex items-center justify-center rounded-full text-ink-muted hover:text-red-500 transition-colors"
              aria-label={t('board.deletePropertyLabel', { name: p.name })}
            >
              ×
            </button>
          </span>
        ))}
        {boardId && (
          <AddPropertyForm
            onAdd={async (name, type, options) => {
              try {
                await createProp.mutateAsync({ boardId: boardId!, name, type, options });
                toast('success', t('board.propertyAdded', { name }));
              } catch {
                toast('error', t('board.issueCreateFailed'));
              }
            }}
          />
        )}
      </div>

      <div className="flex flex-col gap-1">
        <form onSubmit={onCreate} className="flex gap-2 rounded-lg border border-edge-base bg-surface-subtle p-2">
          <Input
            placeholder={t('board.newIssueTitle')}
            value={title}
            onChange={(e) => setTitle(e.target.value)}
            aria-label={t('board.newIssueTitle')}
            className="flex-1 border-0 bg-transparent focus-visible:ring-0"
          />
          <button
            type="button"
            onClick={() => setTemplateOpen(true)}
            className="flex items-center gap-1 rounded-md px-2 py-1 text-xs text-ink-muted hover:bg-surface-muted hover:text-ink-primary transition-colors"
          >
            <FileText size={13} />
            {t('template.title')}
          </button>
          <Button type="submit" size="sm" disabled={createIssue.isPending}>
            {t('board.addIssue')}
          </Button>
        </form>
        {title.trim() && (
          <button
            type="button"
            onClick={() => {
              const name = title.trim();
              const stored = JSON.parse(localStorage.getItem(`mp-templates-${boardId}`) ?? '[]');
              const tmpl = { id: crypto.randomUUID(), name, title: name, body: '', instructions: '' };
              localStorage.setItem(`mp-templates-${boardId}`, JSON.stringify([...stored, tmpl]));
            }}
            className="self-end flex items-center gap-1 rounded-md px-2 py-0.5 text-xs text-ink-muted hover:bg-surface-muted hover:text-ink-primary transition-colors"
          >
            <FileText size={12} />
            {t('template.save')}
          </button>
        )}
      </div>

      {errMsg && (
        <div role="alert" className="flex items-center gap-2 rounded-md border-l-4 border-red-500 bg-red-50 px-4 py-2 text-sm text-red-700 dark:bg-red-950/40 dark:text-red-300">
          <AlertCircle size={16} className="shrink-0" />
          {errMsg}
        </div>
      )}

      {/* filter + view controls */}
      <div className="flex items-center justify-between gap-3">
        {/* filter bar on left */}
        <div className="flex flex-1 items-center gap-2">
          <div className="flex flex-1 items-center gap-2 rounded-md border border-edge-base bg-surface-subtle px-2 py-1.5">
            <Filter size={14} className="text-ink-muted" />
            <input
              type="text"
              value={filterText}
              onChange={(e) => setFilterText(e.target.value)}
              placeholder={t('board.filterPlaceholder', 'Filter issues...')}
              className="flex-1 bg-transparent text-sm text-ink-primary placeholder:text-ink-muted focus:outline-none"
            />
          </div>
          <select
            value={filterStatus}
            onChange={(e) => setFilterStatus(e.target.value as IssueStatus | 'all')}
            className="h-8 rounded-md border border-edge-base bg-surface-subtle px-2 text-xs text-ink-secondary focus:outline-none focus:border-brand-500"
          >
            <option value="all">{t('board.allStatuses', 'All')}</option>
            <option value="Pending">{t('status.Pending')}</option>
            <option value="Approved">{t('status.Approved')}</option>
            <option value="InProgress">{t('status.InProgress')}</option>
            <option value="QA">{t('status.QA')}</option>
            <option value="Done">{t('status.Done')}</option>
            <option value="Rejected">{t('status.Rejected')}</option>
          </select>
        </div>
        {/* view toggle on right */}
        <div className="flex items-center gap-2">
          <button
            type="button"
            onClick={() => setHideDone(h => !h)}
            className={`rounded px-2 py-1 text-xs transition-colors ${hideDone ? 'bg-brand-500/10 text-brand-500' : 'text-ink-muted hover:text-ink-secondary'}`}
          >
            {hideDone ? t('board.showDone', 'Show done') : t('board.hideDone', 'Hide done')}
          </button>
          <button
            type="button"
            onClick={() => setFullscreen(f => !f)}
            className="rounded px-2 py-1 text-xs text-ink-muted hover:text-ink-secondary transition-colors"
            title={fullscreen ? t('board.exitFullscreen', 'Exit fullscreen') : t('board.fullscreen', 'Fullscreen')}
          >
            {fullscreen ? <Minimize2 size={14} /> : <Maximize2 size={14} />}
          </button>
          <div className="flex items-center gap-1 rounded-md bg-surface-muted p-0.5">
            <button
              type="button"
              onClick={() => setViewMode('kanban')}
              className={`rounded px-2 py-1 text-xs transition-colors ${viewMode === 'kanban' ? 'bg-surface-base text-ink-primary shadow-sm' : 'text-ink-muted hover:text-ink-secondary'}`}
            >
              <LayoutGrid size={14} />
            </button>
            <button
              type="button"
              onClick={() => setViewMode('table')}
              className={`rounded px-2 py-1 text-xs transition-colors ${viewMode === 'table' ? 'bg-surface-base text-ink-primary shadow-sm' : 'text-ink-muted hover:text-ink-secondary'}`}
            >
              <List size={14} />
            </button>
          </div>
        </div>
      </div>

      {/* conditionally render kanban or table */}
      {viewMode === 'table' ? (
        <div className="overflow-x-auto rounded-lg border border-edge-base">
          <table className="w-full text-sm">
            <thead>
              <tr className="border-b border-edge-base bg-surface-subtle text-left text-xs font-medium text-ink-muted">
                <th className="px-4 py-2.5">{t('board.colTitle')}</th>
                <th className="px-4 py-2.5">{t('board.colStatus')}</th>
                <th className="px-4 py-2.5">{t('board.colCreated')}</th>
              </tr>
            </thead>
            <tbody>
              {filtered.filter(i => !(hideDone && i.status === 'Done')).map(issue => (
                <tr key={issue.id} className="border-b border-edge-base last:border-0 hover:bg-surface-subtle transition-colors">
                  <td className="px-4 py-2.5">
                    <Link to={`/issues/${issue.id}`} className="font-medium text-ink-primary hover:underline">
                      {issue.title}
                    </Link>
                  </td>
                  <td className="px-4 py-2.5">
                    <StatusBadge status={issue.status} />
                  </td>
                  <td className="px-4 py-2.5 text-xs text-ink-muted">
                    {new Date(issue.createdAt).toLocaleDateString()}
                  </td>
                </tr>
              ))}
              {filtered.length === 0 && (
                <tr>
                  <td colSpan={3} className="px-4 py-8 text-center text-sm text-ink-muted">
                    {t('board.noIssuesFound')}
                  </td>
                </tr>
              )}
            </tbody>
          </table>
        </div>
      ) : (
        <DndContext sensors={sensors} onDragStart={onDragStart} onDragEnd={onDragEnd}>
          <div className="flex gap-4 overflow-x-auto pb-2">
            {COLUMN_KEYS.filter(c => !(hideDone && c.status === 'Done')).map((c) => (
              <KanbanColumn
                key={c.status}
                status={c.status}
                title={t(c.key)}
                issues={grouped[c.status]}
                boardProperties={boardPropsQuery.data}
                onCreateIssue={async (issueTitle, status) => {
                  const created = await createIssue.mutateAsync({ boardId: boardId!, title: issueTitle });
                  if (status !== 'Pending' && created?.id) {
                    await updateIssue.mutateAsync({ id: created.id, patch: { status } });
                  }
                }}
                onReorder={(issueId, direction) => handleReorder(c.status, issueId, direction)}
              />
            ))}
          </div>
          <DragOverlay dropAnimation={null}>
            {activeIssue && (
              <div className="w-[280px] rotate-[2deg] rounded-lg border-2 border-brand-500 bg-surface-base p-3 shadow-2xl">
                <div className="flex items-start gap-2">
                  <div className="flex-1 min-w-0">
                    <span className="font-medium text-ink-primary">{activeIssue.title}</span>
                  </div>
                  <StatusBadge status={activeIssue.status} />
                </div>
              </div>
            )}
          </DragOverlay>
        </DndContext>
      )}

      <TemplateDialog
        boardId={boardId!}
        open={templateOpen}
        onClose={() => setTemplateOpen(false)}
        onSelect={(tmpl) => {
          setTitle(tmpl.title);
          setTemplateOpen(false);
        }}
      />

      <ConfirmDialog
        open={confirmDelete}
        title={t('board.delete')}
        description={t('board.deleteConfirm', { name: board?.name })}
        confirmLabel={t('board.delete')}
        onConfirm={async () => {
          setConfirmDelete(false);
          try {
            await deleteBoard.mutateAsync(boardId!);
            navigate('/boards');
            toast('success', t('board.deleted'));
          } catch (err) {
            toast('error', (err as { message?: string }).message ?? t('board.deleteFailed'));
          }
        }}
        onCancel={() => setConfirmDelete(false)}
      />

      {/* Webhook 설정 */}
      {boardId && <WebhookSettings boardId={boardId} />}
    </section>
  );
}
