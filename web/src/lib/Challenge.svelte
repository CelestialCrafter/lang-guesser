<script>
	import { Highlight } from 'svelte-highlight';
	import { rust } from 'svelte-highlight/languages';
	import { PUBLIC_API_URL } from '$env/static/public';
	let { next, code, id } = $props();

	let response = new Promise(() => {});
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
		fetch(`${PUBLIC_API_URL}/challenge`, {
			method: 'POST',
			body: JSON.stringify({ id, language: selectedLanguage }),
			headers: {
				'Content-Type': 'application/json'
			}
		}).then(res => res.json()).then(data => {
			language = data.language;
			return data;
		});
</script>

<section>
	<Highlight language={languageMap[language]} {code} />
	{#await response}
		<input bind:value={selectedLanguage} list="language-list" />
		<button onclick={() => (response = handleSubmit())}>Submit</button>

		<datalist id="language-list">
			{#each Object.keys(languageMap) as language}
				<option value={language}></option>
			{/each}
		</datalist>
	{:then data}
		<button onclick={next(data.more)}>Next</button>
	{/await}
</section>
