<script>
	import Customization from './Customization.svelte';

	const pages = {
		Customization: Customization
	};

	let currentPage = $state(Object.entries(pages)[0]);
	let CurrentComponent = $derived(currentPage[1]);
</script>

<div class="flex">
	<ul class="menu bg-base-200 w-64 h-screen">
		<li>
			<h2 class="menu-title">Settings</h2>
			<ul>
				{#each Object.entries(pages) as page}
					<li>
						<button class:active={currentPage[0] == page[0]} onclick={() => (currentPage = page)}>
							{page[0]}
						</button>
					</li>
				{/each}
			</ul>
		</li>
	</ul>

	<CurrentComponent />
</div>
