<script lang="ts">
    import { onMount } from 'svelte';
    import PomodoroTimer from "$lib/components/PomodoroTimer.svelte";
    import ShareForm from '$lib/components/ShareForm.svelte';

    let showShareModal = false;
    let shareButtonRef: HTMLButtonElement | null = null;
    let shareFormRef: any = null;
    let shareUrl = '';

    onMount(() => {
        if (typeof window !== 'undefined') shareUrl = window.location.href;
    });

    function openModal() {
        showShareModal = true;
        // focus the form's username input on next tick
        setTimeout(() => shareFormRef?.focus(), 0);
    }

    function closeModal() {
        showShareModal = false;
        // restore focus to the share button
        setTimeout(() => shareButtonRef?.focus(), 0);
    }

    function handleKeydown(e: KeyboardEvent) {
        if (showShareModal && e.key === 'Escape') closeModal();
    }
</script>

<svelte:window on:keydown={handleKeydown} />

<div>
    <button
        bind:this={shareButtonRef}
        on:click={openModal}
        class="top-4 right-4 absolute px-6 py-3 text-base border-none rounded-lg cursor-pointer bg-blue-500 text-white transition-colors duration-200 hover:bg-blue-600"
    >
        share
    </button>

    {#if showShareModal}
                <div
                    class="fixed inset-0 bg-black/50 z-40 flex items-center justify-center"
                    on:click={(e) => { if (e.target === e.currentTarget) closeModal(); }}
                    role="dialog"
                    aria-modal="true"
                    aria-labelledby="shareform-title"
                    tabindex="-1"
                    on:keydown={(e) => { if (e.key === 'Enter' || e.key === ' ') closeModal(); }}
                >
                        <!-- ensure the form content is fully opaque while the backdrop remains semi-transparent -->
                        <ShareForm bind:this={shareFormRef} link={shareUrl} on:close={closeModal} />
                </div>
    {/if}

    <PomodoroTimer />
</div>