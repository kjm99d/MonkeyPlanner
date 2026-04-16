import { Link } from 'react-router-dom';
import { lazy, Suspense, useMemo, useEffect, useState } from 'react';
import { useDayStats, useIssues } from '../../api/hooks';
import { Card } from '../../components/Card';
import Hero3DFallback from './Hero3DFallback';

const Hero3D = lazy(() => import('./Hero3D'));

function todayString(): string {
  const d = new Date();
  const yyyy = d.getFullYear();
  const mm = String(d.getMonth() + 1).padStart(2, '0');
  const dd = String(d.getDate()).padStart(2, '0');
  return `${yyyy}-${mm}-${dd}`;
}

export default function HomePage() {
  const date = useMemo(todayString, []);
  const day = useDayStats(date);
  const recent = useIssues({ status: 'Approved' });

  const created = day.data?.created.length ?? 0;
  const approved = day.data?.approved.length ?? 0;
  const completed = day.data?.completed.length ?? 0;

  const [enable3D, setEnable3D] = useState(false);
  useEffect(() => {
    if (typeof window === 'undefined') return;
    const m = window.matchMedia('(prefers-reduced-motion: reduce)');
    setEnable3D(!m.matches);
    const handler = (e: MediaQueryListEvent) => setEnable3D(!e.matches);
    m.addEventListener('change', handler);
    return () => m.removeEventListener('change', handler);
  }, []);

  return (
    <section className="flex flex-col gap-8">
      <header className="flex flex-col gap-4">
        <div>
          <h1 className="text-4xl font-bold tracking-tight">오늘의 코숭이</h1>
          <p className="mt-2 text-ink-secondary">
            오늘 <time dateTime={date}>{date}</time> 의 활동 요약입니다.
          </p>
        </div>
        {enable3D ? (
          <Suspense fallback={<Hero3DFallback />}>
            <Hero3D />
          </Suspense>
        ) : (
          <Hero3DFallback />
        )}
      </header>

      <div className="grid gap-4 sm:grid-cols-3">
        <StatCard label="생성" value={created} hue="text-brand-500" />
        <StatCard label="승인" value={approved} hue="text-status-approved" />
        <StatCard label="완료" value={completed} hue="text-status-done" />
      </div>

      <section aria-label="최근 승인된 이슈" className="flex flex-col gap-3">
        <h2 className="text-xl font-semibold">최근 승인된 이슈</h2>
        {recent.isLoading && <p className="text-ink-secondary">불러오는 중…</p>}
        {recent.data && recent.data.length === 0 && (
          <p className="text-ink-secondary">아직 승인된 이슈가 없습니다.</p>
        )}
        <ul className="flex flex-col gap-2">
          {recent.data?.slice(0, 5).map((iss) => (
            <li key={iss.id}>
              <Link
                to={`/issues/${iss.id}`}
                className="block rounded-md border border-edge-base bg-surface-subtle px-4 py-3 transition-colors hover:bg-surface-muted focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-brand-500"
              >
                <div className="flex items-center justify-between">
                  <span className="font-medium">{iss.title}</span>
                  <time className="text-xs text-ink-muted" dateTime={iss.approvedAt ?? iss.updatedAt}>
                    {new Date(iss.approvedAt ?? iss.updatedAt).toLocaleString('ko-KR')}
                  </time>
                </div>
              </Link>
            </li>
          ))}
        </ul>
      </section>
    </section>
  );
}

function StatCard({ label, value, hue }: { label: string; value: number; hue: string }) {
  return (
    <Card className="transition-transform hover:[transform:perspective(800px)_rotateX(2deg)_rotateY(-2deg)] motion-reduce:hover:transform-none">
      <div className="flex flex-col gap-1">
        <span className="text-sm font-medium text-ink-secondary">{label}</span>
        <span className={`text-5xl font-extrabold tabular-nums tracking-tight ${hue}`}>{value}</span>
      </div>
    </Card>
  );
}
