import { useState } from 'react';
import { Plus, Trash2, Bell } from 'lucide-react';
import { useWebhooks, useCreateWebhook, useDeleteWebhook } from '../api/hooks';
import type { WebhookEvent } from '../api/types';

const ALL_EVENTS: { value: WebhookEvent; label: string }[] = [
  { value: 'issue.created', label: '이슈 생성' },
  { value: 'issue.approved', label: '이슈 승인' },
  { value: 'issue.status_changed', label: '상태 변경' },
  { value: 'issue.deleted', label: '이슈 삭제' },
];

export function WebhookSettings({ boardId }: { boardId: string }) {
  const webhooks = useWebhooks(boardId);
  const createWh = useCreateWebhook();
  const deleteWh = useDeleteWebhook();
  const [open, setOpen] = useState(false);
  const [name, setName] = useState('');
  const [url, setUrl] = useState('');
  const [events, setEvents] = useState<WebhookEvent[]>(['issue.approved']);

  const submit = () => {
    if (!name.trim() || !url.trim()) return;
    createWh.mutate({ boardId, name: name.trim(), url: url.trim(), events });
    setName('');
    setUrl('');
    setEvents(['issue.approved']);
    setOpen(false);
  };

  const toggleEvent = (e: WebhookEvent) => {
    setEvents((prev) => prev.includes(e) ? prev.filter((v) => v !== e) : [...prev, e]);
  };

  return (
    <section className="flex flex-col gap-3">
      <div className="flex items-center justify-between">
        <h3 className="flex items-center gap-1.5 text-sm font-semibold text-ink-secondary">
          <Bell size={14} /> Webhooks
        </h3>
        {!open && (
          <button
            type="button"
            onClick={() => setOpen(true)}
            className="flex items-center gap-1 rounded-lg border border-dashed border-edge-base px-2.5 py-1 text-xs text-ink-muted transition-colors hover:border-brand-500/30 hover:text-brand-500"
          >
            <Plus size={12} /> 추가
          </button>
        )}
      </div>

      {/* 기존 webhook 목록 */}
      {webhooks.data?.map((wh) => (
        <div key={wh.id} className="flex items-center justify-between rounded-lg border border-edge-base bg-surface-subtle px-3 py-2">
          <div className="flex flex-col gap-0.5">
            <span className="text-sm font-medium text-ink-primary">{wh.name}</span>
            <span className="text-xs text-ink-muted truncate max-w-[300px]">{wh.url}</span>
            <div className="flex gap-1 mt-1">
              {wh.events.map((e) => (
                <span key={e} className="rounded bg-brand-500/10 px-1.5 py-0.5 text-[10px] text-brand-500">
                  {ALL_EVENTS.find((a) => a.value === e)?.label ?? e}
                </span>
              ))}
            </div>
          </div>
          <button
            type="button"
            onClick={() => deleteWh.mutate({ boardId, whId: wh.id })}
            className="rounded p-1.5 text-ink-muted transition-colors hover:bg-red-50 hover:text-red-600 dark:hover:bg-red-950/30"
            aria-label={`${wh.name} 삭제`}
          >
            <Trash2 size={14} />
          </button>
        </div>
      ))}

      {/* 추가 폼 */}
      {open && (
        <div className="flex flex-col gap-2 rounded-lg border border-edge-base bg-surface-subtle p-3">
          <input
            placeholder="이름 (예: Discord 알림)"
            value={name}
            onChange={(e) => setName(e.target.value)}
            className="h-8 rounded-md border border-edge-base bg-surface-base px-2 text-sm focus-visible:border-brand-500 focus-visible:outline-none"
          />
          <input
            placeholder="Webhook URL"
            value={url}
            onChange={(e) => setUrl(e.target.value)}
            className="h-8 rounded-md border border-edge-base bg-surface-base px-2 text-sm focus-visible:border-brand-500 focus-visible:outline-none"
          />
          <div className="flex flex-wrap gap-1.5">
            {ALL_EVENTS.map((ev) => (
              <button
                key={ev.value}
                type="button"
                onClick={() => toggleEvent(ev.value)}
                className={`rounded-full border px-2.5 py-0.5 text-xs font-medium transition-colors ${
                  events.includes(ev.value)
                    ? 'border-brand-500 bg-brand-500/15 text-brand-500'
                    : 'border-edge-base text-ink-muted hover:text-ink-secondary'
                }`}
              >
                {ev.label}
              </button>
            ))}
          </div>
          <div className="flex gap-2">
            <button onClick={submit} className="h-8 flex-1 rounded-md bg-brand-500 text-sm font-medium text-white hover:bg-brand-600 transition-colors">
              추가
            </button>
            <button onClick={() => setOpen(false)} className="h-8 rounded-md border border-edge-base px-3 text-sm text-ink-secondary hover:bg-surface-muted transition-colors">
              취소
            </button>
          </div>
        </div>
      )}
    </section>
  );
}
