import { redirect } from '@sveltejs/kit';
import type { LayoutServerLoad } from './$types';

export const load: LayoutServerLoad = async ({ fetch }) => {
	const res = await fetch('http://localhost:8080/api/me', {
		credentials: 'include'
	});

	if (!res.ok) {
		throw redirect(302, '/login');
	}

	return {};
}