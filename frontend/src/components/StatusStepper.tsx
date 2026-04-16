import { useTranslation } from 'react-i18next';
import type { IssueStatus } from '../api/types';

const STEP_STATUSES: IssueStatus[] = ['Pending', 'Approved', 'InProgress', 'Done', 'Rejected'];

const dotColor: Record<IssueStatus, string> = {
  Pending: 'bg-status-pending',
  Approved: 'bg-status-approved',
  InProgress: 'bg-status-inProgress',
  Done: 'bg-status-done',
  Rejected: 'bg-status-rejected',
};

type Props = {
  current: IssueStatus;
  onSelect: (status: IssueStatus) => void;
  disabled?: boolean;
};

export function StatusStepper({ current, onSelect, disabled }: Props) {
  const { t } = useTranslation();
  const currentIdx = STEP_STATUSES.indexOf(current);

  return (
    <nav aria-label={t('stepper.moveTo', { label: '' }).trim()} className="flex items-center gap-1">
      {STEP_STATUSES.map((status, i) => {
        const label = t(`status.${status}`);
        const isActive = status === current;
        const isPast = i < currentIdx;
        const canClick =
          !disabled &&
          !isActive &&
          status !== 'Pending' &&
          status !== 'Approved' &&
          current !== 'Pending';

        const hint =
          status === 'Approved'
            ? t('stepper.approveOnly')
            : status === 'Pending'
              ? t('stepper.noPending')
              : t('stepper.moveTo', { label });
        const showHint = !canClick && !isActive;

        return (
          <div key={status} className="flex items-center">
            {i > 0 && (
              <div
                className={`mx-1 h-0.5 w-6 rounded ${
                  isPast || isActive ? 'bg-brand-500' : 'bg-edge-base'
                }`}
              />
            )}
            <div className="group relative">
              <button
                type="button"
                disabled={!canClick}
                onClick={() => canClick && onSelect(status)}
                aria-current={isActive ? 'step' : undefined}
                title={hint}
                className={`flex items-center gap-1.5 rounded-full border px-3 py-1 text-xs font-medium transition-colors ${
                  isActive
                    ? `border-current ${dotColor[status]} text-white`
                    : canClick
                      ? 'border-edge-base bg-surface-subtle text-ink-secondary hover:bg-surface-muted hover:text-ink-primary'
                      : 'border-transparent text-ink-muted opacity-50 cursor-not-allowed'
                }`}
              >
                <span
                  className={`inline-block h-2 w-2 rounded-full ${
                    isActive ? 'bg-white' : dotColor[status]
                  }`}
                />
                {label}
              </button>
              {showHint && (
                <span className="pointer-events-none absolute -bottom-7 left-1/2 -translate-x-1/2 whitespace-nowrap rounded bg-ink-primary px-2 py-0.5 text-[11px] text-surface-base opacity-0 transition-opacity group-hover:opacity-100">
                  {hint}
                </span>
              )}
            </div>
          </div>
        );
      })}
    </nav>
  );
}
