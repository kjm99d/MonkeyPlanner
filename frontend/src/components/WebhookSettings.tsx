import { useState } from 'react';
import { useTranslation } from 'react-i18next';
import { Plus, Trash2, Bell, MessageCircle, Hash, Send, Globe } from 'lucide-react';
import { useWebhooks, useCreateWebhook, useDeleteWebhook } from '../api/hooks';
import { ConfirmDialog } from './ConfirmDialog';
import { Input } from './Input';
import { Button } from './Button';
import { useToast } from './Toast';
import type { WebhookEvent } from '../api/types';

const ALL_EVENT_VALUES: WebhookEvent[] = [
  'issue.created',
  'issue.approved',
  'issue.status_changed',
  'issue.deleted',
];

type Platform = 'discord' | 'slack' | 'telegram' | 'custom';

const PLATFORMS: { value: Platform; label: string; icon: typeof Globe; color: string; hint: string }[] = [
  { value: 'discord', label: 'Discord', icon: MessageCircle, color: 'text-[#5865F2] bg-[#5865F2]/10', hint: 'https://discord.com/api/webhooks/...' },
  { value: 'slack', label: 'Slack', icon: Hash, color: 'text-[#E01E5A] bg-[#E01E5A]/10', hint: 'https://hooks.slack.com/services/...' },
  { value: 'telegram', label: 'Telegram', icon: Send, color: 'text-[#26A5E4] bg-[#26A5E4]/10', hint: 'https://api.telegram.org/bot.../sendMessage' },
  { value: 'custom', label: 'Custom', icon: Globe, color: 'text-ink-secondary bg-surface-muted', hint: 'https://your-server.com/webhook' },
];

function detectPlatform(url: string): Platform {
  if (url.includes('discord.com')) return 'discord';
  if (url.includes('slack.com')) return 'slack';
  if (url.includes('telegram.org')) return 'telegram';
  return 'custom';
}

function PlatformBadge({ url }: { url: string }) {
  const p = PLATFORMS.find((pl) => pl.value === detectPlatform(url))!;
  const Icon = p.icon;
  return (
    <span className={`flex items-center gap-1 rounded-md px-2 py-0.5 text-xs font-medium ${p.color}`}>
      <Icon size={12} /> {p.label}
    </span>
  );
}

export function WebhookSettings({ boardId }: { boardId: string }) {
  const { t } = useTranslation();
  const { toast } = useToast();
  const webhooks = useWebhooks(boardId);
  const createWh = useCreateWebhook();
  const deleteWh = useDeleteWebhook();
  const [open, setOpen] = useState(false);
  const [confirmWh, setConfirmWh] = useState<string | null>(null);
  const [platform, setPlatform] = useState<Platform>('discord');
  const [name, setName] = useState('');
  const [url, setUrl] = useState('');
  const [events, setEvents] = useState<WebhookEvent[]>(['issue.approved']);

  const eventKeyMap: Record<WebhookEvent, string> = {
    'issue.created': 'webhook.events.issue_created',
    'issue.approved': 'webhook.events.issue_approved',
    'issue.status_changed': 'webhook.events.issue_status_changed',
    'issue.deleted': 'webhook.events.issue_deleted',
  };
  const ALL_EVENTS: { value: WebhookEvent; label: string }[] = ALL_EVENT_VALUES.map((v) => ({
    value: v,
    label: t(eventKeyMap[v]),
  }));

  const submit = async () => {
    if (!name.trim() || !url.trim()) return;
    try {
      await createWh.mutateAsync({ boardId, name: name.trim(), url: url.trim(), events });
      setName('');
      setUrl('');
      setEvents(['issue.approved']);
      setOpen(false);
    } catch {
      toast('error', t('common.error'));
    }
  };

  const toggleEvent = (e: WebhookEvent) => {
    setEvents((prev) => prev.includes(e) ? prev.filter((v) => v !== e) : [...prev, e]);
  };

  return (
    <section className="flex flex-col gap-3">
      <div className="flex items-center justify-between">
        <h3 className="flex items-center gap-1.5 text-sm font-semibold text-ink-secondary">
          <Bell size={14} /> {t('webhook.title')}
        </h3>
        {!open && (
          <button
            type="button"
            onClick={() => setOpen(true)}
            className="flex items-center gap-1 rounded-lg border border-dashed border-edge-base px-2.5 py-1 text-xs text-ink-muted transition-colors hover:border-brand-500/30 hover:text-brand-500"
          >
            <Plus size={12} /> {t('webhook.add')}
          </button>
        )}
      </div>

      {/* existing webhook list */}
      {webhooks.data?.map((wh) => (
        <div key={wh.id} className="flex items-center justify-between rounded-lg border border-edge-base bg-surface-subtle px-4 py-3 shadow-sm">
          <div className="flex flex-col gap-1">
            <div className="flex items-center gap-2">
              <PlatformBadge url={wh.url} />
              <span className="text-sm font-semibold text-ink-primary">{wh.name}</span>
            </div>
            <span className="text-xs text-ink-muted truncate max-w-[400px]">{wh.url}</span>
            <div className="flex gap-1 mt-1">
              {wh.events.map((e) => (
                <span key={e} className="rounded bg-brand-500/10 px-1.5 py-0.5 text-[11px] text-brand-500">
                  {ALL_EVENTS.find((a) => a.value === e)?.label ?? e}
                </span>
              ))}
            </div>
          </div>
          <button
            type="button"
            onClick={() => setConfirmWh(wh.id)}
            className="rounded p-1.5 text-ink-muted transition-colors hover:bg-red-50 hover:text-red-600 dark:hover:bg-red-950/30"
            aria-label={t('webhook.deleteLabel', { name: wh.name })}
          >
            <Trash2 size={14} />
          </button>
        </div>
      ))}

      <ConfirmDialog
        open={confirmWh !== null}
        title={t('webhook.deleteConfirm', { name: webhooks.data?.find(w => w.id === confirmWh)?.name ?? '' })}
        description={t('webhook.deleteConfirm', { name: webhooks.data?.find(w => w.id === confirmWh)?.name ?? '' })}
        onConfirm={async () => {
          try {
            if (confirmWh) await deleteWh.mutateAsync({ boardId, whId: confirmWh });
          } catch {
            toast('error', t('common.error'));
          }
          setConfirmWh(null);
        }}
        onCancel={() => setConfirmWh(null)}
      />

      {/* add form */}
      {open && (
        <div className="flex flex-col gap-2.5 rounded-lg border border-edge-base bg-surface-subtle p-4 shadow-sm">
          {/* platform selector */}
          <div className="flex gap-2">
            {PLATFORMS.map((p) => {
              const Icon = p.icon;
              return (
                <button
                  key={p.value}
                  type="button"
                  onClick={() => { setPlatform(p.value); setName(t('webhook.notification', { name: p.label })); }}
                  className={`flex items-center gap-1.5 rounded-lg border px-3 py-1.5 text-xs font-medium transition-all ${
                    platform === p.value
                      ? `border-current ${p.color} shadow-sm`
                      : 'border-edge-base text-ink-muted hover:text-ink-secondary'
                  }`}
                >
                  <Icon size={14} /> {p.label}
                </button>
              );
            })}
          </div>
          <Input
            placeholder={t('webhook.name')}
            value={name}
            onChange={(e) => setName(e.target.value)}
          />
          <Input
            placeholder={PLATFORMS.find((p) => p.value === platform)?.hint ?? 'Webhook URL'}
            value={url}
            onChange={(e) => setUrl(e.target.value)}
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
            <Button onClick={submit} size="sm" className="flex-1">
              {t('webhook.add')}
            </Button>
            <Button onClick={() => setOpen(false)} size="sm" variant="ghost">
              {t('webhook.cancel')}
            </Button>
          </div>
        </div>
      )}
    </section>
  );
}
