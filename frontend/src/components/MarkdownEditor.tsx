import { useState } from 'react';
import ReactMarkdown from 'react-markdown';
import rehypeSanitize from 'rehype-sanitize';

type Props = {
  value: string;
  onChange: (v: string) => void;
  label?: string;
};

export function MarkdownEditor({ value, onChange, label = '본문 (Markdown)' }: Props) {
  const [tab, setTab] = useState<'edit' | 'preview' | 'split'>('split');

  return (
    <section aria-label="마크다운 에디터" className="flex flex-col gap-2">
      <div className="flex items-center justify-between">
        <span className="text-sm font-medium text-ink-secondary">{label}</span>
        <div role="tablist" className="flex gap-1 rounded-md bg-surface-muted p-0.5 text-xs">
          {(['edit', 'split', 'preview'] as const).map((t) => (
            <button
              key={t}
              role="tab"
              aria-selected={tab === t}
              onClick={() => setTab(t)}
              className={`rounded px-2 py-1 transition-colors ${
                tab === t ? 'bg-surface-base text-ink-primary shadow-sm' : 'text-ink-secondary hover:text-ink-primary'
              }`}
            >
              {t === 'edit' ? '편집' : t === 'preview' ? '미리보기' : '분할'}
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
            aria-label="마크다운 본문 편집기"
          />
        )}
        {tab !== 'edit' && (
          <article
            aria-label="마크다운 미리보기"
            className="prose prose-sm max-w-none rounded-md border border-edge-base bg-surface-subtle p-3 text-ink-primary dark:prose-invert"
          >
            <ReactMarkdown rehypePlugins={[rehypeSanitize]}>
              {value || '_미리볼 내용이 없습니다._'}
            </ReactMarkdown>
          </article>
        )}
      </div>
    </section>
  );
}
