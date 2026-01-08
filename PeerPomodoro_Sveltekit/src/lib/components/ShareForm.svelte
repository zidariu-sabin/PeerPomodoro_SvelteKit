  
  <script lang="ts">
    import { createEventDispatcher } from 'svelte';

    interface userData {
      username?: string;
    }

    // exported prop for the pre-filled link to display (set by parent)
    export let link: string = '';

    const dispatch = createEventDispatcher();

    let userData: userData = {};
    let usernameInput: HTMLInputElement | null = null;

    // allow parent to set focus to the username input via component API
    export function focus() {
      usernameInput?.focus();
    }

    function StopSession() {
      // emit a close event; parent will close the modal
      dispatch('close');
    }
  </script>

  <!-- stopPropagation so clicks inside the form don't close the modal backdrop -->
  <form class="space-y-8">
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
            class="text-xl sm:text-3xl w-full font-semibold text-center border-2 border-border bg-background h-12 sm:h-16 rounded-2xl focus-visible:ring-work focus-visible:ring-4"
            bind:value={userData.username}
          />
        </div>

        <div class="space-y-3">
          <label for="link" class="text-lg sm:text-2xl font-semibold text-foreground block">Link</label>
          <input
            id="link"
            type="url"
            readonly
            value={link}
            class="text-xl sm:text-3xl w-full font-semibold text-center border-2 border-border bg-background h-12 sm:h-16 rounded-2xl focus-visible:ring-break focus-visible:ring-4"
          />
        </div>
      </div>

      <button
        type="button"
        on:click={StopSession}
        class="w-full mt-6 sm:mt-8 h-12 sm:h-16 text-xl sm:text-3xl font-bold bg-start hover:bg-start/90 text-start-foreground rounded-2xl border-2 border-primary shadow-[4px_4px_0px_0px_rgba(0,0,0,0.2)] hover:shadow-[2px_2px_0px_0px_rgba(0,0,0,0.2)] transform hover:-rotate-1 transition-all"
      >
        Stop Session!
      </button>
    </div>
  </form>