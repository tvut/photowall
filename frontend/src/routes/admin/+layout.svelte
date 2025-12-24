<script lang="ts">
	import { goto } from '$app/navigation';

	let { children } = $props();

	import { apiFetch } from '$lib/api';

	async function logout() {
		await apiFetch('/logout', { method: 'POST' });
		location.href = '/login';
	}
</script>

<div class="navbar bg-primary text-primary-content">
	<span class="text-xl ml-6">Admin</span>
	<button class="text-ml ml-6 cursor-pointer" onclick={() => goto('/')}>Blog</button>
	<button class="text-ml ml-6 cursor-pointer" onclick={() => goto('/admin')}>Home</button>
	<button onclick={logout} class="cursor-pointer btn btn-secondary ml-auto mr-6">Logout</button>
</div>

<div class="container mx-auto p-4">
	{@render children()}
</div>
