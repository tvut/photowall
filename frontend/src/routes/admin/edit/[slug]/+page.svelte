<script lang="ts">
import type { PageProps } from './$types';

let { data }: PageProps = $props();

</script>

<h1>Edit Post: {data.post.title}</h1>
<pre>{JSON.stringify(data, null, 2)}</pre>  