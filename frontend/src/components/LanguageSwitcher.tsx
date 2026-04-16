import { useState, useRef, useEffect } from 'react';
import { useTranslation } from 'react-i18next';
import { Globe } from 'lucide-react';
import { LANGUAGES } from '../i18n';

export function LanguageSwitcher() {
  const { i18n } = useTranslation();
  const [open, setOpen] = useState(false);
  const ref = useRef<HTMLDivElement>(null);

  useEffect(() => {
    if (!open) return;
    const handleClickOutside = (e: MouseEvent) => {
      if (ref.current && !ref.current.contains(e.target as Node)) setOpen(false);
    };
    document.addEventListener('mousedown', handleClickOutside);
    return () => document.removeEventListener('mousedown', handleClickOutside);
  }, [open]);

  return (
    <div className="relative" ref={ref}>
      <button
        type="button"
        onClick={() => setOpen((v) => !v)}
        className="flex h-9 w-9 items-center justify-center rounded-lg border border-edge-base text-ink-secondary transition-all duration-200 hover:bg-surface-muted hover:text-brand-500 active:scale-95"
        aria-label="Language"
      >
        <Globe size={18} />
      </button>
      {open && (
        <div className="absolute left-0 bottom-full mb-1 min-w-[120px] rounded-lg border border-edge-base bg-surface-base py-1 shadow-lg z-50">
          {LANGUAGES.map((lang) => (
            <button
              key={lang.code}
              type="button"
              onClick={() => { i18n.changeLanguage(lang.code); setOpen(false); }}
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
      )}
    </div>
  );
}
