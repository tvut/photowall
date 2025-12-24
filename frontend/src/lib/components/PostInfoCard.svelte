<script lang="ts">
    export let post: {
        slug: string;
        created_at: string;
        display_time: string | null;
    };
</script>

<div class="card bg-base-100 shadow-xl h-full">
    <div class="card-body">
        <h2 class="card-title text-2xl mb-4">Post Information</h2>
        <div class="space-y-4">
            <div>
                <h3 class="text-lg font-semibold">Slug</h3>
                <p class="text-gray-600 font-mono text-sm bg-base-200 p-2 rounded">
                    {post.slug}
                </p>
            </div>
            <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
                <div>
                    <h3 class="text-lg font-semibold">Created</h3>
                    <p class="text-gray-600">{new Date(post.created_at).toLocaleString()}</p>
                </div>
                <div>
                    <h3 class="text-lg font-semibold">Display Time</h3>
                    <p class="text-gray-600">
                        {post.display_time
                            ? new Date(post.display_time).toLocaleString()
                            : 'Never'}
                    </p>
                </div>
            </div>
        </div>
    </div>
</div>
