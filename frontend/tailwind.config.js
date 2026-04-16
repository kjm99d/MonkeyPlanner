/** @type {import('tailwindcss').Config} */
export default {
  darkMode: 'class',
  content: ['./index.html', './src/**/*.{ts,tsx}'],
  theme: {
    extend: {
      colors: {
        // tokens는 M2에서 src/design/tokens.ts로 이동 예정
        brand: {
          50: '#fff7ed',
          500: '#f97316',
          600: '#ea580c',
          900: '#7c2d12',
        },
      },
    },
  },
  plugins: [],
};
