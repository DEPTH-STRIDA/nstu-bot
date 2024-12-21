/** @type {import('tailwindcss').Config} */
export default {
  content: ["./index.html", "./src/**/*.{js,ts,jsx,tsx}"],
  theme: {
    extend: {
      colors: {
        primary: "#222222",
        secondary: "#FFFFFF",
      },
      fontFamily: {
        inter: ["Inter", "sans-serif"],
        roboto: ["Roboto", "sans-serif"],
      },
      boxShadow: {
        activeMenu:
          "0px 0px 8px 0 rgba(255,255,255,1), 0px -1px 6px 0 rgba(0,0,0,0.25)",
      },
    },
  },
  plugins: [],
};
