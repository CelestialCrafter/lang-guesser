<script>
	import { Highlight, LineNumbers } from 'svelte-highlight';
	import './highlight.css';
	import { languageToHighlight } from '$lib/languages.js';

	let { onnext, onsubmit, code, submission } = $props();
	let duration = $state(0);
	let languageInput = $state('');

	const focus = el => el.focus();
	const setstart = () => {
		$effect(() => {
			let last_time = performance.now();
			const update = time => {
				frame = requestAnimationFrame(update);
				duration += (time - last_time) / 1000;
				last_time = time
			};

			let frame = requestAnimationFrame(update);
			return () => cancelAnimationFrame(frame);
		});
	};
</script>

<section>
	<div class="code">
	<Highlight language={languageToHighlight[submission?.challenge.language ?? '']} {code} let:highlighted>
		<LineNumbers {highlighted} hideBorder wrapLines />
	</Highlight>
	</div>

	<div class="inline">
	{#if !submission}
		<form onsubmit={event => event.preventDefault()}>
			<input class="input input-bordered" bind:value={languageInput} list="language-list" use:focus />
			<button class="btn btn-primary" onclick={() => onsubmit(languageInput)}>Submit</button>
			<span use:setstart>{duration.toFixed(2)}s</span>
		</form>

		<datalist id="language-list">
			{#each Object.keys(languageToHighlight) as language}
				<option value={language}></option>
			{/each}
		</datalist>
	{:else}
		<button class="btn btn-primary" use:focus onclick={onnext}>Next</button>
		<!-- d / 1e+9 = ns -> s -->
		<span>{submission.challenge.language == submission.guessed ? 'correct!' : 'wrong.'} {(submission.duration / 1e+9).toFixed(2)}s</span>
	{/if}
	</div>
</section>

<style>
section {
display: flex;
height: 70vh;
flex-direction: column;
}

.code {
overflow: scroll;
flex: 1;
}
</style>
