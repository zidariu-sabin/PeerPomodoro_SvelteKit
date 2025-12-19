  
  <script lang="ts">
    import { timer, setTimerData } from "$lib/stores/pomodoroStore.svelte";
    import type { TimerConfigurableData } from "$lib/stores/pomodoroStore.svelte";
    const formData: TimerConfigurableData = {
      workTime: timer.workTime,
      breakTime: timer.breakTime,
      totalRounds: timer.totalRounds,
    } 

    function StartTimer() {
      setTimerData(formData)
    }

  </script>
  <form class="space-y-8" >
    <div
      class="bg-card px-4 sm:px-6 md:px-8 py-6 sm:py-8 rounded-3xl border-2 border-border shadow-[4px_4px_0px_0px_rgba(0,0,0,0.1)] transform hover:rotate-1 transition-transform"
    >
    <div class="space-y-3" >
      <div class="space-y-6">
        <label for="workMinutes" class="text-lg sm:text-2xl font-semibold text-foreground block"
        >Work Time (minutes)</label
        >
        <input
        id="workMinutes"
        type="number"
        min="1"
        step="1"
        class="text-xl sm:text-3xl w-full font-semibold text-center border-2 border-border bg-background h-12 sm:h-16 rounded-2xl focus-visible:ring-work focus-visible:ring-4"
        bind:value="{formData.workTime}"  
        />
        </div>

        <div class="space-y-3">
          <label for="breakMinutes" class="text-lg sm:text-2xl font-semibold text-foreground block"
            >Break Time (minutes)</label
          >
          <input
            id="breakMinutes"
            type="number"
            min="1"
            step="1"
            class="text-xl sm:text-3xl w-full font-semibold text-center border-2 border-border bg-background h-12 sm:h-16 rounded-2xl focus-visible:ring-break focus-visible:ring-4"
            bind:value="{formData.breakTime}" 
            />
        </div>

        <div class="space-y-3">
          <label for="totalRounds" class="text-lg sm:text-2xl font-semibold text-foreground block"
            >Number of Cycles</label
          >
          <input
            id="totalRounds"
            type="number"
            min="1"
            step="1"
            class="text-xl sm:text-3xl w-full font-semibold text-center border-2 border-border bg-background h-12 sm:h-16 rounded-2xl focus-visible:ring-accent focus-visible:ring-4"
            bind:value="{formData.totalRounds}"           
            />
        </div>
      </div>
       <a href="/timer">
         <button
         on:click={() => setTimerData(formData)}
         class="w-full mt-6 sm:mt-8 h-12 sm:h-16 text-xl sm:text-3xl font-bold bg-start hover:bg-start/90 text-start-foreground rounded-2xl border-2 border-primary shadow-[4px_4px_0px_0px_rgba(0,0,0,0.2)] hover:shadow-[2px_2px_0px_0px_rgba(0,0,0,0.2)] transform hover:-rotate-1 transition-all"
         >
         <!-- <ProiconsPlay class="mr-3 h-8 w-8" /> -->
         Start Timer!
        </button>
      </a> 
      <button
        type="submit"
        class="w-full mt-6 sm:mt-8 h-12 sm:h-16 text-xl sm:text-3xl font-bold bg-session hover:bg-session/90 text-start-foreground rounded-2xl border-2 border-primary shadow-[4px_4px_0px_0px_rgba(0,0,0,0.2)] hover:shadow-[2px_2px_0px_0px_rgba(0,0,0,0.2)] transform hover:-rotate-1 transition-all"
      >
        <!-- <ProiconsPlay class="mr-3 h-8 w-8" /> -->
        Create a Session!
      </button>
    </div>
  </form>