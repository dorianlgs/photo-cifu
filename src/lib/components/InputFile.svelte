<script>
	let avatarInput = $state();

	let { errors } = $props();
</script>

<input
	id={'avatar'}
	name={'avatar'}
	type={'file'}
	bind:this={avatarInput}
	placeholder={'Upload an avatar'}
	class="{errors['avatar'] ? 'input-error' : ''} input-md"
/>
{#if avatarInput}
	<h2>Selected files:</h2>
	{#each Array.from(avatarInput) as file}
		<p>{file.name} ({file.size} bytes)</p>
	{/each}
{/if}

<style>
	/* file upload button */
	input[type='file']::file-selector-button {
		border-radius: 4px;
		padding: 0 16px;
		height: 40px;
		cursor: pointer;
		background-color: white;
		border: 1px solid rgba(0, 0, 0, 0.16);
		box-shadow: 0px 1px 0px rgba(0, 0, 0, 0.05);
		margin-right: 16px;
		transition: background-color 200ms;
	}

	/* file upload button hover state */
	input[type='file']::file-selector-button:hover {
		background-color: #f3f4f6;
	}

	/* file upload button active state */
	input[type='file']::file-selector-button:active {
		background-color: #e5e7eb;
	}
</style>
