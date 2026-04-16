import { useTranslation } from 'react-i18next';
import { Globe } from 'lucide-react';
import { LANGUAGES } from '../i18n';

export function LanguageSwitcher() {
  const { i18n } = useTranslation();

  return (
    <div className="relative group">
      <button
        type="button"
        className="flex h-9 w-9 items-center justify-center rounded-lg border border-edge-base text-ink-secondary transition-all duration-200 hover:bg-surface-muted hover:text-brand-500 active:scale-95"
        aria-label="Language"
      >
        <Globe size={18} />
      </button>
      <div className="invisible absolute right-0 top-full mt-1 min-w-[120px] rounded-lg border border-edge-base bg-surface-base py-1 shadow-lg opacity-0 transition-all group-hover:visible group-hover:opacity-100">
        {LANGUAGES.map((lang) => (
          <button
            key={lang.code}
            type="button"
            onClick={() => i18n.changeLanguage(lang.code)}
            className={`block w-full px-3 py-1.5 text-left text-sm transition-colors ${
              i18n.language === lang.code
                ? 'bg-brand-500/10 font-medium text-brand-500'
                : 'text-ink-secondary hover:bg-surface-muted'
            }`}
          >
            {lang.label}
          </button>
        ))}
      </div>
    </div>
  );
}
