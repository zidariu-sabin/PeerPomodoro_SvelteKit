<script lang="ts">
	import './layout.css';
	import favicon from '$lib/assets/favicon.svg';

	let { children } = $props();
</script>

<svelte:head><link rel="icon" href={favicon} /></svelte:head>
<main
	class="h-screen w-screen flex flex-col items-center justify-center flex-1 font-hand text-center"
>
{@render children()}
</main>