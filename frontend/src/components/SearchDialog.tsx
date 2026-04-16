import { useState, useEffect, useRef, useMemo } from 'react';
import { useNavigate } from 'react-router-dom';
import { Search, FileText, LayoutDashboard } from 'lucide-react';
import { useBoards, useIssues } from '../api/hooks';

export function SearchDialog() {
  const [open, setOpen] = useState(false);
  const [query, setQuery] = useState('');
  const inputRef = useRef<HTMLInputElement>(null);
  const navigate = useNavigate();
  const boards = useBoards();
  const issues = useIssues({});

  // Cmd+K / Ctrl+K to open
  useEffect(() => {
    const handler = (e: KeyboardEvent) => {
      if ((e.metaKey || e.ctrlKey) && e.key === 'k') {
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

  // Auto-focus input when opened
  useEffect(() => {
    if (open) {
      setQuery('');
      setTimeout(() => inputRef.current?.focus(), 50);
    }
  }, [open]);

  // Filter results
  const results = useMemo(() => {
    if (!query.trim()) return { boards: [], issues: [] };
    const q = query.toLowerCase();
    return {
      boards: (boards.data ?? []).filter(b => b.name.toLowerCase().includes(q)).slice(0, 3),
      issues: (issues.data ?? []).filter(i => i.title.toLowerCase().includes(q)).slice(0, 8),
    };
  }, [query, boards.data, issues.data]);

  const hasResults = results.boards.length > 0 || results.issues.length > 0;

  function go(path: string) {
    navigate(path);
    setOpen(false);
  }

  if (!open) return null;

  return (
    <div className="fixed inset-0 z-50 flex items-start justify-center pt-[20vh] p-4">
      <div className="absolute inset-0 bg-black/40" onClick={() => setOpen(false)} aria-hidden />
      <div
        role="dialog"
        aria-modal="true"
        aria-label="Search"
        className="relative z-10 flex w-full max-w-lg flex-col overflow-hidden rounded-xl border border-edge-base bg-surface-base shadow-2xl animate-in"
      >
        {/* Search input */}
        <div className="flex items-center gap-3 border-b border-edge-base px-4 py-3">
          <Search size={18} className="shrink-0 text-ink-muted" />
          <input
            ref={inputRef}
            value={query}
            onChange={e => setQuery(e.target.value)}
            placeholder="Search issues and boards..."
            className="flex-1 bg-transparent text-sm text-ink-primary placeholder:text-ink-muted focus:outline-none"
          />
          <kbd className="hidden sm:inline-flex items-center rounded border border-edge-base bg-surface-muted px-1.5 py-0.5 text-[10px] font-medium text-ink-muted">
            ESC
          </kbd>
        </div>

        {/* Results */}
        {query.trim() && (
          <div className="max-h-[300px] overflow-y-auto p-2">
            {!hasResults && (
              <p className="py-6 text-center text-sm text-ink-muted">No results found</p>
            )}
            {results.boards.length > 0 && (
              <div className="mb-2">
                <p className="px-2 py-1 text-[11px] font-semibold uppercase tracking-wider text-ink-muted">Boards</p>
                {results.boards.map(b => (
                  <button
                    key={b.id}
                    onClick={() => go(`/boards/${b.id}`)}
                    className="flex w-full items-center gap-2 rounded-md px-2 py-2 text-sm text-ink-primary hover:bg-surface-muted transition-colors"
                  >
                    <LayoutDashboard size={14} className="shrink-0 text-ink-muted" />
                    {b.name}
                  </button>
                ))}
              </div>
            )}
            {results.issues.length > 0 && (
              <div>
                <p className="px-2 py-1 text-[11px] font-semibold uppercase tracking-wider text-ink-muted">Issues</p>
                {results.issues.map(i => (
                  <button
                    key={i.id}
                    onClick={() => go(`/issues/${i.id}`)}
                    className="flex w-full items-center gap-2 rounded-md px-2 py-2 text-sm text-ink-primary hover:bg-surface-muted transition-colors"
                  >
                    <FileText size={14} className="shrink-0 text-ink-muted" />
                    <span className="truncate">{i.title}</span>
                  </button>
                ))}
              </div>
            )}
          </div>
        )}

        {/* Footer */}
        {!query.trim() && (
          <div className="border-t border-edge-base px-4 py-3">
            <p className="text-xs text-ink-muted">Type to search issues and boards</p>
          </div>
        )}
      </div>
    </div>
  );
}
