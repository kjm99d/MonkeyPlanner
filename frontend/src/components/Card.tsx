import type { HTMLAttributes, ReactNode } from 'react';

type Props = HTMLAttributes<HTMLDivElement> & {
  children: ReactNode;
};

export function Card({ className = '', children, ...rest }: Props) {
  return (
    <div
      {...rest}
      className={`rounded-xl border border-edge-base bg-surface-subtle p-4 shadow-sm transition-all duration-200 hover:shadow-md hover:border-brand-500/40 ${className}`}
    >
      {children}
    </div>
  );
}
