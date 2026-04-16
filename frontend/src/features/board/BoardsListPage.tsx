import { useState, memo } from 'react';
import { Link } from 'react-router-dom';
import { useTranslation } from 'react-i18next';
import { Plus, AlertCircle } from 'lucide-react';
import { useBoards, useCreateBoard, useIssues } from '../../api/hooks';
import { Button } from '../../components/Button';
import { Card } from '../../components/Card';
import { Input } from '../../components/Input';
import { Skeleton } from '../../components/Skeleton';
import { useToast } from '../../components/Toast';

const BoardCardWithCount = memo(function BoardCardWithCount({ id, name, viewType, createdAt }: { id: string; name: string; viewType: string; createdAt: string }) {
  const { t } = useTranslation();
  const issues = useIssues({ boardId: id });
  const total = issues.data?.length ?? 0;
  const done = issues.data?.filter((i) => i.status === 'Done').length ?? 0;
  const inProg = issues.data?.filter((i) => i.status === 'InProgress').length ?? 0;
  const pct = total > 0 ? Math.round((done / total) * 100) : 0;

  return (
    <Card className="relative overflow-hidden">
      <div className="flex items-start justify-between">
        <h2 className="text-lg font-semibold">{name}</h2>
        {total > 0 && (
          <span className="rounded-full bg-brand-500/10 px-2 py-0.5 text-xs font-medium tabular-nums text-brand-500">
            {total}
          </span>
        )}
      </div>
      <p className="mt-1 text-xs text-ink-muted">
        {new Date(createdAt).toLocaleDateString()} · {viewType}
      </p>
      {total > 0 && (
        <div className="mt-3 flex flex-col gap-1.5">
          <div className="flex items-center justify-between text-xs">
            <span className="text-ink-secondary">
              {t('board.progress', { done, total })}
              {inProg > 0 && <span className="ml-1.5 text-status-inProgress">{t('board.inProgressCount', { count: inProg })}</span>}
            </span>
            <span className="font-medium tabular-nums text-ink-primary">{pct}%</span>
          </div>
          <div className="h-1.5 w-full overflow-hidden rounded-full bg-surface-muted">
            <div
              className="h-full rounded-full bg-status-done transition-all duration-500"
              style={{ width: `${pct}%` }}
            />
          </div>
        </div>
      )}
      {total === 0 && (
        <p className="mt-3 text-xs text-ink-muted">{t('board.noIssues')}</p>
      )}
    </Card>
  );
});

export default function BoardsListPage() {
  const { t } = useTranslation();
  const [name, setName] = useState('');
  const boards = useBoards();
  const createBoard = useCreateBoard();
  const { toast } = useToast();

  const submit = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!name.trim()) return;
    try {
      await createBoard.mutateAsync({ name: name.trim() });
      setName('');
      toast('success', t('board.issueCreated'));
    } catch {
      toast('error', t('board.issueCreateFailed'));
    }
  };

  return (
    <section className="flex flex-col gap-6">
      <header className="fade-up flex items-end justify-between">
        <div>
          <h1 className="text-2xl font-bold">{t('board.title')}</h1>
          <p className="mt-1 text-ink-secondary">{t('board.description')}</p>
        </div>
      </header>

      <form onSubmit={submit} className="fade-up flex gap-2" style={{ animationDelay: '60ms' }}>
        <Input
          placeholder={t('board.newName')}
          value={name}
          onChange={(e) => setName(e.target.value)}
          aria-label={t('board.newName')}
          className="flex-1"
        />
        <Button type="submit" disabled={createBoard.isPending}>
          <Plus size={16} /> {t('board.create')}
        </Button>
      </form>

      {boards.isLoading && (
        <div className="grid gap-3 sm:grid-cols-2 lg:grid-cols-3">
          <Skeleton className="h-32 rounded-lg" />
          <Skeleton className="h-32 rounded-lg" />
          <Skeleton className="h-32 rounded-lg" />
        </div>
      )}
      {boards.isError && (
        <div role="alert" className="flex items-center gap-2 rounded-md border-l-4 border-red-500 bg-red-50 px-4 py-3 text-sm text-red-700 dark:bg-red-950/40 dark:text-red-300">
          <AlertCircle size={16} className="shrink-0" />
          <span>{t('common.errorLoad')}</span>
        </div>
      )}
      {boards.data && boards.data.length === 0 && (
        <div className="fade-up flex flex-col items-center gap-4 py-16 text-center" style={{ animationDelay: '120ms' }}>
          <div className="flex h-16 w-16 items-center justify-center rounded-2xl bg-brand-500/10">
            <Plus size={28} className="text-brand-500" />
          </div>
          <div>
            <p className="text-lg font-semibold text-ink-primary">{t('board.noBoards')}</p>
            <p className="mt-1 text-sm text-ink-secondary">{t('board.noBoardsHint')}</p>
          </div>
        </div>
      )}
      <ul className="grid gap-3 sm:grid-cols-2 lg:grid-cols-3">
        {boards.data?.map((b, i) => (
          <li key={b.id} className="fade-up" style={{ animationDelay: `${120 + i * 50}ms` }}>
            <Link
              to={`/boards/${b.id}`}
              className="block rounded-lg focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-brand-500"
            >
              <BoardCardWithCount id={b.id} name={b.name} viewType={b.viewType} createdAt={b.createdAt} />
            </Link>
          </li>
        ))}
      </ul>
    </section>
  );
}
