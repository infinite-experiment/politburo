import type { Config } from 'tailwindcss'

export default {
  content: [
    'vizburo/ui/templates/**/*.html',
    'vizburo/ui/templates/**/**/*.html',
  ],
  theme: {
    extend: {
      colors: {
        // Vizburo brand colors
        'flight': {
          'primary': '#0f766e',      // Teal
          'secondary': '#1f2937',    // Dark gray
          'accent': '#14b8a6',       // Light teal
          'error': '#dc2626',        // Red
          'warning': '#f59e0b',      // Amber
          'success': '#10b981',      // Green
        },
      },
      spacing: {
        'sidebar': '16rem',          // Sidebar width
      },
      animation: {
        'slide-in': 'slideIn 0.3s ease-in-out',
      },
      keyframes: {
        slideIn: {
          'from': { transform: 'translateX(-100%)' },
          'to': { transform: 'translateX(0)' },
        },
      },
    },
  },
  safelist: [
    // Common utilities
    'flex', 'grid', 'block', 'hidden', 'inline-block',
    'px-4', 'py-3', 'py-2', 'px-6', 'py-3',
    'text-gray-600', 'text-gray-400', 'text-gray-900', 'text-white',
    'bg-gray-100', 'bg-gray-700', 'bg-gray-800', 'bg-white',
    'dark:text-gray-400', 'dark:bg-gray-700', 'dark:bg-gray-800', 'dark:bg-gray-900',
    'border-b', 'border-b-2', 'border-gray-200', 'border-gray-700',
    'dark:border-gray-700',
    'flight-primary', 'flight-secondary',
    'bg-flight-primary', 'text-flight-primary', 'border-flight-primary',
    'hover:bg-white', 'hover:opacity-80', 'hover:text-flight-primary',
    'dark:hover:bg-gray-600',
    'transition', 'transition-all', 'transition-colors', 'transition-opacity',
    'rounded', 'rounded-lg', 'rounded-md', 'rounded-full',
    'shadow-lg', 'opacity-50',
    'active:scale-95',
  ],
  plugins: [],
} satisfies Config
