// 몽키 플래너 디자인 토큰
// 원본 정의: .omc/design/DESIGN.md
// Tailwind config 와 컴포넌트가 공통으로 import 합니다.

export const lightColors = {
  bg: {
    base: '#ffffff',
    subtle: '#f8fafc',
    muted: '#f1f5f9',
  },
  border: {
    base: '#e2e8f0',
  },
  text: {
    primary: '#0f172a',
    secondary: '#475569',
    muted: '#94a3b8',
  },
  brand: {
    50: '#f0fdfa',
    500: '#0d9488',
    600: '#0f766e',
    900: '#134e4a',
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
    base: '#0f172a',
    subtle: '#1e293b',
    muted: '#334155',
  },
  border: {
    base: '#475569',
  },
  text: {
    primary: '#f8fafc',
    secondary: '#cbd5e1',
    muted: '#64748b',
  },
  brand: {
    50: '#042f2e',
    500: '#14b8a6',
    600: '#0d9488',
    900: '#ccfbf1',
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
    '"Plus Jakarta Sans"',
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
  focus: '0 0 0 4px rgb(13 148 136 / 0.35)',
} as const;

/** 상태별 뱃지 색상을 한 자리에서 쓰기 위한 맵. */
export const statusColor = {
  Pending: lightColors.status.pending,
  Approved: lightColors.status.approved,
  InProgress: lightColors.status.inProgress,
  Done: lightColors.status.done,
} as const;

export type StatusKey = keyof typeof statusColor;
