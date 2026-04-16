import { NavLink, Outlet } from 'react-router-dom';
import { useTranslation } from 'react-i18next';
import { Home, LayoutDashboard, Calendar, Squirrel } from 'lucide-react';
import { ThemeToggle } from './ThemeToggle';
import { LanguageSwitcher } from './LanguageSwitcher';

const linkClass = ({ isActive }: { isActive: boolean }) =>
  `flex items-center gap-1.5 rounded-lg px-3 py-2 text-sm font-medium transition-all duration-150 ${
    isActive
      ? 'bg-brand-500 text-white shadow-sm'
      : 'text-ink-secondary hover:bg-surface-muted hover:text-ink-primary'
  }`;

export default function Layout() {
  const { t } = useTranslation();
  return (
    <div className="min-h-screen bg-surface-base text-ink-primary">
      <header className="sticky top-0 z-10 border-b border-edge-base bg-surface-base/95 backdrop-blur-md">
        <div className="mx-auto flex max-w-7xl items-center justify-between px-4 py-2.5 sm:px-6">
          {/* 로고 — 네비와 시각적으로 분리 */}
          <a href="/" className="flex items-center gap-2 font-bold tracking-tight text-ink-primary">
            <div className="flex h-8 w-8 items-center justify-center rounded-lg bg-brand-500 text-white">
              <Squirrel size={18} />
            </div>
            <span className="hidden sm:inline">{t('app.name')}</span>
          </a>

          {/* 네비 — 중앙 정렬 */}
          <nav aria-label="주 메뉴" className="flex items-center gap-1">
            <NavLink to="/" end className={linkClass}>
              <Home size={16} /> <span className="hidden xs:inline">{t('nav.home')}</span>
            </NavLink>
            <NavLink to="/boards" className={linkClass}>
              <LayoutDashboard size={16} /> <span className="hidden xs:inline">{t('nav.boards')}</span>
            </NavLink>
            <NavLink to="/calendar" className={linkClass}>
              <Calendar size={16} /> <span className="hidden xs:inline">{t('nav.calendar')}</span>
            </NavLink>
          </nav>

          {/* 우측 유틸 */}
          <div className="flex items-center gap-2">
            <LanguageSwitcher />
            <ThemeToggle />
          </div>
        </div>
      </header>
      <main className="mx-auto max-w-7xl px-4 py-6 sm:px-6 sm:py-8">
        <Outlet />
      </main>
    </div>
  );
}
