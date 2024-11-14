<script>
	import { base } from '$app/paths';
	import { token } from '$lib/auth.js';
	const routes = {
		settings: '⚙️',
		auth : token ? auth : '🔒',
		challenge: '⚔️'
	};
</script>

{#snippet auth()}
	<div class="avatar">
		<img
			class="rounded-full"
			alt={token.username}
			src={token.picture}
		/>
	</div>
{/snippet}

<div class="absolute bottom-0 right-0 m-4 dropdown dropdown-hover dropdown-top">
	<button class="btn mx-2 shadow-xl rounded-full">🌎</button>
	<ul class="dropdown-content menu z-[1]">
		{#each Object.entries(routes) as [route, item]}
			<li>
				<a class="btn my-1 rounded-full shadow-xl" href="{base}/{route}">
					{#if typeof item == typeof ''}
						{item}
					{:else}
						{@render item()}
					{/if}
				</a>
			</li>
		{/each}
	</ul>
</div>
