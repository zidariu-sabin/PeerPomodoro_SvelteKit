<script lang="ts">
import { page } from '$app/stores';
import { onMount } from 'svelte';
import { connection } from "$lib/stores/connectionStore.svelte";
import { timer } from "$lib/stores/pomodoroStore.svelte";
import PomodoroTimer from "$lib/components/PomodoroTimer.svelte";

const sessionId = $page.params.slug;
let joined = false;

onMount(() => {
    connection.connect();
});

$effect(() => {
    if (connection.state === 'connected' && !joined && sessionId) {
        connection.joinSession(sessionId, "User-" + Math.floor(Math.random() * 1000));
        joined = true;
    }
});
</script>

<div class="flex flex-col items-center gap-4">
    <div class="text-sm text-muted-foreground">
        Session: {sessionId}
    </div>
    <div class="text-xs text-muted-foreground">
        Status: {connection.state}
    </div>
    <PomodoroTimer 
        timerState={timer}
        onStart={() => connection.startTimer()}
        onPause={() => connection.pauseTimer()}
        onReset={() => connection.resetTimer()}
    />
</div>