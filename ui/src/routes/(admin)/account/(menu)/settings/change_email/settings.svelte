<script lang="ts">
	import { pb } from '$lib/pocketbase';
	import { currentUser } from '$lib/stores/user';
	import { onMount } from 'svelte';

	let title: string = $state('Change Email');
	let loading: boolean = $state(false);
	let showSuccess: boolean = $state(false);
	let message: string = $state('mensaje');
	let dangerous: boolean = $state(false);
	let errors: { [fieldName: string]: string } = $state({});
	let newEmailInput: HTMLInputElement | undefined = $state();


	const handleSubmit = async (e: SubmitEvent) => {
	e.preventDefault();
	errors = {};


	const formData = new FormData(e.target as HTMLFormElement);

	const newEmail = formData.get('newEmail')?.toString() ?? '';

	if (newEmail.length < 6) {
		errors['newEmail'] = 'Email is required';
	} else if (newEmail.length > 500) {
		errors['newEmail'] = 'Email too long';
	} else if (!newEmail.includes('@') || !newEmail.includes('.')) {
		errors['newEmail'] = 'Invalid email';
	}

	if (Object.keys(errors).length > 0) {
		return;
	}


	loading = true;
	await pb.collection('users').requestEmailChange(newEmail);
	loading = false;
	showSuccess = true

	}

	onMount(() => {
		if (newEmailInput) newEmailInput.focus();
	});

</script>

<div class="card p-6 pb-7 mt-8 max-w-xl flex flex-col md:flex-row shadow">
	{#if title}
		<div class="text-xl font-bold mb-3 w-48 md:pr-8 flex-none">{title}</div>
	{/if}

	<div class="w-full min-w-48">
		{#if !showSuccess}
			{#if message}
				<div class="mb-6 {dangerous ? 'alert alert-warning' : ''}">
					{#if dangerous}
						<svg
							xmlns="http://www.w3.org/2000/svg"
							class="stroke-current shrink-0 h-6 w-6"
							fill="none"
							viewBox="0 0 24 24"
							><path
								stroke-linecap="round"
								stroke-linejoin="round"
								stroke-width="2"
								d="M12 9v2m0 4h.01m-6.938 4h13.856c1.54 0 2.502-1.667 1.732-3L13.732 4c-.77-1.333-2.694-1.333-3.464 0L3.34 16c-.77 1.333.192 3 1.732 3z"
							/></svg
						>
						<span>{message}</span>
					{/if}
				</div>
			{/if}
			<form class="form-widget flex flex-col" onsubmit={handleSubmit}>
			<label for={'currentEmail'}>
				<span class="text-sm text-gray-500">{'Current Email'}</span>
			</label>
			<input
				id={"currentEmail"}
				name={"currentEmail"}
				type={'email'}
				placeholder={"CurrentEmail"}
				class="input-sm mt-1 input input-bordered w-full max-w-xs mb-3 text-base py-4"
				value={$currentUser?.email}
				disabled={true}
			/>
				<label for={'newEmail'}>
					<span class="text-sm text-gray-500">{'New Email'}</span>
					{#if errors['newEmail']}
						<div class="text-red-600 flex-grow text-sm ml-2 text-right">
							{errors['newEmail']}
						</div>
					{/if}
				</label>
				<input
				    bind:this={newEmailInput}
					id={"newEmail"}
					name={"newEmail"}
					type={'text'}
					placeholder={"New Email"}
					class="input-sm mt-1 input input-bordered w-full max-w-xs mb-3 text-base py-4"
				/>
				<button
					type="submit"
					class="ml-auto btn btn-sm mt-3 min-w-[145px] {dangerous
						? 'btn-error'
						: 'btn-primary btn-outline'}"
					disabled={loading}
				>
					{#if loading}
						<span class="loading loading-spinner loading-md align-middle mx-3"></span>
					{:else}
						{'Change'}
					{/if}
				</button>
			</form>
		{:else}
			<div>
				<div class="text-base">{'Check your inbox in order to verify your new email'}</div>
			</div>
			<a href="/account/settings">
				<button class="btn btn-outline btn-sm mt-3 min-w-[145px]"> Return to Settings </button>
			</a>
		{/if}
	</div>
</div>
