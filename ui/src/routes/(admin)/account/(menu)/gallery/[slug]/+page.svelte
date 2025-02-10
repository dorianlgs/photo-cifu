<script lang="ts">
	import { PUBLIC_POCKETBASE_URL } from '$env/static/public';
	import { pb } from '$lib/pocketbase';
	import type { RecordModel } from 'pocketbase';
	import { onMount, onDestroy } from 'svelte';
	import { page } from '$app/state';
	import LikeButton from '$lib/components/LikeButton.svelte';

	interface Image {
		collectionId: string;
		id: string;
		image: string;
		likes: number;
	}

	let { slug: galleryId } = page.params;

	let gallery: RecordModel | undefined = $state();
	let images: Image[] = $state([]);
	let unsubscribe: () => void;

	onMount(async () => {
		const _gallery = await pb.collection('galleries').getOne(galleryId, { expand: 'images' });
		gallery = _gallery;
		images = gallery?.expand?.images;

		unsubscribe = await pb.collection('images').subscribe('*', async ({ action, record }) => {
			if (action === 'update') {
				var foundIndex = images.findIndex((x) => x.id == record.id);
				images[foundIndex] = { ...images[foundIndex], likes: record.likes };
			}
		});
	});

	onDestroy(() => {
		unsubscribe?.();
	});
</script>

<svelte:head>
	<title>{gallery ? gallery.name : ''}</title>
</svelte:head>

{#if gallery}
	<div class="carousel rounded-box">
		{#each images as image, i}
			<div>
				<div class="carousel-item">
					<img
						src="{PUBLIC_POCKETBASE_URL}/api/files/{image.collectionId}/{image.id}/{image.image}"
						alt="Burger"
						height="400px"
						width="450px"
					/>
				</div>
				<LikeButton image={images[i]} />
			</div>
		{/each}
	</div>
	<a href="/account"><button class="btn btn-outline btn-primary mt-3 btn-wide">Back</button></a>
{/if}

<style>
	img {
		padding: 20px;
	}
</style>
