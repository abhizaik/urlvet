<script lang="ts">
  export let screenshotUrl: string | null = null;

  let showModal = false;
</script>

{#if screenshotUrl}
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

{#if showModal}
  <div class="fixed inset-0 bg-black/80 flex items-center justify-center z-50">
    <button
      class="absolute top-4 right-4 text-gray-300 hover:text-white text-2xl"
      on:click={() => (showModal = false)}>×</button
    >
    <img
      src={screenshotUrl}
      alt="Full screenshot"
      class="max-h-[90vh] max-w-[90vw] rounded-lg shadow-lg"
    />
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
