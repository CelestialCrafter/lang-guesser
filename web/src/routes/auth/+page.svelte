<script>
	import { token } from '$lib/auth.js';
	import LoggedIn from './LoggedIn.svelte';
	import LoggedOut from './LoggedOut.svelte';
</script>

<div class="flex h-screen justify-center items-center text-center">
	<section class="w-96 gap-4 flex flex-col bg-base-200 p-4 rounded-box">
		{#if token}
			<LoggedIn />
		{:else}
			<LoggedOut />
		{/if}
	</section>
</div>
