<script lang="ts">
    
    let title = $state('');
    let error = $state('');
    import { apiFetch } from '$lib/api';
    import { goto } from '$app/navigation';

    async function createPost() {

		const res = await apiFetch('/admin/add-post', {
			method: 'POST',
			body: JSON.stringify({ "title": title })
		});

		if (res.ok) {
            let slug = await res.text();
			goto('/admin/edit/' + slug);
		} else {
			error = 'Failed to create post';
		}
	}
</script>

<form onsubmit={(e) => { e.preventDefault(); createPost(); }}>
    <input type="text" placeholder="Title" bind:value={title} />
    <button type="submit">Create Post</button>
    {#if error}
        <p style="color: red;">{error}</p>
    {/if}
</form>