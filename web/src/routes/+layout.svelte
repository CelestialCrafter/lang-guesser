<script>
	import { onMount } from 'svelte';
	import { setTheme, setFont, setDefaultFont } from '$lib';

	import '../app.css';
	import Navigation from './Navigation.svelte';

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

<Navigation />
{@render children()}
