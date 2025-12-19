<script lang="ts">
  //methods
    import {pause, startTimer} from "$lib/stores/pomodoroStore.svelte"
  //vars
    import { timer } from "$lib/stores/pomodoroStore.svelte"
  //types
    import type { TimerDisplayData } from "$lib/stores/pomodoroStore.svelte"

    const displayData: TimerDisplayData = $derived(timer)
    const periodType = $derived(displayData.isWorkPeriod ? 'Work' : 'Break')
    const minutes = $derived(Math.floor(displayData.secondsRemaining / 60))
    const seconds = $derived(displayData.secondsRemaining % 60)
    const formattedTime = $derived(`${String(minutes).padStart(2, '0')}:${String(seconds).padStart(2, '0')}`)
    const progress = $derived(() => {
    const totalSeconds = displayData.isWorkPeriod ? timer.workTime * 60 : timer.breakTime * 60
    return ((totalSeconds - displayData.secondsRemaining) / totalSeconds) * 60
  })

  function handleReset(){
    pause()
  }
</script>

<div class="max-w-125 p-8 text-center">
    <!-- Timer Display -->
    <div class="p-8 rounded-xl mb-8 transition-colors duration-300">
      <h2 class="text-xl">{ periodType }</h2>
      <div class="text-6xl font-bold my-4">{ formattedTime }</div>
      <div class="text-xl mb-4">Round { displayData.currentRound } / { displayData.totalRounds }</div>
      <div class="w-full h-2 bg-black/10 rounded overflow-hidden">
        <div
          class="h-full bg-current transition-[width] duration-1000 linear"
          style="width:{progress}%"
        ></div>
      </div>
    </div>

    <!-- Controls -->
    <div class="flex gap-4 justify-center mb-8">
    {#if !displayData.isRunning}
      <button
        onclick={startTimer}
        disabled="{displayData.isCompleted}"
        class="px-6 py-3 text-base border-none rounded-lg cursor-pointer bg-blue-500 text-white transition-colors duration-200 hover:bg-blue-600 disabled:opacity-50 disabled:cursor-not-allowed"
      >
      Continue
      </button>
    {:else}
      <button
      onclick={pause}
      class="px-6 py-3 text-base border-none rounded-lg cursor-pointer bg-blue-500 text-white transition-colors duration-200 hover:bg-blue-600"
      >
        Pause
      </button>
    {/if }
    <a href="/">
      <button
      onclick={handleReset}
      class="px-6 py-3 text-base border-none rounded-lg cursor-pointer bg-blue-500 text-white transition-colors duration-200 hover:bg-blue-600"
      >
      Reset
    </button>
  </a>
  </div>
  
  <!-- Completion Message -->
  {#if displayData.isCompleted}
  <div class="text-2xl text-green-600 mb-4">
    🎉 All rounds completed! Great work!
    </div>
    {/if }
  </div>