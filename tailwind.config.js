/** @type {import('tailwindcss').Config} */
export default {
  content: ["./src/**/*.{html,js,svelte,ts}"],
  theme: {
    extend: {
      colors: {
        ess: {
          400: "#E88419",
          500: "#ff931e",
          600: "#de7a0f",
        },
      },
      spacing: {
        logo: "6rem",
      },
      fontFamily: {
        sans: ["Roboto", "sans-serif"],
      },
    },
  },
  plugins: [],
};
