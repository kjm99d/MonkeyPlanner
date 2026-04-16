import { useMemo } from 'react';
import { useTranslation } from 'react-i18next';
import {
  BarChart,
  Bar,
  XAxis,
  YAxis,
  Tooltip,
  ResponsiveContainer,
  CartesianGrid,
  Legend,
} from 'recharts';
import { useMonthStats } from '../../api/hooks';
import { BarChart3 } from 'lucide-react';

function getLast7Days(): string[] {
  const days: string[] = [];
  const now = new Date();
  for (let i = 6; i >= 0; i--) {
    const d = new Date(now);
    d.setDate(d.getDate() - i);
    days.push(
      `${d.getFullYear()}-${String(d.getMonth() + 1).padStart(2, '0')}-${String(d.getDate()).padStart(2, '0')}`,
    );
  }
  return days;
}

function shortDay(dateStr: string, locale: string): string {
  const d = new Date(dateStr + 'T00:00:00');
  const dayName = new Intl.DateTimeFormat(locale, { weekday: 'short' }).format(d);
  return `${d.getDate()}(${dayName})`;
}

/* Colors chosen for maximum contrast between the three series */
const COLORS = {
  created: '#0d9488', // teal-600
  approved: '#f59e0b', // amber-500 (was green — too close to teal)
  done: '#8b5cf6', // violet-500
};

export function WeeklyChart() {
  const { t, i18n } = useTranslation();
  const now = new Date();
  const stats = useMonthStats(now.getFullYear(), now.getMonth() + 1);

  const days = useMemo(getLast7Days, []);

  const data = useMemo(() => {
    const map = new Map<string, { created: number; approved: number; completed: number }>();
    (stats.data ?? []).forEach((d) => {
      const key = d.date.slice(0, 10);
      map.set(key, { created: d.created, approved: d.approved, completed: d.completed });
    });
    return days.map((day) => ({
      name: shortDay(day, i18n.language),
      [t('home.created')]: map.get(day)?.created ?? 0,
      [t('home.approved')]: map.get(day)?.approved ?? 0,
      [t('home.done')]: map.get(day)?.completed ?? 0,
    }));
  }, [stats.data, days, t, i18n.language]);

  const hasData = data.some((d) =>
    Object.values(d).some((v) => typeof v === 'number' && v > 0),
  );

  /* Empty state — show a placeholder instead of hiding entirely */
  if (!hasData && !stats.isLoading) {
    return (
      <section className="fade-up flex flex-col gap-3" style={{ animationDelay: '80ms' }}>
        <h2 className="text-base font-semibold">{t('home.weeklyActivity')}</h2>
        <div className="flex flex-col items-center gap-2 rounded-xl border border-dashed border-edge-base bg-surface-subtle py-10 text-center">
          <BarChart3 size={24} className="text-ink-muted opacity-30" />
          <p className="text-xs text-ink-muted">
            {t('chart.noData', '아직 활동 데이터가 없습니다')}
          </p>
        </div>
      </section>
    );
  }

  return (
    <section
      className="fade-up flex flex-col gap-3"
      style={{ animationDelay: '80ms' }}
      aria-label={t('home.weeklyActivity')}
    >
      <h2 className="text-base font-semibold">{t('home.weeklyActivity')}</h2>
      <div
        className="rounded-xl border border-edge-base bg-surface-subtle p-4"
        style={{ height: 280 }}
      >
        <ResponsiveContainer width="100%" height="100%">
          <BarChart data={data} barGap={2} barCategoryGap="25%">
            <CartesianGrid
              strokeDasharray="3 3"
              stroke="var(--border-base)"
              strokeOpacity={0.5}
              vertical={false}
            />
            <XAxis
              dataKey="name"
              tick={{ fill: 'var(--text-muted)', fontSize: 11 }}
              axisLine={{ stroke: 'var(--border-base)' }}
              tickLine={false}
            />
            <YAxis
              allowDecimals={false}
              tick={{ fill: 'var(--text-muted)', fontSize: 11 }}
              axisLine={false}
              tickLine={false}
              width={28}
            />
            <Tooltip
              cursor={{ fill: 'var(--bg-muted)', opacity: 0.4, radius: 4 }}
              contentStyle={{
                background: 'var(--bg-subtle)',
                border: '1px solid var(--border-base)',
                borderRadius: 8,
                fontSize: 12,
                color: 'var(--text-primary)',
                boxShadow: '0 4px 6px -1px rgb(0 0 0 / 0.1)',
              }}
            />
            <Legend
              verticalAlign="top"
              align="right"
              height={32}
              iconType="circle"
              iconSize={8}
              wrapperStyle={{ fontSize: 11, color: 'var(--text-secondary)' }}
            />
            <Bar
              dataKey={t('home.created')}
              fill={COLORS.created}
              radius={[4, 4, 0, 0]}
              maxBarSize={32}
            />
            <Bar
              dataKey={t('home.approved')}
              fill={COLORS.approved}
              radius={[4, 4, 0, 0]}
              maxBarSize={32}
            />
            <Bar
              dataKey={t('home.done')}
              fill={COLORS.done}
              radius={[4, 4, 0, 0]}
              maxBarSize={32}
            />
          </BarChart>
        </ResponsiveContainer>
      </div>
    </section>
  );
}
