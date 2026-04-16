import { useState, useEffect, useRef } from 'react';
import { useTranslation } from 'react-i18next';
import { X, Plus, FileText, Trash2 } from 'lucide-react';
import { Button } from './Button';

type Template = {
  id: string;
  name: string;
  title: string;
  body: string;
  instructions: string;
};

type Props = {
  boardId: string;
  open: boolean;
  onClose: () => void;
  onSelect: (template: Template) => void;
};

function getTemplates(boardId: string): Template[] {
  try {
    return JSON.parse(localStorage.getItem(`mp-templates-${boardId}`) ?? '[]');
  } catch { return []; }
}

function saveTemplates(boardId: string, templates: Template[]) {
  localStorage.setItem(`mp-templates-${boardId}`, JSON.stringify(templates));
}

export function TemplateDialog({ boardId, open, onClose, onSelect }: Props) {
  const { t } = useTranslation();
  const [templates, setTemplates] = useState<Template[]>([]);
  const [showAdd, setShowAdd] = useState(false);
  const [newName, setNewName] = useState('');
  const [newTitle, setNewTitle] = useState('');
  const [newBody, setNewBody] = useState('');
  const [newInstructions, setNewInstructions] = useState('');
  const closeRef = useRef<HTMLButtonElement>(null);

  useEffect(() => {
    if (!open) return;
    setTemplates(getTemplates(boardId));
    setShowAdd(false);
    setNewName('');
    setNewTitle('');
    setNewBody('');
    setNewInstructions('');
    closeRef.current?.focus();
    const handler = (e: KeyboardEvent) => {
      if (e.key === 'Escape') onClose();
    };
    document.addEventListener('keydown', handler);
    return () => document.removeEventListener('keydown', handler);
  }, [open, boardId, onClose]);

  if (!open) return null;

  function handleDelete(id: string) {
    const updated = templates.filter((t) => t.id !== id);
    setTemplates(updated);
    saveTemplates(boardId, updated);
  }

  function handleAdd() {
    if (!newName.trim()) return;
    const tmpl: Template = {
      id: crypto.randomUUID(),
      name: newName.trim(),
      title: newTitle.trim(),
      body: newBody.trim(),
      instructions: newInstructions.trim(),
    };
    const updated = [...templates, tmpl];
    setTemplates(updated);
    saveTemplates(boardId, updated);
    setShowAdd(false);
    setNewName('');
    setNewTitle('');
    setNewBody('');
    setNewInstructions('');
  }

  return (
    <div className="fixed inset-0 z-50 flex items-center justify-center p-4">
      <div className="absolute inset-0 bg-black/40" onClick={onClose} aria-hidden />
      <div
        role="dialog"
        aria-modal="true"
        aria-labelledby="template-dialog-title"
        className="relative z-10 flex w-full max-w-md flex-col gap-4 rounded-xl border border-edge-base bg-surface-base p-6 shadow-lg animate-in"
      >
        {/* Header */}
        <div className="flex items-center justify-between">
          <div className="flex items-center gap-2">
            <FileText size={16} className="text-brand-500" />
            <h2 id="template-dialog-title" className="text-sm font-semibold text-ink-primary">
              {t('template.title')}
            </h2>
          </div>
          <button
            ref={closeRef}
            type="button"
            onClick={onClose}
            aria-label={t('common.close')}
            className="rounded-md p-1 text-ink-muted transition-colors hover:bg-surface-muted hover:text-ink-primary"
          >
            <X size={16} />
          </button>
        </div>

        {/* Template list */}
        <div className="flex flex-col gap-2">
          {templates.length === 0 && !showAdd && (
            <p className="py-4 text-center text-sm text-ink-muted">
              {t('template.empty')}
            </p>
          )}
          {templates.map((tmpl) => (
            <div
              key={tmpl.id}
              className="group flex items-center justify-between rounded-lg border border-edge-base bg-surface-subtle px-3 py-2 transition-colors hover:border-brand-500/50 hover:bg-surface-muted"
            >
              <button
                type="button"
                className="flex-1 text-left"
                onClick={() => { onSelect(tmpl); onClose(); }}
              >
                <span className="block text-sm font-medium text-ink-primary">{tmpl.name}</span>
                {tmpl.title && (
                  <span className="block truncate text-xs text-ink-muted">{tmpl.title}</span>
                )}
              </button>
              <button
                type="button"
                onClick={() => handleDelete(tmpl.id)}
                aria-label={t('template.delete')}
                className="ml-2 shrink-0 rounded p-1 text-ink-muted opacity-0 transition-all group-hover:opacity-100 hover:text-red-500"
              >
                <Trash2 size={14} />
              </button>
            </div>
          ))}
        </div>

        {/* Add new template form */}
        {showAdd ? (
          <div className="flex flex-col gap-2 rounded-lg border border-edge-base bg-surface-subtle p-3">
            <input
              autoFocus
              type="text"
              value={newName}
              onChange={(e) => setNewName(e.target.value)}
              placeholder={t('template.name')}
              className="w-full rounded-md border border-edge-base bg-surface-base px-2.5 py-1.5 text-sm text-ink-primary placeholder:text-ink-muted focus:border-brand-500 focus:outline-none"
            />
            <input
              type="text"
              value={newTitle}
              onChange={(e) => setNewTitle(e.target.value)}
              placeholder={t('board.newIssueTitle')}
              className="w-full rounded-md border border-edge-base bg-surface-base px-2.5 py-1.5 text-sm text-ink-primary placeholder:text-ink-muted focus:border-brand-500 focus:outline-none"
            />
            <textarea
              value={newBody}
              onChange={(e) => setNewBody(e.target.value)}
              placeholder={t('issue.body')}
              rows={2}
              className="w-full resize-none rounded-md border border-edge-base bg-surface-base px-2.5 py-1.5 text-sm text-ink-primary placeholder:text-ink-muted focus:border-brand-500 focus:outline-none"
            />
            <textarea
              value={newInstructions}
              onChange={(e) => setNewInstructions(e.target.value)}
              placeholder={t('issue.instructionsPlaceholder')}
              rows={2}
              className="w-full resize-none rounded-md border border-edge-base bg-surface-base px-2.5 py-1.5 text-sm text-ink-primary placeholder:text-ink-muted focus:border-brand-500 focus:outline-none"
            />
            <div className="flex justify-end gap-2">
              <Button variant="ghost" size="sm" onClick={() => setShowAdd(false)}>
                {t('webhook.cancel')}
              </Button>
              <Button size="sm" disabled={!newName.trim()} onClick={handleAdd}>
                {t('property.add')}
              </Button>
            </div>
          </div>
        ) : (
          <button
            type="button"
            onClick={() => setShowAdd(true)}
            className="flex items-center gap-1.5 rounded-lg border border-dashed border-edge-base px-3 py-2 text-sm text-ink-muted transition-colors hover:border-brand-500/50 hover:text-brand-500"
          >
            <Plus size={14} />
            {t('template.save')}
          </button>
        )}
      </div>
    </div>
  );
}
