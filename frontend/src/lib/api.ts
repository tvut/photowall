const API = 'http://localhost:8080/api';

export type Post = {
	title: string;
	slug: string;
	display_time: Date;
	photos: string[];
};

export async function apiFetchJson(path: string, options: RequestInit = {}) {
	const defaultHeaders: Record<string, string> = {
		'Content-Type': 'application/json'
	};

	return fetch(`${API}${path}`, {
		credentials: 'include',
		headers: defaultHeaders,
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
