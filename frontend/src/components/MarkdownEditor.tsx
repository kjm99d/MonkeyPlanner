import { useState } from 'react';
import ReactMarkdown from 'react-markdown';
import rehypeSanitize from 'rehype-sanitize';
import { useTranslation } from 'react-i18next';

type Props = {
  value: string;
  onChange: (v: string) => void;
  label?: string;
};

export function MarkdownEditor({ value, onChange, label }: Props) {
  const { t } = useTranslation();
  const [tab, setTab] = useState<'edit' | 'preview' | 'split'>('split');
  const resolvedLabel = label ?? t('issue.body');

  return (
    <section aria-label={t('issue.body')} className="flex flex-col gap-2">
      <div className="flex items-center justify-between">
        <span className="text-sm font-medium text-ink-secondary">{resolvedLabel}</span>
        <div role="tablist" className="flex gap-1 rounded-md bg-surface-muted p-0.5 text-xs">
          {(['edit', 'split', 'preview'] as const).map((tab_) => (
            <button
              key={tab_}
              role="tab"
              aria-selected={tab === tab_}
              onClick={() => setTab(tab_)}
              className={`rounded px-2 py-1 transition-colors ${
                tab === tab_ ? 'bg-surface-base text-ink-primary shadow-sm' : 'text-ink-secondary hover:text-ink-primary'
              }`}
            >
              {tab_ === 'edit' ? t('editor.edit') : tab_ === 'preview' ? t('editor.preview') : t('editor.split')}
            </button>
          ))}
        </div>
      </div>

      <div className={`grid gap-3 ${tab === 'split' ? 'md:grid-cols-2' : 'grid-cols-1'}`}>
        {tab !== 'preview' && (
          <textarea
            value={value}
            onChange={(e) => onChange(e.target.value)}
            spellCheck={false}
            className="min-h-[16rem] rounded-md border border-edge-base bg-surface-base p-3 font-mono text-sm text-ink-primary focus-visible:border-brand-500 focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-brand-500/40"
            aria-label={t('editor.placeholder')}
          />
        )}
        {tab !== 'edit' && (
          <article
            aria-label="마크다운 미리보기"
            className="prose prose-sm max-w-none rounded-md border border-edge-base bg-surface-subtle p-3 text-ink-primary dark:prose-invert"
          >
            <ReactMarkdown rehypePlugins={[rehypeSanitize]}>
              {value || `_${t('issue.noPreview')}_`}
            </ReactMarkdown>
          </article>
        )}
      </div>
    </section>
  );
}
