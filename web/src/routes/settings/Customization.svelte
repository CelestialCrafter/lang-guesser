<script>
	const themes = ['', 'rosepine', 'rosepine-moon', 'rosepine-dawn', 'latte', 'frappe', 'macchiato', 'mocha'];
</script>

{#each themes as theme}
	<div class="form-control">
		<label class="label cursor-pointer gap-4">
			<span class="label-text">{theme || 'Default'}</span>
			<input type="radio" name="theme-radios" class="radio theme-controller" onclick={() => localStorage.setItem('theme', theme)} value={theme} />
		</label>
	</div>
{/each}
