<script lang="ts">
import { page } from '$app/stores';
import { onMount } from 'svelte';
import { connection } from "$lib/stores/connectionStore.svelte";
import { timer } from "$lib/stores/pomodoroStore.svelte";
import { userPersistentStore } from "$lib/stores/userPersistentStore.svelte";
import PomodoroTimer from "$lib/components/PomodoroTimer.svelte";
import UserList from "$lib/components/UserList.svelte";

let { data } = $props();

const sessionId = $page.params.slug;
let joined = false;
let isDrawerOpen = $state(false);
let shareUrl = $state('');
let myName = $state(userPersistentStore.userName);

onMount(() => {
    if (data.isValid) {
        connection.connect();
    }
    shareUrl = window.location.href;
});

$effect(() => {
    if (data.isValid && connection.state === 'connected' && !joined && sessionId) {
        connection.joinSession(sessionId, myName);
        joined = true;
    }
});

function updateUserName() {
    if (myName.trim()) {
        userPersistentStore.setUserName(myName.trim());
        connection.updateName(myName.trim());
    }
}

const displayError = $derived(data.isValid ? connection.session.error : data.error);
</script>

{#if displayError}
    <div class="fixed inset-0 z-50 flex items-center justify-center bg-background/95 backdrop-blur-sm p-4">
        <div class="bg-card border-2 border-destructive p-8 rounded-3xl shadow-[8px_8px_0px_0px_rgba(0,0,0,0.1)] max-w-md w-full text-center space-y-6 transform hover:-rotate-1 transition-transform">
            <div class="space-y-2">
                <h2 class="text-3xl font-bold text-destructive uppercase tracking-widest font-hand">Session Error</h2>
                <p class="text-muted-foreground text-lg font-medium">{displayError}</p>
            </div>
            <div class="pt-4">
                <a href="/" class="block w-full">
                    <button class="w-full py-4 text-xl font-bold bg-primary text-primary-foreground rounded-2xl border-2 border-border shadow-[4px_4px_0px_0px_rgba(0,0,0,0.2)] hover:shadow-[2px_2px_0px_0px_rgba(0,0,0,0.2)] hover:translate-y-0.5 transition-all uppercase tracking-wide">
                        Return to Home
                    </button>
                </a>
            </div>
        </div>
    </div>
{/if}

<!-- Top Right Controls -->

<div class="fixed top-4 right-4 z-30 flex items-center gap-3">

    <div 

        class="w-3 h-3 rounded-full shadow-sm ring-2 ring-background transition-colors duration-300"

        class:bg-green-500={connection.state === 'connected'}

        class:bg-red-500={connection.state === 'disconnected' || connection.state === 'error'}

        class:bg-yellow-500={connection.state === 'connecting'}

        title="Connection status: {connection.state}"

    ></div>

    <button

        onclick={() => isDrawerOpen = !isDrawerOpen}

        class="px-5 py-2.5 bg-session hover:bg-session/90 text-session-foreground border-2 border-primary shadow-[4px_4px_0px_0px_rgba(0,0,0,0.2)] hover:shadow-[2px_2px_0px_0px_rgba(0,0,0,0.2)] transform hover:-translate-y-0.5 rounded-2xl font-bold text-sm flex items-center gap-2 transition-all"

    >

        <span>👥</span>

        <span class="hidden sm:inline">Users</span>

        <span class="bg-white/30 text-session-foreground text-xs px-2 py-0.5 rounded-full ml-1 font-mono">

            {connection.users.length}

        </span>

    </button>

</div>



<!-- Drawer Backdrop -->

{#if isDrawerOpen}

    <!-- svelte-ignore a11y_click_events_have_key_events -->

    <!-- svelte-ignore a11y_no_static_element_interactions -->

    <div 

        onclick={() => isDrawerOpen = false} 

        class="fixed inset-0 bg-black/20 backdrop-blur-sm z-40 transition-opacity"

    ></div>

{/if}



<!-- Drawer Panel -->

<div class={`fixed top-0 right-0 h-full w-80 bg-card border-l-2 border-border shadow-2xl z-50 transform transition-transform duration-300 ease-in-out ${isDrawerOpen ? 'translate-x-0' : 'translate-x-full'}`}>

    <div class="p-6 h-full flex flex-col">

        <div class="flex justify-between items-center mb-8">

            <h2 class="text-xl font-bold text-foreground">Session Dashboard</h2>

            <button 

                onclick={() => isDrawerOpen = false}

                class="text-muted-foreground hover:text-foreground transition-colors p-1 hover:bg-black/5 rounded-full"

            >

                ✕

            </button>

        </div>

        <div class="mb-8 space-y-2">

            <div class="text-xs font-bold text-muted-foreground uppercase tracking-wider">Invite Link</div>

            <div class="relative group">

                <input 

                    readonly 

                    value={shareUrl}

                    onclick={(e) => e.currentTarget.select()}

                    class="w-full text-sm font-mono bg-background border-2 border-border rounded-xl p-3 pr-12 focus:outline-none focus:ring-4 focus:ring-accent/20 text-foreground truncate shadow-sm cursor-text"

                />

                <button

                    onclick={() => {

                        navigator.clipboard.writeText(shareUrl);

                    }}

                    class="absolute right-2 top-1/2 -translate-y-1/2 opacity-0 group-hover:opacity-100 transition-opacity bg-primary text-primary-foreground text-xs font-bold px-3 py-1.5 rounded-lg shadow-sm hover:bg-primary/90"

                    title="Copy to clipboard"

                >

                    Copy

                </button>

            </div>

            <div class="text-[10px] text-muted-foreground text-center">Share this link to invite others!</div>

        </div>

        <div class="mb-6 space-y-2">
            <div class="text-xs font-bold text-muted-foreground uppercase tracking-wider">Your Name</div>
            <input 
                type="text"
                bind:value={myName}
                onblur={updateUserName}
                onkeydown={(e) => e.key === 'Enter' && e.currentTarget.blur()}
                class="w-full text-sm font-semibold bg-background border-2 border-border rounded-xl p-3 focus:outline-none focus:ring-4 focus:ring-accent/20 text-foreground shadow-sm placeholder:font-normal"
                placeholder="Enter your name"
            />
        </div>

                <div class="flex-1 overflow-hidden">

                    <UserList users={connection.users} />

                </div>

        

    </div>

</div>




<div class="flex flex-col items-center gap-4 w-full px-4 pt-16">    

    <div class="w-full max-w-2xl">

        <PomodoroTimer 

            timerState={timer}

            onStart={() => connection.startTimer()}

            onPause={() => connection.pauseTimer()}

            onReset={() => connection.resetTimer()}

        />

    </div>

</div>
