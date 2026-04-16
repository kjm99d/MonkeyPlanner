import { useMemo, useState } from 'react';
import { useParams } from 'react-router-dom';
import { DndContext, type DragEndEvent, PointerSensor, useSensor, useSensors } from '@dnd-kit/core';
import { useBoards, useCreateIssue, useIssues, useUpdateIssue } from '../../api/hooks';
import { Button } from '../../components/Button';
import { Input } from '../../components/Input';
import { Breadcrumb } from '../../components/Breadcrumb';
import { useToast } from '../../components/Toast';
import { KanbanColumn } from './KanbanColumn';
import type { Issue, IssueStatus } from '../../api/types';

const COLUMNS: { status: IssueStatus; title: string }[] = [
  { status: 'Pending', title: '대기' },
  { status: 'Approved', title: '승인됨' },
  { status: 'InProgress', title: '진행 중' },
  { status: 'Done', title: '완료' },
];

export default function BoardPage() {
  const { boardId } = useParams<{ boardId: string }>();
  const boards = useBoards();
  const board = boards.data?.find((b) => b.id === boardId);
  const issues = useIssues({ boardId });
  const createIssue = useCreateIssue();
  const updateIssue = useUpdateIssue();

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
    toast('success', '이슈가 생성되었습니다');
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
      toast('success', `상태가 "${toStatus}"로 변경되었습니다`);
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

      <form onSubmit={onCreate} className="flex gap-2">
        <Input
          placeholder="새 이슈 제목"
          value={title}
          onChange={(e) => setTitle(e.target.value)}
          aria-label="새 이슈 제목"
          className="flex-1"
        />
        <Button type="submit" disabled={createIssue.isPending}>
          이슈 추가
        </Button>
      </form>

      {errMsg && (
        <div role="alert" className="rounded-md border border-red-200 bg-red-50 px-4 py-2 text-sm text-red-700 dark:border-red-900/50 dark:bg-red-950/40 dark:text-red-300">
          {errMsg}
        </div>
      )}

      <DndContext sensors={sensors} onDragEnd={onDragEnd}>
        <div className="grid gap-4 lg:grid-cols-4 md:grid-cols-2">
          {COLUMNS.map((c) => (
            <KanbanColumn key={c.status} status={c.status} title={c.title} issues={grouped[c.status]} />
          ))}
        </div>
      </DndContext>
    </section>
  );
}
