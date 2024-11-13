<script>
	import { setFont } from '$lib';

	const { font } = $props();
</script>

<button
	class="btn"
	onclick={(event) => setFont(event.target.dataset.setFont)}
	data-set-font={font}
>
	{font || 'default'}
</button>
