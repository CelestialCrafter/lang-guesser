<script>
	import { Highlight, LineNumbers } from 'svelte-highlight';
	import { languageToHighlight } from '$lib/languages.js';
	import './highlight.css';
	import ChallengeForm from './ChallengeForm.svelte';

	let { onnext, onsubmit, more, code, submission, duration = $bindable() } = $props();

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

		<ChallengeForm {onnext} {onsubmit} {more} {submission} bind:duration />
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
