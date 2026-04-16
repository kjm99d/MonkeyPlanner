import { useMemo } from 'react';
import { useTranslation } from 'react-i18next';
import { BarChart, Bar, XAxis, YAxis, Tooltip, ResponsiveContainer, CartesianGrid } from 'recharts';
import { useMonthStats } from '../../api/hooks';

function getLast7Days(): string[] {
  const days: string[] = [];
  const now = new Date();
  for (let i = 6; i >= 0; i--) {
    const d = new Date(now);
    d.setDate(d.getDate() - i);
    days.push(`${d.getFullYear()}-${String(d.getMonth() + 1).padStart(2, '0')}-${String(d.getDate()).padStart(2, '0')}`);
  }
  return days;
}

function shortDay(dateStr: string): string {
  const d = new Date(dateStr + 'T00:00:00');
  const weekday = ['일', '월', '화', '수', '목', '금', '토'];
  return `${d.getDate()}(${weekday[d.getDay()]})`;
}

export function WeeklyChart() {
  const { t } = useTranslation();
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
      name: shortDay(day),
      [t('home.created')]: map.get(day)?.created ?? 0,
      [t('home.approved')]: map.get(day)?.approved ?? 0,
      [t('home.done')]: map.get(day)?.completed ?? 0,
    }));
  }, [stats.data, days, t]);

  const hasData = data.some((d) => Object.values(d).some((v) => typeof v === 'number' && v > 0));

  if (!hasData && !stats.isLoading) {
    return null;
  }

  return (
    <section className="fade-up flex flex-col gap-3" style={{ animationDelay: '80ms' }}>
      <h2 className="text-lg font-semibold">{t('home.weeklyActivity', '주간 활동')}</h2>
      <div className="rounded-xl border border-edge-base bg-surface-subtle p-4" style={{ height: 220 }}>
        <ResponsiveContainer width="100%" height="100%">
          <BarChart data={data} barGap={2} barCategoryGap="20%">
            <CartesianGrid strokeDasharray="3 3" stroke="var(--border-base)" vertical={false} />
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
              width={24}
            />
            <Tooltip
              contentStyle={{
                background: 'var(--bg-subtle)',
                border: '1px solid var(--border-base)',
                borderRadius: 8,
                fontSize: 12,
                color: 'var(--text-primary)',
              }}
            />
            <Bar dataKey={t('home.created')} fill="#0d9488" radius={[3, 3, 0, 0]} />
            <Bar dataKey={t('home.approved')} fill="#16a34a" radius={[3, 3, 0, 0]} />
            <Bar dataKey={t('home.done')} fill="#8b5cf6" radius={[3, 3, 0, 0]} />
          </BarChart>
        </ResponsiveContainer>
      </div>
    </section>
  );
}
