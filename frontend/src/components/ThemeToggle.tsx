import { useTheme } from '../lib/theme';

export function ThemeToggle() {
  const { mode, toggle } = useTheme();
  return (
    <button
      type="button"
      onClick={toggle}
      aria-label={mode === 'dark' ? '라이트 모드로 전환' : '다크 모드로 전환'}
      className="inline-flex h-9 w-9 items-center justify-center rounded-md border border-edge-base text-ink-secondary transition-colors hover:bg-surface-muted focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-brand-500"
    >
      <span aria-hidden>{mode === 'dark' ? '☀' : '◐'}</span>
    </button>
  );
}
