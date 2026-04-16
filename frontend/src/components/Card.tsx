import type { HTMLAttributes, ReactNode } from 'react';

type Props = HTMLAttributes<HTMLDivElement> & {
  children: ReactNode;
};

export function Card({ className = '', children, ...rest }: Props) {
  return (
    <div
      {...rest}
      className={`rounded-xl border border-[var(--glass-border)] bg-[var(--glass-bg)] p-4 shadow-sm backdrop-blur-sm transition-all duration-200 hover:shadow-md hover:border-brand-500/30 ${className}`}
    >
      {children}
    </div>
  );
}
