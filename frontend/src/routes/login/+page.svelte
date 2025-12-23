<script lang="ts">
	import { goto } from '$app/navigation';
	import { apiFetch } from '$lib/api';

	let username = '';
	let password = '';
	let error = '';
	let isLoading = false;

	async function login() {
		error = '';
		isLoading = true;
		try {
			const res = await apiFetch('/login', {
				method: 'POST',
				body: JSON.stringify({ username, password })
			});

			if (res.ok) {
				goto('/admin');
			} else {
				error = 'Invalid credentials';
			}
		} catch (e) {
			error = 'Login failed. Please try again.';
		} finally {
			isLoading = false;
		}
	}
</script>

<div class="min-h-screen flex items-center justify-center bg-base-200 p-4">
	<div class="card w-full max-w-md bg-base-100 shadow-xl">
		<div class="card-body">
			<h1 class="card-title text-2xl font-bold mb-2">Login</h1>

			<form on:submit|preventDefault={login} class="space-y-4">
				<div class="form-control">
					<label class="label">
						<span class="label-text mb-1">Username</span>
					</label>
					<input 
						type="text" 
						bind:value={username} 
						placeholder="Enter your username" 
						class="input input-bordered w-full"
						required
					/>
				</div>

				<div class="form-control">
					<label class="label">
						<span class="label-text mb-1">Password</span>
					</label>
					<input 
						type="password" 
						bind:value={password} 
						placeholder="••••••••" 
						class="input input-bordered w-full"
						required
					/>
				</div>

				{#if error}
					<div class="alert alert-error mt-4">
						<svg xmlns="http://www.w3.org/2000/svg" class="stroke-current shrink-0 h-6 w-6" fill="none" viewBox="0 0 24 24">
							<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M10 14l2-2m0 0l2-2m-2 2l-2-2m2 2l2 2m7-2a9 9 0 11-18 0 9 9 0 0118 0z" />
						</svg>
						<span>{error}</span>
					</div>
				{/if}

				<div class="card-actions justify-end mt-6">
					<button 
						class="btn btn-primary w-full"
						class:loading={isLoading}
						disabled={isLoading}
					>
						{isLoading ? 'Logging in...' : 'Login'}
					</button>
				</div>
			</form>
		</div>
	</div>
</div>
