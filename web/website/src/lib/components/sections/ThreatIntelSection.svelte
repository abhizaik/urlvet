<script lang="ts">
  import type { ThreatIntel } from "../../types";
  export let threatIntel: ThreatIntel | undefined;
</script>

{#if threatIntel?.phishtank?.in_database}
  <section
    id="section-threat"
    class="bg-gray-900/80 border border-gray-800 rounded-lg p-5 shadow-md hover:shadow-lg hover:scale-[1.01] transition-all scroll-mt-20"
  >
    <div class="flex items-center justify-between mb-4">
      <h3 class="text-base font-semibold text-white">Threat Intelligence</h3>
      <span class="text-[10px] text-gray-400 uppercase tracking-wide px-2 py-0.5 bg-gray-800 rounded"
        >External Feeds</span
      >
    </div>

    <div class="space-y-4">
      <div class="p-4 rounded-lg bg-gray-800/50 border border-gray-700">
        <div class="flex items-center gap-3 mb-2">
          <img src="https://phishtank.com/favicon.ico" alt="" class="w-4 h-4" />
          <h4 class="text-sm font-medium text-white">PhishTank Result</h4>
        </div>

        <div class="grid grid-cols-1 md:grid-cols-2 gap-x-8 gap-y-2 text-sm">
          <div
            class="flex justify-between md:justify-start md:gap-4 border-b border-gray-700/50 pb-1"
          >
            <span class="text-gray-400">Status:</span>
            {#if threatIntel.phishtank.verified}
              <span class="text-red-400 font-semibold">Verified Phishing</span>
            {:else}
              <span class="text-yellow-400 font-semibold">Suspected</span>
            {/if}
          </div>
          <div
            class="flex justify-between md:justify-start md:gap-4 border-b border-gray-700/50 pb-1"
          >
            <span class="text-gray-400">Online:</span>
            <span class={threatIntel.phishtank.is_online ? "text-red-400" : "text-gray-400"}>
              {threatIntel.phishtank.is_online ? "Yes" : "No"}
            </span>
          </div>
          {#if threatIntel.phishtank.target}
            <div
              class="flex justify-between md:justify-start md:gap-4 border-b border-gray-700/50 pb-1 md:col-span-2"
            >
              <span class="text-gray-400">Reported Target:</span>
              <span class="text-white font-medium">{threatIntel.phishtank.target}</span>
            </div>
          {/if}
        </div>
      </div>
    </div>
  </section>
{/if}
