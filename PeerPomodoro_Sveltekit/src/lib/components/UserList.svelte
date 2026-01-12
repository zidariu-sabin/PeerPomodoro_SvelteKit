<script lang="ts">
    import type { User } from "$lib/stores/connectionStore.svelte";
    
    interface Props {
        users: User[];
    }

    let { users }: Props = $props();
</script>

<div class="w-full h-full">
    <h3 class="font-bold mb-4 text-sm uppercase tracking-wider text-muted-foreground">
        Connected Users ({users.length})
    </h3>
    <ul class="space-y-3 overflow-y-auto max-h-[calc(100vh-150px)] px-1">
        {#each users as user (user.id)}
            <li class="flex items-center gap-3 p-3 bg-background border-2 border-border rounded-xl shadow-[2px_2px_0px_0px_rgba(0,0,0,0.1)] transform transition-transform hover:-translate-y-0.5">
                <div class="w-8 h-8 rounded-full bg-start border-2 border-border flex items-center justify-center shrink-0 shadow-sm">
                    <span class="text-xs font-bold text-start-foreground">
                        {user.name.substring(0, 2).toUpperCase()}
                    </span>
                </div>
                <span class="text-base font-semibold text-foreground truncate flex-1">
                    {user.name}
                </span>
                <div 
                    class="w-2.5 h-2.5 rounded-full bg-green-500 shadow-sm ring-1 ring-black/5"
                    title="Status: Connected"
                ></div>
            </li>
        {/each}
        {#if users.length === 0}
            <li class="text-sm text-muted-foreground italic p-2 text-center">Waiting for users...</li>
        {/if}
    </ul>
</div>
