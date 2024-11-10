<script>
	import { Highlight, LineNumbers } from 'svelte-highlight';
	import rust from 'svelte-highlight/languages/rust';
	import "svelte-highlight/styles/github.css";

	import { PUBLIC_API_URL } from '$env/static/public';

	let { next, code, id } = $props();
	let response = $state(new Promise(() => {}));
	let language = $state('');
	let duration = $state(0);
	let languageInput = $state('');

	const languageMap = {
		rust: rust,
		'': {
			name: "none",
			register: () => ({}),
		}
	};

	const handleSubmit = () =>
		response = fetch(`${PUBLIC_API_URL}/api/challenge`, {
			method: 'POST',
			body: JSON.stringify({ id, language: languageInput }),
			headers: {
				'Content-Type': 'application/json'
			}
		})
		.then(res => res.json())
		.then(data => {
			language = data.language;
			duration = data.duration;
			return data;
		});

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
	<Highlight language={languageMap[language]} {code} let:highlighted>
		<LineNumbers {highlighted} hideBorder wrapLines />
	</Highlight>
	</div>

	<div class="controls">
	{#await response}
		<form>
			<input bind:value={languageInput} list="language-list" use:focus />
			<button onclick={handleSubmit}>Submit</button>
			<span use:setstart>{duration.toFixed(2)}s</span>
		</form>

		<datalist id="language-list">
			{#each Object.keys(languageMap) as language}
				<option value={language}></option>
			{/each}
		</datalist>
	{:then { more } }
		<button use:focus onclick={() => next(more)}>Next</button>
		<span>{languageInput == language ? 'correct!' : 'wrong.'} {duration.toFixed(2)}s</span>
	{:catch error}
		<span>could not submit challenge: {error.toString()}</span>
	{/await}
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
