<script lang="ts">
	let title = $state('');
	let error = $state('');
	import { apiFetch } from '$lib/api';
	import { goto } from '$app/navigation';

	async function createPost() {
		const res = await apiFetch('/admin/add-post', {
			method: 'POST',
			body: JSON.stringify({ title: title })
		});

		if (res.ok) {
			let slug = await res.text();
			goto('/admin/edit/' + slug);
		} else {
			error = 'Failed to create post';
		}
	}
</script>

<div class="card bg-base-100 shadow-sm mb-6">
	<div class="card-body">
		<form
			class="flex flex-row gap-4"
			onsubmit={(e) => {
				e.preventDefault();
				createPost();
			}}
		>
			<div class="text-2xl font-medium whitespace-nowrap">New Post</div>
			<input
				type="text"
				class="input input-bordered w-full"
				placeholder="Title"
				bind:value={title}
			/>
			<button class="btn btn-soft" type="submit">Create Post</button>
			{#if error}
				<p style="color: red;">{error}</p>
			{/if}
		</form>
	</div>
</div>
