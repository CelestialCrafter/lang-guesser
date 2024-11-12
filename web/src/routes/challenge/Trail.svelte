<script>
	import { convert, deserialize, OKLCH, RGBToHex, sRGB } from '@texel/color';
	import { languageToIcon } from '$lib/languages.js';
	import { formatDuration } from '$lib';

	const { submissions, currentDuration } = $props();

	const createIconUrl = (language) => {
		if (!language) return '';

		let { svg } = languageToIcon[language];
		const oklch = getComputedStyle(document.querySelector(':root')).getPropertyValue('--b3');
		const { coords } = deserialize(`oklch(${oklch})`);
		const rgb = convert(coords, OKLCH, sRGB);
		svg = svg.replace('<svg', `<svg fill="${RGBToHex(rgb)}"`);
		const base64 = btoa(svg);
		return `url("data:image/svg+xml;base64,${base64}")`;
	};

	const scrollEnd = (node) => {
		$effect(() => {
			$state.snapshot(submissions);
			node.scroll(node.scrollWidth, 0);
		});
	};
</script>

<section use:scrollEnd class="overflow-x-scroll scroll-smooth">
	<ul class="steps">
		{#each submissions as submission}
			{#if submission}
				{@const {
					duration,
					guessed,
					challenge: { language }
				} = submission}
				{@const stepClass = guessed == language ? 'step-success' : 'step-error'}
				<li
					style="--icon: {createIconUrl(submission?.challenge.language)};"
					class="step use-icon {stepClass}"
					data-content=""
				>
					{formatDuration(duration)}
				</li>
			{:else}
				<li class="step" data-content="">{formatDuration(currentDuration)}</li>
			{/if}
		{/each}
	</ul>
</section>

<style>
	.use-icon::after {
		background-image: var(--icon);
		background-position: center;
		background-repeat: no-repeat;
		background-size: 1.7rem;
	}

	section {
		scrollbar-width: none;
	}
</style>
