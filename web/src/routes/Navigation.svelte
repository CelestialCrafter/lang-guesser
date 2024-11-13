<script>
	import { base } from '$app/paths';
	const routes = {
		settings: '⚙️',
		challenge: '⚔️'
	};
</script>

<div class="absolute bottom-0 right-0 m-4 dropdown dropdown-hover dropdown-top">
	<button class="btn mx-2 shadow-xl rounded-full">🌎</button>
	<ul class="dropdown-content menu z-[1]">
		{#each Object.entries(routes) as [route, name]}
			<li>
				<a class="btn my-1 rounded-full shadow-xl" href="{base}/{route}">{name}</a>
			</li>
		{/each}
	</ul>
</div>
