/** @type {import('tailwindcss').Config} */
module.exports = {
	content: ['./src/**/*.{html,js,svelte,ts}'],
	theme: {
		extend: {
			colors: {
				primary: '#c7f860',
				secondary: '#ff814b',
				bg: '#f8f8f7',
				'bg-dimmed': '#eeeeec',
				'text-dimmed': '#5f5f5f',
				'text-dimmed-dark': '#252525'
			}
		}
	},
	plugins: []
};
