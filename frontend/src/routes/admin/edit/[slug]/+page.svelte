<script lang="ts">
	import { toast } from 'svelte-sonner';
	import { updatePostStatus, updatePostDisplayTime } from '$lib/api/posts';
	import type { PageProps } from './$types';

	let { data }: PageProps = $props();

	let status = $derived(data.post.status);
	let displayTime = $derived(
		data.post.display_time ? new Date(data.post.display_time).toISOString().slice(0, 16) : ''
	);

	const updateStatus = async (newStatus: string) => {
		try {
			await updatePostStatus(data.post.slug, newStatus);
			toast.success(`Status updated to ${newStatus}`);
		} catch (error) {
			console.error('Error updating status:', error);
			toast.error('Failed to update status');
		}
	};

	const updateDisplayTime = async (newDisplayTime: string) => {
		try {
			await updatePostDisplayTime(data.post.slug, newDisplayTime);
			toast.success('Display time updated');
		} catch (error) {
			console.error('Error updating display time:', error);
			toast.error('Failed to update display time');
		}
	};
</script>

<div class="container mx-auto p-4">
	<div class="flex flex-row gap-4">
		<div>
			<h1 class="text-2xl font-bold">Post: {data.post.title}</h1>
			<h2 class="text-xl">Slug: {data.post.slug}</h2>
			<h3 class="text-lg">Created at: {data.post.created_at}</h3>
			<h3 class="text-lg">Display time: {data.post.display_time}</h3>
		</div>

		<div class="flex-auto bg-base-100 rounded-lg shadow p-6 mb-6">
			<div class="form-control mb-4">
				<label class="label">
					<span class="label-text">Status</span>
				</label>
				<select
					class="select select-bordered w-full max-w-xs"
					bind:value={status}
					onchange={() => updateStatus(status)}
				>
					<option value="draft">Draft</option>
					<option value="published">Published</option>
				</select>
			</div>

			<div class="form-control mb-4">
				<label class="label">
					<span class="label-text">Display Time</span>
				</label>
				<div class="flex gap-2">
					<input
						type="datetime-local"
						class="input input-bordered flex-1"
						bind:value={displayTime}
					/>
					<button class="btn btn-primary" onclick={() => updateDisplayTime(displayTime)}>
						Update Time
					</button>
				</div>
			</div>
		</div>
	</div>
</div>
