import { Sun, Moon } from 'lucide-react';
import { useTranslation } from 'react-i18next';
import { useTheme } from '../lib/theme';

export function ThemeToggle() {
  const { mode, toggle } = useTheme();
  const { t } = useTranslation();
  return (
    <button
      type="button"
      onClick={toggle}
      aria-label={mode === 'dark' ? t('theme.toLight') : t('theme.toDark')}
      className="inline-flex h-9 w-9 cursor-pointer items-center justify-center rounded-lg border border-edge-base text-ink-secondary transition-all duration-200 hover:bg-surface-muted hover:text-brand-500 hover:border-brand-500/30 focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-brand-500 active:scale-95"
    >
      {mode === 'dark' ? <Sun size={18} /> : <Moon size={18} />}
    </button>
  );
}
