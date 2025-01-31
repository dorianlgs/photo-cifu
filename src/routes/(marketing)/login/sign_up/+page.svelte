<script lang="ts">
	import { pb } from '$lib/pocketbase';
	import GitHubButton from '$lib/components/GitHubButton.svelte';
	import InputFile from '$lib/components/InputFile.svelte';
	import { goto } from '$app/navigation';
	import { onMount } from 'svelte';

	let errors: { [fieldName: string]: string } = $state({});
	let loading = $state(false);
	let emailInput: HTMLInputElement | undefined = $state();

	const handleSubmit = async (e: SubmitEvent) => {
		e.preventDefault();

		errors = {};

		const formData = new FormData(e.target as HTMLFormElement);

		await pb.collection('users').create({
			email: formData.get('email'),
			name: formData.get('name'),
			password: formData.get('password'),
			passwordConfirm: formData.get('password'),
			avatar: formData.get('avatar')
		});

		goto(`/login/sign_in?not_verified=true`);
	};

	onMount(() => {
		if (emailInput) emailInput.focus();
	});
</script>

<svelte:head>
	<title>Sign up</title>
</svelte:head>

<h1 class="text-2xl font-bold mb-6">Sign Up</h1>
<GitHubButton />
<br />
<hr class="solid" />
<br />
<form class="form-widget flex flex-col" onsubmit={handleSubmit}>
	<label for={'email'}>
		<div class="flex flex-row">
			<div class="text-base font-bold">{'Email address'}</div>
			{#if errors['email']}
				<div class="text-red-600 flex-grow text-sm ml-2 text-right">
					{errors['email']}
				</div>
			{/if}
		</div>
		<input
			bind:this={emailInput}
			id={'email'}
			name={'email'}
			type={'email'}
			autocomplete={'email'}
			placeholder={'Your email address'}
			class="{errors['email']
				? 'input-error'
				: ''} input-md mt-1 input input-bordered w-full mb-3 text-base py-4"
		/>
	</label>
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
			id={'name'}
			name={'name'}
			autocomplete={'off'}
			placeholder={'Your name'}
			class="{errors['name']
				? 'input-error'
				: ''} input-md mt-1 input input-bordered w-full mb-3 text-base py-4"
		/>
	</label>
	<label for={'password'}>
		<div class="flex flex-row">
			<div class="text-base font-bold">{'Create a Password'}</div>
			{#if errors['password']}
				<div class="text-red-600 flex-grow text-sm ml-2 text-right">
					{errors['password']}
				</div>
			{/if}
		</div>
		<input
			id={'password'}
			name={'password'}
			type={'password'}
			autocomplete={'off'}
			placeholder={'Your password'}
			class="{errors['email']
				? 'input-error'
				: ''} input-md mt-1 input input-bordered w-full mb-3 text-base py-4"
		/>
	</label>
	<label for={'avatar'}>
		<div class="flex flex-row">
			<div class="text-base font-bold">{'Avatar'}</div>
			{#if errors['avatar']}
				<div class="text-red-600 flex-grow text-sm ml-2 text-right">
					{errors['avatar']}
				</div>
			{/if}
		</div>
		<InputFile {errors} />

		<br />
	</label>

	{#if Object.keys(errors).length > 0}
		{#if errors['loginResult']}
			<p class="text-red-600 text-sm mb-2">{errors['loginResult']}</p>
		{:else}
			<p class="text-red-600 text-sm mb-2">Please resolve above issues.</p>
		{/if}
	{/if}

	<button type={'submit'} class="btn btn-primary {loading ? 'btn-disabled' : ''}">Sign up</button>
</form>
<div class="text-l text-slate-800 mt-4 mb-2">
	Have an account? <a class="underline" href="/login/sign_in">Sign in</a>
</div>
