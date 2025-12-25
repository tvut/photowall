import { apiFetchJson } from '../api';

export async function uploadImagesBackend(files: File[]): Promise<string[]> {
	const formData = new FormData();

	files.forEach((file) => {
		formData.append('images', file);
	});

	const response = await fetch('http://localhost:8080/api/admin/upload-images', {
		method: 'POST',
		body: formData,
		credentials: 'include'
	});

	if (!response.ok) {
		throw new Error('Upload failed');
	}

	const result = await response.json();
	return result.image_urls;
}

export async function attachImagesToPost(postSlug: string, imageUrls: string[]): Promise<void> {
	const res = await apiFetchJson('/admin/attach-images', {
		method: 'POST',
		body: JSON.stringify({
			post_slug: postSlug,
			image_urls: imageUrls
		})
	});

	if (!res.ok) {
		const errorText = await res.text();
		throw new Error(`Failed to attach images: ${errorText}`);
	}
}
