{#snippet authButton(url, name)}
	<button class="btn btn-primary">
		<img class="h-6 w-6" src={url} alt={name} />
		<span>{name}</span>
	</button>
{/snippet}

<span class="text-xl">Login</span>
<ul class="menu gap-4 p-0">
	{@render authButton('static/github.svg', 'Github')}
	{@render authButton('static/google.svg', 'Google')}
</ul>
