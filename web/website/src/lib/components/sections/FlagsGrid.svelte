<script lang="ts">
  export let reasons: { bad_reasons?: string[]; good_reasons?: string[] } | undefined;
</script>

{#if reasons}
  <div class="grid grid-cols-1 md:grid-cols-2 gap-6">
    <!-- Red Flags -->
    <div
      class="rounded-xl border border-red-200 dark:border-red-700 bg-red-50 dark:bg-red-900/20 p-5 shadow-md hover:shadow-lg hover:scale-[1.01] transition-all"
    >
      <div class="flex items-center gap-3 mb-3">
        <svg class="w-5 h-5 text-red-500 flex-shrink-0" fill="currentColor" viewBox="0 0 20 20">
          <path
            fill-rule="evenodd"
            d="M8.257 3.099c.765-1.36 2.721-1.36 3.486 0l6.518 11.59c.75 1.335-.213 3.011-1.743 3.011H3.482c-1.53 0-2.493-1.676-1.743-3.01L8.257 3.1zM11 14a1 1 0 10-2 0 1 1 0 002 0zm-1-2a.75.75 0 01-.75-.75V8a.75.75 0 011.5 0v3.25A.75.75 0 0110 12z"
            clip-rule="evenodd"
          />
        </svg>
        <h3 class="text-sm font-semibold text-red-600 dark:text-red-400 uppercase tracking-wide">
          Red Flags
        </h3>
      </div>
      {#if reasons.bad_reasons?.length}
        <ul class="space-y-2 text-red-700 dark:text-red-200 text-sm">
          {#each reasons.bad_reasons as r}
            <li class="flex items-start gap-2" title="Potential risk">
              <span class="mt-1 h-2 w-2 rounded-full bg-red-500 flex-shrink-0"></span>
              <span class="break">{r}</span>
            </li>
          {/each}
        </ul>
      {:else}
        <p class="text-gray-500 dark:text-gray-400 text-sm">No red flags found.</p>
      {/if}
    </div>

    <!-- Green Flags -->
    <div
      class="rounded-xl border border-emerald-200 dark:border-emerald-700 bg-emerald-50 dark:bg-emerald-900/20 p-5 shadow-md hover:shadow-lg hover:scale-[1.01] transition-all"
    >
      <div class="flex items-center gap-3 mb-3">
        <svg
          class="w-5 h-5 text-emerald-500 dark:text-emerald-400 flex-shrink-0"
          fill="currentColor"
          viewBox="0 0 20 20"
        >
          <path
            fill-rule="evenodd"
            d="M16.704 5.29a1 1 0 010 1.42l-7.388 7.388a1 1 0 01-1.42 0L3.296 9.498a1 1 0 111.408-1.42L8.5 11.874l6.796-6.795a1 1 0 011.408 0z"
            clip-rule="evenodd"
          />
        </svg>
        <h3
          class="text-sm font-semibold text-emerald-600 dark:text-emerald-400 uppercase tracking-wide"
        >
          Green Flags
        </h3>
      </div>
      {#if reasons.good_reasons?.length}
        <ul class="space-y-2 text-emerald-700 dark:text-emerald-200 text-sm">
          {#each reasons.good_reasons as r}
            <li class="flex items-start gap-2" title="Positive sign">
              <span class="mt-1 h-2 w-2 rounded-full bg-emerald-500 flex-shrink-0"></span>
              <span class="break">{r}</span>
            </li>
          {/each}
        </ul>
      {:else}
        <p class="text-gray-500 dark:text-gray-400 text-sm">No green flags found.</p>
      {/if}
    </div>
  </div>
{/if}
