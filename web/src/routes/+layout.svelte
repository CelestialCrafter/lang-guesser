<script>
	import { onMount } from 'svelte';
	import { setTheme, setFont, setDefaultFont } from '$lib';

	import '../app.css';
	import ThemeButton from './ThemeButton.svelte';
	import FontButton from './FontButton.svelte';

	let { children } = $props();

	onMount(() => {
		window.addEventListener(
			'theme',
			(event) => (document.documentElement.dataset.theme = event.detail)
		);
		window.addEventListener(
			'font',
			(event) => (document.documentElement.style.fontFamily = event.detail)
		);

		setTheme(localStorage.getItem('theme'));
		setDefaultFont(getComputedStyle(document.documentElement).fontFamily);
		setFont(localStorage.getItem('font'));
	});
</script>

<section class="p-4">
	<div class="flex flex-wrap gap-2">
		{#each ['', 'rosepine', 'rosepine-moon', 'rosepine-dawn'] as theme}
			<ThemeButton {theme} />
		{/each}
		{#each ['', 'JetBrainsMono', 'Monaspace Neon', 'FiraCode'] as font}
			<FontButton {font} />
		{/each}
	</div>
	
	

	<div class="divider"></div>

	{@render children()}
</section>
