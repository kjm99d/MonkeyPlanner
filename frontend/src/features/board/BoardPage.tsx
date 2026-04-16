import { useMemo, useState } from 'react';
import { useParams } from 'react-router-dom';
import { useTranslation } from 'react-i18next';
import { DndContext, type DragEndEvent, PointerSensor, useSensor, useSensors } from '@dnd-kit/core';
import { useBoards, useBoardProperties, useCreateBoardProperty, useCreateIssue, useIssues, useUpdateIssue } from '../../api/hooks';
import { Button } from '../../components/Button';
import { Input } from '../../components/Input';
import { Breadcrumb } from '../../components/Breadcrumb';
import { AddPropertyForm } from '../../components/PropertyEditor';
import { useToast } from '../../components/Toast';
import { WebhookSettings } from '../../components/WebhookSettings';
import { KanbanColumn } from './KanbanColumn';
import type { Issue, IssueStatus } from '../../api/types';

const COLUMN_KEYS: { status: IssueStatus; key: string }[] = [
  { status: 'Pending', key: 'kanban.pending' },
  { status: 'Approved', key: 'kanban.approved' },
  { status: 'InProgress', key: 'kanban.inProgress' },
  { status: 'Done', key: 'kanban.done' },
];

export default function BoardPage() {
  const { boardId } = useParams<{ boardId: string }>();
  const boards = useBoards();
  const board = boards.data?.find((b) => b.id === boardId);
  const issues = useIssues({ boardId });
  const createIssue = useCreateIssue();
  const boardPropsQuery = useBoardProperties(boardId);
  const createProp = useCreateBoardProperty();
  const updateIssue = useUpdateIssue();

  const { t } = useTranslation();
  const { toast } = useToast();
  const [title, setTitle] = useState('');
  const [errMsg, setErrMsg] = useState<string | null>(null);

  const sensors = useSensors(
    useSensor(PointerSensor, { activationConstraint: { distance: 6 } }),
  );

  const grouped = useMemo(() => {
    const map: Record<IssueStatus, Issue[]> = {
      Pending: [],
      Approved: [],
      InProgress: [],
      Done: [],
    };
    (issues.data ?? []).forEach((i) => map[i.status].push(i));
    return map;
  }, [issues.data]);

  async function onCreate(e: React.FormEvent) {
    e.preventDefault();
    if (!title.trim() || !boardId) return;
    await createIssue.mutateAsync({ boardId, title: title.trim() });
    setTitle('');
    toast('success', t('board.issueCreated'));
  }

  async function onDragEnd(e: DragEndEvent) {
    setErrMsg(null);
    if (!e.over) return;
    const toStatus = e.over.id as IssueStatus;
    const issueId = String(e.active.id);
    const current = (issues.data ?? []).find((i) => i.id === issueId);
    if (!current || current.status === toStatus) return;
    try {
      await updateIssue.mutateAsync({ id: issueId, patch: { status: toStatus } });
      toast('success', t('board.statusChanged', { status: toStatus }));
    } catch (err) {
      const msg =
        (err as { message?: string; code?: string })?.message ??
        (err as { code?: string })?.code ??
        '상태 변경 실패';
      setErrMsg(msg);
      toast('error', msg);
    }
  }

  return (
    <section className="flex flex-col gap-6">
      <header className="flex flex-col gap-2">
        <Breadcrumb items={[
          { label: '보드', to: '/boards' },
          { label: board?.name ?? '...' },
        ]} />
        <h1 className="text-3xl font-bold">{board?.name ?? '보드'}</h1>
      </header>

      {/* 속성 관리 */}
      <div className="flex flex-wrap items-center gap-2">
        {boardPropsQuery.data?.map((p) => (
          <span key={p.id} className="rounded-full border border-edge-base bg-surface-subtle px-2.5 py-0.5 text-xs text-ink-secondary">
            {p.name} · {p.type}
          </span>
        ))}
        {boardId && (
          <AddPropertyForm
            onAdd={(name, type, options) => {
              createProp.mutate({ boardId: boardId!, name, type, options });
              toast('success', `속성 "${name}" 추가됨`);
            }}
          />
        )}
      </div>

      <form onSubmit={onCreate} className="flex gap-2">
        <Input
          placeholder={t('board.newIssueTitle')}
          value={title}
          onChange={(e) => setTitle(e.target.value)}
          aria-label={t('board.newIssueTitle')}
          className="flex-1"
        />
        <Button type="submit" disabled={createIssue.isPending}>
          {t('board.addIssue')}
        </Button>
      </form>

      {errMsg && (
        <div role="alert" className="rounded-md border border-red-200 bg-red-50 px-4 py-2 text-sm text-red-700 dark:border-red-900/50 dark:bg-red-950/40 dark:text-red-300">
          {errMsg}
        </div>
      )}

      <DndContext sensors={sensors} onDragEnd={onDragEnd}>
        <div className="grid gap-4 lg:grid-cols-4 md:grid-cols-2">
          {COLUMN_KEYS.map((c) => (
            <KanbanColumn key={c.status} status={c.status} title={t(c.key)} issues={grouped[c.status]} />
          ))}
        </div>
      </DndContext>

      {/* Webhook 설정 */}
      {boardId && <WebhookSettings boardId={boardId} />}
    </section>
  );
}
