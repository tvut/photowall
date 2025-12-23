<script lang="ts">
	import { apiFetch } from '$lib/api';

	async function logout() {
		await apiFetch('/logout', { method: 'POST' });
		location.href = '/login';
	}
</script>

<h1>Admin Dashboard</h1>

<button on:click={logout}>Logout</button>