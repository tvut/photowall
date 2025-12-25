import { apiFetchJson } from '../api';
import { invalidateAll } from '$app/navigation';
import { toast } from 'svelte-sonner';

export async function updatePostStatus(slug: string, status: string): Promise<void> {
	try {
		const res = await apiFetchJson(`/admin/posts/${slug}/status`, {
			method: 'PUT',
			body: JSON.stringify({ status })
		});

		if (!res.ok) {
			throw new Error('Failed to update status');
		}

		await invalidateAll();
		return;
	} catch (error) {
		console.error('Error updating status:', error);
		throw error;
	}
}

export async function updatePostDisplayTime(slug: string, displayTime: string): Promise<void> {
	try {
		const formattedTime = new Date(displayTime).toISOString();
		const res = await apiFetchJson(`/admin/posts/${slug}/display-time`, {
			method: 'PUT',
			body: JSON.stringify({ display_time: formattedTime })
		});

		if (!res.ok) {
			throw new Error('Failed to update display time');
		}

		await invalidateAll();
		return;
	} catch (error) {
		console.error('Error updating display time:', error);
		throw error;
	}
}

export async function deletePost(slug: string): Promise<void> {
	try {
		const res = await apiFetchJson(`/admin/posts/${slug}`, {
			method: 'DELETE'
		});

		if (!res.ok) {
			throw new Error('Failed to delete post');
		}

		await invalidateAll();
		return;
	} catch (error) {
		console.error('Error deleting post:', error);
		throw error;
	}
}

export async function createPost(title: string): Promise<string> {
	const res = await apiFetchJson('/admin/add-post', {
		method: 'POST',
		body: JSON.stringify({ title: title })
	});

	if (res.ok) {
		const slug = await res.text();
		return slug;
	} else {
		console.error('Failed to create post, response: ' + res);
		toast.error('Failed to create post');
		throw new Error('Failed to create post');
	}
}
