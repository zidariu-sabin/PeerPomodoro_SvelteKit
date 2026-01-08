
<script lang="ts">
    import TimerForm from "$lib/components/TimerForm.svelte";
</script>
    <div class="h-full py-6 sm:py-8 md:py-12 mx-auto px-4 max-w-2xl">
      <header class="text-center mb-8 sm:mb-12">
        <h1
          class="text-3xl sm:text-5xl md:text-7xl font-bold text-foreground mb-2 sm:mb-3 transform -rotate-1"
        >
          Pomodoro
        </h1>
        <p class="text-md sm:text-xl md:text-2xl text-muted-foreground transform rotate-1">
          ~ stay focused, work smart ~
        </p>
      </header>
      <div>
        <TimerForm />
      </div>
    </div>