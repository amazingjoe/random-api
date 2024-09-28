/** @type {import('tailwindcss').Config} */
module.exports = {
  content: ["view/*.tsx", "view/**/*.tsx"],
  theme: {
    extend: {},
  },
  plugins: [require("daisyui")],
};
