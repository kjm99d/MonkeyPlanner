import typography from '@tailwindcss/typography';
import {
  lightColors,
  fontFamily,
  fontSize,
  spacing,
  radius,
  shadow,
} from './src/design/tokens.ts';

/** @type {import('tailwindcss').Config} */
export default {
  darkMode: 'class',
  content: ['./index.html', './src/**/*.{ts,tsx}'],
  theme: {
    extend: {
      colors: {
        brand: {
          50: '#f0fdfa',
          100: '#ccfbf1',
          500: 'var(--brand-500)',
          600: 'var(--brand-600)',
          900: '#134e4a',
        },
        accent: {
          DEFAULT: 'var(--accent)',
          ring: lightColors.accent.ring,
        },
        status: {
          pending: 'var(--status-pending)',
          approved: 'var(--status-approved)',
          inProgress: 'var(--status-inProgress)',
          done: 'var(--status-done)',
        },
        surface: {
          base: 'var(--bg-base)',
          subtle: 'var(--bg-subtle)',
          muted: 'var(--bg-muted)',
        },
        ink: {
          primary: 'var(--text-primary)',
          secondary: 'var(--text-secondary)',
          muted: 'var(--text-muted)',
        },
        edge: {
          base: 'var(--border-base)',
        },
      },
      fontFamily: {
        sans: fontFamily.sans,
        mono: fontFamily.mono,
      },
      fontSize,
      spacing,
      borderRadius: radius,
      boxShadow: shadow,
    },
  },
  plugins: [typography],
};
