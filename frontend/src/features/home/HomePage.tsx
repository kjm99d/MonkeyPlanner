import { Link } from 'react-router-dom';
import { useMemo } from 'react';
import { useTranslation } from 'react-i18next';
import { Clock, CheckCircle2, ArrowRight } from 'lucide-react';
import { useBoards, useDayStats, useIssues } from '../../api/hooks';
import { StatusBadge } from '../../components/StatusBadge';
import { Skeleton } from '../../components/Skeleton';
import { WeeklyChart } from './WeeklyChart';
import { AgentMetrics } from './AgentMetrics';

function todayString(): string {
  const d = new Date();
  return `${d.getFullYear()}-${String(d.getMonth() + 1).padStart(2, '0')}-${String(d.getDate()).padStart(2, '0')}`;
}

function relativeTime(dateStr: string): string {
  const now = Date.now();
  const then = new Date(dateStr).getTime();
  const diff = Math.floor((now - then) / 1000);
  if (diff < 60) return 'just now';
  if (diff < 3600) return `${Math.floor(diff / 60)}m ago`;
  if (diff < 86400) return `${Math.floor(diff / 3600)}h ago`;
  if (diff < 604800) return `${Math.floor(diff / 86400)}d ago`;
  return new Date(dateStr).toLocaleDateString();
}

export default function HomePage() {
  const date = useMemo(todayString, []);
  const { t } = useTranslation();
  const day = useDayStats(date);
  const approved = useIssues({ status: 'Approved' });
  const inProgress = useIssues({ status: 'InProgress' });
  const boards = useBoards();

  const created = day.data?.created.length ?? 0;
  const approvedCount = day.data?.approved.length ?? 0;
  const completed = day.data?.completed.length ?? 0;

  return (
    <section className="flex flex-col gap-6">
      <header className="fade-up">
        <h1 className="text-2xl font-bold tracking-tight">{t('home.title')}</h1>
        <p className="mt-1 text-sm text-ink-secondary">
          {t('home.summary', { date })}
        </p>
      </header>

      {/* 통계 카드 3분할 */}
      <div className="grid gap-3 sm:grid-cols-3">
        {day.isLoading ? (
          <>
            <Skeleton className="h-24 rounded-lg" />
            <Skeleton className="h-24 rounded-lg" />
            <Skeleton className="h-24 rounded-lg" />
          </>
        ) : (
          <>
            <StatCard label={t('home.created')} value={created} hue="text-brand-500" accent="bg-brand-500" delay={0} />
            <StatCard label={t('home.approved')} value={approvedCount} hue="text-status-approved" accent="bg-status-approved" delay={1} />
            <StatCard label={t('home.done')} value={completed} hue="text-status-done" accent="bg-status-done" delay={2} />
          </>
        )}
      </div>

      {/* 주간 차트 */}
      <WeeklyChart />

      {/* 에이전트 메트릭 */}
      <AgentMetrics />

      {/* 2열 레이아웃: 진행 중 이슈 + 보드 요약 */}
      <div className="grid gap-6 lg:grid-cols-2">
        {/* 좌: 진행 중 + 승인 대기 이슈 */}
        <section className="flex flex-col gap-4 fade-up" style={{ animationDelay: '100ms' }}>
          <div className="flex items-center justify-between">
            <h2 className="flex items-center gap-2 text-lg font-semibold">
              <Clock size={18} className="text-status-inProgress" /> {t('home.inProgress')}
            </h2>
            <span className="rounded-full bg-status-inProgress/10 px-2 py-0.5 text-xs font-medium text-status-inProgress">
              {t('home.items', { count: inProgress.data?.length ?? 0 })}
            </span>
          </div>
          {inProgress.isLoading ? (
            <Skeleton className="h-14 rounded-lg" count={3} />
          ) : inProgress.data && inProgress.data.length > 0 ? (
            <ul className="flex flex-col gap-2">
              {inProgress.data.slice(0, 5).map((iss, i) => (
                <li key={iss.id} className="fade-up" style={{ animationDelay: `${150 + i * 40}ms` }}>
                  <Link
                    to={`/issues/${iss.id}`}
                    className="flex items-center justify-between rounded-lg border border-edge-base bg-surface-subtle px-4 py-3 transition-all duration-200 hover:bg-surface-muted hover:shadow-sm"
                  >
                    <span className="font-medium">{iss.title}</span>
                    <StatusBadge status={iss.status} />
                  </Link>
                </li>
              ))}
            </ul>
          ) : (
            <p className="py-6 text-center text-sm text-ink-muted">{t('home.noInProgress')}</p>
          )}

          {/* 승인 대기 */}
          {approved.data && approved.data.length > 0 && (
            <div className="flex flex-col gap-2 fade-up" style={{ animationDelay: '250ms' }}>
              <h3 className="flex items-center gap-2 text-sm font-medium text-ink-secondary">
                <CheckCircle2 size={14} className="text-status-approved" /> {t('home.approvedWaiting')}
              </h3>
              <ul className="flex flex-col gap-1.5">
                {approved.data.slice(0, 3).map((iss) => (
                  <li key={iss.id}>
                    <Link
                      to={`/issues/${iss.id}`}
                      className="flex items-center justify-between rounded-md border border-edge-base bg-surface-base px-3 py-2 text-sm transition-colors hover:bg-surface-muted"
                    >
                      <span>{iss.title}</span>
                      <StatusBadge status={iss.status} />
                    </Link>
                  </li>
                ))}
              </ul>
            </div>
          )}
        </section>

        {/* 우: 보드 요약 */}
        <section className="flex flex-col gap-4 fade-up" style={{ animationDelay: '200ms' }}>
          <div className="flex items-center justify-between">
            <h2 className="text-lg font-semibold">{t('home.boardStatus')}</h2>
            <Link to="/boards" className="flex items-center gap-1 text-xs font-medium text-brand-500 hover:underline">
              {t('home.viewAll')} <ArrowRight size={12} />
            </Link>
          </div>
          {boards.isLoading ? (
            <Skeleton className="h-20 rounded-lg" count={2} />
          ) : boards.data && boards.data.length > 0 ? (
            <ul className="flex flex-col gap-2">
              {boards.data.slice(0, 4).map((b, i) => (
                <li key={b.id} className="fade-up" style={{ animationDelay: `${250 + i * 40}ms` }}>
                  <Link
                    to={`/boards/${b.id}`}
                    className="flex items-center justify-between rounded-lg border border-edge-base bg-surface-subtle px-4 py-3 transition-all duration-200 hover:bg-surface-muted hover:shadow-sm"
                  >
                    <div>
                      <span className="font-medium">{b.name}</span>
                      <span className="ml-2 text-xs text-ink-muted">{b.viewType}</span>
                    </div>
                    <ArrowRight size={14} className="text-ink-muted" />
                  </Link>
                </li>
              ))}
            </ul>
          ) : (
            <Link
              to="/boards"
              className="flex flex-col items-center gap-2 rounded-lg border-2 border-dashed border-edge-base py-8 text-center transition-colors hover:border-brand-500/30 hover:bg-brand-500/5"
            >
              <span className="text-2xl">+</span>
              <span className="text-sm text-ink-secondary">{t('home.createFirst')}</span>
            </Link>
          )}
        </section>
      </div>
    </section>
  );
}

function StatCard({ label, value, hue, accent, delay }: { label: string; value: number; hue: string; accent: string; delay: number }) {
  return (
    <div
      className="fade-up relative overflow-hidden rounded-lg border border-edge-base bg-gradient-to-br from-surface-subtle to-surface-muted p-4 shadow-sm transition-shadow hover:shadow-md"
      style={{ animationDelay: `${delay * 60}ms` }}
    >
      <div className={`absolute left-0 top-0 h-full w-1 rounded-l-lg ${accent}`} />
      <div className="flex items-end justify-between pl-2">
        <div className="flex flex-col gap-0.5">
          <span className="text-xs font-medium text-ink-muted">{label}</span>
          <span className={`text-3xl font-bold tabular-nums tracking-tight ${hue}`}>{value}</span>
        </div>
        <span className="text-xs text-ink-muted tabular-nums">{relativeTime(new Date().toISOString())}</span>
      </div>
    </div>
  );
}
