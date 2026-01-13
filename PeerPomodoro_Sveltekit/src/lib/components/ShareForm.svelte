<script lang="ts">
    import { createEventDispatcher } from 'svelte';
    import { userPersistentStore } from '$lib/stores/userPersistentStore.svelte';
    import { timer } from '$lib/stores/pomodoroStore.svelte';
    import { createSession } from '$lib/services/sessionService';

    let usernameInput: HTMLInputElement | null = null;
    let userName = $state(userPersistentStore.userName);
    let isLoading = $state(false);
    let { closeModal }: { closeModal: () => void } = $props();
    export function focus() {
        usernameInput?.focus();
    }

    async function handleStartSession() {
        if (!userName.trim()) return;
        
        isLoading = true;
        try {
            // Update persistent store
            userPersistentStore.setUserName(userName);
            
            // Create session using current timer state
            await createSession({
                workTime: timer.workTime,
                breakTime: timer.breakTime,
                totalRounds: timer.totalRounds
            });
            
            // Modal closing is handled by navigation or parent, but we can emit close just in case
           closeModal();
        } catch (e) {
            console.error(e);
            isLoading = false;
        } finally {
            isLoading = false;
        }
    }
</script>

<!-- stopPropagation so clicks inside the form don't close the modal backdrop -->
<form class="space-y-8" onsubmit={(e) => { e.preventDefault(); handleStartSession(); }}>
    <div
        class="bg-card px-4 sm:px-6 md:px-8 py-6 sm:py-8 rounded-3xl border-2 border-border shadow-[4px_4px_0px_0px_rgba(0,0,0,0.1)] transform hover:rotate-1 transition-transform"
    >
        <div class="space-y-3">
            <div class="space-y-6">
                <label id="shareform-title" for="username" class="text-lg sm:text-2xl font-semibold text-foreground block">Your Name</label>
                <input
                    id="username"
                    bind:this={usernameInput}
                    type="text"
                    required
                    class="text-xl sm:text-3xl w-full font-semibold text-center border-2 border-border bg-background h-12 sm:h-16 rounded-2xl focus-visible:ring-work focus-visible:ring-4"
                    bind:value={userName}
                />
            </div>
        </div>

        <button
            type="submit"
            disabled={isLoading}
            class="w-full mt-6 sm:mt-8 h-12 sm:h-16 text-xl sm:text-3xl font-bold bg-start hover:bg-start/90 text-start-foreground rounded-2xl border-2 border-primary shadow-[4px_4px_0px_0px_rgba(0,0,0,0.2)] hover:shadow-[2px_2px_0px_0px_rgba(0,0,0,0.2)] transform hover:-rotate-1 transition-all disabled:opacity-50 disabled:cursor-not-allowed"
        >
            {isLoading ? 'Starting...' : 'Start a Session!'}
        </button>
    </div>
</form>