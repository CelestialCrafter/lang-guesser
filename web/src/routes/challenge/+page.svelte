<script>
	import { onMount } from 'svelte';

	import { PUBLIC_API_URL } from '$env/static/public';
	import Challenge from './Challenge.svelte';
	import Trail from './Trail.svelte';

	let code = $state(null);
	let submissions = $state([]);
	let currentDuration = $state(0);

	const onnext = async () => {
		code = await (await fetch(`${PUBLIC_API_URL}/challenge`)).text();
		submissions.push(null);
	};

	const loadSession = async () =>
		(submissions = await (await fetch(`${PUBLIC_API_URL}/session`)).json());

	onMount(() => {
		onnext();
		loadSession();
	});

	const onsubmit = async (language) =>
		(submissions = await (
			await fetch(`${PUBLIC_API_URL}/challenge`, {
				method: 'POST',
				body: JSON.stringify({ language }),
				headers: {
					'Content-Type': 'application/json'
				}
			})
		).json());
</script>

<section class="p-4">
	<Trail {submissions} {currentDuration} />

	<div class="divider"></div>

	{#if !code}
		<div role="alert" class="alert">
			<span>loading challenge...</span>
		</div>
	{:else}
		<Challenge
			bind:duration={currentDuration}
			{onsubmit}
			{onnext}
			{code}
			submission={submissions[submissions.length - 1]}
		/>
	{/if}
</section>
