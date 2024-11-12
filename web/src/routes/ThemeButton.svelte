<script>
	const { theme } = $props();

	const fireThemeChangeEvent = () => {
		const event = new CustomEvent('theme', { detail: localStorage.getItem('theme') });
		window.dispatchEvent(event);
	};
</script>

<button class="btn" onclick={fireThemeChangeEvent} data-set-theme={theme}
	>{theme || 'default'}</button
>
