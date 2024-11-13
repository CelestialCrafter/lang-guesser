import { browser } from '$app/environment';

class Customization {
	#theme = $state(localStorage.getItem('theme'));
	#font = $state(localStorage.getItem('font'));

	constructor() {
		window.addEventListener('storage', () => {
			this.#theme = localStorage.getItem('theme');
			this.#font = localStorage.getItem('font');
		});
	}

	get themeRaw() {
		return this.#theme;
	}

	get theme() {
		return this.#theme || '';
	}

	set theme(value) {
		localStorage.setItem('theme', value);
		this.#theme = value;
	}

	get fontRaw() {
		return this.#font;
	}

	get font() {
		const defaultFont = getComputedStyle(document.documentElement).fontFamily;
		return this.#font ? `${this.#font}, ${defaultFont}` : defaultFont;
	}

	set font(value) {
		localStorage.setItem('font', value);
		this.#font = value;
	}
}

export default browser && new Customization();
