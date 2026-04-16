import { useEffect, useRef, type ReactNode } from 'react';

type MenuItem = {
  label: string;
  icon?: ReactNode;
  onClick: () => void;
  danger?: boolean;
  divider?: boolean;
};

type Props = {
  items: MenuItem[];
  position: { x: number; y: number } | null;
  onClose: () => void;
};

export function ContextMenu({ items, position, onClose }: Props) {
  const ref = useRef<HTMLDivElement>(null);

  useEffect(() => {
    if (!position) return;
    const handler = (e: MouseEvent) => {
      if (ref.current && !ref.current.contains(e.target as Node)) onClose();
    };
    const keyHandler = (e: KeyboardEvent) => {
      if (e.key === 'Escape') onClose();
    };
    document.addEventListener('mousedown', handler);
    document.addEventListener('keydown', keyHandler);
    return () => {
      document.removeEventListener('mousedown', handler);
      document.removeEventListener('keydown', keyHandler);
    };
  }, [position, onClose]);

  if (!position) return null;

  return (
    <div
      ref={ref}
      className="fixed z-50 min-w-[160px] rounded-lg border border-edge-base bg-surface-base py-1 shadow-lg animate-in"
      style={{ left: position.x, top: position.y }}
    >
      {items.map((item, i) => (
        item.divider ? (
          <hr key={i} className="my-1 border-edge-base" />
        ) : (
          <button
            key={i}
            type="button"
            onClick={() => { item.onClick(); onClose(); }}
            className={`flex w-full items-center gap-2 px-3 py-1.5 text-sm transition-colors ${
              item.danger
                ? 'text-red-600 hover:bg-red-50 dark:hover:bg-red-950/30'
                : 'text-ink-primary hover:bg-surface-muted'
            }`}
          >
            {item.icon}
            {item.label}
          </button>
        )
      ))}
    </div>
  );
}
