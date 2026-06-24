import type { Config } from 'tailwindcss';

export default {
  content: ['./index.html', './src/**/*.{ts,tsx}'],
  theme: {
    extend: {
      colors: {
        primary: '#644599',
        secondary: '#495599',
        accent: '#29DBF2',
        bg: '#0F1020',
        surface: '#1C1F3A',
        surface2: '#262A4D',
        muted: '#8A8FB5',
        'text-accent': '#80DEEA',
        border: 'rgba(255,255,255,0.06)',
        positive: '#3DDC97',
        negative: '#FF6B7A',
        donate: '#FF8FB3',
        buy: '#FFB347',
      },
      fontFamily: {
        sans: ['Manrope', 'Helvetica Neue', 'system-ui', 'sans-serif'],
        display: ['Space Grotesk', 'Manrope', 'system-ui', 'sans-serif'],
        mono: ['JetBrains Mono', 'ui-monospace', 'SF Mono', 'Menlo', 'monospace'],
      },
      borderRadius: {
        '4xl': '2rem',
      },
      animation: {
        'pulse-dot': 'pulse 1.2s infinite',
        'spin-slow': 'spin 0.8s linear infinite',
        'pop': 'pop 0.3s cubic-bezier(.2,.8,.3,1.4)',
        'sheet-in': 'sheet-in 0.22s cubic-bezier(.2,.7,.3,1)',
        'toast-in': 'toast-in 0.22s cubic-bezier(.2,.7,.3,1)',
      },
      keyframes: {
        pop: {
          from: { transform: 'scale(0.5)', opacity: '0' },
          to: { transform: 'scale(1)', opacity: '1' },
        },
        'sheet-in': {
          from: { transform: 'translateY(100%)' },
          to: { transform: 'translateY(0)' },
        },
        'toast-in': {
          from: { transform: 'translateY(20px)', opacity: '0' },
          to: { transform: 'translateY(0)', opacity: '1' },
        },
      },
    },
  },
} satisfies Config;
