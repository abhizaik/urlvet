<script lang="ts">
  export let performance: any;
</script>

{#if performance}
  <section
    id="section-performance"
    class="bg-gray-900/80 border border-gray-800 rounded-lg p-5 shadow-md hover:shadow-lg hover:scale-[1.01] transition-all scroll-mt-20"
  >
    <div class="flex items-center justify-between mb-4">
      <h3 class="text-base font-semibold text-white">Performance</h3>
      <span
        class="text-[10px] text-gray-400 uppercase tracking-wide px-2 py-0.5 bg-gray-800 rounded"
        >Timings</span
      >
    </div>

    <div class="space-y-0 divide-y divide-gray-800 text-sm text-gray-200 max-w-4xl w-full mx-auto">
      <div
        class="flex flex-col md:grid md:grid-cols-[350px,1fr] md:items-center gap-2 md:gap-4 py-2 first:pt-0 last:pb-0"
      >
        <span class="text-gray-400">Total time:</span>
        <span class="font-medium text-white">{performance.total_time}</span>
      </div>

      {#if performance.timings && performance.timings.length > 0}
        <ul class="text-sm text-gray-100 space-y-0 divide-y divide-gray-800">
          {#each performance.timings as timing}
            <li
              class="flex flex-col md:grid md:grid-cols-[350px,1fr] md:items-center gap-2 md:gap-4 py-2 first:pt-0 last:pb-0"
            >
              <span class="text-gray-400">{timing.task}</span>
              <span>{timing.time}</span>
            </li>
          {/each}
        </ul>
      {/if}
    </div>
  </section>
{/if}
