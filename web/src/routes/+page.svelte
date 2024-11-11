<script>
	import { PUBLIC_API_URL } from '$env/static/public';
	import Challenge from '$lib/Challenge.svelte';
	import { languageToIcon } from '$lib/languages.js';
	import { convert, deserialize, OKLCH, RGBToHex, sRGB } from '@texel/color';
	import { onMount } from 'svelte';

	let response = $state(new Promise(() => {}));
	let submissions = $state([]);

	const onnext = () => {
		response = fetch(`${PUBLIC_API_URL}/api/challenge`).then(res => res.text());
		submissions.push(null);
	};

	onMount(onnext);

	const onsubmit = language =>
		fetch(`${PUBLIC_API_URL}/api/challenge`, {
			method: 'POST',
			body: JSON.stringify({ language }),
			headers: {
				'Content-Type': 'application/json'
			}
		})
		.then(res => res.json())
		.then(data => (submissions = data));

	const submissionToClass = submission => {
		if (!submission) return '';
		if (submission.guessed != submission.challenge.language) return 'use-icon step-error';
		return 'use-icon step-success';
	};

	const createIconUrl = language => {
		if (!language) return '';

		let { svg } = languageToIcon[language];
		const oklch = getComputedStyle(document.querySelector(':root')).getPropertyValue('--b3');
		const { coords } = deserialize(`oklch(${oklch})`);
		const rgb = convert(coords, OKLCH, sRGB);
		svg = svg.replace('<svg', `<svg fill="${RGBToHex(rgb)}"`);
		const base64 = btoa(svg);
		return `url("data:image/svg+xml;base64,${base64}")`;
	};

	const scrollEnd = node => {
		$effect(() => {
			$state.snapshot(submissions);
			node.scroll(node.scrollWidth, 0);
		});
	};
</script>

<section>
	{#await response}
		<span>loading challenge...</span>
	{:then code}
		<Challenge {onsubmit} {onnext} {code} submission={submissions[submissions.length - 1]} />
	{:catch error}
		<span>could not load challenge: {error.toString()}</span>
	{/await}


	<div use:scrollEnd class="overflow-x-hidden scroll-smooth">
		<ul class="steps">
		{#each submissions as submission}
			{@const language = submission?.challenge.language}
			<li style="--icon: {createIconUrl(language)};" data-content="" class="step {submissionToClass(submission)}">{language}</li>
		{/each}
		</ul>
	</div>
</section>

<style>
.use-icon::after {
	background-image: var(--icon);
	background-position: center;
	background-repeat: no-repeat;
	background-size: 1.7rem;
}
</style>
