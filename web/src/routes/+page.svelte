<script>
	import { PUBLIC_API_URL } from '$env/static/public';
	import Challenge from '$lib/Challenge.svelte';
	import { onMount } from 'svelte';

	let response = $state(new Promise(() => {}));
	let submissions = $state([]);

	const onnext = () => {
		response = fetch(`${PUBLIC_API_URL}/api/challenge`).then(res => res.text());
		submissions.push(null);
	};

	onMount(onnext);

	const onsubmit = language =>
		fetch(`${PUBLIC_API_URL}/api/challenge`, {
			method: 'POST',
			body: JSON.stringify({ language }),
			headers: {
				'Content-Type': 'application/json'
			}
		})
		.then(res => res.json())
		.then(data => (submissions = data));
</script>

<section>
	{#await response}
		<span>loading challenge...</span>
	{:then code}
		<Challenge {onsubmit} {onnext} {code} submission={submissions[submissions.length - 1]} />
	{:catch error}
		<span>could not load challenge: {error.toString()}</span>
	{/await}
</section>
