<script lang="ts">
	import { pb } from '$lib/pocketbase';
	import { goto } from '$app/navigation';
	import { onMount } from 'svelte';
	import { page } from '$app/state';

	let errors: { [fieldName: string]: string } = $state({});
	let loading = $state(false);
	let passwordInput: HTMLInputElement | undefined = $state();

	const handleSubmit = async (e: SubmitEvent) => {
		e.preventDefault();
		errors = {};

		const formData = new FormData(e.target as HTMLFormElement);

		const currentPassword = formData.get('currentPassword')?.toString() ?? '';
		if (currentPassword.length > 500) {
			errors['currentPassword'] = 'Current Password too long';
		}

		const emailChangeToken = page.url.searchParams.get('email_change_token') as string;

		if (!emailChangeToken) {
			errors['currentPassword'] = 'Empty token';
		}

		if (Object.keys(errors).length > 0) {
			return;
		}

		try {
			loading = true;
			await pb.collection('users').confirmEmailChange(emailChangeToken, currentPassword);
			loading = false;
			goto(`/login/sign_in?email_changed=true`);
		} catch (err: any) {
		  console.log({err})
		  	errors['changeResult'] = err.toString();
			loading = false;
		}
	};

	onMount(() => {
		if (passwordInput) passwordInput.focus();
	});
</script>

<svelte:head>
	<title>Change Email</title>
</svelte:head>
<h1 class="text-2xl font-bold mb-6">Change Email</h1>
<form class="form-widget flex flex-col" onsubmit={handleSubmit}>
	<label for={'currentPassword'}>
		<div class="flex flex-row">
			<div class="text-base font-bold">{'Type your password to confirm changing your email'}</div>
			{#if errors['currentPassword']}
				<div class="text-red-600 flex-grow text-sm ml-2 text-right">
					{errors['currentPassword']}
				</div>
			{/if}
		</div>
		<input
			bind:this={passwordInput}
			id={'currentPassword'}
			name={'currentPassword'}
			type={'password'}
			autocomplete={'off'}
			placeholder={'Password'}
			class="{errors['currentPassword']
				? 'input-error'
				: ''} input-md mt-1 input input-bordered w-full mb-3 text-base py-4"
		/>
	</label>
	{#if Object.keys(errors).length > 0}
		{#if errors['changeResult']}
			<p class="text-red-600 text-sm mb-2">{errors['changeResult']}</p>
		{:else}
			<p class="text-red-600 text-sm mb-2">Please resolve above issues.</p>
		{/if}
	{/if}

	<button disabled={loading} class="btn btn-primary {loading ? 'btn-disabled' : ''}"
		>Confirm New Password</button
	>
</form>
