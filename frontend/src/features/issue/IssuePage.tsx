import { useEffect, useState } from 'react';
import { Link, useParams } from 'react-router-dom';
import {
  useApproveIssue,
  useDeleteIssue,
  useIssue,
  useUpdateIssue,
} from '../../api/hooks';
import { Button } from '../../components/Button';
import { Input } from '../../components/Input';
import { MarkdownEditor } from '../../components/MarkdownEditor';
import { StatusBadge } from '../../components/StatusBadge';
import type { IssueStatus } from '../../api/types';

export default function IssuePage() {
  const { issueId } = useParams<{ issueId: string }>();
  const query = useIssue(issueId);
  const update = useUpdateIssue();
  const approve = useApproveIssue();
  const remove = useDeleteIssue();

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

  const nextButtons: Array<{ status: IssueStatus; label: string }> = [];
  if (iss.status === 'Pending') nextButtons.push({ status: 'Approved', label: '승인' });
  if (iss.status === 'Approved') nextButtons.push({ status: 'InProgress', label: '진행' });
  if (iss.status === 'InProgress') nextButtons.push({ status: 'Done', label: '완료' });

  return (
    <section className="flex flex-col gap-6">
      <header className="flex flex-wrap items-center justify-between gap-3">
        <div className="flex items-center gap-3">
          <StatusBadge status={iss.status} />
          <span className="text-xs text-ink-muted">
            {new Date(iss.createdAt).toLocaleString('ko-KR')} 생성
          </span>
        </div>
        <div className="flex flex-wrap gap-2">
          {nextButtons.map((b) => (
            <Button key={b.status} size="sm" onClick={() => setStatus(b.status)}>
              {b.label}
            </Button>
          ))}
          <Button
            size="sm"
            variant="danger"
            onClick={async () => {
              if (!issueId) return;
              if (!window.confirm('이슈와 모든 자식 이슈가 삭제됩니다. 계속할까요?')) return;
              await remove.mutateAsync(issueId);
              window.history.back();
            }}
          >
            삭제
          </Button>
        </div>
      </header>

      <Input
        value={title}
        onChange={(e) => setTitle(e.target.value)}
        label="제목"
        className="text-lg font-semibold"
      />

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
    </section>
  );
}
