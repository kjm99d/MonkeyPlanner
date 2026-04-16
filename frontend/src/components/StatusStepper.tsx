import type { IssueStatus } from '../api/types';

const STEPS: { status: IssueStatus; label: string }[] = [
  { status: 'Pending', label: '대기' },
  { status: 'Approved', label: '승인됨' },
  { status: 'InProgress', label: '진행 중' },
  { status: 'Done', label: '완료' },
];

const dotColor: Record<IssueStatus, string> = {
  Pending: 'bg-status-pending',
  Approved: 'bg-status-approved',
  InProgress: 'bg-status-inProgress',
  Done: 'bg-status-done',
};

type Props = {
  current: IssueStatus;
  onSelect: (status: IssueStatus) => void;
  disabled?: boolean;
};

export function StatusStepper({ current, onSelect, disabled }: Props) {
  const currentIdx = STEPS.findIndex((s) => s.status === current);

  return (
    <nav aria-label="상태 전이" className="flex items-center gap-1">
      {STEPS.map((step, i) => {
        const isActive = step.status === current;
        const isPast = i < currentIdx;
        const canClick =
          !disabled &&
          !isActive &&
          step.status !== 'Pending' &&
          step.status !== 'Approved' &&
          current !== 'Pending';

        return (
          <div key={step.status} className="flex items-center">
            {i > 0 && (
              <div
                className={`mx-1 h-0.5 w-6 rounded ${
                  isPast || isActive ? 'bg-brand-500' : 'bg-edge-base'
                }`}
              />
            )}
            <button
              type="button"
              disabled={!canClick}
              onClick={() => canClick && onSelect(step.status)}
              aria-current={isActive ? 'step' : undefined}
              title={
                step.status === 'Approved'
                  ? 'Approve 버튼으로만 전환 가능'
                  : step.status === 'Pending'
                    ? '대기 상태로 되돌릴 수 없음'
                    : `${step.label}(으)로 전환`
              }
              className={`flex items-center gap-1.5 rounded-full border px-3 py-1 text-xs font-medium transition-colors ${
                isActive
                  ? `border-current ${dotColor[step.status]} text-white`
                  : canClick
                    ? 'border-edge-base bg-surface-subtle text-ink-secondary hover:bg-surface-muted hover:text-ink-primary'
                    : 'border-transparent text-ink-muted opacity-50 cursor-not-allowed'
              }`}
            >
              <span
                className={`inline-block h-2 w-2 rounded-full ${
                  isActive ? 'bg-white' : dotColor[step.status]
                }`}
              />
              {step.label}
            </button>
          </div>
        );
      })}
    </nav>
  );
}
