// 몽키 플래너 디자인 토큰
// 원본 정의: .omc/design/DESIGN.md
// Tailwind config 와 컴포넌트가 공통으로 import 합니다.

export const lightColors = {
  bg: {
    base: '#ffffff',
    subtle: '#fafafa',
    muted: '#f4f4f5',
  },
  border: {
    base: '#e4e4e7',
  },
  text: {
    primary: '#18181b',
    secondary: '#52525b',
    muted: '#a1a1aa',
  },
  brand: {
    50: '#fff7ed',
    500: '#f97316',
    600: '#ea580c',
    900: '#7c2d12',
  },
  status: {
    pending: '#a1a1aa',
    approved: '#16a34a',
    inProgress: '#2563eb',
    done: '#8b5cf6',
  },
  accent: {
    ring: '#fbbf24',
  },
} as const;

export const darkColors = {
  bg: {
    base: '#0a0a0a',
    subtle: '#18181b',
    muted: '#27272a',
  },
  border: {
    base: '#3f3f46',
  },
  text: {
    primary: '#fafafa',
    secondary: '#a1a1aa',
    muted: '#71717a',
  },
  brand: {
    50: '#431407',
    500: '#fb923c',
    600: '#f97316',
    900: '#fed7aa',
  },
  status: {
    pending: '#71717a',
    approved: '#22c55e',
    inProgress: '#3b82f6',
    done: '#a78bfa',
  },
  accent: {
    ring: '#facc15',
  },
} as const;

export const fontFamily = {
  sans: [
    '"Pretendard Variable"',
    'Pretendard',
    '-apple-system',
    'BlinkMacSystemFont',
    '"Segoe UI"',
    'Roboto',
    '"Noto Sans KR"',
    'sans-serif',
  ],
  mono: [
    '"JetBrains Mono"',
    '"D2Coding"',
    'ui-monospace',
    'SFMono-Regular',
    'Menlo',
    'Consolas',
    'monospace',
  ],
} as const;

export const fontSize = {
  xs: ['0.75rem', { lineHeight: '1rem' }],
  sm: ['0.875rem', { lineHeight: '1.25rem' }],
  base: ['1rem', { lineHeight: '1.5rem' }],
  lg: ['1.125rem', { lineHeight: '1.5rem' }],
  xl: ['1.25rem', { lineHeight: '1.4' }],
  '2xl': ['1.5rem', { lineHeight: '1.3' }],
  '3xl': ['1.875rem', { lineHeight: '1.2' }],
  '4xl': ['2.25rem', { lineHeight: '1.1' }],
} as const;

export const spacing = {
  0: '0',
  1: '0.25rem',
  2: '0.5rem',
  3: '0.75rem',
  4: '1rem',
  6: '1.5rem',
  8: '2rem',
  12: '3rem',
  16: '4rem',
  24: '6rem',
} as const;

export const radius = {
  none: '0',
  sm: '0.25rem',
  md: '0.5rem',
  lg: '0.75rem',
  xl: '1rem',
  full: '9999px',
} as const;

export const shadow = {
  sm: '0 1px 2px 0 rgb(0 0 0 / 0.05)',
  md: '0 4px 6px -1px rgb(0 0 0 / 0.1), 0 2px 4px -2px rgb(0 0 0 / 0.1)',
  lg: '0 10px 15px -3px rgb(0 0 0 / 0.1), 0 4px 6px -4px rgb(0 0 0 / 0.1)',
  focus: '0 0 0 4px rgb(249 115 22 / 0.35)',
} as const;

/** 상태별 뱃지 색상을 한 자리에서 쓰기 위한 맵. */
export const statusColor = {
  Pending: lightColors.status.pending,
  Approved: lightColors.status.approved,
  InProgress: lightColors.status.inProgress,
  Done: lightColors.status.done,
} as const;

export type StatusKey = keyof typeof statusColor;
