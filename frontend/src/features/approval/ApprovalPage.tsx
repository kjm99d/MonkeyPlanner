import { useState, useMemo } from 'react';
import { Link } from 'react-router-dom';
import { useTranslation } from 'react-i18next';
import { CheckCircle2, XCircle, Inbox, AlertCircle } from 'lucide-react';
import { useIssues, useApproveIssue, useBoards } from '../../api/hooks';
import { Button } from '../../components/Button';
import { StatusBadge } from '../../components/StatusBadge';
import { Skeleton } from '../../components/Skeleton';

function relativeTime(iso: string, locale: string): string {
  const diff = Date.now() - new Date(iso).getTime();
  const minutes = Math.floor(diff / 60_000);
  const hours = Math.floor(diff / 3_600_000);
  const days = Math.floor(diff / 86_400_000);
  try {
    const rtf = new Intl.RelativeTimeFormat(locale, { numeric: 'auto' });
    if (days > 0) return rtf.format(-days, 'day');
    if (hours > 0) return rtf.format(-hours, 'hour');
    return rtf.format(-minutes, 'minute');
  } catch {
    if (days > 0) return `${days}d ago`;
    if (hours > 0) return `${hours}h ago`;
    return `${minutes}m ago`;
  }
}

export default function ApprovalPage() {
  const { t, i18n } = useTranslation();
  const issues = useIssues({ status: 'Pending' });
  const boards = useBoards();
  const approveIssue = useApproveIssue();

  const [selected, setSelected] = useState<Set<string>>(new Set());
  const [approving, setApproving] = useState<Set<string>>(new Set());

  const boardMap = useMemo(() => {
    const m = new Map<string, string>();
    boards.data?.forEach((b) => m.set(b.id, b.name));
    return m;
  }, [boards.data]);

  const pendingIssues = issues.data ?? [];

  const toggleSelect = (id: string) => {
    setSelected((prev) => {
      const next = new Set(prev);
      if (next.has(id)) next.delete(id);
      else next.add(id);
      return next;
    });
  };

  const toggleAll = () => {
    if (selected.size === pendingIssues.length) {
      setSelected(new Set());
    } else {
      setSelected(new Set(pendingIssues.map((i) => i.id)));
    }
  };

  const handleApprove = async (id: string) => {
    setApproving((prev) => new Set(prev).add(id));
    try {
      await approveIssue.mutateAsync(id);
    } finally {
      setApproving((prev) => {
        const next = new Set(prev);
        next.delete(id);
        return next;
      });
      setSelected((prev) => {
        const next = new Set(prev);
        next.delete(id);
        return next;
      });
    }
  };

  const handleApproveSelected = async () => {
    const ids = Array.from(selected);
    await Promise.allSettled(ids.map((id) => handleApprove(id)));
  };

  const clearSelection = () => setSelected(new Set());

  const isLoading = issues.isLoading || boards.isLoading;
  const isError = issues.isError || boards.isError;
  const allSelected = pendingIssues.length > 0 && selected.size === pendingIssues.length;
  const someSelected = selected.size > 0;

  return (
    <section className="flex flex-col gap-6">
      <header className="fade-up flex items-center gap-3">
        <div className="flex-1">
          <div className="flex items-center gap-2">
            <h1 className="text-2xl font-bold">{t('approval.title')}</h1>
            {!isLoading && pendingIssues.length > 0 && (
              <span className="rounded-full bg-status-pending/15 px-2.5 py-0.5 text-xs font-semibold tabular-nums text-status-pending border border-status-pending/30">
                {pendingIssues.length}
              </span>
            )}
          </div>
          <p className="mt-1 text-sm text-ink-secondary">{t('approval.description')}</p>
        </div>
      </header>

      {isLoading && (
        <div className="fade-up flex flex-col gap-2" style={{ animationDelay: '60ms' }}>
          {[...Array(3)].map((_, i) => (
            <Skeleton key={i} className="h-14 rounded-lg" />
          ))}
        </div>
      )}

      {isError && (
        <div role="alert" className="flex items-center gap-2 rounded-md border-l-4 border-red-500 bg-red-50 px-4 py-3 text-sm text-red-700 dark:bg-red-950/40 dark:text-red-300">
          <AlertCircle size={16} className="shrink-0" />
          <span>{t('common.errorLoad')}</span>
        </div>
      )}

      {!isLoading && !isError && pendingIssues.length === 0 && (
        <div className="fade-up flex flex-col items-center gap-4 py-20 text-center" style={{ animationDelay: '60ms' }}>
          <div className="flex h-16 w-16 items-center justify-center rounded-2xl bg-surface-subtle border border-edge-base">
            <Inbox size={28} className="text-ink-muted" />
          </div>
          <div>
            <p className="text-lg font-semibold text-ink-primary">{t('approval.empty')}</p>
            <p className="mt-1 text-sm text-ink-secondary">{t('approval.emptyHint')}</p>
          </div>
        </div>
      )}

      {!isLoading && !isError && pendingIssues.length > 0 && (
        <div className="fade-up flex flex-col gap-1" style={{ animationDelay: '60ms' }}>
          <div className="mb-1 flex items-center gap-3 rounded-lg border border-edge-base bg-surface-subtle px-4 py-2">
            <input
              type="checkbox"
              id="select-all"
              checked={allSelected}
              onChange={toggleAll}
              aria-label={t('approval.selectAll')}
              className="h-4 w-4 rounded accent-brand-500"
            />
            <label htmlFor="select-all" className="text-xs font-medium text-ink-secondary select-none cursor-pointer">
              {t('approval.selectAll')}
            </label>
          </div>

          <ul className="flex flex-col gap-1" role="list" aria-label={t('approval.queueLabel')}>
            {pendingIssues.map((issue, i) => {
              const boardName = boardMap.get(issue.boardId) ?? issue.boardId;
              const isChecked = selected.has(issue.id);
              const isApprovingThis = approving.has(issue.id);
              return (
                <li
                  key={issue.id}
                  className="fade-up flex items-center gap-3 rounded-lg border border-edge-base bg-surface-subtle px-4 py-3 transition-colors hover:bg-surface-muted"
                  style={{ animationDelay: `${80 + i * 40}ms` }}
                >
                  <input
                    type="checkbox"
                    checked={isChecked}
                    onChange={() => toggleSelect(issue.id)}
                    aria-label={t('approval.selectIssue', { title: issue.title })}
                    className="h-4 w-4 shrink-0 rounded accent-brand-500"
                  />

                  <div className="min-w-0 flex-1">
                    <Link
                      to={`/issues/${issue.id}`}
                      className="truncate text-sm font-medium text-ink-primary hover:text-brand-600 focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-brand-500 rounded"
                    >
                      {issue.title}
                    </Link>
                  </div>

                  <span className="shrink-0 rounded-full border border-edge-base bg-surface-muted px-2 py-0.5 text-xs text-ink-secondary">
                    {boardName}
                  </span>

                  <StatusBadge status={issue.status} />

                  <span className="shrink-0 text-xs text-ink-muted tabular-nums">
                    {relativeTime(issue.createdAt, i18n.language)}
                  </span>

                  <Button
                    variant="ghost"
                    size="sm"
                    onClick={() => handleApprove(issue.id)}
                    disabled={isApprovingThis}
                    aria-label={t('approval.approveIssue', { title: issue.title })}
                    className="shrink-0 border-green-300 text-green-700 hover:border-green-400 hover:bg-green-50 hover:text-green-800 dark:border-green-700 dark:text-green-400 dark:hover:bg-green-950/40"
                  >
                    <CheckCircle2 size={15} />
                    {t('approval.approve')}
                  </Button>
                </li>
              );
            })}
          </ul>
        </div>
      )}

      {someSelected && (
        <div className="fixed bottom-6 left-1/2 z-50 flex -translate-x-1/2 items-center gap-3 rounded-xl border border-edge-base bg-surface-base px-5 py-3 shadow-lg">
          <span className="text-sm font-medium text-ink-primary">
            {t('approval.selectedCount', { count: selected.size })}
          </span>
          <Button
            size="sm"
            onClick={handleApproveSelected}
            disabled={approveIssue.isPending}
            className="border-green-300 bg-green-600 text-white hover:bg-green-700 shadow-none"
          >
            <CheckCircle2 size={15} />
            {t('approval.approveSelected', { count: selected.size })}
          </Button>
          <Button variant="ghost" size="sm" onClick={clearSelection}>
            <XCircle size={15} />
            {t('approval.clearSelection')}
          </Button>
        </div>
      )}
    </section>
  );
}
