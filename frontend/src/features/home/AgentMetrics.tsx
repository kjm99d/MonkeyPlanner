import { useMemo } from 'react';
import { useTranslation } from 'react-i18next';
import { Timer, TrendingUp, CheckCircle2, BarChart3 } from 'lucide-react';
import { useIssues } from '../../api/hooks';

export function AgentMetrics() {
  const { t } = useTranslation();
  const issues = useIssues({});

  const metrics = useMemo(() => {
    const all = issues.data ?? [];
    const total = all.length;
    const done = all.filter(i => i.status === 'Done').length;
    const rejected = all.filter(i => i.status === 'Rejected').length;
    const approved = all.filter(i => i.status === 'Approved').length;
    const inProgress = all.filter(i => i.status === 'InProgress').length;
    const pending = all.filter(i => i.status === 'Pending').length;

    // Completion rate
    const completionRate = total > 0 ? Math.round((done / total) * 100) : 0;

    // Approval rate (approved + inProgress + done out of total non-pending)
    const reviewed = total - pending;
    const approvalRate = reviewed > 0 ? Math.round(((approved + inProgress + done) / reviewed) * 100) : 0;

    // Average cycle time (created → done, for completed issues)
    const completedIssues = all.filter(i => i.status === 'Done');
    let avgCycleHours = 0;
    if (completedIssues.length > 0) {
      const totalMs = completedIssues.reduce((sum, i) => {
        return sum + (new Date(i.updatedAt).getTime() - new Date(i.createdAt).getTime());
      }, 0);
      avgCycleHours = Math.round(totalMs / completedIssues.length / 3600000 * 10) / 10;
    }

    return { total, done, rejected, approved, inProgress, pending, completionRate, approvalRate, avgCycleHours };
  }, [issues.data]);

  if (!issues.data || issues.data.length === 0) return null;

  return (
    <section className="fade-up flex flex-col gap-3" style={{ animationDelay: '120ms' }}>
      <h2 className="flex items-center gap-2 text-base font-semibold">
        <BarChart3 size={16} className="text-brand-500" />
        {t('home.agentMetrics', 'Agent Metrics')}
      </h2>
      <div className="grid gap-3 sm:grid-cols-2 lg:grid-cols-4">
        {/* Completion Rate */}
        <MetricCard
          icon={<CheckCircle2 size={16} className="text-status-done" />}
          label={t('home.completionRate', 'Completion Rate')}
          value={`${metrics.completionRate}%`}
          sub={`${metrics.done}/${metrics.total}`}
        />
        {/* Approval Rate */}
        <MetricCard
          icon={<TrendingUp size={16} className="text-status-approved" />}
          label={t('home.approvalRate', 'Approval Rate')}
          value={`${metrics.approvalRate}%`}
          sub={`${metrics.rejected} rejected`}
        />
        {/* Avg Cycle Time */}
        <MetricCard
          icon={<Timer size={16} className="text-status-inProgress" />}
          label={t('home.avgCycleTime', 'Avg Cycle Time')}
          value={metrics.avgCycleHours < 24 ? `${metrics.avgCycleHours}h` : `${Math.round(metrics.avgCycleHours / 24 * 10) / 10}d`}
          sub={`${metrics.done} completed`}
        />
        {/* Pipeline */}
        <div className="flex flex-col gap-2 rounded-lg border border-edge-base bg-gradient-to-br from-surface-subtle to-surface-muted p-3">
          <span className="text-xs font-medium text-ink-muted">{t('home.pipeline', 'Pipeline')}</span>
          <div className="flex items-center gap-1 text-xs">
            <span className="rounded bg-status-pending/15 px-1.5 py-0.5 text-status-pending">{metrics.pending}</span>
            <span className="text-ink-muted">→</span>
            <span className="rounded bg-status-approved/15 px-1.5 py-0.5 text-status-approved">{metrics.approved}</span>
            <span className="text-ink-muted">→</span>
            <span className="rounded bg-status-inProgress/15 px-1.5 py-0.5 text-status-inProgress">{metrics.inProgress}</span>
            <span className="text-ink-muted">→</span>
            <span className="rounded bg-status-done/15 px-1.5 py-0.5 text-status-done">{metrics.done}</span>
          </div>
          {/* Progress bar */}
          <div className="flex h-2 w-full overflow-hidden rounded-full bg-surface-muted">
            {metrics.total > 0 && (
              <>
                <div className="bg-status-done h-full" style={{ width: `${(metrics.done / metrics.total) * 100}%` }} />
                <div className="bg-status-inProgress h-full" style={{ width: `${(metrics.inProgress / metrics.total) * 100}%` }} />
                <div className="bg-status-approved h-full" style={{ width: `${(metrics.approved / metrics.total) * 100}%` }} />
                <div className="bg-status-pending h-full" style={{ width: `${(metrics.pending / metrics.total) * 100}%` }} />
              </>
            )}
          </div>
        </div>
      </div>
    </section>
  );
}

function MetricCard({ icon, label, value, sub }: { icon: React.ReactNode; label: string; value: string; sub: string }) {
  return (
    <div className="flex flex-col gap-1 rounded-lg border border-edge-base bg-gradient-to-br from-surface-subtle to-surface-muted p-3">
      <div className="flex items-center gap-1.5">
        {icon}
        <span className="text-xs font-medium text-ink-muted">{label}</span>
      </div>
      <span className="text-2xl font-bold tabular-nums text-ink-primary">{value}</span>
      <span className="text-xs text-ink-muted tabular-nums">{sub}</span>
    </div>
  );
}
