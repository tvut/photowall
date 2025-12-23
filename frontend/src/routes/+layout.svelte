<script lang="ts">
	import { Toaster } from 'svelte-sonner';
</script>

<slot />
<Toaster position="bottom-right" richColors closeButton />
