/** @type {import('tailwindcss').Config} */
export default {
  content: [
    "./index.html",
    "./src/**/*.{js,ts,jsx,tsx}",
  ],
  darkMode: 'class', // Enable dark mode with class strategy
  theme: {
    extend: {
      width: {
        '65': '16.25rem',
      },
      spacing: {
        '65': '16.25rem',
      }
    },
  },
  plugins: [],
} 