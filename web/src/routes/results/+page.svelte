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

<section>
	<SubmissionTrail {submissions} />

	<div class="divider"></div>

	<div class="bg-base-200 stats shadow-xl">
		<div class="stat text-primary">
			<span class="stat-title">Points</span>
			<span class="stat-value">x</span>
			<span class="stat-desc">#X on leaderboard</span>
		</div>
		<div class="stat">
			<span class="stat-title">Mode</span>
			<span class="stat-value">X</span>
		</div>
	</div>

	<div class="bg-base-200 stats shadow-xl text-secondary">
		<div class="stat">
			<span class="stat-title">Correct</span>
			<span class="stat-value">{correct}/{submissions.length}</span>
			<span class="stat-desc">top X% of all users</span>
		</div>
		<div class="stat">
			<span class="stat-title">Time</span>
			<span class="stat-value">Xs</span>
			<span class="stat-desc">top X% of all users</span>
		</div>
	</div>
</section>
