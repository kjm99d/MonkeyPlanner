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
          50: '#F0FDFA',
          500: 'var(--brand-500)',
          600: 'var(--brand-600)',
          900: '#134E4A',
        },
        accent: {
          DEFAULT: 'var(--accent)',
        },
        status: lightColors.status,
        accent: lightColors.accent,
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
