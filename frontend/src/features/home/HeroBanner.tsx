import { useTranslation } from 'react-i18next';
import { Rocket, ArrowRight } from 'lucide-react';
import { Link } from 'react-router-dom';

export default function HeroBanner() {
  const { t } = useTranslation();
  return (
    <div className="flex items-center justify-between gap-4 rounded-xl border border-edge-base bg-gradient-to-r from-brand-500/8 to-surface-subtle px-5 py-4">
      <div className="flex items-center gap-3">
        <div className="flex h-10 w-10 shrink-0 items-center justify-center rounded-lg bg-brand-500/15 text-brand-500">
          <Rocket size={20} />
        </div>
        <div>
          <p className="text-sm font-semibold text-ink-primary">{t('app.tagline')}</p>
          <p className="text-xs text-ink-secondary">{t('app.description')}</p>
        </div>
      </div>
      <Link
        to="/boards"
        className="hidden sm:flex items-center gap-1 rounded-lg bg-brand-500 px-3 py-1.5 text-xs font-medium text-white shadow-sm transition-colors hover:bg-brand-600"
      >
        {t('nav.boards')} <ArrowRight size={12} />
      </Link>
    </div>
  );
}
