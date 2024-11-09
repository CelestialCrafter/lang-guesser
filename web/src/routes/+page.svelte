<script>
	import { PUBLIC_API_URL } from '$env/static/public';
	import Challenge from '$lib/Challenge.svelte';
	import { onMount } from 'svelte';

	let response = new Promise(() => {});
	onMount(() => (response = fetch(`${PUBLIC_API_URL}/challenge`).then(res => res.json())));
</script>

<section>
	{#await response then { code, id }}
		<Challenge code={atob(code)} {id} />
	{/await}
</section>
