import type { Config } from "tailwindcss";

export default {
	content: ["./app/**/*.{ts,tsx}", "./components/**/*.{ts,tsx}"],
	theme: {
		extend: {},
	},
	plugins: [require("@tailwindcss/typography")],
} satisfies Config;
