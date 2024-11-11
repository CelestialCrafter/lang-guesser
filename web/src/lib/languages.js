import { siRust, siPython, siGo } from 'simple-icons';
import { rust, python, go } from 'svelte-highlight/languages';

export const languageToHighlight = {
	rust: rust,
	python: python,
	go: go,
	'': {
		name: "none",
		register: () => ({}),
	}
};

export const languageToIcon = {
	rust: siRust,
	python: siPython,
	go: siGo
};

