<script lang="ts">
	import { PUBLIC_POCKETBASE_URL } from '$env/static/public';
	import { pb } from '$lib/pocketbase';
	import type { RecordModel } from 'pocketbase';
	import { onMount } from 'svelte';

	import { page } from '$app/state';

	let { slug: galleryId } = page.params;

	let gallery: RecordModel | undefined = $state();
	let images: string[] = $state([]);

	onMount(async () => {
		const _gallery = await pb.collection('galleries').getOne(galleryId, { expand: 'images' });
		gallery = _gallery;
		images = gallery?.expand?.images;
	});
</script>

<svelte:head>
	<title>{gallery ? gallery.name : ''}</title>
</svelte:head>

{#if gallery}
	<div class="carousel rounded-box">
		{#each images as item}
			<div class="carousel-item">
				<img
					src="{PUBLIC_POCKETBASE_URL}/api/files/{item.collectionId}/{item.id}/{item.image}"
					alt="Burger"
					height="400px"
					width="450px"
				/>
			</div>
		{/each}
	</div>
{/if}

<style>
	img {
		padding: 20px;
	}
</style>
