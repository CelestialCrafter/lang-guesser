<script>
	import { Highlight, LineNumbers } from 'svelte-highlight';
	import { languageToHighlight } from '$lib/languages.js';
	import './highlight.css';

	let { onnext, onsubmit, more, code, submission, duration = $bindable() } = $props();
	let submitDisabled = $state(false);
	let language = $state('');

	$effect(() => {
		let last_time = performance.now();
		const update = (time) => {
			frame = requestAnimationFrame(update);
			duration += (time - last_time) * 1e6;
			last_time = time;
		};

		let frame = requestAnimationFrame(update);
		return () => {
			duration = 0;
			cancelAnimationFrame(frame);
		};
	});

	$effect(() => {
		if (!submission) language = '';
	});

	const guessedCorrectly = $derived(
		submission ? submission.challenge.language == submission.guessed : null
	);
	const dotsClass = $derived(
		guessedCorrectly !== null ? (guessedCorrectly ? 'success-dots' : 'error-dots') : ''
	);

	const focus = el => el.focus();
</script>

<section class="card card-compact bg-base-200 shadow-xl">
	<div class="card-body">
		<div class="overflow-y-scroll mockup-code bg-base-300 {dotsClass}">
			<Highlight
				language={languageToHighlight[submission?.challenge.language ?? '']}
				{code}
				let:highlighted
			>
				<LineNumbers {highlighted} hideBorder wrapLines />
			</Highlight>
		</div>

		<form
			class="card-actions items-center justify-between"
			onsubmit={(event) => event.preventDefault()}
		>
			<input
				class="input input-bordered"
				disabled={submitDisabled}
				bind:value={language}
				use:focus
				list="language-list"
			/>

			{#if !submission}
				<button
					class="btn btn-primary"
					disabled={submitDisabled}
					onclick={event => {
						event.target.setAttribute('disabled', '');
						onsubmit(language);
					}}
				>
					Submit
				</button>
			{:else}
				<button
					class="btn btn-primary"
					onclick={onnext}
				>
					{more ? 'Next' : 'Finish'}
				</button>
			{/if}
		</form>
		<datalist id="language-list">
			{#each Object.keys(languageToHighlight) as language}
				<option value={language}></option>
			{/each}
		</datalist>
	</div>
</section>

<style>
	.mockup-code {
		height: 38rem;
	}

	.error-dots::before {
		color: oklch(var(--er));
		opacity: 1;
	}

	.success-dots::before {
		color: oklch(var(--su));
		opacity: 1;
	}
</style>
