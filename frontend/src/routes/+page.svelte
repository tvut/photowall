<script lang="ts">
	import { onMount } from 'svelte';
	import { getPosts, type Post } from '$lib/api';

	let posts: Post[] = [];

	onMount(async () => {
		posts = await getPosts();
	});
</script>

<h1>Blog</h1>

{#each posts as post}
	<article>
		<h2>{post.title}</h2>
	</article>
{/each}
