<script lang="ts">
	import type { FullAutoFill } from 'svelte/elements';

	let errors: { [fieldName: string]: string } = $state({});
	let loading = $state(false);
	let showSuccess = $state(false);

	interface FormField {
		id: string;
		label: string;
		inputType: string;
		autocomplete: FullAutoFill;
	}

	const formFields: FormField[] = [
		{
			id: 'first_name',
			label: 'First Name *',
			inputType: 'text',
			autocomplete: 'given-name'
		},
		{
			id: 'last_name',
			label: 'Last Name *',
			inputType: 'text',
			autocomplete: 'family-name'
		},
		{
			id: 'email',
			label: 'Email *',
			inputType: 'email',
			autocomplete: 'email'
		},
		{
			id: 'phone',
			label: 'Phone Number',
			inputType: 'tel',
			autocomplete: 'tel'
		},
		{
			id: 'company',
			label: 'Company Name',
			inputType: 'text',
			autocomplete: 'organization'
		},
		{
			id: 'message',
			label: 'Message',
			inputType: 'textarea',
			autocomplete: 'off'
		}
	];

	const handleSubmit = async (e: SubmitEvent) => {
		try {
			e.preventDefault();
			errors = {};

			const formData = new FormData(e.target as HTMLFormElement);

			const firstName = formData.get('first_name')?.toString() ?? '';
			if (firstName.length < 2) {
				errors['first_name'] = 'First name is required';
			}
			if (firstName.length > 500) {
				errors['first_name'] = 'First name too long';
			}

			const lastName = formData.get('last_name')?.toString() ?? '';
			if (lastName.length < 2) {
				errors['last_name'] = 'Last name is required';
			}
			if (lastName.length > 500) {
				errors['last_name'] = 'Last name too long';
			}

			const email = formData.get('email')?.toString() ?? '';
			if (email.length < 6) {
				errors['email'] = 'Email is required';
			} else if (email.length > 500) {
				errors['email'] = 'Email too long';
			} else if (!email.includes('@') || !email.includes('.')) {
				errors['email'] = 'Invalid email';
			}

			const company = formData.get('company')?.toString() ?? '';
			if (company.length > 500) {
				errors['company'] = 'Company too long';
			}

			const phone = formData.get('phone')?.toString() ?? '';
			if (phone.length > 100) {
				errors['phone'] = 'Phone number too long';
			}

			const message = formData.get('message')?.toString() ?? '';
			if (message.length > 2000) {
				errors['message'] = 'Message too long (' + message.length + ' of 2000)';
			}

			if (Object.keys(errors).length > 0) {
				return;
			}

			loading = true;
		} catch (err) {
			loading = false;
		}
	};
</script>

<div
	class="flex flex-col lg:flex-row mx-auto my-4 min-h-[70vh] place-items-center lg:place-items-start place-content-center"
>
	<div
		class="max-w-[400px] lg:max-w-[500px] flex flex-col place-content-center p-4 lg:mr-8 lg:mb-8 lg:min-h-[70vh]"
	>
		<div class="px-6">
			<h1 class="text-2xl lg:text-4xl font-bold mb-4">Contact Us</h1>
			<p class="text-lg">Talk to one of our service professionals to:</p>
			<ul class="list-disc list-outside pl-6 py-4 space-y-1">
				<li class="">Get a live demo</li>
				<li class="">Discuss your specific needs</li>
				<li>Get a quote</li>
				<li>Answer any technical questions you have</li>
			</ul>
			<p>Once you complete the form, we'll reach out to you! *</p>
			<p class="text-sm pt-8">
				*Not really for this demo page, but you should say something like that 😉
			</p>
		</div>
	</div>

	<div
		class="flex flex-col flex-grow m-4 lg:ml-10 min-w-[300px] stdphone:min-w-[360px] max-w-[400px] place-content-center lg:min-h-[70vh]"
	>
		{#if showSuccess}
			<div class="flex flex-col place-content-center lg:min-h-[70vh]">
				<div class="card card-bordered shadow-lg py-6 px-6 mx-2 lg:mx-0 lg:p-6 mb-10">
					<div class="text-2xl font-bold mb-4">Thank you!</div>
					<p class="">We've received your message and will be in touch soon.</p>
				</div>
			</div>
		{:else}
			<div class="card card-bordered shadow-lg p-4 pt-6 mx-2 lg:mx-0 lg:p-6">
				<form class="form-widget flex flex-col" onsubmit={handleSubmit}>
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
							{#if field.inputType === 'textarea'}
								<textarea
									id={field.id}
									name={field.id}
									autocomplete={field.autocomplete}
									rows={4}
									class="{errors[field.id]
										? 'input-error'
										: ''} h-24 input-sm mt-1 input input-bordered w-full mb-3 text-base py-4"
								></textarea>
							{:else}
								<input
									id={field.id}
									name={field.id}
									type={field.inputType}
									autocomplete={field.autocomplete}
									class="{errors[field.id]
										? 'input-error'
										: ''} input-sm mt-1 input input-bordered w-full mb-3 text-base py-4"
								/>
							{/if}
						</label>
					{/each}

					{#if Object.keys(errors).length > 0}
						<p class="text-red-600 text-sm mb-2">Please resolve above issues.</p>
					{/if}

					<button
						disabled={loading}
						type="submit"
						class="btn btn-primary {loading ? 'btn-disabled' : ''}"
						>{loading ? 'Submitting' : 'Submit'}</button
					>
				</form>
			</div>
		{/if}
	</div>
</div>
