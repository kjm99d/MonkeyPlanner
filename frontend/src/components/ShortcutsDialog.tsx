import { useState, useEffect } from 'react';
import { Keyboard } from 'lucide-react';

interface Shortcut {
  keys: string[];
  description: string;
}

const SHORTCUTS: Shortcut[] = [
  { keys: ['⌘', 'K'], description: 'Open search' },
  { keys: ['⌘', 'S'], description: 'Save changes' },
  { keys: ['Esc'], description: 'Close / cancel' },
  { keys: ['?'], description: 'Show keyboard shortcuts' },
];

function ShortcutRow({ shortcut }: { shortcut: Shortcut }) {
  return (
    <div className="flex items-center justify-between py-2">
      <span className="text-sm text-ink-secondary">{shortcut.description}</span>
      <div className="flex items-center gap-1">
        {shortcut.keys.map((key, idx) => (
          <kbd
            key={idx}
            className="inline-flex items-center rounded border border-edge-base bg-surface-muted px-2 py-0.5 text-xs font-medium text-ink-secondary"
          >
            {key}
          </kbd>
        ))}
      </div>
    </div>
  );
}

export function ShortcutsDialog() {
  const [open, setOpen] = useState(false);

  useEffect(() => {
    const handler = (e: KeyboardEvent) => {
      const target = e.target as HTMLElement;
      const isInputFocused =
        target.tagName === 'INPUT' ||
        target.tagName === 'TEXTAREA' ||
        target.isContentEditable;

      if (e.key === '?' && !isInputFocused) {
        e.preventDefault();
        setOpen(prev => !prev);
      }
      if (e.key === 'Escape' && open) {
        setOpen(false);
      }
    };
    document.addEventListener('keydown', handler);
    return () => document.removeEventListener('keydown', handler);
  }, [open]);

  if (!open) return null;

  return (
    <div className="fixed inset-0 z-50 flex items-start justify-center pt-[20vh] p-4">
      <div className="absolute inset-0 bg-black/40" onClick={() => setOpen(false)} aria-hidden />
      <div
        role="dialog"
        aria-modal="true"
        aria-label="Keyboard shortcuts"
        className="relative z-10 flex w-full max-w-sm flex-col overflow-hidden rounded-xl border border-edge-base bg-surface-base shadow-2xl"
      >
        {/* Header */}
        <div className="flex items-center gap-2.5 border-b border-edge-base px-4 py-3">
          <Keyboard size={16} className="shrink-0 text-ink-muted" />
          <h2 className="text-sm font-semibold text-ink-primary">Keyboard Shortcuts</h2>
        </div>

        {/* Shortcuts list */}
        <div className="divide-y divide-edge-base px-4">
          {SHORTCUTS.map((shortcut, idx) => (
            <ShortcutRow key={idx} shortcut={shortcut} />
          ))}
        </div>

        {/* Footer */}
        <div className="border-t border-edge-base px-4 py-3">
          <p className="text-xs text-ink-muted">
            Press <kbd className="inline-flex items-center rounded border border-edge-base bg-surface-muted px-1.5 py-0.5 text-[10px] font-medium text-ink-muted">Esc</kbd> to close
          </p>
        </div>
      </div>
    </div>
  );
}
