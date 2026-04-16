import { Link } from 'react-router-dom';
import { useTranslation } from 'react-i18next';

type Crumb = { label: string; to?: string };

export function Breadcrumb({ items }: { items: Crumb[] }) {
  const { t } = useTranslation();
  return (
    <nav aria-label={t('common.location')} className="flex items-center gap-1 text-sm text-ink-muted">
      {items.map((item, i) => (
        <span key={i} className="flex items-center gap-1">
          {i > 0 && <span aria-hidden>/</span>}
          {item.to ? (
            <Link to={item.to} className="hover:text-ink-primary hover:underline">
              {item.label}
            </Link>
          ) : (
            <span className="text-ink-secondary">{item.label}</span>
          )}
        </span>
      ))}
    </nav>
  );
}
