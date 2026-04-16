import { useTranslation } from 'react-i18next';
import { Banana, Brain, Zap } from 'lucide-react';

export default function HeroBanner() {
  const { t } = useTranslation();
  return (
    <div className="relative overflow-hidden rounded-xl border border-edge-base bg-gradient-to-br from-brand-500/5 via-surface-subtle to-surface-muted px-6 py-8">
      {/* 배경 장식 원 */}
      <div className="absolute -right-8 -top-8 h-32 w-32 rounded-full bg-brand-500/10 blur-2xl" />
      <div className="absolute -bottom-6 -left-6 h-24 w-24 rounded-full bg-accent/10 blur-2xl" />

      <div className="relative flex items-center justify-between gap-6">
        <div className="flex flex-col gap-2">
          <h2 className="text-lg font-bold text-ink-primary">
            {t('app.tagline')}
          </h2>
          <p className="text-sm text-ink-secondary">
            {t('app.description')}
          </p>
        </div>
        <div className="hidden sm:flex items-center gap-3">
          <div className="flex h-12 w-12 items-center justify-center rounded-xl bg-brand-500/15 text-brand-500">
            <Brain size={24} />
          </div>
          <div className="flex h-12 w-12 items-center justify-center rounded-xl bg-accent/15 text-accent">
            <Banana size={24} />
          </div>
          <div className="flex h-12 w-12 items-center justify-center rounded-xl bg-status-inProgress/15 text-status-inProgress">
            <Zap size={24} />
          </div>
        </div>
      </div>
    </div>
  );
}
