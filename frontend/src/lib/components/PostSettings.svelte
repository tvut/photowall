<script lang="ts">
  import { toast } from 'svelte-sonner';
  import { updatePostStatus, updatePostDisplayTime, deletePost } from '$lib/api/posts';
  import { goto } from '$app/navigation';
  
  const { 
    post
  } = $props<{
    post: {
      slug: string;
      title: string;
      status: string;
      displayTime: string;
    };
  }>();
  
  let status = $state(post.status);
  let displayTime = $state(post.displayTime);
  let isDeleting = $state(false);
  
  $effect(() => {
    status = post.status;
    displayTime = post.displayTime;
  });
  
  async function updateStatus(newStatus: string) {
    try {
      await updatePostStatus(post.slug, newStatus);
      toast.success(`Status updated to ${newStatus}`);
    } catch (error) {
      console.error('Error updating status:', error);
      toast.error('Failed to update status');
    }
  }
  
  async function updateDisplayTime(newDisplayTime: string) {
    try {
      await updatePostDisplayTime(post.slug, newDisplayTime);
      toast.success('Display time updated');
    } catch (error) {
      console.error('Error updating display time:', error);
      toast.error('Failed to update display time');
    }
  }

  async function handleDelete() {
    if (confirm(`Are you sure you want to delete "${post.title}"?`)) {
      try {
        isDeleting = true;
        await deletePost(post.slug);
        toast.success('Post deleted successfully');
        await goto('/admin');
      } catch (error) {
        console.error('Error deleting post:', error);
        toast.error('Failed to delete post');
      } finally {
        isDeleting = false;
      }
    }
  }
</script>

<div class="card bg-base-100 shadow-xl">
  <div class="card-body">
    <h2 class="card-title text-xl mb-4">Post Settings</h2>
    <div class="space-y-6">
      <div class="form-control w-full">
        <label for="status-select" class="label">
          <span class="label-text font-semibold">Status</span>
        </label>
        <select
          id="status-select"
          class="select select-bordered w-full"
          bind:value={status}
          onchange={() => updateStatus(status)}
        >
          <option value="draft" class={status === 'draft' ? 'bg-primary-content' : ''}
            >Draft</option
          >
          <option value="published" class={status === 'published' ? 'bg-primary-content' : ''}
            >Published</option
          >
        </select>
      </div>
      <div class="form-control w-full">
        <label for="display-time" class="label">
          <span class="label-text font-semibold">Display Time</span>
        </label>
        <div class="join w-full">
          <input
            id="display-time"
            type="datetime-local"
            class="input input-bordered join-item flex-1"
            bind:value={displayTime}
          />
          <button
            class="btn btn-primary join-item"
            onclick={() => updateDisplayTime(displayTime)}
            disabled={!displayTime}
          >
            Update
          </button>
        </div>
      </div>
      <div class="space-x-2">
        <a href="/posts/{post.slug}" class="btn btn-outline w-40" target="_blank">
          <svg
            xmlns="http://www.w3.org/2000/svg"
            class="h-5 w-5 mr-2"
            fill="none"
            viewBox="0 0 24 24"
            stroke="currentColor"
          >
            <path
              stroke-linecap="round"
              stroke-linejoin="round"
              stroke-width="2"
              d="M10 6H6a2 2 0 00-2 2v10a2 2 0 002 2h10a2 2 0 002-2v-4M14 4h6m0 0v6m0-6L10 14"
            />
          </svg>
          View Post
        </a>
        <button
          class="btn btn-error w-40"
          onclick={handleDelete}
          disabled={isDeleting}
        >
          {isDeleting ? 'Deleting...' : 'Delete Post'}
        </button>
      </div>
    </div>
  </div>
</div>
