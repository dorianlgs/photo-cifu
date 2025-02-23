<script lang="ts">
	import { getContext } from 'svelte';
	import type { Writable } from 'svelte/store';
	import ChangeEmailModule from './settings.svelte';

	let adminSection: Writable<string> = getContext('adminSection');
	adminSection.set('settings');
</script> 

<svelte:head>
	<title>Change Email</title>
</svelte:head>

<h1 class="text-2xl font-bold mb-6">Settings</h1>

<ChangeEmailModule />
