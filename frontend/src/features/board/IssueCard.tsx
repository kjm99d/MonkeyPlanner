import { useDraggable } from '@dnd-kit/core';
import { Link } from 'react-router-dom';
import type { Issue } from '../../api/types';
import { StatusBadge } from '../../components/StatusBadge';
import { Button } from '../../components/Button';
import { useApproveIssue } from '../../api/hooks';

type Props = {
  issue: Issue;
};

export function IssueCard({ issue }: Props) {
  const { attributes, listeners, setNodeRef, transform, isDragging } = useDraggable({
    id: issue.id,
  });
  const approve = useApproveIssue();

  const style: React.CSSProperties = transform
    ? {
        transform: `translate3d(${transform.x}px, ${transform.y}px, 0) perspective(800px) rotateX(3deg) rotateY(-3deg)`,
      }
    : undefined as unknown as React.CSSProperties;

  return (
    <article
      ref={setNodeRef}
      style={style}
      {...attributes}
      {...listeners}
      aria-grabbed={isDragging}
      className={`group flex cursor-grab flex-col gap-2 rounded-lg border border-edge-base bg-surface-base p-3 shadow-sm transition-shadow hover:shadow-md motion-reduce:transition-none ${
        isDragging ? 'shadow-lg opacity-90 cursor-grabbing' : ''
      }`}
    >
      <div className="flex items-start justify-between gap-2">
        <Link
          to={`/issues/${issue.id}`}
          onClick={(e) => e.stopPropagation()}
          className="font-medium text-ink-primary hover:underline focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-brand-500 rounded"
        >
          {issue.title}
        </Link>
        <StatusBadge status={issue.status} />
      </div>
      {issue.status === 'Pending' && (
        <Button
          size="sm"
          onPointerDown={(e) => e.stopPropagation()}
          onClick={(e) => {
            e.stopPropagation();
            approve.mutate(issue.id);
          }}
          disabled={approve.isPending}
          aria-label={`이슈 ${issue.title} 승인`}
        >
          ✓ Approve
        </Button>
      )}
    </article>
  );
}
