<script lang="ts">
	import { page } from '$app/state';
	import { enhance, applyAction } from '$app/forms';
	import { redirect, type SubmitFunction } from '@sveltejs/kit';
	import type { FullAutoFill } from 'svelte/elements';

	let errors: { [fieldName: string]: string } = $state({});
	let loading = $state(false);

	const handleSubmit: SubmitFunction = () => {
		loading = true;
		errors = {};
		return async ({ update, result }) => {
			await update({ reset: false });
			await applyAction(result);
			loading = false;

			if (result.type === 'success') {
				redirect(303, '/account');
			} else if (result.type === 'failure') {
				errors = result.data?.errors ?? {};
			} else if (result.type === 'error') {
				errors = { _: 'An error occurred. Please check inputs and try again.' };
			}
		};
	};

	interface FormField {
		id: string;
		label: string;
		inputType: string;
		autocomplete: FullAutoFill;
		placeholder: string;
	}

	const formFields: FormField[] = [
		{
			id: 'email',
			label: 'Email address',
			inputType: 'email',
			autocomplete: 'email',
			placeholder: 'Your email address'
		},
		{
			id: 'password',
			label: 'Your Password',
			inputType: 'password',
			autocomplete: 'off',
			placeholder: 'Your Password'
		}
	];
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

<form class="form-widget flex flex-col" method="POST" action="?/signIn" use:enhance={handleSubmit}>
	{#each formFields as field}
		<label for={field.id}>
			<div class="flex flex-row">
				<div class="text-base font-bold">{field.label}</div>
				{#if errors[field.id]}
					<div class="text-red-600 flex-grow text-sm ml-2 text-right">
						{errors[field.id]}
					</div>
				{/if}
			</div>
			<input
				id={field.id}
				name={field.id}
				type={field.inputType}
				autocomplete={field.autocomplete}
				placeholder={field.placeholder}
				class="{errors[field.id]
					? 'input-error'
					: ''} input-sm mt-1 input input-bordered w-full mb-3 text-base py-4"
			/>
		</label>
	{/each}

	{#if Object.keys(errors).length > 0}
		<p class="text-red-600 text-sm mb-2">Please resolve above issues.</p>
	{/if}

	<button class="btn btn-primary {loading ? 'btn-disabled' : ''}">Sign in</button>
</form>

<div class="text-l text-slate-800 mt-4">
	<a class="underline" href="/login/forgot_password">Forgot password?</a>
</div>
<div class="text-l text-slate-800 mt-3">
	Don't have an account? <a class="underline" href="/login/sign_up">Sign up</a>.
</div>
