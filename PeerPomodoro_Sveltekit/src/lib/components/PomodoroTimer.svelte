<script lang="ts">
    import type { TimerDisplayData, TimerConfigurableData } from "$lib/stores/pomodoroStore.svelte"

    type TimerState = TimerDisplayData & TimerConfigurableData;

    let { timerState, onStart, onPause, onReset }: {
      timerState: TimerState,
      onStart: () => void,
      onPause: () => void,
      onReset: () => void
    } = $props();

    const periodType = $derived(timerState.isWorkPeriod ? 'Work' : 'Break')
    const minutes = $derived(Math.floor(timerState.secondsRemaining / 60))
    const seconds = $derived(timerState.secondsRemaining % 60)
    const formattedTime = $derived(`${String(minutes).padStart(2, '0')}:${String(seconds).padStart(2, '0')}`)
    
    const progress = $derived(
      ((timerState.isWorkPeriod ? timerState.workTime * 60 : timerState.breakTime * 60) - timerState.secondsRemaining) / 
      (timerState.isWorkPeriod ? timerState.workTime * 60 : timerState.breakTime * 60) * 100
    )
</script>

<div class="max-w-125 p-8 text-center">
    <!-- Timer Display -->
    <div class="p-8 rounded-xl mb-8 transition-colors duration-300">
      <h2 class="text-xl">{ periodType }</h2>
      <div class="text-6xl font-bold my-4">{ formattedTime }</div>
      <div class="text-xl mb-4">Round { timerState.currentRound } / { timerState.totalRounds }</div>
      <div class="w-full h-2 bg-black/10 rounded overflow-hidden">
        <div
          class="h-full bg-current transition-[width] duration-1000 linear"
          style="width:{progress}%"
        ></div>
      </div>
    </div>

    <!-- Controls -->
    <div class="flex gap-4 justify-center mb-8">
    {#if !timerState.isRunning}
      <button
        onclick={onStart}
        disabled="{timerState.isCompleted}"
        class="px-6 py-3 text-base border-none rounded-lg cursor-pointer bg-blue-500 text-white transition-colors duration-200 hover:bg-blue-600 disabled:opacity-50 disabled:cursor-not-allowed"
      >
      Continue
      </button>
    {:else}
      <button
      onclick={onPause}
      class="px-6 py-3 text-base border-none rounded-lg cursor-pointer bg-blue-500 text-white transition-colors duration-200 hover:bg-blue-600"
      >
        Pause
      </button>
    {/if }
    
    <button
      onclick={onReset}
      class="px-6 py-3 text-base border-none rounded-lg cursor-pointer bg-blue-500 text-white transition-colors duration-200 hover:bg-blue-600"
    >
      Reset
    </button>
  </div>
  
  <!-- Completion Message -->
  {#if timerState.isCompleted}
  <div class="text-2xl text-green-600 mb-4">
    🎉 All rounds completed! Great work!
    </div>
    {/if }
  </div>