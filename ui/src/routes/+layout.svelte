<script lang="ts">
	import '../app.css';
	import { navigating } from '$app/state';
	import { expoOut } from 'svelte/easing';
	import { slide } from 'svelte/transition';
	interface Props {
		children?: import('svelte').Snippet;
	}

	let { children }: Props = $props();
</script>

{#if navigating}
	<!-- 
	  Loading animation for next page since svelte doesn't show any indicator. 
	   - delay 100ms because most page loads are instant, and we don't want to flash 
	   - long 12s duration because we don't actually know how long it will take
	   - exponential easing so fast loads (>100ms and <1s) still see enough progress,
		 while slow networks see it moving for a full 12 seconds
	-->
	<div
		class="bg-primary fixed left-0 right-0 top-0 z-50 h-1 w-full"
		in:slide={{ delay: 100, duration: 12000, axis: 'x', easing: expoOut }}
	></div>
{/if}
{@render children?.()}
