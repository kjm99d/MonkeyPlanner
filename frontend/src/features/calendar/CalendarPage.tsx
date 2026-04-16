import { useMemo, useState } from 'react';
import { Link } from 'react-router-dom';
import { useTranslation } from 'react-i18next';
import { useDayStats, useMonthStats } from '../../api/hooks';
import { Button } from '../../components/Button';
import { StatusBadge } from '../../components/StatusBadge';
import type { DayCount, Issue } from '../../api/types';

function pad(n: number): string {
  return n.toString().padStart(2, '0');
}

function ymd(d: Date): string {
  return `${d.getFullYear()}-${pad(d.getMonth() + 1)}-${pad(d.getDate())}`;
}

function startOfMonth(y: number, m: number): Date {
  return new Date(y, m - 1, 1);
}

function daysInMonth(y: number, m: number): number {
  return new Date(y, m, 0).getDate();
}

export default function CalendarPage() {
  const { t } = useTranslation();
  const today = useMemo(() => new Date(), []);
  const [year, setYear] = useState(today.getFullYear());
  const [month, setMonth] = useState(today.getMonth() + 1);
  const [selected, setSelected] = useState<string>(ymd(today));

  const monthStats = useMonthStats(year, month);
  const dayStats = useDayStats(selected);

  const countsByDate = useMemo(() => {
    const map = new Map<string, DayCount>();
    (monthStats.data ?? []).forEach((d) => {
      const key = d.date.slice(0, 10);
      map.set(key, d);
    });
    return map;
  }, [monthStats.data]);

  const firstOfMonth = startOfMonth(year, month);
  const firstWeekday = firstOfMonth.getDay(); // 0=일
  const total = daysInMonth(year, month);

  const cells: Array<{ day: number | null; dateStr?: string }> = [];
  for (let i = 0; i < firstWeekday; i++) cells.push({ day: null });
  for (let d = 1; d <= total; d++) {
    cells.push({ day: d, dateStr: `${year}-${pad(month)}-${pad(d)}` });
  }
  while (cells.length % 7 !== 0) cells.push({ day: null });

  function shift(delta: number) {
    const d = new Date(year, month - 1 + delta, 1);
    setYear(d.getFullYear());
    setMonth(d.getMonth() + 1);
  }

  return (
    <section className="grid gap-6 lg:grid-cols-[2fr_1fr]">
      <div className="flex flex-col gap-4">
        <header className="flex items-center justify-between">
          <h1 className="text-3xl font-bold">{t('calendar.title', { year, month })}</h1>
          <div className="flex gap-2">
            <Button size="sm" variant="ghost" onClick={() => shift(-1)} aria-label={t('calendar.prevMonth')}>
              ←
            </Button>
            <Button size="sm" variant="ghost" onClick={() => { setYear(today.getFullYear()); setMonth(today.getMonth()+1); setSelected(ymd(today)); }}>
              {t('calendar.today')}
            </Button>
            <Button size="sm" variant="ghost" onClick={() => shift(1)} aria-label={t('calendar.nextMonth')}>
              →
            </Button>
          </div>
        </header>

        <div role="grid" aria-label={t('nav.calendar')} className="grid grid-cols-7 gap-1">
          {['일', '월', '화', '수', '목', '금', '토'].map((d) => (
            <div key={d} className="pb-2 text-center text-xs font-medium text-ink-muted">
              {d}
            </div>
          ))}
          {cells.map((c, idx) => {
            if (c.day === null) return <div key={`empty-${idx}`} />;
            const isToday = c.dateStr === ymd(today);
            const isSelected = c.dateStr === selected;
            const stat = c.dateStr ? countsByDate.get(c.dateStr) : undefined;
            return (
              <button
                key={c.dateStr}
                role="gridcell"
                aria-selected={isSelected}
                aria-label={t('calendar.selectDate', { date: c.dateStr })}
                onClick={() => c.dateStr && setSelected(c.dateStr)}
                className={`flex min-h-[5.5rem] flex-col items-start gap-1 rounded-md border p-2 text-left transition-transform hover:-translate-y-0.5 motion-reduce:hover:transform-none focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-brand-500 ${
                  isSelected
                    ? 'border-brand-500 bg-brand-500/10'
                    : 'border-edge-base bg-surface-subtle hover:bg-surface-muted'
                }`}
              >
                <span
                  className={`text-sm font-medium ${
                    isToday ? 'rounded-full bg-brand-500 px-2 text-white' : 'text-ink-primary'
                  }`}
                >
                  {c.day}
                </span>
                {stat && (stat.created + stat.approved + stat.completed > 0) && (
                  <div className="flex flex-wrap gap-1 text-[10px] tabular-nums">
                    {stat.created > 0 && (
                      <span className="rounded bg-ink-primary/10 px-1 text-ink-primary">+{stat.created}</span>
                    )}
                    {stat.approved > 0 && (
                      <span className="rounded bg-status-approved/15 px-1 text-status-approved">✓{stat.approved}</span>
                    )}
                    {stat.completed > 0 && (
                      <span className="rounded bg-status-done/15 px-1 text-status-done">●{stat.completed}</span>
                    )}
                  </div>
                )}
              </button>
            );
          })}
        </div>
      </div>

      <aside aria-label={t('calendar.daily', { date: selected })} className="flex flex-col gap-4 rounded-lg border border-edge-base bg-surface-subtle p-4">
        <h2 className="text-xl font-semibold">{t('calendar.daily', { date: selected })}</h2>
        {dayStats.isLoading && <p className="text-ink-secondary">{t('common.loading')}</p>}
        {dayStats.data && (
          <div className="flex flex-col gap-4">
            <DaySection title={t('calendar.created')} issues={dayStats.data.created} />
            <DaySection title={t('calendar.approved')} issues={dayStats.data.approved} />
            <DaySection title={t('calendar.done')} issues={dayStats.data.completed} />
          </div>
        )}
      </aside>
    </section>
  );
}

function DaySection({ title, issues }: { title: string; issues: Issue[] }) {
  const { t } = useTranslation();
  return (
    <section aria-label={`${title} 이슈`}>
      <h3 className="mb-2 text-sm font-medium text-ink-secondary">
        {title} <span className="text-ink-muted">({issues.length})</span>
      </h3>
      {issues.length === 0 ? (
        <p className="text-xs text-ink-muted">{t('calendar.none')}</p>
      ) : (
        <ul className="flex flex-col gap-1">
          {issues.map((iss) => (
            <li key={iss.id}>
              <Link
                to={`/issues/${iss.id}`}
                className="flex items-center justify-between gap-2 rounded-md bg-surface-base px-2 py-1.5 text-sm hover:bg-surface-muted focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-brand-500"
              >
                <span className="truncate">{iss.title}</span>
                <StatusBadge status={iss.status} />
              </Link>
            </li>
          ))}
        </ul>
      )}
    </section>
  );
}
