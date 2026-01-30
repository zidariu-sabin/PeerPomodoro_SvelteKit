<script lang="ts">
	import './layout.css';
	import favicon from '$lib/assets/favicon.svg';
	import { timer, refreshTimer } from '$lib/stores/pomodoroStore.svelte';
	import { page } from '$app/stores';
	import { onMount } from 'svelte';
	import type { WorkerMessage } from '$lib/types/workerMessages';

	let { children } = $props();

	const isTimerPage = $derived($page.url.pathname.startsWith('/timer'));

	onMount(() => {
		const sw = self as unknown as ServiceWorkerContainer;
		const handleVisibilityChange = () => {
			if (document.visibilityState === 'visible') {
				refreshTimer();
			}
		};

		document.addEventListener('visibilitychange', handleVisibilityChange);
		sw.addEventListener('message', (event) => {
			const message = event.data as WorkerMessage;
    		const { type, payload } = message;
		});
		

		return () => {
			document.removeEventListener('visibilitychange', handleVisibilityChange);
		};
	});
</script>

<svelte:head><link rel="icon" href={favicon} /></svelte:head>
<main
	class="min-h-screen w-screen flex flex-col items-center justify-center flex-1 font-hand text-center relative transition-colors duration-500"
	class:work-period={isTimerPage && timer.isWorkPeriod}
	class:break-period={isTimerPage && !timer.isWorkPeriod}
>
{@render children()}
</main>

<style>
	.work-period {
		background-color: #e0f2fe;
		color: #0c4a6e;
	}

	.break-period {
		background-color: #dcfce7;
		color: #14532d;
	}
</style>