import { useEffect, useState } from 'react';
import { Link, useParams, useNavigate } from 'react-router-dom';
import { useTranslation } from 'react-i18next';
import {
  useAddDependency,
  useApproveIssue,
  useBoards,
  useBoardProperties,
  useDeleteIssue,
  useIssue,
  useRemoveDependency,
  useUpdateIssue,
  useUpdateIssueProperties,
} from '../../api/hooks';
import { Breadcrumb } from '../../components/Breadcrumb';
import { Button } from '../../components/Button';
import { ConfirmDialog } from '../../components/ConfirmDialog';
import { Input } from '../../components/Input';
import { MarkdownEditor } from '../../components/MarkdownEditor';
import { PropertyEditor } from '../../components/PropertyEditor';
import { CommentSection } from '../../components/CommentSection';
import { Skeleton } from '../../components/Skeleton';
import { StatusBadge } from '../../components/StatusBadge';
import { StatusStepper } from '../../components/StatusStepper';
import type { Criterion, IssueStatus } from '../../api/types';

export default function IssuePage() {
  const { t } = useTranslation();
  const { issueId } = useParams<{ issueId: string }>();
  const navigate = useNavigate();
  const query = useIssue(issueId);
  const boards = useBoards();
  const update = useUpdateIssue();
  const approve = useApproveIssue();
  const remove = useDeleteIssue();
  const updateProps = useUpdateIssueProperties();
  const boardProps = useBoardProperties(query.data?.issue.boardId);

  const addDep = useAddDependency();
  const removeDep = useRemoveDependency();

  const [title, setTitle] = useState('');
  const [body, setBody] = useState('');
  const [instructions, setInstructions] = useState('');
  const [criteria, setCriteria] = useState<Criterion[]>([]);
  const [saveErr, setSaveErr] = useState<string | null>(null);
  const [confirmDelete, setConfirmDelete] = useState(false);
  const [depInput, setDepInput] = useState('');

  useEffect(() => {
    if (query.data?.issue) {
      setTitle(query.data.issue.title);
      setBody(query.data.issue.body);
      setInstructions(query.data.issue.instructions ?? '');
      setCriteria(query.data.issue.criteria ?? []);
    }
  }, [query.data?.issue?.id, query.data?.issue?.title, query.data?.issue?.body, query.data?.issue?.instructions]);

  // Track dirty state for unsaved changes warning
  const isDirty =
    query.data?.issue != null &&
    (title !== query.data.issue.title || body !== query.data.issue.body);

  useEffect(() => {
    if (!isDirty) return;
    const handler = (e: BeforeUnloadEvent) => {
      e.preventDefault();
    };
    window.addEventListener('beforeunload', handler);
    return () => window.removeEventListener('beforeunload', handler);
  }, [isDirty]);

  if (query.isLoading || !query.data) {
    return (
      <section className="flex flex-col gap-6">
        <Skeleton className="h-4 w-48" />
        <Skeleton className="h-8 w-3/4" />
        <Skeleton className="h-4 w-32" />
        <Skeleton className="h-40 w-full" />
      </section>
    );
  }

  const iss = query.data.issue;

  async function onSave() {
    if (!issueId) return;
    setSaveErr(null);
    try {
      await update.mutateAsync({ id: issueId, patch: { title, body, instructions, criteria } });
    } catch (err) {
      setSaveErr((err as { message?: string }).message ?? t('issue.saveFailed'));
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
      setSaveErr((err as { message?: string }).message ?? t('issue.statusChangeFailed'));
    }
  }

  const board = boards.data?.find((b) => b.id === iss.boardId);

  return (
    <section className="flex flex-col gap-6">
      <Breadcrumb items={[
        { label: t('nav.boards'), to: '/boards' },
        { label: board?.name ?? '...', to: `/boards/${iss.boardId}` },
        { label: iss.title },
      ]} />

      <header className="flex flex-col gap-4">
        <div className="flex items-center justify-between">
          <div className="flex items-center gap-3">
            <StatusBadge status={iss.status} />
            <span className="text-xs text-ink-muted">
              {t('issue.createdAt', { date: new Date(iss.createdAt).toLocaleString() })}
            </span>
          </div>
          <div className="flex gap-2">
            {iss.status === 'Pending' && (
              <Button size="sm" onClick={() => setStatus('Approved')}>
                {t('kanban.approve')}
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
        label={t('issue.title')}
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

      <section className="flex flex-col gap-2">
        <div className="flex items-center gap-2">
          <span className="text-sm font-medium text-ink-secondary">{t('issue.instructions')}</span>
          <span className="rounded bg-brand-500/10 px-1.5 py-0.5 text-[10px] text-brand-500">MCP</span>
        </div>
        <textarea
          value={instructions}
          onChange={(e) => setInstructions(e.target.value)}
          placeholder={t('issue.instructionsPlaceholder')}
          className="min-h-[80px] rounded-md border border-edge-base bg-surface-subtle p-3 font-mono text-sm text-ink-primary focus-visible:border-brand-500 focus-visible:outline-none"
        />
      </section>

      <section className="flex flex-col gap-2">
        <h3 className="text-sm font-medium text-ink-secondary">{t('issue.criteria')}</h3>
        {criteria.map((c, i) => (
          <label key={i} className="flex items-center gap-2 text-sm">
            <input
              type="checkbox"
              checked={c.done}
              onChange={() => {
                const next = [...criteria];
                next[i] = { ...next[i], done: !next[i].done };
                setCriteria(next);
              }}
              className="h-4 w-4 rounded border-edge-base accent-brand-500"
            />
            <span className={c.done ? 'line-through text-ink-muted' : 'text-ink-primary'}>{c.text}</span>
            <button type="button" onClick={() => setCriteria(criteria.filter((_, j) => j !== i))}
              className="ml-auto text-ink-muted hover:text-red-500 text-xs">×</button>
          </label>
        ))}
        <div className="flex gap-1">
          <input
            placeholder={t('issue.addCriterion')}
            onKeyDown={(e) => {
              if (e.key === 'Enter' && e.currentTarget.value.trim()) {
                setCriteria([...criteria, { text: e.currentTarget.value.trim(), done: false }]);
                e.currentTarget.value = '';
              }
            }}
            className="flex-1 rounded-md border border-edge-base bg-surface-base px-2 py-1 text-sm focus:outline-none focus:border-brand-500"
          />
        </div>
      </section>

      <MarkdownEditor value={body} onChange={setBody} />

      <div className="flex items-center gap-3">
        <Button onClick={onSave} disabled={update.isPending}>
          {update.isPending ? t('issue.saving') : t('issue.save')}
        </Button>
        {saveErr && (
          <span role="alert" className="text-sm text-red-600">
            {saveErr}
          </span>
        )}
      </div>

      {query.data.children.length > 0 && (
        <section aria-label={t('issue.children')} className="flex flex-col gap-2">
          <h2 className="text-lg font-semibold">{t('issue.children')}</h2>
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

      <section aria-label={t('issue.dependencies')} className="flex flex-col gap-2">
        <h2 className="text-lg font-semibold">{t('issue.dependencies')}</h2>
        <p className="text-sm text-ink-secondary">{t('issue.blockedBy')}</p>
        {(iss.blockedBy ?? []).length > 0 ? (
          <ul className="flex flex-col gap-2">
            {(iss.blockedBy ?? []).map((blockerId) => (
              <li key={blockerId} className="flex items-center justify-between rounded-md border border-edge-base bg-surface-subtle px-3 py-2">
                <Link
                  to={`/issues/${blockerId}`}
                  className="text-sm text-brand-500 hover:underline focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-brand-500"
                >
                  {blockerId}
                </Link>
                <button
                  type="button"
                  onClick={() => issueId && removeDep.mutate({ issueId, blockerId })}
                  className="rounded px-2 py-0.5 text-xs text-ink-muted transition-colors hover:bg-red-50 hover:text-red-600 dark:hover:bg-red-950/30 dark:hover:text-red-400"
                >
                  {t('issue.removeDependency')}
                </button>
              </li>
            ))}
          </ul>
        ) : null}
        <div className="flex gap-2">
          <input
            value={depInput}
            onChange={(e) => setDepInput(e.target.value)}
            placeholder={t('issue.addDependency')}
            className="flex-1 rounded-md border border-edge-base bg-surface-subtle px-3 py-1.5 text-sm text-ink-primary focus-visible:border-brand-500 focus-visible:outline-none"
          />
          <Button
            size="sm"
            disabled={!depInput.trim() || addDep.isPending}
            onClick={() => {
              if (!issueId || !depInput.trim()) return;
              addDep.mutate({ issueId, blockerId: depInput.trim() }, {
                onSuccess: () => setDepInput(''),
              });
            }}
          >
            {t('issue.addDependency')}
          </Button>
        </div>
      </section>

      <CommentSection issueId={issueId!} />

      <hr className="border-edge-base" />
      <div className="flex justify-end">
        <button
          type="button"
          onClick={() => setConfirmDelete(true)}
          className="rounded-md px-3 py-1.5 text-xs text-ink-muted transition-colors hover:bg-red-50 hover:text-red-600 dark:hover:bg-red-950/30 dark:hover:text-red-400"
        >
          {t('issue.delete')}
        </button>
      </div>
      <ConfirmDialog
        open={confirmDelete}
        title={t('issue.delete')}
        description={t('issue.deleteConfirm')}
        confirmLabel={t('issue.delete')}
        onConfirm={async () => {
          setConfirmDelete(false);
          if (!issueId) return;
          await remove.mutateAsync(issueId);
          navigate(`/boards/${iss.boardId}`);
        }}
        onCancel={() => setConfirmDelete(false)}
      />
    </section>
  );
}
