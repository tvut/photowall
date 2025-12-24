const API = 'http://localhost:8080/api';

export type Post = {
	id: string;
	title: string;
	content: string;
	createdAt: string;
	updatedAt: string;
};

export async function apiFetch(path: string, options: RequestInit = {}) {
	return fetch(`${API}${path}`, {
		credentials: 'include',
		headers: {
			'Content-Type': 'application/json'
		},
		...options
	});
}

export async function getPosts(): Promise<Post[]> {
	const res = await fetch(`${API}/posts`);

	if (!res.ok) {
		throw new Error(`Failed to fetch posts: ${res.status}`);
	}

	return (await res.json()) as Post[];
}
