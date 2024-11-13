<script>
	import { onMount } from 'svelte';
	import { setTheme, setFont, setDefaultFont } from '$lib';

	import '../app.css';
	import Navigation from './Navigation.svelte';

	let { children } = $props();

	onMount(() => {
		window.addEventListener(
			'font',
			(event) => (document.documentElement.style.fontFamily = event.detail)
		);

		window.addEventListener(
			'theme',
			(event) => (document.documentElement.dataset.theme = event.detail)
		);

		window.addEventListener(
			'storage',
			() => setCustomization()
		);

		setDefaultFont(getComputedStyle(document.documentElement).fontFamily);
		const setCustomization = () => {
			setTheme(localStorage.getItem('theme'));
			setFont(localStorage.getItem('font'));
		};

		setCustomization();
		window.onstorage = setCustomization;
	});
</script>

<Navigation />
{@render children()}
