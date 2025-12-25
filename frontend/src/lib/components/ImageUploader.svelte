<script lang="ts">
	import { uploadImagesBackend } from '$lib/api/images';
	import { toast } from 'svelte-sonner';

	let { onFileUpload } = $props<{ onFileUpload: (imageUrls: string[]) => void }>();

	let isDragging = $state(false);
	let isUploading = $state(false);
	let selectedFiles: File[] = $state([]);
	let imagePreviews: string[] = $state([]);

	function handleDragOver(e: DragEvent) {
		e.preventDefault();
		isDragging = true;
	}

	function handleDragLeave(e: DragEvent) {
		e.preventDefault();
		isDragging = false;
	}

	function handleDrop(e: DragEvent) {
		e.preventDefault();
		isDragging = false;

		if (e.dataTransfer?.files) {
			addFiles(Array.from(e.dataTransfer.files));
		}
	}

	function handleFileSelect(e: Event) {
		const target = e.target as HTMLInputElement;
		if (target.files) {
			addFiles(Array.from(target.files));
		}
	}

	function createPreview(file: File): Promise<string> {
		return new Promise((resolve) => {
			const reader = new FileReader();
			reader.onload = (e) => {
				resolve(e.target?.result as string);
			};
			reader.readAsDataURL(file);
		});
	}

	async function addFiles(files: File[]) {
		const validFiles = files.filter((file) => {
			if (!file.type.startsWith('image/')) {
				toast.error(`${file.name} is not an image file`);
				return false;
			}
			if (file.size > 30 * 1024 * 1024) {
				toast.error(`${file.name} is too large (max 30MB)`);
				return false;
			}
			return true;
		});

		const newPreviews = await Promise.all(validFiles.map((file) => createPreview(file)));

		selectedFiles = [...selectedFiles, ...validFiles];
		imagePreviews = [...imagePreviews, ...newPreviews];
	}

	function removeFile(index: number) {
		selectedFiles = selectedFiles.filter((_, i) => i !== index);
		imagePreviews = imagePreviews.filter((_, i) => i !== index);
	}

	function clearFiles() {
		selectedFiles = [];
		imagePreviews = [];
	}

	async function uploadImages() {
		if (selectedFiles.length === 0) {
			toast.error('Please select images to upload');
			return;
		}

		isUploading = true;
		try {
			let image_urls = await uploadImagesBackend(selectedFiles);

			if (image_urls && image_urls.length > 0) {
				toast.success(`Successfully uploaded ${image_urls.length} images`);
				onFileUpload(image_urls);
				clearFiles();
			} else {
				toast.error('No images were uploaded');
			}
		} catch (error) {
			console.error('Upload error:', error);
			toast.error('Upload failed');
		} finally {
			isUploading = false;
		}
	}

	function formatFileSize(bytes: number): string {
		if (bytes === 0) return '0 Bytes';
		const k = 1024;
		const sizes = ['Bytes', 'KB', 'MB', 'GB'];
		const i = Math.floor(Math.log(bytes) / Math.log(k));
		return parseFloat((bytes / Math.pow(k, i)).toFixed(2)) + ' ' + sizes[i];
	}
</script>

<div class="w-full">
	<!-- Upload Area -->
	<div
		class="border-2 border-dashed rounded-lg p-8 text-center transition-colors {isDragging
			? 'border-primary bg-primary/5'
			: 'border-base-300 hover:border-base-400'}"
		role="button"
		tabindex="0"
		ondragover={handleDragOver}
		ondragleave={handleDragLeave}
		ondrop={handleDrop}
	>
		<svg
			class="w-12 h-12 mx-auto mb-4 text-base-content/50"
			fill="none"
			stroke="currentColor"
			viewBox="0 0 24 24"
		>
			<path
				stroke-linecap="round"
				stroke-linejoin="round"
				stroke-width="2"
				d="M7 16a4 4 0 01-.88-7.903A5 5 0 1115.9 6L16 6a5 5 0 011 9.9M15 13l-3-3m0 0l-3 3m3-3v12"
			/>
		</svg>

		<div class="mb-4">
			<p class="text-lg font-medium mb-2">Drop images here or click to browse</p>
			<p class="text-sm text-base-content/60">
				Support for JPG, PNG, GIF, WebP • Max 30MB per file
			</p>
		</div>

		<input
			type="file"
			accept="image/*"
			multiple
			onchange={handleFileSelect}
			class="hidden"
			id="file-input"
		/>

		<label for="file-input" class="btn btn-primary cursor-pointer"> Choose Images </label>
	</div>

	<!-- File List -->
	{#if selectedFiles.length > 0}
		<div class="mt-6">
			<div class="flex items-center justify-between mb-4">
				<h3 class="text-lg font-medium">
					{selectedFiles.length}
					{selectedFiles.length === 1 ? 'Image' : 'Images'} Selected
				</h3>
				<button class="btn btn-ghost btn-sm" onclick={clearFiles}> Clear All </button>
			</div>

			<!-- Image Grid Table -->
			<div class="grid grid-cols-3 gap-4">
				{#each selectedFiles as file, index (index)}
					<div class="card bg-base-100 w-96 shadow-sm">
						<figure>
							<div class="aspect-3/2 overflow-hidden justify-center items-center flex">
								<img class="h-full display-block" src={imagePreviews[index]} alt={file.name} />
							</div>
						</figure>
						<div class="card-body">
							<div class="flex justify-between">
								<div>
									<p class="text-sm font-medium truncate max-w-xs">{file.name}</p>
									<p class="text-sm text-base-content/60">{formatFileSize(file.size)}</p>
								</div>
								<button
									class="btn btn-ghost btn-sm btn-circle"
									onclick={() => removeFile(index)}
									disabled={isUploading}
								>
									<svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
										<path
											stroke-linecap="round"
											stroke-linejoin="round"
											stroke-width="2"
											d="M6 18L18 6M6 6l12 12"
										/>
									</svg>
								</button>
							</div>
						</div>
					</div>
				{/each}
			</div>

			<!-- Upload Button -->
			<div class="mt-6 flex justify-end">
				<button
					class="btn btn-primary"
					onclick={uploadImages}
					disabled={isUploading || selectedFiles.length === 0}
				>
					{#if isUploading}
						<span class="loading loading-spinner loading-sm"></span>
						Uploading...
					{:else}
						Upload {selectedFiles.length} {selectedFiles.length === 1 ? 'Image' : 'Images'}
					{/if}
				</button>
			</div>
		</div>
	{/if}
</div>
