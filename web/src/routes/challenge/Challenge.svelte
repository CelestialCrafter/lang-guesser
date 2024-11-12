<script>
	import { Highlight, LineNumbers } from 'svelte-highlight';
	import './highlight.css';
	import { languageToHighlight } from '$lib/languages.js';

	let { onnext, onsubmit, more, code, submission, duration = $bindable() } = $props();
	let disabled = $state(false);
	let languageInput = $state('');

	const focus = (el) => el.focus();
	const setstart = () => {
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
	};

	const guessedCorrectly = $derived(
		submission ? submission.challenge.language == submission.guessed : null
	);
	const dotsClass = $derived(
		guessedCorrectly !== null ? (guessedCorrectly ? 'success-dots' : 'error-dots') : ''
	);
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
			class="card-actions items-center {submission ? 'justify-end' : 'justify-between'}"
			use:setstart
			onsubmit={(event) => event.preventDefault()}
		>
			{#if !submission}
				<input
					class="input input-bordered"
					bind:value={languageInput}
					list="language-list"
					use:focus
				/>
				<button
					class="btn btn-primary"
					{disabled}
					onclick={() => {
						disabled = true;
						onsubmit(languageInput);
						languageInput = '';
						disabled = false;
					}}
				>
					Submit
				</button>
			{:else}
				<button class="btn btn-primary" use:focus onclick={onnext}
					>{more ? 'Next' : 'Finish'}</button
				>
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
