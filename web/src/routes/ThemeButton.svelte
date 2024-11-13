<script>
	import { setTheme } from '$lib';

	const { theme } = $props();
</script>

<button
	class="btn"
	onclick={(event) => setTheme(event.target.dataset.setTheme)}
	data-set-theme={theme}
>
	{theme || 'default'}
</button>
