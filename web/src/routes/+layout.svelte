<script>
	import customization from '$lib/customization.svelte.js';
	import Navigation from './Navigation.svelte';
	import '../app.css';

	let { children } = $props();

	$effect(() => {
		document.documentElement.style.fontFamily = customization.font;
		document.documentElement.dataset.theme = customization.theme;
	});
</script>

<Navigation />
{@render children()}
