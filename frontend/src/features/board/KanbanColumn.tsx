import { useDroppable } from '@dnd-kit/core';
import type { Issue, IssueStatus } from '../../api/types';
import { IssueCard } from './IssueCard';

type Props = {
  status: IssueStatus;
  title: string;
  issues: Issue[];
};

export function KanbanColumn({ status, title, issues }: Props) {
  const { setNodeRef, isOver } = useDroppable({ id: status });
  return (
    <section
      ref={setNodeRef}
      aria-label={`${title} 컬럼`}
      className={`flex min-h-[24rem] flex-col gap-3 rounded-lg border border-edge-base bg-surface-subtle p-3 transition-colors ${
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
          <p className="py-8 text-center text-xs text-ink-muted">비어 있음</p>
        ) : (
          issues.map((iss) => <IssueCard key={iss.id} issue={iss} />)
        )}
      </div>
    </section>
  );
}
