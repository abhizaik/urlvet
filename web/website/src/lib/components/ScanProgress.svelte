<script lang="ts">
  import { onDestroy, onMount } from "svelte";

  export let loading = false;
  export let done = false;

  const STEPS = [
    "Resolving domain & IP addresses",
    "Checking URL structure & patterns",
    "Analyzing TLD & DNS records",
    "Checking redirections & HTTP status",
    "Looking up WHOIS & domain age",
    "Checking homoglyphs & typosquatting",
    "Analyzing SSL & TLS certificate",
    "Scanning threat intelligence feeds",
    "Inspecting page content & forms",
    "Computing trust score",
  ];

  // -1 = not started, 0..N-1 = current active step
  let activeStep = -1;
  let completedSteps: Set<number> = new Set();
  let timer: ReturnType<typeof setInterval>;

  function startProgress() {
    activeStep = 0;
    completedSteps = new Set();
    let stepIdx = 0;
    // Advance a step roughly every 1.4s, leaving last step for when done arrives
    timer = setInterval(() => {
      if (stepIdx < STEPS.length - 2) {
        completedSteps = new Set([...completedSteps, stepIdx]);
        stepIdx++;
        activeStep = stepIdx;
      }
    }, 1000);
  }

  function finishProgress() {
    clearInterval(timer);
    // Complete all steps quickly
    let i = activeStep;
    const flush = setInterval(() => {
      if (i < STEPS.length) {
        completedSteps = new Set([...completedSteps, i]);
        i++;
        activeStep = i;
      } else {
        clearInterval(flush);
        activeStep = STEPS.length; // all done
      }
    }, 80);
  }

  $: if (loading) startProgress();
  $: if (done) finishProgress();

  onDestroy(() => clearInterval(timer));
</script>

{#if loading || done}
  <div class="max-w-sm mx-auto mt-6 space-y-2">
    {#each STEPS as step, i}
      {@const isCompleted = completedSteps.has(i)}
      {@const isActive = activeStep === i && loading}
      <div class="flex items-center gap-3 transition-opacity duration-300 {i > activeStep && !isCompleted ? 'opacity-30' : 'opacity-100'}">
        <!-- Status icon -->
        <div class="w-5 h-5 flex-shrink-0 flex items-center justify-center">
          {#if isCompleted}
            <svg class="w-4 h-4 text-emerald-400" fill="none" stroke="currentColor" stroke-width="2.5" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" d="M5 13l4 4L19 7" />
            </svg>
          {:else if isActive}
            <svg class="w-4 h-4 text-blue-400 animate-spin" fill="none" viewBox="0 0 24 24">
              <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="3"/>
              <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8v4a4 4 0 00-4 4H4z"/>
            </svg>
          {:else}
            <span class="w-1.5 h-1.5 rounded-full bg-gray-700 inline-block"></span>
          {/if}
        </div>
        <!-- Label -->
        <span class="text-sm {isCompleted ? 'text-gray-400 line-through decoration-gray-600' : isActive ? 'text-white font-medium' : 'text-gray-600'}">
          {step}
        </span>
      </div>
    {/each}
  </div>
{/if}
