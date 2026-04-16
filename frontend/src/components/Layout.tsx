import { NavLink, Outlet } from 'react-router-dom';
import { Home, LayoutDashboard, Calendar } from 'lucide-react';
import { ThemeToggle } from './ThemeToggle';

const linkClass = ({ isActive }: { isActive: boolean }) =>
  `rounded-md px-3 py-1.5 text-sm font-medium transition-colors ${
    isActive
      ? 'bg-brand-500 text-white'
      : 'text-ink-secondary hover:bg-surface-muted hover:text-ink-primary'
  }`;

export default function Layout() {
  return (
    <div className="min-h-screen bg-surface-base text-ink-primary">
      <header className="sticky top-0 z-10 border-b border-b-brand-500/15 bg-surface-base/80 backdrop-blur">
        <div className="mx-auto flex max-w-7xl items-center justify-between px-6 py-3">
          <div className="flex items-center gap-6">
            <a href="/" className="font-bold tracking-tight">
              🐒 몽키 플래너
            </a>
            <nav aria-label="주 메뉴" className="flex gap-1">
              <NavLink to="/" end className={linkClass}>
                <Home size={16} /> 홈
              </NavLink>
              <NavLink to="/boards" className={linkClass}>
                <LayoutDashboard size={16} /> 보드
              </NavLink>
              <NavLink to="/calendar" className={linkClass}>
                <Calendar size={16} /> 캘린더
              </NavLink>
            </nav>
          </div>
          <ThemeToggle />
        </div>
      </header>
      <main className="mx-auto max-w-7xl px-6 py-8">
        <Outlet />
      </main>
    </div>
  );
}
