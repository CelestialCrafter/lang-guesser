{#snippet authButton(url, provider, name)}
	<a href="/api/auth/{provider}" class="btn btn-primary">
		<img class="h-6 w-6" src={url} alt={name} />
		<span>{name}</span>
	</a>
{/snippet}

<span class="text-xl">Login</span>
<ul class="menu gap-4 p-0">
	{@render authButton('static/github.svg', 'github', 'Github')}
	{@render authButton('static/google.svg', 'google', 'Google')}
</ul>
