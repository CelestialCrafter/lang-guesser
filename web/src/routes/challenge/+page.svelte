<script>
	import { onMount } from 'svelte';

	import { PUBLIC_API_URL } from '$env/static/public';
	import Challenge from './Challenge.svelte';
	import SubmissionTrail from '$lib/SubmissionTrail.svelte';
	import { loadSession } from '$lib/session.js';
	import { goto } from '$app/navigation';
	import { base } from '$app/paths';

	let code = $state(null);
	let more = $state(true);
	let submissions = $state([]);
	let currentDuration = $state(0);

	const onnext = async () => {
		if (!more) return await goto(base + '/results');

		code = await (await fetch(`${PUBLIC_API_URL}/challenge`)).text();
		submissions.push(null);
	};

	onMount(() => loadSession().then(() => onnext()));

	const onsubmit = async (language) => {
		let data = await (
			await fetch(`${PUBLIC_API_URL}/challenge`, {
				method: 'POST',
				body: JSON.stringify({ language }),
				headers: {
					'Content-Type': 'application/json'
				}
			})
		).json();

		submissions = data.past;
		more = data.more;
	};
</script>

<section class="p-4">
	<SubmissionTrail {submissions} {currentDuration} />

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
			{more}
			submission={submissions[submissions.length - 1]}
		/>
	{/if}
</section>
