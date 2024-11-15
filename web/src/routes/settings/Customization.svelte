<script>
	import customization from '$lib/customization.svelte.js';

	const themes = {
		Default: '',
		'Rosé Pine': 'rosepine',
		'Rosé Pine Moon': 'rosepine-moon',
		'Rosé Pine Dawn': 'rosepine-dawn',
		'Catppuccin Latte': 'latte',
		'Catppuccin Frappé': 'frappe',
		'Catppuccin Macchiato': 'macchiato',
		'Catppuccin Mocha': 'mocha',
		'Qtile Rice': 'qtile-rice'
	};

	const fonts = {
		Default: '',
		'JetBrains Mono': 'JetBrainsMono',
		'Monaspace Neon': 'Monaspace Neon',
		'Fira Code': 'FiraCode'
	};

	const uppercase = (str) => str.charAt(0).toUpperCase() + str.slice(1);
</script>

{#snippet control(type, name, value)}
{@const active = value == customization[type + 'Raw']}
	<div class="form-control">
		<label class="label cursor-pointer gap-4">
			<span class="label-text">{name}</span>
			<input
				checked={active}
				type="radio"
				name="{type}-radios"
				class="radio"
				class:radio-secondary={active}
				onclick={() => (customization[type] = value)}
			/>
		</label>
	</div>
{/snippet}

{#snippet section(type, data)}
	<div class="collapse collapse-arrow bg-base-200">
		<input type="radio" name="section-radios" checked="checked" />
		<div class="collapse-title text-xl font-medium">{uppercase(type)}</div>
		<div class="collapse-content">
			{#each Object.entries(data) as [name, value]}
				{@render control(type, name, value)}
			{/each}
		</div>
	</div>
{/snippet}

<section class="p-4 grow flex flex-col gap-4">
	{@render section('theme', themes)}
	{@render section('font', fonts)}
</section>
