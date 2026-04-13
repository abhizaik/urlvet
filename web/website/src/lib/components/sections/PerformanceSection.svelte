<script lang="ts">
  export let performance: any;

  let expanded = false;
</script>

{#if performance}
  <div class="border border-gray-800 rounded-lg overflow-hidden">
    <button
      type="button"
      class="w-full flex items-center justify-between px-5 py-3 bg-gray-900/50 hover:bg-gray-900/80 transition-colors text-left focus:outline-none focus:ring-2 focus:ring-inset focus:ring-emerald-400"
      on:click={() => (expanded = !expanded)}
      aria-expanded={expanded}
      aria-controls="section-performance"
    >
      <div class="flex items-center gap-2">
        <span class="text-sm font-medium text-gray-400">Performance Timings</span>
      </div>
      <svg
        class="w-4 h-4 text-gray-500 transition-transform duration-200 {expanded
          ? 'rotate-180'
          : ''}"
        viewBox="0 0 20 20"
        fill="currentColor"
      >
        <path
          fill-rule="evenodd"
          d="M5.23 7.21a.75.75 0 011.06.02L10 11.188l3.71-3.958a.75.75 0 111.08 1.04l-4.25 4.53a.75.75 0 01-1.08 0l-4.25-4.53a.75.75 0 01.02-1.06z"
          clip-rule="evenodd"
        />
      </svg>
    </button>

    <div
      id="section-performance"
      class="transition-all duration-300 ease-in-out {expanded
        ? 'max-h-[2000px] opacity-100'
        : 'max-h-0 opacity-0 overflow-hidden'}"
    >
      <div class="space-y-0 divide-y divide-gray-800 text-sm text-gray-200 px-5 py-4">
        <div
          class="flex flex-col md:grid md:grid-cols-[350px,1fr] md:items-center gap-2 md:gap-4 py-2 first:pt-0 last:pb-0"
        >
          <span class="text-gray-400">Total time:</span>
          <span class="font-medium text-white">{performance.total_time}</span>
        </div>

        {#if performance.timings && performance.timings.length > 0}
          {#each performance.timings as timing}
            <div
              class="flex flex-col md:grid md:grid-cols-[350px,1fr] md:items-center gap-2 md:gap-4 py-2 first:pt-0 last:pb-0"
            >
              <span class="text-gray-400">{timing.task}</span>
              <span class="font-medium text-white">{timing.time}</span>
            </div>
          {/each}
        {/if}
      </div>
    </div>
  </div>
{/if}
