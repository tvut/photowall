<script lang="ts">
	import type { PageProps } from './$types';
	import PostSettings from '$lib/components/PostSettings.svelte';
	import PostInfoCard from '$lib/components/PostInfoCard.svelte';

	let { data }: PageProps = $props();

	let status = $derived(data.post.status);
	let displayTime = $derived(
		data.post.display_time ? new Date(data.post.display_time).toISOString().slice(0, 16) : ''
	);
</script>

<div class="mb-8">
	<h1 class="text-3xl font-bold text-gray-900">Edit Post: {data.post.title}</h1>
</div>
<div class="grid grid-cols-1 lg:grid-cols-3 gap-6 h-full">
	<div class="lg:col-span-2">
		<PostInfoCard post={data.post} />
	</div>
	<div class="lg:col-span-1">
		<PostSettings
			post={{
				slug: data.post.slug,
				title: data.post.title,
				status: status,
				displayTime: displayTime
			}}
		/>
	</div>
</div>
