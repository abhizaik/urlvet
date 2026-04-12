<script lang="ts">
  export let verdict: string | undefined;
  export let finalScore: number | undefined;
</script>

<div
  class="flex flex-col md:flex-row gap-6 p-6 bg-gray-900/80 rounded-xl shadow-md hover:shadow-lg transition-transform hover:scale-[1.01]"
>
  <div class="flex-1">
    <div class="text-sm font-medium text-gray-300 uppercase tracking-wide mb-1">Verdict</div>
    <div class="flex items-center gap-3">
      <span class="text-2xl font-bold text-white">{verdict ?? "-"}</span>
      {#if verdict === "Safe"}
        <span
          class="px-3 py-1 rounded-full bg-green-700 text-white font-medium text-xs uppercase tracking-wide"
        >
          Trusted
        </span>
      {:else if verdict === "Suspicious"}
        <span
          class="px-3 py-1 rounded-full bg-yellow-500 text-black font-medium text-xs uppercase tracking-wide"
        >
          Be Cautious
        </span>
      {:else if verdict === "Risky"}
        <span
          class="px-3 py-1 rounded-full bg-red-700 text-white font-medium text-xs uppercase tracking-wide"
        >
          High Risk
        </span>
      {:else if verdict === "Unclear"}
        <span
          class="px-3 py-1 rounded-full bg-gray-500 text-white font-medium text-xs uppercase tracking-wide"
        >
          Not Enough Data
        </span>
      {:else}
        <span
          class="px-3 py-1 rounded-full bg-red-600 text-white font-medium text-xs uppercase tracking-wide"
        >
          Dangerous
        </span>
      {/if}
    </div>
  </div>

  <div class="flex-1 md:text-right flex flex-col justify-center">
    <div class="text-sm font-medium text-gray-300 uppercase tracking-wide mb-1">Trust Score</div>
    <span class="text-3xl font-extrabold text-white-50">{finalScore ?? "-"} / 100</span>
  </div>
</div>
