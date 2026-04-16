import { useTranslation } from 'react-i18next';
import type { IssueStatus } from '../api/types';

const classes: Record<IssueStatus, string> = {
  Pending: 'bg-status-pending/15 text-status-pending border-status-pending/30',
  Approved: 'bg-status-approved/15 text-status-approved border-status-approved/30',
  InProgress: 'bg-status-inProgress/15 text-status-inProgress border-status-inProgress/40',
  Done: 'bg-status-done/15 text-status-done border-status-done/40',
  Rejected: 'bg-status-rejected/15 text-status-rejected border-status-rejected/30',
};

export function StatusBadge({ status }: { status: IssueStatus }) {
  const { t } = useTranslation();
  const label = t(`status.${status}`);
  return (
    <span
      className={`inline-flex items-center gap-1 rounded-full border px-2 py-0.5 text-xs font-medium ${classes[status]}`}
      role="status"
      aria-label={label}
    >
      <span aria-hidden className="h-1.5 w-1.5 rounded-full bg-current" />
      {label}
    </span>
  );
}
