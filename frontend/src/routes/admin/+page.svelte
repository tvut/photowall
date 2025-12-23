<script lang="ts">
	import { goto, invalidateAll } from '$app/navigation';
	import NewPost from '$components/admin/NewPost.svelte';
	import { apiFetch } from '$lib/api';
	import type { PageProps } from './$types';

	let { data }: PageProps = $props();

	async function deletePost(slug: string) {
		await apiFetch(`/admin/posts/${slug}`, { method: 'DELETE' });
		await invalidateAll();
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
			<button
				onclick={() => deletePost(post.slug)}
				class="btn btn-square btn-ghost hover:bg-error"
				aria-label={`Delete ${post.title}`}
			>
				<svg class="size-[1.2em]" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24">
					<path
						stroke-linecap="round"
						stroke-linejoin="round"
						stroke-width="2"
						stroke="currentColor"
						fill="none"
						d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16"
					></path>
				</svg>
			</button>
		</li>
	{/each}
</ul>
