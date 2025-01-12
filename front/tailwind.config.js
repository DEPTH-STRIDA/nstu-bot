/** @type {import('tailwindcss').Config} */
export default {
  content: ["./index.html", "./src/**/*.{js,ts,jsx,tsx}"],
  theme: {
    extend: {
      colors: {
        primary: "#1E1E1E",
        secondary: "#FFFFFF",
        inputBg: "#F2F2F2",
        weekDayBg:"#E2E2E2",
      },
      boxShadow: {
        activeMenu: "0px 0px 8px rgba(255, 255, 255, 0.5)",
      },
      keyframes: {
        fadeIn: {
          "0%": { opacity: "0" },
          "100%": { opacity: "1" },
        },
        popIn: {
          "0%": {
            transform: "translate(-50%, -45%) scale(0.95)",
            opacity: "0",
          },
          "100%": { transform: "translate(-50%, -50%) scale(1)", opacity: "1" },
        },
      },
      animation: {
        fadeIn: "fadeIn 0.2s ease-out",
        popIn: "popIn 0.3s ease-out forwards",
      },
    },
    fontFamily: {
      roboto: ["Roboto", "sans-serif"],
    },
  },
  plugins: [],
};
