import type { Issue } from '../../api/types';

type Props = {
  issues: Issue[];
};

/**
 * AgentPresenceBar — shows which issues AI agents are actively working on.
 * Only rendered when there are InProgress or QA issues so it does not take
 * up space on quiet boards.
 *
 * Data comes from the already-loaded board issues; updates appear in real
 * time because BoardPage invalidates the issues query on every SSE event.
 */
export function AgentPresenceBar({ issues }: Props) {
  const active = issues.filter((i) => i.status === 'InProgress');
  const inQA = issues.filter((i) => i.status === 'QA');

  if (active.length === 0 && inQA.length === 0) return null;

  return (
    <div
      className="flex items-center gap-2 rounded-lg border border-brand-500/20 bg-brand-500/5 px-3 py-1.5 text-[11px] text-ink-secondary"
      role="status"
      aria-live="polite"
      aria-label={`${active.length} agent${active.length !== 1 ? 's' : ''} active`}
    >
      {/* live pulse dot */}
      <span className="relative flex h-2 w-2 shrink-0">
        <span className="absolute inline-flex h-full w-full animate-ping rounded-full bg-brand-500 opacity-60" />
        <span className="relative inline-flex h-2 w-2 rounded-full bg-brand-500" />
      </span>

      {active.length > 0 && (
        <span className="font-medium text-brand-500">
          {active.length} working
        </span>
      )}

      {active.length > 0 && (
        <span className="flex min-w-0 gap-2 overflow-hidden">
          {active.slice(0, 3).map((i) => (
            <span
              key={i.id}
              className="max-w-[12rem] truncate rounded bg-brand-500/10 px-1.5 py-0.5 text-brand-500"
              title={i.title}
            >
              {i.title}
            </span>
          ))}
          {active.length > 3 && (
            <span className="text-ink-muted">+{active.length - 3} more</span>
          )}
        </span>
      )}

      {inQA.length > 0 && (
        <>
          {active.length > 0 && <span className="text-ink-muted">·</span>}
          <span className="font-medium text-amber-500">
            {inQA.length} in QA
          </span>
          <span className="flex min-w-0 gap-2 overflow-hidden">
            {inQA.slice(0, 2).map((i) => (
              <span
                key={i.id}
                className="max-w-[12rem] truncate rounded bg-amber-500/10 px-1.5 py-0.5 text-amber-600 dark:text-amber-400"
                title={i.title}
              >
                {i.title}
              </span>
            ))}
            {inQA.length > 2 && (
              <span className="text-ink-muted">+{inQA.length - 2} more</span>
            )}
          </span>
        </>
      )}
    </div>
  );
}
