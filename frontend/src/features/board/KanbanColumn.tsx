import { useDroppable } from '@dnd-kit/core';
import { useTranslation } from 'react-i18next';
import { ClipboardList } from 'lucide-react';
import type { Issue, IssueStatus } from '../../api/types';
import { IssueCard } from './IssueCard';

const topAccent: Record<IssueStatus, string> = {
  Pending: 'border-t-amber-400',
  Approved: 'border-t-emerald-500',
  InProgress: 'border-t-blue-500',
  Done: 'border-t-violet-400',
};

type Props = {
  status: IssueStatus;
  title: string;
  issues: Issue[];
};

export function KanbanColumn({ status, title, issues }: Props) {
  const { t } = useTranslation();
  const { setNodeRef, isOver } = useDroppable({ id: status });
  return (
    <section
      ref={setNodeRef}
      aria-label={`${title} 컬럼`}
      className={`flex min-h-[24rem] flex-col gap-3 rounded-lg border border-t-[3px] ${topAccent[status]} border-edge-base bg-surface-subtle p-3 transition-colors ${
        isOver ? 'border-brand-500 bg-brand-500/5' : ''
      }`}
    >
      <header className="flex items-center justify-between">
        <h2 className="text-sm font-semibold text-ink-secondary">{title}</h2>
        <span className="rounded-full bg-surface-muted px-2 py-0.5 text-xs tabular-nums text-ink-secondary">
          {issues.length}
        </span>
      </header>
      <div className="flex flex-col gap-2">
        {issues.length === 0 ? (
          <div className="flex flex-col items-center gap-2 py-10 text-center">
            <ClipboardList size={28} className="text-ink-muted opacity-25" />
            <p className="text-xs text-ink-muted">{t('kanban.dragHere')}</p>
          </div>
        ) : (
          issues.map((iss) => <IssueCard key={iss.id} issue={iss} />)
        )}
      </div>
    </section>
  );
}
