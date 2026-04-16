import { useState } from 'react';
import { useDroppable } from '@dnd-kit/core';
import { useTranslation } from 'react-i18next';
import { ClipboardList, Plus } from 'lucide-react';
import type { Issue, IssueStatus, BoardProperty } from '../../api/types';
import { IssueCard } from './IssueCard';

const topAccent: Record<IssueStatus, string> = {
  Pending: 'border-t-amber-400',
  Approved: 'border-t-emerald-500',
  InProgress: 'border-t-blue-500',
  Done: 'border-t-violet-400',
  Rejected: 'border-t-red-500',
};

type Props = {
  status: IssueStatus;
  title: string;
  issues: Issue[];
  boardProperties?: BoardProperty[];
  onCreateIssue?: (title: string, status: IssueStatus) => void;
  onReorder?: (issueId: string, direction: 'up' | 'down') => void;
};

export function KanbanColumn({ status, title, issues, boardProperties, onCreateIssue, onReorder }: Props) {
  const { t } = useTranslation();
  const { setNodeRef, isOver } = useDroppable({ id: status });
  const [adding, setAdding] = useState(false);
  const [newTitle, setNewTitle] = useState('');

  return (
    <section
      ref={setNodeRef}
      aria-label={t('kanban.columnLabel', { title })}
      className={`min-w-[280px] w-[280px] shrink-0 lg:min-w-0 lg:w-auto lg:shrink lg:flex-1 flex min-h-[24rem] flex-col gap-3 rounded-lg border border-t-[3px] ${topAccent[status]} border-edge-base bg-surface-subtle p-3 transition-colors ${
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
          <div className="flex flex-1 flex-col items-center justify-center gap-2 rounded-md border-2 border-dashed border-edge-base/60 py-8 text-center">
            <ClipboardList size={24} className="text-ink-muted opacity-30" />
            <p className="text-xs text-ink-muted">{t('kanban.dragHere')}</p>
          </div>
        ) : (
          issues.map((iss, idx) => (
            <div key={iss.id} className="group/card relative">
              <IssueCard issue={iss} boardProperties={boardProperties} />
              {onReorder && (
                <div className="absolute right-1 top-1 flex flex-col gap-0.5 opacity-0 group-hover/card:opacity-100 transition-opacity">
                  <button
                    type="button"
                    disabled={idx === 0}
                    onClick={() => onReorder(iss.id, 'up')}
                    className="flex h-5 w-5 items-center justify-center rounded bg-surface-base/90 text-ink-muted shadow-sm hover:bg-brand-500 hover:text-white disabled:cursor-not-allowed disabled:opacity-30 transition-colors"
                    aria-label="Move up"
                  >
                    ▲
                  </button>
                  <button
                    type="button"
                    disabled={idx === issues.length - 1}
                    onClick={() => onReorder(iss.id, 'down')}
                    className="flex h-5 w-5 items-center justify-center rounded bg-surface-base/90 text-ink-muted shadow-sm hover:bg-brand-500 hover:text-white disabled:cursor-not-allowed disabled:opacity-30 transition-colors"
                    aria-label="Move down"
                  >
                    ▼
                  </button>
                </div>
              )}
            </div>
          ))
        )}
      </div>
      {adding ? (
        <div className="flex gap-1">
          <input
            autoFocus
            value={newTitle}
            onChange={(e) => setNewTitle(e.target.value)}
            onKeyDown={(e) => {
              if (e.key === 'Enter' && newTitle.trim()) {
                onCreateIssue?.(newTitle.trim(), status);
                setNewTitle('');
                setAdding(false);
              }
              if (e.key === 'Escape') { setNewTitle(''); setAdding(false); }
            }}
            placeholder={t('board.newIssueTitle')}
            className="flex-1 rounded-md border border-edge-base bg-surface-base px-2 py-1 text-sm focus:outline-none focus:border-brand-500"
          />
        </div>
      ) : (
        <button
          type="button"
          onClick={() => setAdding(true)}
          className="flex items-center gap-1 rounded-md px-2 py-1.5 text-xs text-ink-muted hover:bg-surface-muted hover:text-ink-primary transition-colors w-full"
        >
          <Plus size={14} /> {t('board.addIssue')}
        </button>
      )}
    </section>
  );
}
