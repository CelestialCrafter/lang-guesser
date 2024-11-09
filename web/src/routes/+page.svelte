<script>
	import { PUBLIC_API_URL } from '$env/static/public';
	import Challenge from '$lib/Challenge.svelte';
	import { onMount } from 'svelte';

	let response = $state(new Promise(() => {}));
	const loadNewChallenge = () => (response = fetch(`${PUBLIC_API_URL}/challenge`).then(res => res.json()));
	onMount(loadNewChallenge);
</script>

<section>
	{#await response}
		<span>loading challenge...</span>
	{:then { code, id }}
		<Challenge next={more => more && loadNewChallenge()} code={atob(code)} {id} />
	{:catch error}
		<span>could not load challenge: {error.toString()}</span>
	{/await}
</section>
