<script lang="ts">
	import { ClientResponseError } from 'pocketbase';
	import { pb } from '$lib/pocketbase';
	import { goto } from '$app/navigation';
	import { page } from '$app/state';
	import { onMount } from 'svelte';

	let errors: { [fieldName: string]: string } = $state({});
	let loading = $state(false);
	let emailInput: HTMLInputElement | undefined = $state();

	const handleSubmit = async (e: SubmitEvent) => {
		const formData = new FormData(e.target as HTMLFormElement);

		const email = formData.get('email')?.toString() ?? '';
		if (email.length < 6) {
			errors['email'] = 'Email is required';
		} else if (email.length > 500) {
			errors['email'] = 'Email too long';
		} else if (!email.includes('@') || !email.includes('.')) {
			errors['email'] = 'Invalid email';
		}
		const password = formData.get('password')?.toString() ?? '';
		if (password.length > 500) {
			errors['password'] = 'Password too long';
		}

		try {
			await pb.collection('users').authWithPassword(email, password);
			if (!pb?.authStore?.record) {
				pb.authStore.clear();
				return {
					notVerified: true
				};
			}
		} catch (err: any) {
			if (err instanceof ClientResponseError) {
				if (err.response.message === 'Failed to authenticate.') {
					errors['loginResult'] = 'The Email or Password is Incorrect. Try again.';
				}
			}

			return;
		}

		goto('/account');
	};

	onMount(() => {
		if (emailInput) emailInput.focus();
	});
</script>

<svelte:head>
	<title>Sign in</title>
</svelte:head>

{#if page.url.searchParams.get('verified') === 'true'}
	<div role="alert" class="alert alert-success mb-5">
		<svg
			xmlns="http://www.w3.org/2000/svg"
			class="stroke-current shrink-0 h-6 w-6"
			fill="none"
			viewBox="0 0 24 24"
			><path
				stroke-linecap="round"
				stroke-linejoin="round"
				stroke-width="2"
				d="M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z"
			/></svg
		>
		<span>Email verified! Please sign in.</span>
	</div>
{/if}
<h1 class="text-2xl font-bold mb-6">Sign In</h1>
<button
	aria-label="Sign in with Github"
	class="btn btn-github {loading ? 'btn-disabled' : ''}"
	onclick={async () => {
		try {
			await pb.collection('users').authWithOAuth2({ provider: 'github' });
			goto('/account');
		} catch (err) {}
	}}
>
	<svg
		fill="gray"
		xmlns="http://www.w3.org/2000/svg"
		viewBox="0 0 30 30"
		width="21px"
		height="21px"
		class="svelte-10a6av0"
	>
		<path
			d="M15,3C8.373,3,3,8.373,3,15c0,5.623,3.872,10.328,9.092,11.63C12.036,26.468,12,26.28,12,26.047v-2.051 c-0.487,0-1.303,0-1.508,0c-0.821,0-1.551-0.353-1.905-1.009c-0.393-0.729-0.461-1.844-1.435-2.526 c-0.289-0.227-0.069-0.486,0.264-0.451c0.615,0.174,1.125,0.596,1.605,1.222c0.478,0.627,0.703,0.769,1.596,0.769 c0.433,0,1.081-0.025,1.691-0.121c0.328-0.833,0.895-1.6,1.588-1.962c-3.996-0.411-5.903-2.399-5.903-5.098 c0-1.162,0.495-2.286,1.336-3.233C9.053,10.647,8.706,8.73,9.435,8c1.798,0,2.885,1.166,3.146,1.481C13.477,9.174,14.461,9,15.495,9 c1.036,0,2.024,0.174,2.922,0.483C18.675,9.17,19.763,8,21.565,8c0.732,0.731,0.381,2.656,0.102,3.594 c0.836,0.945,1.328,2.066,1.328,3.226c0,2.697-1.904,4.684-5.894,5.097C18.199,20.49,19,22.1,19,23.313v2.734 c0,0.104-0.023,0.179-0.035,0.268C23.641,24.676,27,20.236,27,15C27,8.373,21.627,3,15,3z"
		></path>
	</svg>
</button>
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
	<label for={'password'}>
		<div class="flex flex-row">
			<div class="text-base font-bold">{'Your password'}</div>
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

	{#if Object.keys(errors).length > 0}
		{#if errors['loginResult']}
			<p class="text-red-600 text-sm mb-2">{errors['loginResult']}</p>
		{:else}
			<p class="text-red-600 text-sm mb-2">Please resolve above issues.</p>
		{/if}
	{/if}

	<button class="btn btn-primary {loading ? 'btn-disabled' : ''}">Sign in</button>
</form>

<div class="text-l text-slate-800 mt-4">
	<a class="underline" href="/login/forgot_password">Forgot password?</a>
</div>
<div class="text-l text-slate-800 mt-3">
	Don't have an account? <a class="underline" href="/login/sign_up">Sign up</a>.
</div>

<style>
	.btn-github {
		background-color: white;
		border-color: rgb(224, 224, 224);
	}
</style>
