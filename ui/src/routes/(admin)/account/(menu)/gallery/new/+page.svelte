<script lang="ts">
	import { pb } from '$lib/pocketbase';
	import InputFile from '$lib/components/InputFile.svelte';
	import { goto } from '$app/navigation';
	import { onMount } from 'svelte';

	let errors: { [fieldName: string]: string } = $state({});
	let loading = $state(false);
	let nameInput: HTMLInputElement | undefined = $state();

	const handleSubmit = async (e: SubmitEvent) => {
		try {
			loading = true;
			e.preventDefault();
			errors = {};

			const formData = new FormData(e.target as HTMLFormElement);

			const name = formData.get('name')?.toString() ?? '';
			if (name.length < 2) {
				errors['name'] = 'Name is required';
			}
			if (name.length > 500) {
				errors['name'] = 'Name too long';
			}

			const location = formData.get('location')?.toString() ?? '';
			if (location.length < 2) {
				errors['location'] = 'Location is required';
			}
			if (location.length > 500) {
				errors['location'] = 'Location too long';
			}

			const imagesZip = formData.get('imagesZip') as File;

			if (imagesZip?.size === 0) {
				errors['imagesZip'] = 'Images Zip is required';
			}

			const thumbnail = formData.get('thumbnail') as File;

			if (thumbnail?.size === 0) {
				errors['thumbnail'] = 'Thumbnail is required';
			}

			if (Object.keys(errors).length > 0) {
				return;
			}

			loading = true;
			const result = await pb.send('/api/photocifu/gallery/create', {
				method: 'POST',
				body: formData
			});
			loading = false;

			goto(`/account/gallery/${result.galleryId}`);
		} catch (err) {
			loading = false;
		}
	};

	onMount(() => {
		if (nameInput) nameInput.focus();
	});
</script>

<svelte:head>
	<title>New Gallery</title>
</svelte:head>

<h1 class="text-2xl font-bold mb-6">New Gallery</h1>
<div class="max-w-lg min-h-[70vh] pb-12 flex">
	<div class="">
		<form class="form-widget flex flex-col" onsubmit={handleSubmit}>
			<label for={'name'}>
				<div class="flex flex-row">
					<div class="text-base font-bold">{'Name'}</div>
					{#if errors['name']}
						<div class="text-red-600 flex-grow text-sm ml-2 text-right">
							{errors['name']}
						</div>
					{/if}
				</div>
				<input
					bind:this={nameInput}
					id={'name'}
					name={'name'}
					type={'name'}
					placeholder={'Gallery name'}
					class="{errors['name']
						? 'input-error'
						: ''} input-md mt-1 input input-bordered w-full mb-3 text-base py-4"
				/>
			</label>
			<label for={'location'}>
				<div class="flex flex-row">
					<div class="text-base font-bold">{'Location'}</div>
					{#if errors['location']}
						<div class="text-red-600 flex-grow text-sm ml-2 text-right">
							{errors['location']}
						</div>
					{/if}
				</div>
				<input
					id={'location'}
					name={'location'}
					type={'location'}
					placeholder={'Location'}
					class="{errors['location']
						? 'input-error'
						: ''} input-md mt-1 input input-bordered w-full mb-3 text-base py-4"
				/>
			</label>
			<label for={'imagesZip'}>
				<div class="flex flex-row">
					<div class="text-base font-bold">{'Images Zip'}</div>
					{#if errors['imagesZip']}
						<div class="text-red-600 flex-grow text-sm ml-2 text-right">
							{errors['imagesZip']}
						</div>
					{/if}
				</div>
				<InputFile {errors} name={'imagesZip'} placeholder={'Upload your zip images'} />
				<br />
			</label>

			<label for={'thumbnail'}>
				<div class="flex flex-row">
					<div class="text-base font-bold">{'Thumbnail'}</div>
					{#if errors['thumbnail']}
						<div class="text-red-600 flex-grow text-sm ml-2 text-right">
							{errors['thumbnail']}
						</div>
					{/if}
				</div>
				<InputFile {errors} name={'thumbnail'} placeholder={'Upload your gallery thumbnail'} />
				<br />
			</label>

			{#if Object.keys(errors).length > 0}
				{#if errors['submitResult']}
					<p class="text-red-600 text-sm mb-2">{errors['submitResult']}</p>
				{:else}
					<p class="text-red-600 text-sm mb-2">Please resolve above issues.</p>
				{/if}
			{/if}

			<button
				disabled={loading}
				type={'submit'}
				class="btn btn-primary {loading ? 'btn-disabled' : ''}">Create</button
			>
		</form>
	</div>
</div>
