<script lang="ts">
	import { apiFetch } from '$lib/api';
	import NewPost from '$components/admin/NewPost.svelte';

	async function logout() {
		await apiFetch('/logout', { method: 'POST' });
		location.href = '/login';
	}
</script>

<h1>Admin Dashboard</h1>

<NewPost />

<button on:click={logout}>Logout</button>