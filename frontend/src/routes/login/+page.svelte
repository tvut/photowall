<script lang="ts">
	import { goto } from '$app/navigation';
	import { apiFetch } from '$lib/api';

	let username = '';
	let password = '';
	let error = '';

	async function login() {
		error = '';

		const res = await apiFetch('/login', {
			method: 'POST',
			body: JSON.stringify({ username, password })
		});

		if (res.ok) {
			goto('/admin');
		} else {
			error = 'Invalid credentials';
		}
	}
</script>

<h1>Admin Login</h1>

<form on:submit|preventDefault={login}>
	<input bind:value={username} placeholder="Username" />
	<input type="password" bind:value={password} />
	<button>Login</button>
</form>

{#if error}
	<p style="color:red">{error}</p>
{/if}
