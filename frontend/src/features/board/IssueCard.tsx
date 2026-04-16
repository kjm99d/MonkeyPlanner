import { useDraggable } from '@dnd-kit/core';
import { useState } from 'react';
import { Link, useNavigate } from 'react-router-dom';
import { useTranslation } from 'react-i18next';
import { GripVertical, Check, Calendar, Tag, Trash2, ExternalLink, Copy } from 'lucide-react';
import type { Issue, BoardProperty } from '../../api/types';
import { StatusBadge } from '../../components/StatusBadge';
import { useApproveIssue, useUpdateIssue, useDeleteIssue } from '../../api/hooks';
import { ContextMenu } from '../../components/ContextMenu';

type Props = {
  issue: Issue;
  boardProperties?: BoardProperty[];
};

export function IssueCard({ issue, boardProperties = [] }: Props) {
  const { t } = useTranslation();
  const { attributes, listeners, setNodeRef, isDragging } = useDraggable({
    id: issue.id,
  });
  const approve = useApproveIssue();
  const updateIssue = useUpdateIssue();
  const deleteIssue = useDeleteIssue();
  const navigate = useNavigate();
  const [editing, setEditing] = useState(false);
  const [editTitle, setEditTitle] = useState('');
  const [ctxMenu, setCtxMenu] = useState<{ x: number; y: number } | null>(null);

  // DragOverlay handles the visual ghost — hide the original card during drag
  const style: React.CSSProperties | undefined = isDragging
    ? { opacity: 0.3, pointerEvents: 'none' }
    : undefined;

  // 속성 값 추출 (비어있지 않은 것만)
  const propValues = boardProperties
    .map((p) => ({ prop: p, value: issue.properties?.[p.id] }))
    .filter((pv) => pv.value !== undefined && pv.value !== null && pv.value !== '');

  return (
    <article
      ref={setNodeRef}
      style={style}
      {...attributes}
      {...listeners}
      aria-roledescription="draggable item"
      className={`group flex cursor-grab flex-col gap-2 rounded-lg border border-edge-base bg-surface-base p-3 shadow-sm transition-shadow hover:shadow-md motion-reduce:transition-none ${
        isDragging ? 'cursor-grabbing' : ''
      }`}
      onContextMenu={(e) => {
        e.preventDefault();
        e.stopPropagation();
        setCtxMenu({ x: e.clientX, y: e.clientY });
      }}
    >
      <div className="flex items-start gap-2">
        <GripVertical size={14} className="mt-0.5 shrink-0 cursor-grab text-ink-muted opacity-30 group-hover:opacity-70 transition-opacity" aria-hidden />
        <div className="flex-1 min-w-0">
          {editing ? (
            <input
              autoFocus
              value={editTitle}
              onChange={(e) => setEditTitle(e.target.value)}
              onBlur={async () => {
                setEditing(false);
                if (editTitle.trim() && editTitle !== issue.title) {
                  await updateIssue.mutateAsync({ id: issue.id, patch: { title: editTitle.trim() } });
                }
              }}
              onKeyDown={(e) => {
                if (e.key === 'Enter') e.currentTarget.blur();
                if (e.key === 'Escape') { setEditTitle(issue.title); setEditing(false); }
              }}
              onPointerDown={(e) => e.stopPropagation()}
              className="w-full rounded border border-brand-500 bg-surface-base px-1 py-0.5 text-sm font-medium text-ink-primary outline-none"
            />
          ) : (
            <Link
              to={`/issues/${issue.id}`}
              onClick={(e) => e.stopPropagation()}
              onDoubleClick={(e) => {
                e.preventDefault();
                e.stopPropagation();
                setEditTitle(issue.title);
                setEditing(true);
              }}
              className="font-medium text-ink-primary hover:underline focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-brand-500 rounded"
            >
              {issue.title}
            </Link>
          )}
          {/* 속성 값 인라인 표시 */}
          {propValues.length > 0 && (
            <div className="mt-1.5 flex flex-wrap gap-1" onPointerDown={(e) => e.stopPropagation()}>
              {propValues.map(({ prop, value }) => (
                <PropertyPill key={prop.id} prop={prop} value={value} />
              ))}
            </div>
          )}
        </div>
        <StatusBadge status={issue.status} />
      </div>
      {issue.status === 'Pending' && (
        <div className="flex justify-end">
          <button
            type="button"
            onPointerDown={(e) => e.stopPropagation()}
            onClick={(e) => {
              e.stopPropagation();
              approve.mutate(issue.id);
            }}
            disabled={approve.isPending}
            aria-label={t('kanban.approveLabel', { title: issue.title })}
            className="flex items-center gap-1 rounded-md bg-accent px-2.5 py-1 text-xs font-semibold text-white shadow-sm transition-all duration-150 hover:brightness-110 hover:shadow-md active:scale-95 disabled:opacity-50 cursor-pointer"
          >
            <Check size={12} /> {t('kanban.approve')}
          </button>
        </div>
      )}
      <ContextMenu
        position={ctxMenu}
        onClose={() => setCtxMenu(null)}
        items={[
          { label: t('issue.title') + '...', icon: <ExternalLink size={14} />, onClick: () => navigate(`/issues/${issue.id}`) },
          { label: t('board.addIssue'), icon: <Copy size={14} />, onClick: () => { navigator.clipboard.writeText(issue.title); } },
          { divider: true, label: '', onClick: () => {} },
          { label: t('issue.delete'), icon: <Trash2 size={14} />, onClick: () => deleteIssue.mutate(issue.id), danger: true },
        ]}
      />
    </article>
  );
}

function PropertyPill({ prop, value }: { prop: BoardProperty; value: unknown }) {
  if (prop.type === 'checkbox') {
    if (!value) return null;
    return (
      <span className="inline-flex items-center gap-0.5 rounded bg-status-done/10 px-1.5 py-0.5 text-[11px] text-status-done">
        <Check size={10} /> {prop.name}
      </span>
    );
  }

  if (prop.type === 'date') {
    return (
      <span className="inline-flex items-center gap-0.5 rounded bg-status-inProgress/10 px-1.5 py-0.5 text-[11px] text-status-inProgress">
        <Calendar size={10} /> {String(value)}
      </span>
    );
  }

  if (prop.type === 'multi_select' && Array.isArray(value)) {
    return (
      <>
        {value.map((v: string) => (
          <span key={v} className="inline-flex items-center gap-0.5 rounded bg-brand-500/10 px-1.5 py-0.5 text-[11px] text-brand-500">
            <Tag size={9} /> {v}
          </span>
        ))}
      </>
    );
  }

  // select, text, number
  return (
    <span className="rounded bg-surface-muted px-1.5 py-0.5 text-[11px] text-ink-secondary">
      {prop.name}: {String(value)}
    </span>
  );
}
