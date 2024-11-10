<script>
	import { PUBLIC_API_URL } from '$env/static/public';
	import Challenge from '$lib/Challenge.svelte';
	import { onMount } from 'svelte';

	let response = $state(new Promise(() => {}));
	const loadNewChallenge = () => (response = fetch(`${PUBLIC_API_URL}/api/challenge`).then(res => res.text()));
	onMount(loadNewChallenge);
</script>

<section>
	{#await response}
		<span>loading challenge...</span>
	{:then code}
		<Challenge next={more => more && loadNewChallenge()} {code} />
	{:catch error}
		<span>could not load challenge: {error.toString()}</span>
	{/await}
</section>
