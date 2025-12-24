<script lang="ts">
	let title = $state('');
	import { createPost } from '$lib/api/posts';
</script>

<div class="card bg-base-100 shadow-sm mb-6">
	<div class="card-body">
		<form
			class="flex flex-row gap-4"
			onsubmit={(e) => {
				e.preventDefault();
				createPost(title);
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
		</form>
	</div>
</div>
