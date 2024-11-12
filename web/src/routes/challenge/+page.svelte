<script>
	import { onMount } from 'svelte';

	import { PUBLIC_API_URL } from '$env/static/public';
	import Challenge from './Challenge.svelte';
	import Trail from './Trail.svelte';

	let response = $state(new Promise(() => {}));
	let submissions = $state([]);
	let currentDuration = $state(0);

	const onnext = () => {
		response = fetch(`${PUBLIC_API_URL}/api/challenge`).then((res) => res.text());
		submissions.push(null);
	};

	onMount(onnext);

	const onsubmit = (language) =>
		fetch(`${PUBLIC_API_URL}/api/challenge`, {
			method: 'POST',
			body: JSON.stringify({ language }),
			headers: {
				'Content-Type': 'application/json'
			}
		})
			.then((res) => res.json())
			.then((data) => (submissions = data));
</script>

<section class="p-4">
	<Trail {submissions} {currentDuration} />

	<div class="divider"></div>

	{#await response}
		<div role="alert" class="alert">
			<span>loading challenge...</span>
		</div>
	{:then code}
		<Challenge
			bind:duration={currentDuration}
			{onsubmit}
			{onnext}
			{code}
			submission={submissions[submissions.length - 1]}
		/>
	{:catch error}
		<div role="alert" class="alert alert-error">
			<span>could not load challenge: {error.toString()}</span>
		</div>
	{/await}
</section>
