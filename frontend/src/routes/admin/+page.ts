import type { PageLoad } from './$types';
import { apiFetch } from '$lib/api';
import { redirect } from '@sveltejs/kit';

export const load: PageLoad = async () => {
    const res = await apiFetch('/admin/posts', {
			method: 'GET'
		});
        
    if (!res.ok) {
        if (res.status === 401 || res.status === 403) {
            throw redirect(303, '/login');
        }
        throw new Error(`Failed to fetch post: ${res.status}`);
    }
        
	return {
		posts: await res.json()
	};
};