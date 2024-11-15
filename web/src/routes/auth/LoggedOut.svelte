<script>
	import { base } from '$app/paths';
	import { PUBLIC_API_URL } from '$env/static/public';
</script>

{#snippet authButton(url, provider, name)}
	<a href="{PUBLIC_API_URL}/auth/{provider}" class="btn btn-primary">
		<img class="h-6 w-6" src={url} alt={name} />
		<span>{name}</span>
	</a>
{/snippet}

<span class="text-xl">Login</span>
<ul class="menu gap-4 p-0">
	{@render authButton(`${base}/github.svg`, 'github', 'Github')}
	{@render authButton(`${base}/google.svg`, 'google', 'Google')}
</ul>
