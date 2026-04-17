import type { SVGProps } from 'react';

type Props = SVGProps<SVGSVGElement> & { size?: number };

/**
 * MonkeyLogo — minimal monkey face mark for the MonkeyPlanner brand.
 * Uses currentColor so it inherits the surrounding text color.
 */
export function MonkeyLogo({ size = 16, ...rest }: Props) {
  return (
    <svg
      xmlns="http://www.w3.org/2000/svg"
      width={size}
      height={size}
      viewBox="0 0 24 24"
      fill="none"
      stroke="currentColor"
      strokeWidth={2}
      strokeLinecap="round"
      strokeLinejoin="round"
      aria-hidden="true"
      {...rest}
    >
      {/* ears */}
      <circle cx="5" cy="9" r="2.5" />
      <circle cx="19" cy="9" r="2.5" />
      {/* head */}
      <path d="M18 12a6 6 0 0 1-12 0 6 6 0 0 1 12 0z" />
      {/* face mask */}
      <path d="M8.5 13.5c.7 1 1.8 1.8 3.5 1.8s2.8-.8 3.5-1.8" />
      {/* eyes */}
      <circle cx="9.5" cy="11" r="0.6" fill="currentColor" />
      <circle cx="14.5" cy="11" r="0.6" fill="currentColor" />
    </svg>
  );
}
