<script lang="ts">
	import { goto, invalidateAll } from '$app/navigation';
	import NewPost from '$lib/components/NewPost.svelte';
	import { deletePost } from '$lib/api/posts';
	import { toast } from 'svelte-sonner';
	import type { PageProps } from './$types';

	let { data }: PageProps = $props();

	let isDeleting = $state(false);

	async function handleDelete(slug: string, title: string) {
		if (confirm(`Are you sure you want to delete "${title}"?`)) {
			try {
				isDeleting = true;
				await deletePost(slug);
				toast.success('Post deleted successfully');
				await invalidateAll();
			} catch (error) {
				console.error('Error deleting post:', error);
				toast.error('Failed to delete post');
			} finally {
				isDeleting = false;
			}
		}
	}
</script>

<NewPost />

<ul class="list bg-base-100 rounded-box shadow-md">
	<li class="p-4 pb-2 text-xl tracking-wide">Posts</li>

	{#each data.posts as post}
		<li class="list-row flex">
			<button
				onclick={() => goto(`/admin/edit/${post.slug}`)}
				class="flex-1 flex items-center gap-4 p-4 text-left hover:bg-base-200"
			>
				<div class="text-xl tabular-nums">
					{new Date(post.display_time).toLocaleDateString('en-GB')}
				</div>
				<div class="flex-1">
					<div class="text-2xl">{post.title}</div>
				</div>
			</button>
			<button class="btn btn-error w-40" onclick={() => handleDelete(post.slug, post.title)} disabled={isDeleting}>
				{isDeleting ? 'Deleting...' : 'Delete Post'}
			</button>
		</li>
	{/each}
</ul>
