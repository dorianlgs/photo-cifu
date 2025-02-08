<script lang="ts">
	import { PUBLIC_POCKETBASE_URL } from '$env/static/public';
	import { getContext, onMount } from 'svelte';
	import type { Writable } from 'svelte/store';
	import { pb } from '$lib/pocketbase';
	import type { RecordModel } from 'pocketbase';

	let adminSection: Writable<string> = getContext('adminSection');
	adminSection.set('home');

	let galleries: RecordModel[] = $state([]);

	onMount(async () => {
		const _galleries = await pb.collection('galleries').getFullList({
			sort: '-created'
		});

		galleries = _galleries;
	});
</script>

<svelte:head>
	<title>Client Galleries</title>
</svelte:head>

<h1 class="text-2xl font-bold mb-1">Client Galleries</h1>
<a href="/account/gallery/new"
	><button class="btn btn-outline btn-primary mt-3 btn-wide">New</button></a
>

<div class="overflow-x-auto">
	<table class="table">
		<thead>
			<tr>
				<th>
					<label>
						<input type="checkbox" class="checkbox" />
					</label>
				</th>
				<th>Name</th>
				<th>Job</th>
				<th>Favorite Color</th>
				<th></th>
			</tr>
		</thead>
		<tbody>
			{#each galleries as gallery}
				<tr>
					<th>
						<label>
							<input type="checkbox" class="checkbox" />
						</label>
					</th>
					<td>
						<div class="flex items-center gap-3">
							<div class="avatar">
								<div class="mask mask-squircle h-12 w-12">
									<img
										src="{PUBLIC_POCKETBASE_URL}/api/files/{gallery.collectionId}/{gallery.id}/{gallery.thumbnail}"
										alt="Avatar Tailwind CSS Component"
									/>
								</div>
							</div>
							<div>
								<div class="font-bold">{gallery.name}</div>
								<div class="text-sm opacity-50">{gallery.location}</div>
							</div>
						</div>
					</td>
					<td>
						Zemlak, Daniel and Leannon
						<br />
						<span class="badge badge-ghost badge-sm">Desktop Support Technician</span>
					</td>
					<td>Purple</td>
					<th>
						<a href="/account/gallery/{gallery.id}" class="btn btn-ghost btn-xs">details</a>
					</th>
				</tr>
			{/each}
		</tbody>
		<!-- foot -->
		<tfoot>
			<tr>
				<th></th>
				<th>Name</th>
				<th>Job</th>
				<th>Favorite Color</th>
				<th></th>
			</tr>
		</tfoot>
	</table>
</div>
