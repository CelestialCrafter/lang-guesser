<script>
	import { onMount } from 'svelte';

	import SubmissionTrail from '$lib/SubmissionTrail.svelte';
	import { loadSession } from '$lib/session.js';

	let submissions = $state([]);
	onMount(() => loadSession().then((data) => (submissions = data)));

	const correct = $derived(
		submissions.reduce((acc, x) => acc + (x.guessed === x.challenge.language), 0)
	);
</script>

<SubmissionTrail {submissions} />

<div class="divider"></div>

<div class="stats shadow-xl">
	<div class="stat">
		<span class="stat-title">Correct</span>
		<span class="stat-value text-primary">{correct}/{submissions.length}</span>
		<span class="stat-desc">top x% of all users</span>
	</div>
</div>
