import type { InputHTMLAttributes } from 'react';

type Props = InputHTMLAttributes<HTMLInputElement> & {
  label?: string;
  error?: string;
};

export function Input({ label, error, id, className = '', ...rest }: Props) {
  const inputId = id ?? `input-${Math.random().toString(36).slice(2, 9)}`;
  return (
    <div className="flex flex-col gap-1">
      {label && (
        <label htmlFor={inputId} className="text-sm font-medium text-ink-secondary">
          {label}
        </label>
      )}
      <input
        id={inputId}
        aria-invalid={!!error}
        aria-describedby={error ? `${inputId}-err` : undefined}
        {...rest}
        className={`h-10 rounded-md border border-edge-base bg-surface-base px-3 text-ink-primary transition-colors focus-visible:border-brand-500 focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-brand-500/40 ${className}`}
      />
      {error && (
        <span id={`${inputId}-err`} className="text-xs text-red-600">
          {error}
        </span>
      )}
    </div>
  );
}
