<script>
	import { Highlight, LineNumbers } from 'svelte-highlight';
	import { go, python } from 'svelte-highlight/languages';
	import rust from 'svelte-highlight/languages/rust';
	import "svelte-highlight/styles/github.css";

	let { onnext, onsubmit, code, submission } = $props();
	let duration = $state(0);
	let languageInput = $state('');

	const languageMap = {
		rust: rust,
		go: go,
		python: python,
		'': {
			name: "none",
			register: () => ({}),
		}
	};

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
	<Highlight language={languageMap[submission?.language ?? '']} {code} let:highlighted>
		<LineNumbers {highlighted} hideBorder wrapLines />
	</Highlight>
	</div>

	<div class="controls">
	{#if !submission}
		<form>
			<input bind:value={languageInput} list="language-list" use:focus />
			<button onclick={() => onsubmit(languageInput)}>Submit</button>
			<span use:setstart>{duration.toFixed(2)}s</span>
		</form>

		<datalist id="language-list">
			{#each Object.keys(languageMap) as language}
				<option value={language}></option>
			{/each}
		</datalist>
	{:else}
		<button use:focus onclick={onnext}>Next</button>
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

.controls {
display: inherit;
}
</style>
