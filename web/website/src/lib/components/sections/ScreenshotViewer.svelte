<script lang="ts">
  export let screenshotUrl: string | null = null;
  export let loading = false;
  export let failed = false;
  export let unavailableReason: string | null = null;
  export let compact = false;

  let showModal = false;
</script>

{#if compact}
  <!-- Compact inline thumbnail for embedding next to VerdictCard -->
  {#if loading && !screenshotUrl}
    <div class="w-full h-28 md:h-36 animate-pulse bg-gray-800 rounded-lg"></div>
  {:else if screenshotUrl}
    <button
      type="button"
      class="w-full p-0 border-0 bg-transparent cursor-pointer rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500"
      on:click={() => (showModal = true)}
      aria-label="View full-size screenshot"
    >
      <img
        src={screenshotUrl}
        alt="Website screenshot"
        class="w-full rounded-lg border border-gray-800 hover:opacity-90 transition-opacity"
        loading="lazy"
      />
    </button>
  {:else if failed}
    <div class="w-full h-28 md:h-36 flex flex-col items-center justify-center rounded-lg gap-1">
      <svg
        class="w-5 h-5 text-gray-600"
        fill="none"
        viewBox="0 0 24 24"
        stroke="currentColor"
        stroke-width="1.5"
      >
        <path
          stroke-linecap="round"
          stroke-linejoin="round"
          d="M2.25 15.75l5.159-5.159a2.25 2.25 0 013.182 0l5.159 5.159m-1.5-1.5l1.409-1.409a2.25 2.25 0 013.182 0l2.909 2.909M3 3l18 18"
        />
      </svg>
      <span class="text-[10px] text-gray-600">No preview</span>
    </div>
  {/if}
{:else}
  <!-- Full layout -->
  {#if loading && !screenshotUrl}
    <div class="mt-6 rounded-xl border border-gray-800 bg-gray-900/70 p-4 shadow-md">
      <p class="text-sm font-semibold text-gray-300 mb-2">Website Screenshot</p>
      <div class="animate-pulse w-full h-40 bg-gray-800 rounded-lg"></div>
    </div>
  {:else if failed && !screenshotUrl}
    <div class="mt-6 rounded-xl border border-gray-800 bg-gray-900/70 p-4 shadow-md">
      <p class="text-sm font-semibold text-gray-300 mb-2">Website Screenshot</p>
      <div class="flex items-center gap-2 text-xs text-gray-500">
        <svg
          class="w-3.5 h-3.5 flex-shrink-0"
          fill="none"
          viewBox="0 0 24 24"
          stroke="currentColor"
          stroke-width="2"
        >
          <path
            stroke-linecap="round"
            stroke-linejoin="round"
            d="M18.364 18.364A9 9 0 005.636 5.636m12.728 12.728A9 9 0 015.636 5.636m12.728 12.728L5.636 5.636"
          />
        </svg>
        {unavailableReason ?? "Screenshot unavailable."}
      </div>
    </div>
  {:else if screenshotUrl}
    <div
      class="mt-6 rounded-xl border border-gray-800 bg-gray-900/70 p-4 shadow-md hover:shadow-lg transition-all"
    >
      <h4 class="text-sm font-semibold text-gray-300 mb-2">Website Screenshot</h4>
      <button
        type="button"
        class="w-full p-0 border-0 bg-transparent cursor-pointer rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500"
        on:click={() => (showModal = true)}
        aria-label="View full-size screenshot"
      >
        <img
          src={screenshotUrl}
          alt="Website screenshot"
          class="w-full rounded-lg border border-gray-800 hover:opacity-90 transition-opacity"
          loading="lazy"
        />
      </button>
    </div>
  {/if}
{/if}

{#if showModal}
  <div
    class="fixed inset-0 bg-black/80 flex items-center justify-center z-50"
    role="presentation"
    on:click={() => (showModal = false)}
  >
    <button
      class="absolute top-4 right-4 text-gray-300 hover:text-white text-2xl leading-none"
      on:click={() => (showModal = false)}
      aria-label="Close screenshot">×</button
    >
    <div role="presentation" on:click|stopPropagation>
      <img
        src={screenshotUrl}
        alt="Full screenshot"
        class="max-h-[90vh] max-w-[90vw] rounded-lg shadow-lg"
      />
    </div>
  </div>
{/if}

<style>
  img {
    transition: transform 0.2s ease-in-out;
  }
  img:hover {
    transform: scale(1.02);
  }
</style>
