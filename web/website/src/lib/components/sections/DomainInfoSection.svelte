<script lang="ts">
  import TooltipIcon from "../TooltipIcon.svelte";
  export let domainInfo: any;
  export let rank: number | undefined;
</script>

{#if domainInfo}
  <section
    id="section-domain"
    class="bg-gray-900/80 border border-gray-800 rounded-lg p-5 shadow-md hover:shadow-lg hover:scale-[1.01] transition-all scroll-mt-20"
  >
    <div class="flex items-center justify-between mb-4">
      <h3 class="text-base font-semibold text-white">Domain Information</h3>
      <span class="text-[10px] text-gray-400 uppercase tracking-wide px-2 py-0.5 bg-gray-800 rounded"
        >{domainInfo.source}</span
      >
    </div>

    <div class="space-y-0 divide-y divide-gray-800 text-sm text-gray-200 max-w-4xl w-full mx-auto">
      <div
        class="flex flex-col md:grid md:grid-cols-[350px,1fr] md:items-center gap-2 md:gap-4 py-2 first:pt-0 last:pb-0"
      >
        <div class="flex items-center gap-1 text-gray-400">
          <span>Domain:</span>
          <TooltipIcon
            text="The registered name of the website — what users type in the browser to visit it."
          />
        </div>
        <span class="font-medium text-white">{domainInfo.domain}</span>
      </div>

      {#if rank !== undefined}
        <div
          class="flex flex-col md:grid md:grid-cols-[350px,1fr] md:items-center gap-2 md:gap-4 py-2 first:pt-0 last:pb-0"
        >
          <div class="flex items-center gap-1 text-gray-400">
            <span>Global Traffic Rank:</span>
            <TooltipIcon
              text="A rough estimate of the website's global popularity, lower numbers mean more visitors. Derived from traffic and engagement data."
            />
          </div>
          <span class="font-medium text-white">{rank === 0 ? "Unranked" : rank}</span>
        </div>
      {/if}

      <div
        class="flex flex-col md:grid md:grid-cols-[350px,1fr] md:items-center gap-2 md:gap-4 py-2 first:pt-0 last:pb-0"
      >
        <div class="flex items-center gap-1 text-gray-400">
          <span>Registrar:</span>
          <TooltipIcon
            text="The company or organization that manages the registration of this domain (e.g., GoDaddy, Namecheap, Google Domains)."
          />
        </div>
        <span class="font-medium text-white">{domainInfo.registrar || "-"}</span>
      </div>

      <div
        class="flex flex-col md:grid md:grid-cols-[350px,1fr] md:items-center gap-2 md:gap-4 py-2 first:pt-0 last:pb-0"
      >
        <div class="flex items-center gap-1 text-gray-400">
          <span>Domain Age:</span>
          <TooltipIcon
            text="How long ago the domain was first registered. Older domains often suggest more established or legitimate websites."
          />
        </div>
        <span class="font-medium text-white">{domainInfo.age_human}</span>
      </div>

      <div
        class="flex flex-col md:grid md:grid-cols-[350px,1fr] md:items-center gap-2 md:gap-4 py-2 first:pt-0 last:pb-0"
      >
        <div class="flex items-center gap-1 text-gray-400">
          <span>DNSSEC Enabled:</span>
          <TooltipIcon
            text="A security feature that helps protect against DNS tampering and redirection attacks by digitally signing DNS data."
          />
        </div>
        {#if domainInfo.dnssec}
          <span class="text-green-400 font-medium flex items-center gap-1">✅ Yes</span>
        {:else}
          <span class="text-red-400 font-medium flex items-center gap-1">❌ No</span>
        {/if}
      </div>

      <div
        class="flex flex-col md:grid md:grid-cols-[350px,1fr] md:items-center gap-2 md:gap-4 py-2 first:pt-0 last:pb-0"
      >
        <div class="flex items-center gap-1 text-gray-400">
          <span>Created:</span>
          <TooltipIcon
            text="The date when this domain was first registered and became active on the internet."
          />
        </div>
        <span class="font-medium text-white">{domainInfo.created}</span>
      </div>

      <div
        class="flex flex-col md:grid md:grid-cols-[350px,1fr] md:items-center gap-2 md:gap-4 py-2 first:pt-0 last:pb-0"
      >
        <div class="flex items-center gap-1 text-gray-400">
          <span>Updated:</span>
          <TooltipIcon
            text="The last date the domain registration information was modified (e.g., contact change or nameserver update)."
          />
        </div>
        <span class="font-medium text-white">{domainInfo.updated}</span>
      </div>

      <div
        class="flex flex-col md:grid md:grid-cols-[350px,1fr] md:items-center gap-2 md:gap-4 py-2 first:pt-0 last:pb-0"
      >
        <div class="flex items-center gap-1 text-gray-400">
          <span>Expiry:</span>
          <TooltipIcon
            text="The date when this domain's registration will expire unless renewed by the owner."
          />
        </div>
        <span class="font-medium text-white">{domainInfo.expiry}</span>
      </div>

      {#if domainInfo.nameservers?.length}
        <div
          class="flex flex-col md:grid md:grid-cols-[350px,1fr] gap-2 md:gap-4 py-2 first:pt-0 last:pb-0"
        >
          <div class="flex items-center gap-1 text-gray-400">
            <span>Nameservers:</span>
            <TooltipIcon
              text="Servers responsible for directing internet traffic to the correct web host. They link your domain to its hosting service."
            />
          </div>
          <div class="flex flex-wrap gap-2">
            {#each domainInfo.nameservers as ns}
              <span class="px-2 py-1 bg-gray-800 text-white rounded text-xs">{ns}</span>
            {/each}
          </div>
        </div>
      {/if}

      {#if domainInfo.status?.length}
        <div
          class="flex flex-col md:grid md:grid-cols-[350px,1fr] gap-2 md:gap-4 py-2 first:pt-0 last:pb-0"
        >
          <div class="flex items-center gap-1 text-gray-400">
            <span>Status:</span>
            <TooltipIcon
              text="Domain lifecycle or control states — such as 'active', 'clientTransferProhibited', or 'pendingDelete'. They indicate administrative or operational restrictions."
            />
          </div>
          <div class="flex flex-wrap gap-2">
            {#each domainInfo.status as st}
              <span class="px-2 py-1 bg-gray-800 text-white rounded text-xs">{st}</span>
            {/each}
          </div>
        </div>
      {/if}
    </div>
  </section>
{/if}
