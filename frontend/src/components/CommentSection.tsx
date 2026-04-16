import { useState } from 'react';
import { useTranslation } from 'react-i18next';
import { Trash2, MessageCircle, Send } from 'lucide-react';
import { useComments, useCreateComment, useDeleteComment } from '../api/hooks';
import { Button } from './Button';
import { Skeleton } from './Skeleton';

function relativeTime(dateStr: string): string {
  const now = Date.now();
  const then = new Date(dateStr).getTime();
  const diff = Math.floor((now - then) / 1000);
  if (diff < 60) return 'just now';
  if (diff < 3600) return `${Math.floor(diff / 60)}m ago`;
  if (diff < 86400) return `${Math.floor(diff / 3600)}h ago`;
  if (diff < 604800) return `${Math.floor(diff / 86400)}d ago`;
  return new Date(dateStr).toLocaleDateString();
}

interface Props {
  issueId: string;
}

export function CommentSection({ issueId }: Props) {
  const { t } = useTranslation();
  const { data: comments, isLoading } = useComments(issueId);
  const createComment = useCreateComment();
  const deleteComment = useDeleteComment();
  const [body, setBody] = useState('');

  async function handleSubmit(e: React.FormEvent) {
    e.preventDefault();
    const trimmed = body.trim();
    if (!trimmed) return;
    await createComment.mutateAsync({ issueId, body: trimmed });
    setBody('');
  }

  return (
    <section aria-label={t('comments.section')} className="flex flex-col gap-4">
      <h2 className="text-base font-semibold text-ink-primary">{t('comments.title')}</h2>

      {isLoading && (
        <div className="flex flex-col gap-3">
          <Skeleton className="h-16 w-full" />
          <Skeleton className="h-16 w-full" />
        </div>
      )}

      {!isLoading && comments && comments.length === 0 && (
        <div className="flex flex-col items-center gap-2 rounded-lg border border-edge-base bg-surface-subtle py-8 text-ink-muted">
          <MessageCircle className="h-8 w-8 opacity-40" />
          <span className="text-sm">{t('comments.empty')}</span>
        </div>
      )}

      {!isLoading && comments && comments.length > 0 && (
        <ul className="flex flex-col gap-2">
          {comments.map((comment) => (
            <li
              key={comment.id}
              className="group flex flex-col gap-1 rounded-lg border border-edge-base bg-surface-subtle px-4 py-3"
            >
              <p className="whitespace-pre-wrap text-sm text-ink-primary">{comment.body}</p>
              <div className="flex items-center justify-between">
                <span className="text-xs text-ink-muted">{relativeTime(comment.createdAt)}</span>
                <button
                  type="button"
                  aria-label={t('comments.delete')}
                  onClick={() => deleteComment.mutate({ issueId, commentId: comment.id })}
                  className="rounded p-1 text-ink-muted opacity-0 transition-opacity hover:text-red-500 group-hover:opacity-100 focus-visible:opacity-100 focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-brand-500"
                >
                  <Trash2 className="h-3.5 w-3.5" />
                </button>
              </div>
            </li>
          ))}
        </ul>
      )}

      <form onSubmit={handleSubmit} className="flex flex-col gap-2">
        <textarea
          value={body}
          onChange={(e) => setBody(e.target.value)}
          placeholder={t('comments.placeholder')}
          rows={3}
          className="min-h-[80px] w-full resize-y rounded-lg border border-edge-base bg-surface-subtle px-3 py-2 text-sm text-ink-primary placeholder:text-ink-muted focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-brand-500"
        />
        <div className="flex justify-end">
          <Button
            type="submit"
            size="sm"
            disabled={!body.trim() || createComment.isPending}
          >
            <Send className="h-3.5 w-3.5" />
            {createComment.isPending ? t('comments.submitting') : t('comments.submit')}
          </Button>
        </div>
      </form>
    </section>
  );
}
