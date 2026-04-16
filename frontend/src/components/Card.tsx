import type { HTMLAttributes, ReactNode } from 'react';

type Props = HTMLAttributes<HTMLDivElement> & {
  children: ReactNode;
};

export function Card({ className = '', children, ...rest }: Props) {
  return (
    <div
      {...rest}
      className={`rounded-lg border border-edge-base bg-surface-subtle p-4 shadow-sm transition-shadow hover:shadow-md ${className}`}
    >
      {children}
    </div>
  );
}
