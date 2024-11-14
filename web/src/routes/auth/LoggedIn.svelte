<script>
	import { token } from '$lib/auth.js';
</script>

<div>
	<div class="w-24 avatar">
		<img class="rounded-full" alt={token.username} src={token.picture} />
	</div>
</div>

<span class="text-xl">Logged in as {token.username}</span>
<button
	onclick={() => {
		localStorage.removeItem('token');
		location.reload();
	}}
	class="btn btn-primary">Log Out</button
>
