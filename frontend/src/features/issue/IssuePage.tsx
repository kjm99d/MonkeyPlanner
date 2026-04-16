import { useEffect, useState } from 'react';
import { Link, useParams } from 'react-router-dom';
import {
  useApproveIssue,
  useBoards,
  useBoardProperties,
  useDeleteIssue,
  useIssue,
  useUpdateIssue,
  useUpdateIssueProperties,
} from '../../api/hooks';
import { Breadcrumb } from '../../components/Breadcrumb';
import { Button } from '../../components/Button';
import { Input } from '../../components/Input';
import { MarkdownEditor } from '../../components/MarkdownEditor';
import { PropertyEditor } from '../../components/PropertyEditor';
import { StatusBadge } from '../../components/StatusBadge';
import { StatusStepper } from '../../components/StatusStepper';
import type { IssueStatus } from '../../api/types';

export default function IssuePage() {
  const { issueId } = useParams<{ issueId: string }>();
  const query = useIssue(issueId);
  const boards = useBoards();
  const update = useUpdateIssue();
  const approve = useApproveIssue();
  const remove = useDeleteIssue();
  const updateProps = useUpdateIssueProperties();
  const boardProps = useBoardProperties(query.data?.issue.boardId);

  const [title, setTitle] = useState('');
  const [body, setBody] = useState('');
  const [saveErr, setSaveErr] = useState<string | null>(null);

  useEffect(() => {
    if (query.data?.issue) {
      setTitle(query.data.issue.title);
      setBody(query.data.issue.body);
    }
  }, [query.data?.issue?.id, query.data?.issue?.title, query.data?.issue?.body]);

  if (query.isLoading || !query.data) {
    return <p className="text-ink-secondary">이슈를 불러오는 중…</p>;
  }

  const iss = query.data.issue;

  async function onSave() {
    if (!issueId) return;
    setSaveErr(null);
    try {
      await update.mutateAsync({ id: issueId, patch: { title, body } });
    } catch (err) {
      setSaveErr((err as { message?: string }).message ?? '저장 실패');
    }
  }

  async function setStatus(status: IssueStatus) {
    if (!issueId) return;
    setSaveErr(null);
    try {
      if (status === 'Approved') {
        await approve.mutateAsync(issueId);
      } else {
        await update.mutateAsync({ id: issueId, patch: { status } });
      }
    } catch (err) {
      setSaveErr((err as { message?: string }).message ?? '상태 변경 실패');
    }
  }

  const board = boards.data?.find((b) => b.id === iss.boardId);

  return (
    <section className="flex flex-col gap-6">
      <Breadcrumb items={[
        { label: '보드', to: '/boards' },
        { label: board?.name ?? '...', to: `/boards/${iss.boardId}` },
        { label: iss.title },
      ]} />

      <header className="flex flex-col gap-4">
        <div className="flex items-center justify-between">
          <div className="flex items-center gap-3">
            <StatusBadge status={iss.status} />
            <span className="text-xs text-ink-muted">
              {new Date(iss.createdAt).toLocaleString('ko-KR')} 생성
            </span>
          </div>
          <div className="flex gap-2">
            {iss.status === 'Pending' && (
              <Button size="sm" onClick={() => setStatus('Approved')}>
                ✓ Approve
              </Button>
            )}
          </div>
        </div>
        <StatusStepper
          current={iss.status}
          onSelect={setStatus}
          disabled={update.isPending || approve.isPending}
        />
      </header>

      <Input
        value={title}
        onChange={(e) => setTitle(e.target.value)}
        label="제목"
        className="text-lg font-semibold"
      />

      {/* 커스텀 속성 — 본문 위에 배치 */}
      {boardProps.data && boardProps.data.length > 0 && (
        <PropertyEditor
          properties={boardProps.data}
          values={query.data?.issue.properties ?? {}}
          onChange={(propId, value) => {
            if (!issueId) return;
            updateProps.mutate({ id: issueId, properties: { [propId]: value } });
          }}
        />
      )}

      <MarkdownEditor value={body} onChange={setBody} />

      <div className="flex items-center gap-3">
        <Button onClick={onSave} disabled={update.isPending}>
          {update.isPending ? '저장 중…' : '저장'}
        </Button>
        {saveErr && (
          <span role="alert" className="text-sm text-red-600">
            {saveErr}
          </span>
        )}
      </div>

      {query.data.children.length > 0 && (
        <section aria-label="자식 이슈" className="flex flex-col gap-2">
          <h2 className="text-lg font-semibold">자식 이슈</h2>
          <ul className="flex flex-col gap-2">
            {query.data.children.map((c) => (
              <li key={c.id}>
                <Link
                  to={`/issues/${c.id}`}
                  className="flex items-center justify-between rounded-md border border-edge-base bg-surface-subtle px-3 py-2 hover:bg-surface-muted focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-brand-500"
                >
                  <span>{c.title}</span>
                  <StatusBadge status={c.status} />
                </Link>
              </li>
            ))}
          </ul>
        </section>
      )}

      <hr className="border-edge-base" />
      <div className="flex justify-end">
        <button
          type="button"
          onClick={async () => {
            if (!issueId) return;
            if (!window.confirm('이슈와 모든 자식 이슈가 삭제됩니다. 계속할까요?')) return;
            await remove.mutateAsync(issueId);
            window.history.back();
          }}
          className="rounded-md px-3 py-1.5 text-xs text-ink-muted transition-colors hover:bg-red-50 hover:text-red-600 dark:hover:bg-red-950/30 dark:hover:text-red-400"
        >
          이 이슈 삭제
        </button>
      </div>
    </section>
  );
}
