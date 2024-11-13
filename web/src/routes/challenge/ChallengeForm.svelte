<script>
	let { onnext, onsubmit, more, submission, duration = $bindable() } = $props();
	let language = null;

	$effect(() => {
		let last_time = performance.now();
		const update = (time) => {
			frame = requestAnimationFrame(update);
			duration += (time - last_time) * 1e6;
			last_time = time;
		};

		let frame = requestAnimationFrame(update);
		return () => {
			duration = 0;
			cancelAnimationFrame(frame);
		};
	});

	const focus = (el) => el.focus();
</script>

<form
	class="card-actions items-center justify-between"
	onsubmit={(event) => event.preventDefault()}
>
	<input class="input input-bordered" bind:this={language} use:focus list="language-list" />

	{#if !submission}
		<button
			class="btn btn-primary"
			onclick={(event) => {
				language.setAttribute('disabled', '');
				event.target.setAttribute('disabled', '');
				onsubmit(language.value);
			}}
		>
			Submit
		</button>
	{:else}
		<button
			use:focus
			class="btn btn-primary"
			onclick={() => {
				language.value = '';
				language.removeAttribute('disabled');
				language.focus();
				duration = 0;
				onnext();
			}}
		>
			{more ? 'Next' : 'Finish'}
		</button>
	{/if}
</form>
