<script>
	import { Highlight, LineNumbers } from 'svelte-highlight';
	import rust from 'svelte-highlight/languages/rust';
	import "svelte-highlight/styles/github.css";

	import { PUBLIC_API_URL } from '$env/static/public';

	let { next, code, id } = $props();
	let response = $state(new Promise(() => {}));
	let language = $state('');
	let selectedLanguage = $state('');

	const languageMap = {
		rust: rust,
		'': {
			name: "none",
			register: () => ({}),
		}
	};

	const handleSubmit = () =>
		response = fetch(`${PUBLIC_API_URL}/challenge`, {
			method: 'POST',
			body: JSON.stringify({ id, language: selectedLanguage }),
			headers: {
				'Content-Type': 'application/json'
			}
		})
		.then(res => res.json())
		.then(data => {
			language = data.language ?? '';
			return data;
		});
</script>

<section>
	<div class="code">
	<Highlight language={languageMap[language]} {code} let:highlighted>
		<LineNumbers {highlighted} hideBorder wrapLines />
	</Highlight>
	</div>

	<div class="controls">
	{#await response}
		<input bind:value={selectedLanguage} list="language-list" />
		<button onclick={handleSubmit}>Submit</button>

		<datalist id="language-list">
			{#each Object.keys(languageMap) as language}
				<option value={language}></option>
			{/each}
		</datalist>
	{:then { more } }
		<button onclick={() => next(more)}>Next</button>
		<span>{selectedLanguage == language ? 'correct!' : 'wrong.'}</span>
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
