<script lang="ts">
  import TooltipIcon from "../TooltipIcon.svelte";
  import type { SSLInfo, TLSInfo } from "../../types";
  export let sslInfo: SSLInfo | undefined;
  export let tlsInfo: TLSInfo | undefined;
</script>

{#if sslInfo || tlsInfo}
  <section
    id="section-security"
    class="bg-gray-900/80 border border-gray-800 rounded-lg p-5 shadow-md hover:shadow-lg hover:scale-[1.01] transition-all scroll-mt-20"
  >
    <div class="flex items-center justify-between mb-4">
      <h3 class="text-base font-semibold text-white">Security & Encryption</h3>
      <span
        class="text-[10px] text-gray-400 uppercase tracking-wide px-2 py-0.5 bg-gray-800 rounded"
        >SSL / TLS</span
      >
    </div>

    <div class="space-y-0 divide-y divide-gray-800 text-sm text-gray-200 max-w-4xl w-full mx-auto">
      {#if sslInfo}
        <div
          class="flex flex-col md:grid md:grid-cols-[350px,1fr] md:items-center gap-2 md:gap-4 py-2 first:pt-0"
        >
          <div class="flex items-center gap-1 text-gray-400">
            <span>SSL Support:</span>
            <TooltipIcon text="Checks if the website supports secure HTTPS connections." />
          </div>
          {#if sslInfo.HasTLS}
            <span class="text-green-400 font-medium flex items-center gap-1">✅ Enabled</span>
          {:else}
            <span class="text-red-400 font-medium flex items-center gap-1">❌ Disabled</span>
          {/if}
        </div>

        {#if sslInfo.HasTLS}
          <div
            class="flex flex-col md:grid md:grid-cols-[350px,1fr] md:items-center gap-2 md:gap-4 py-2"
          >
            <div class="flex items-center gap-1 text-gray-400">
              <span>Certificate Chain:</span>
              <TooltipIcon
                text="Verifies if the SSL certificate is issued by a trusted authority and the full chain is valid."
              />
            </div>
            {#if sslInfo.ChainValid}
              <span class="text-green-400 font-medium flex items-center gap-1">✅ Valid</span>
            {:else}
              <span class="text-red-400 font-medium flex items-center gap-1"
                >❌ Invalid / Self-signed</span
              >
            {/if}
          </div>

          <div
            class="flex flex-col md:grid md:grid-cols-[350px,1fr] md:items-center gap-2 md:gap-4 py-2"
          >
            <div class="flex items-center gap-1 text-gray-400">
              <span>Certificate Issuer:</span>
              <TooltipIcon text="The organization that issued the SSL certificate." />
            </div>
            <span class="font-medium text-white">{sslInfo.Issuer || "-"}</span>
          </div>

          <div
            class="flex flex-col md:grid md:grid-cols-[350px,1fr] md:items-center gap-2 md:gap-4 py-2"
          >
            <div class="flex items-center gap-1 text-gray-400">
              <span>Certificate Age:</span>
              <TooltipIcon
                text="How many days ago the certificate was issued. Recently issued certificates on new domains can be suspicious."
              />
            </div>
            <span class="font-medium text-white">{sslInfo.AgeDays} days</span>
          </div>

          <div
            class="flex flex-col md:grid md:grid-cols-[350px,1fr] md:items-center gap-2 md:gap-4 py-2"
          >
            <div class="flex items-center gap-1 text-gray-400">
              <span>Valid From:</span>
              <TooltipIcon text="The date this certificate first became active." />
            </div>
            <span class="font-medium text-white">{sslInfo.NotBefore}</span>
          </div>

          <div
            class="flex flex-col md:grid md:grid-cols-[350px,1fr] md:items-center gap-2 md:gap-4 py-2"
          >
            <div class="flex items-center gap-1 text-gray-400">
              <span>Expiry Date:</span>
              <TooltipIcon text="When the current SSL certificate will expire." />
            </div>
            <span class="font-medium text-white">{sslInfo.NotAfter}</span>
          </div>

          <div
            class="flex flex-col md:grid md:grid-cols-[350px,1fr] md:items-center gap-2 md:gap-4 py-2"
          >
            <div class="flex items-center gap-1 text-gray-400">
              <span>Certificate Risk Level:</span>
              <TooltipIcon text="Overall assessment of the certificate's technical integrity." />
            </div>
            {#if !sslInfo.IsSuspicious}
              <span class="text-green-400 font-medium flex items-center gap-1">✅ Low Risk</span>
            {:else}
              <span class="text-yellow-400 font-medium flex items-center gap-1"
                >⚠️ Suspicious</span
              >
            {/if}
          </div>

          {#if sslInfo.Reasons && sslInfo.Reasons.length > 0}
            <div
              class="flex flex-col md:grid md:grid-cols-[350px,1fr] gap-2 md:gap-4 py-2"
            >
              <div class="flex items-center gap-1 text-gray-400">
                <span>Technical Warnings:</span>
                <TooltipIcon text="Specific technical reasons why this certificate is flagged." />
              </div>
              <ul class="text-xs text-yellow-400 list-disc list-inside">
                {#each sslInfo.Reasons as reason}
                  <li>{reason}</li>
                {/each}
              </ul>
            </div>
          {/if}

          <div
            class="flex flex-col md:grid md:grid-cols-[350px,1fr] md:items-center gap-2 md:gap-4 py-2"
          >
            <div class="flex items-center gap-1 text-gray-400">
              <span>Certificate Fingerprint:</span>
              <TooltipIcon text="A unique identifier (SHA-256 hash) for this specific certificate." />
            </div>
            <span class="font-mono text-[12px] text-gray-300 break-all">{sslInfo.Fingerprint}</span>
          </div>
        {/if}
      {/if}

      {#if tlsInfo}
        <div
          class="flex flex-col md:grid md:grid-cols-[350px,1fr] md:items-center gap-2 md:gap-4 py-2"
        >
          <div class="flex items-center gap-1 text-gray-400">
            <span>TLS Issuer (Connection):</span>
            <TooltipIcon text="The certificate issuer detected during the live connection." />
          </div>
          <span class="font-medium text-white">{tlsInfo.Issuer || "-"}</span>
        </div>

        <div
          class="flex flex-col md:grid md:grid-cols-[350px,1fr] md:items-center gap-2 md:gap-4 py-2 last:pb-0"
        >
          <div class="flex items-center gap-1 text-gray-400">
            <span>Hostname Match:</span>
            <TooltipIcon
              text="Ensures the certificate is actually issued for the domain you are visiting."
            />
          </div>
          {#if !tlsInfo.HostnameMismatch}
            <span class="text-green-400 font-medium flex items-center gap-1">✅ Match</span>
          {:else}
            <span class="text-red-400 font-medium flex items-center gap-1">❌ Mismatch</span>
          {/if}
        </div>
      {/if}
    </div>
  </section>
{/if}
