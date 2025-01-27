<script lang="ts">
	import { onMount, onDestroy } from 'svelte';
	import { pb } from '$lib/pocketbase';

	let unsubscribe: () => void;

	onMount(async () => {
		unsubscribe = await pb.collection('messages').subscribe('*', (newData) => {
			console.log({ newData });
		});
	});

	onDestroy(() => {
		if (unsubscribe) unsubscribe();
	});
</script>

<svelte:head>
	<title>Sign up</title>
</svelte:head>

<h1 class="text-2xl font-bold mb-6">Sign Up</h1>

<div class="text-l text-slate-800 mt-4 mb-2">
	Have an account? <a class="underline" href="/login/sign_in">Sign in</a>.
</div>
