const remToPx = require('tailwindcss-rem-to-px');

/** @type {import('tailwindcss').Config} */
module.exports = {
  important: ".tailwind",
  corePlugins: {
    preflight: true
  },
  content: [
    "./public/*.html",
    "./src/**/*.{vue,js,ts,jsx,tsx}", // 扩展更多文件类型
    "./src/components/**/*.{vue,js,ts}", // 确保覆盖所有组件
  ],
  plugins: [
    // remToPx(),
  ],
  theme: {
    extend: {
      colors: {primary: '#1565C0', secondary: '#0277BD'},
      borderRadius: {
        'none': '0px',
        'sm': '0px',
        DEFAULT: '0px',
        'md': '0px',
        'lg': '0px',
        'xl': '0px',
        '2xl': '0px',
        '3xl': '0px',
        'full': '9999px',
        'button': '0px'
      }
    }
  }
}
