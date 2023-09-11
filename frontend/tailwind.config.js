/** @type {import('tailwindcss').Config} */
module.exports = {
  content: ["../app/public/*.{html,js}"],
  theme: {
    extend: {},
  },
  plugins: [require("daisyui")],
}

