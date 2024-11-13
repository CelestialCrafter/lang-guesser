export const formatDuration = (ns) => {
	const s = ns / 1e9;
	if (s > 60) {
		return (s / 60).toFixed(2) + 'm';
	}

	return s.toFixed(2) + 's';
};

export const setTheme = (theme) => {
	theme = theme || '';
	localStorage.setItem('theme', theme);
	window.dispatchEvent(new CustomEvent('theme', { detail: theme }));
};

let defaultFont = '';
export const setDefaultFont = (font) => (defaultFont = font);

export const setFont = (font) => {
	localStorage.setItem('font', font || '');
	font = font ? `${font}, ${defaultFont}` : defaultFont;
	window.dispatchEvent(new CustomEvent('font', { detail: font }));
};
