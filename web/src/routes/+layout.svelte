<script>
	import { onMount } from 'svelte';
	import { themeChange } from 'theme-change';
	import ThemeButton from './ThemeButton.svelte';

	let { children } = $props();

	onMount(() => themeChange(false));
</script>

<section class="p-4">
	<div class="flex gap-2">
		{#each ['', 'rosepine', 'rosepine-moon', 'rosepine-dawn'] as theme}
			<ThemeButton {theme} />
		{/each}
	</div>

	<div class="divider"></div>

	{@render children()}
</section>
