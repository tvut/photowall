<script lang="ts">
	let title = $state('');
	import { createPost } from '$lib/api/posts';
	import ImageUploader from './ImageUploader.svelte';
	import { attachImagesToPost } from '$lib/api/images';
	import { toast } from 'svelte-sonner';
	import { goto } from '$app/navigation';
	import { resolve } from '$app/paths';

	let uploadedImages: string[] = $state([]);
	let isCreatingPost = $state(false);

	let addingImages = $state(true);

	async function handleCreatePost() {
		if (!title.trim()) {
			toast.error('Please enter a title');
			return;
		}
		isCreatingPost = true;

		try {
			const slug = await createPost(title);

			if (uploadedImages.length > 0) {
				await attachImagesToPost(slug, uploadedImages);
				toast.success(`Created post "${title}" with ${uploadedImages.length} image(s)`);
			} else {
				toast.success(`Created post "${title}"`);
			}

			goto(resolve('/admin/edit/' + slug));
		} catch (error) {
			console.error('Failed to create post:', error);
			toast.error('Failed to create post');
		} finally {
			isCreatingPost = false;
		}
	}
</script>

<div class="card bg-base-100 shadow-sm mb-6">
	<div class="card-body">
		<h2 class="card-title">New Post</h2>

		<div class="form-control">
			<label class="label" for="post-title">
				<span class="label-text">Post Title</span>
			</label>
			<input
				id="post-title"
				type="text"
				class="input input-bordered"
				placeholder="Enter post title"
				bind:value={title}
			/>
		</div>

		<div class="form-control">
			<label class="cursor-pointer label flex items-center" for="images">
				<input
					type="checkbox"
					class="checkbox checkbox-primary"
					name="images"
					bind:checked={addingImages}
				/>
				<span class="ml-2 label-text">Adding Images</span>
			</label>
		</div>

		{#if addingImages && uploadedImages.length === 0}
			<div class="mt-4">
				<ImageUploader onFileUpload={(imageUrls) => (uploadedImages = imageUrls)} />
			</div>
		{:else if addingImages && uploadedImages.length > 0}
			<div class="mt-4">
				{uploadedImages.length} images uploaded
			</div>
		{/if}

		<div class="card-actions justify-end mt-6">
			<button
				class="btn btn-primary"
				class:loading={isCreatingPost}
				onclick={handleCreatePost}
				disabled={isCreatingPost || !title.trim() || (addingImages && uploadedImages.length === 0)}
			>
				{#if isCreatingPost}
					Creating...
				{:else}
					Create Post
				{/if}
			</button>
		</div>
	</div>
</div>
