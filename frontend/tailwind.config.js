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
        brand: lightColors.brand,
        status: lightColors.status,
        accent: lightColors.accent,
        surface: {
          base: lightColors.bg.base,
          subtle: lightColors.bg.subtle,
          muted: lightColors.bg.muted,
        },
        ink: {
          primary: lightColors.text.primary,
          secondary: lightColors.text.secondary,
          muted: lightColors.text.muted,
        },
        edge: {
          base: lightColors.border.base,
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
  plugins: [],
};
