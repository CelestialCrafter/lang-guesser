<script>
	import { Highlight, HighlightAuto } from 'svelte-highlight';
	import { rust } from 'svelte-highlight/languages';
	import { API_URL } from '$env/static/public';
	let { next, code, id } = $props();

	let response = $state(new Promise(() => {}));
	let selectedLanguage = $state('');
	const languageMap = {
		rust: rust
	};

	const handleSubmit = () => fetch({
		method: 'POST',
		url: `${API_URL}/challenge`,
		body: JSON.stringify({ id, language: selectedLanguage })
	});
</script>

<section>
	{#await response}
		<HighlightAuto {code} />
		<input bind:value={selectedLanguage} list="language-list" />
		<button onclick={handleSubmit}>Submit</button>

		<datalist id="language-list">
			{#each Object.keys(languageMap) as language}
				<option value={language}></option>
			{/each}
		</datalist>
	{:then data}
		<Highlight language={languageMap[data.language]} {code} />
		<button onclick={next(data.more)}>Next</button>
	{/await}
</section>
