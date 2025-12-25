<script lang="ts">
	import { goto } from '$app/navigation';
	import { resolve } from '$app/paths';
	import { apiFetchJson } from '$lib/api';

	let { children } = $props();

	async function logout() {
		await apiFetchJson('/logout', { method: 'POST' });
		goto(resolve('/login'));
	}
</script>

<div class="navbar bg-primary text-primary-content">
	<span class="text-xl ml-6">Admin</span>
	<button class="text-ml ml-6 cursor-pointer" onclick={() => goto(resolve('/'))}>Blog</button>
	<button class="text-ml ml-6 cursor-pointer" onclick={() => goto(resolve('/admin'))}>Home</button>
	<button onclick={logout} class="cursor-pointer btn btn-secondary ml-auto mr-6">Logout</button>
</div>

<div class="container mx-auto p-4">
	{@render children()}
</div>
