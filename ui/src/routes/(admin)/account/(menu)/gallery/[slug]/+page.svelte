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
		let _images = gallery?.expand?.images;
		images = _images.sort((a: Image, b: Image) => b.likes - a.likes);

		unsubscribe = await pb.collection('images').subscribe('*', async ({ action, record }) => {
			if (action === 'update') {
				var foundIndex = images.findIndex((x) => x.id == record.id);
				images[foundIndex] = { ...images[foundIndex], likes: record.likes };
				images = _images.sort((a: Image, b: Image) => b.likes - a.likes);
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
	<div class="row">
		{#each images as image, i}
			<div class="column">
				<div class="container">
					<img
						src="{PUBLIC_POCKETBASE_URL}/api/files/{image.collectionId}/{image.id}/{image.image}"
						alt="Norway"
						style="width:100%;"
					/>
					<div class="top-right">
						<LikeButton image={images[i]} />
					</div>
				</div>
			</div>
		{/each}
	</div>
	<a href="/account"><button class="btn btn-outline btn-primary mt-3 btn-wide">Back</button></a>
{/if}

<style>
	/* Container holding the image and the text */
	.container {
		position: relative;
	}

	.top-right {
		position: absolute;
		top: 8px;
		right: 16px;
	}

	.row {
		display: -ms-flexbox; /* IE10 */
		display: flex;
		-ms-flex-wrap: wrap; /* IE10 */
		flex-wrap: wrap;
		padding: 0 4px;
	}

	/* Create four equal columns that sits next to each other */
	.column {
		-ms-flex: 25%; /* IE10 */
		flex: 25%;
		max-width: 25%;
		padding: 0 4px;
	}

	.column img {
		margin-top: 8px;
		vertical-align: middle;
		width: 100%;
	}

	/* Responsive layout - makes a two column-layout instead of four columns */
	@media screen and (max-width: 800px) {
		.column {
			-ms-flex: 50%;
			flex: 50%;
			max-width: 50%;
		}
	}

	/* Responsive layout - makes the two columns stack on top of each other instead of next to each other */
	@media screen and (max-width: 600px) {
		.column {
			-ms-flex: 100%;
			flex: 100%;
			max-width: 100%;
		}
	}
</style>
