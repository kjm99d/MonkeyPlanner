import { useState, useEffect, useMemo } from 'react';
import { NavLink, Outlet, Link, useNavigate } from 'react-router-dom';
import { useTranslation } from 'react-i18next';
import { Home, LayoutDashboard, Calendar, CheckCircle2, Squirrel, Plus, Menu, X, PanelLeftClose, PanelLeft, Search } from 'lucide-react';
import { ThemeToggle } from './ThemeToggle';
import { LanguageSwitcher } from './LanguageSwitcher';
import { useBoards, useIssues } from '../api/hooks';
import { SearchDialog } from '../components/SearchDialog';
import { ShortcutsDialog } from '../components/ShortcutsDialog';

const navLinkClass = ({ isActive }: { isActive: boolean }) =>
  `flex items-center gap-2 rounded-md px-2.5 py-1.5 text-sm transition-colors ${
    isActive
      ? 'bg-brand-500/10 font-medium text-brand-500'
      : 'text-ink-secondary hover:bg-surface-muted hover:text-ink-primary'
  }`;

const boardLinkClass = ({ isActive }: { isActive: boolean }) =>
  `flex items-center gap-2 rounded-md px-2.5 py-1.5 text-sm transition-colors ${
    isActive
      ? 'bg-surface-muted font-medium text-ink-primary'
      : 'text-ink-secondary hover:bg-surface-muted hover:text-ink-primary'
  }`;

function Sidebar({ onNavigate, hideHeader, collapsed }: { onNavigate?: () => void; hideHeader?: boolean; collapsed?: boolean }) {
  const { t } = useTranslation();
  const boards = useBoards();
  const allIssues = useIssues({});

  const issueCounts = useMemo(() => {
    const map: Record<string, number> = {};
    (allIssues.data ?? []).forEach(i => {
      map[i.boardId] = (map[i.boardId] ?? 0) + 1;
    });
    return map;
  }, [allIssues.data]);

  const recentIssues = useMemo(() => {
    return (allIssues.data ?? [])
      .sort((a, b) => new Date(b.createdAt).getTime() - new Date(a.createdAt).getTime())
      .slice(0, 3);
  }, [allIssues.data]);

  const pendingCount = useMemo(() => (allIssues.data ?? []).filter(i => i.status === 'Pending').length, [allIssues.data]);

  return (
    <div className="flex h-full flex-col">
      {/* Logo — hidden in mobile drawer (drawer renders its own header) */}
      {!hideHeader && (
        <div className={`flex h-14 items-center border-b border-edge-base ${collapsed ? 'justify-center px-2' : 'gap-2.5 px-4'}`}>
          <div className="flex h-7 w-7 shrink-0 items-center justify-center rounded-lg bg-brand-500 text-white">
            <Squirrel size={16} />
          </div>
          {!collapsed && (
            <span className="font-semibold text-ink-primary">{t('app.name')}</span>
          )}
        </div>
      )}

      {/* Search hint — only shown when expanded */}
      {!collapsed && (
        <div className="px-2 pt-2">
          <button
            type="button"
            onClick={() => {/* SearchDialog handles its own open state via Cmd+K */}}
            className="mx-2 mb-2 flex items-center gap-2 rounded-md border border-edge-base bg-surface-muted px-2.5 py-1.5 text-xs text-ink-muted transition-colors hover:text-ink-primary w-full"
          >
            <Search size={14} />
            <span className="flex-1 text-left">Search...</span>
            <kbd className="rounded border border-edge-base bg-surface-base px-1 py-0.5 text-[10px]">⌘K</kbd>
          </button>
        </div>
      )}

      {/* Collapsed search icon */}
      {collapsed && (
        <div className="flex justify-center px-2 pt-2 pb-1">
          <button
            type="button"
            onClick={() => {/* SearchDialog handles its own open state via Cmd+K */}}
            className="rounded-md p-1.5 text-ink-muted hover:bg-surface-muted hover:text-ink-primary transition-colors"
          >
            <Search size={16} />
          </button>
        </div>
      )}

      {/* Recent issues section */}
      {!collapsed && recentIssues.length > 0 && (
        <div className="flex flex-col gap-0.5 border-t border-edge-base px-2 pt-2 pb-1">
          <span className="px-2.5 py-1 text-[11px] font-semibold uppercase tracking-wider text-ink-muted">
            Recent
          </span>
          {recentIssues.map(issue => (
            <NavLink
              key={issue.id}
              to={`/issues/${issue.id}`}
              onClick={onNavigate}
              className={({ isActive }) =>
                `flex items-center gap-2 rounded-md px-2.5 py-1 text-xs transition-colors truncate ${
                  isActive ? 'bg-surface-muted text-ink-primary' : 'text-ink-muted hover:bg-surface-muted hover:text-ink-secondary'
                }`
              }
            >
              <span className="truncate">{issue.title}</span>
            </NavLink>
          ))}
        </div>
      )}

      {/* Main nav */}
      <nav aria-label={t('nav.menu')} className="flex flex-col gap-0.5 px-2 pt-1 pb-2">
        <NavLink
          to="/"
          end
          className={navLinkClass}
          onClick={onNavigate}
          title={collapsed ? t('nav.home') : undefined}
        >
          <Home size={16} className="shrink-0" />
          {!collapsed && t('nav.home')}
        </NavLink>
        <NavLink
          to="/calendar"
          className={navLinkClass}
          onClick={onNavigate}
          title={collapsed ? t('nav.calendar') : undefined}
        >
          <Calendar size={16} className="shrink-0" />
          {!collapsed && t('nav.calendar')}
        </NavLink>
        <NavLink to="/approve" className={navLinkClass} onClick={onNavigate}
          title={collapsed ? t('approval.title') : undefined}>
          <CheckCircle2 size={16} className="shrink-0" />
          {!collapsed && t('approval.title')}
          {!collapsed && pendingCount > 0 && (
            <span className="ml-auto rounded-full bg-accent/15 px-1.5 py-0.5 text-[10px] tabular-nums text-accent font-medium">
              {pendingCount}
            </span>
          )}
        </NavLink>
      </nav>

      {/* Boards section — hidden when collapsed */}
      {!collapsed && (
        <div className="flex flex-col gap-0.5 border-t border-edge-base px-2 pt-3">
          <div className="flex items-center justify-between px-2.5 py-1">
            <span className="text-[11px] font-semibold uppercase tracking-wider text-ink-muted">
              {t('nav.boards')}
            </span>
            <Link
              to="/boards"
              onClick={onNavigate}
              className="rounded p-0.5 text-ink-muted transition-colors hover:bg-surface-muted hover:text-ink-primary"
              title={t('board.create')}
            >
              <Plus size={14} />
            </Link>
          </div>
          <div className="flex flex-col gap-0.5 overflow-y-auto">
            {boards.data?.map((b) => (
              <NavLink
                key={b.id}
                to={`/boards/${b.id}`}
                className={boardLinkClass}
                onClick={onNavigate}
              >
                <LayoutDashboard size={14} className="shrink-0 opacity-50" />
                <span className="truncate">{b.name}</span>
                {!collapsed && issueCounts[b.id] > 0 && (
                  <span className="ml-auto rounded-full bg-surface-muted px-1.5 py-0.5 text-[10px] tabular-nums text-ink-muted">
                    {issueCounts[b.id]}
                  </span>
                )}
              </NavLink>
            ))}
            {boards.data?.length === 0 && (
              <p className="px-2.5 py-2 text-xs text-ink-muted">{t('board.noBoards')}</p>
            )}
          </div>
        </div>
      )}

      {/* Spacer */}
      <div className="flex-1" />

      {/* Footer */}
      <div className={`flex items-center border-t border-edge-base px-4 py-3 ${collapsed ? 'flex-col gap-2 px-2' : 'gap-2'}`}>
        <LanguageSwitcher />
        <ThemeToggle />
      </div>
    </div>
  );
}

export default function Layout() {
  const { t } = useTranslation();
  const navigate = useNavigate();
  const [sidebarOpen, setSidebarOpen] = useState(false);
  const [collapsed, setCollapsed] = useState(() => localStorage.getItem('sidebar-collapsed') === 'true');

  useEffect(() => {
    localStorage.setItem('sidebar-collapsed', String(collapsed));
  }, [collapsed]);

  useEffect(() => {
    if (!sidebarOpen) return;
    const handler = (e: KeyboardEvent) => {
      if (e.key === 'Escape') setSidebarOpen(false);
    };
    document.addEventListener('keydown', handler);
    return () => document.removeEventListener('keydown', handler);
  }, [sidebarOpen]);

  useEffect(() => {
    const handler = (e: KeyboardEvent) => {
      const tag = (e.target as HTMLElement)?.tagName;
      if (tag === 'INPUT' || tag === 'TEXTAREA' || tag === 'SELECT') return;
      if ((e.target as HTMLElement)?.isContentEditable) return;

      if (e.key === 'a' && !e.metaKey && !e.ctrlKey) {
        e.preventDefault();
        navigate('/approve');
      }
      if (e.key === 'h' && !e.metaKey && !e.ctrlKey) {
        e.preventDefault();
        navigate('/');
      }
    };
    document.addEventListener('keydown', handler);
    return () => document.removeEventListener('keydown', handler);
  }, [navigate]);

  return (
    <div className="flex h-screen overflow-hidden bg-surface-base text-ink-primary">
      {/* Desktop sidebar */}
      <aside className={`hidden lg:flex ${collapsed ? 'w-14' : 'w-60'} shrink-0 flex-col border-r border-edge-base bg-surface-subtle h-full transition-[width] duration-200`}>
        <div className="flex flex-col h-full overflow-hidden">
          <div className="flex-1 overflow-hidden">
            <Sidebar collapsed={collapsed} />
          </div>
          <button
            type="button"
            onClick={() => setCollapsed(c => !c)}
            className="mx-auto my-2 rounded-md p-1.5 text-ink-muted hover:bg-surface-muted hover:text-ink-primary transition-colors"
            aria-label={collapsed ? 'Expand sidebar' : 'Collapse sidebar'}
          >
            {collapsed ? <PanelLeft size={16} /> : <PanelLeftClose size={16} />}
          </button>
        </div>
      </aside>

      {/* Mobile sidebar overlay */}
      {sidebarOpen && (
        <div className="fixed inset-0 z-40 lg:hidden">
          <div
            className="absolute inset-0 bg-black/40 transition-opacity"
            onClick={() => setSidebarOpen(false)}
            aria-hidden
          />
          <aside
            className="relative z-50 flex h-full w-60 flex-col bg-surface-subtle shadow-lg animate-in slide-in-from-left"
            role="dialog"
            aria-modal="true"
          >
            <div className="flex h-14 shrink-0 items-center justify-between border-b border-edge-base px-4">
              <div className="flex items-center gap-2.5">
                <div className="flex h-7 w-7 items-center justify-center rounded-lg bg-brand-500 text-white">
                  <Squirrel size={16} />
                </div>
                <span className="font-semibold text-ink-primary">{t('app.name')}</span>
              </div>
              <button
                type="button"
                onClick={() => setSidebarOpen(false)}
                className="rounded-md p-1 text-ink-secondary hover:bg-surface-muted"
                aria-label={t('common.close')}
              >
                <X size={18} />
              </button>
            </div>
            <div className="flex-1 overflow-y-auto">
              <Sidebar onNavigate={() => setSidebarOpen(false)} hideHeader />
            </div>
          </aside>
        </div>
      )}

      {/* Main content */}
      <div className="flex flex-1 flex-col overflow-hidden">
        {/* Mobile top bar */}
        <header className="flex h-14 items-center gap-3 border-b border-edge-base px-4 lg:hidden">
          <button
            type="button"
            onClick={() => setSidebarOpen(true)}
            className="rounded-md p-1.5 text-ink-secondary hover:bg-surface-muted"
            aria-label={t('nav.menu')}
          >
            <Menu size={20} />
          </button>
          <span className="font-semibold">{t('app.name')}</span>
        </header>

        <main className="flex-1 overflow-auto min-h-0">
          <div className="mx-auto max-w-6xl px-4 py-6 sm:px-6 lg:px-8 lg:py-8">
            <Outlet />
          </div>
        </main>
      </div>
      <SearchDialog />
      <ShortcutsDialog />
    </div>
  );
}
