<script lang="ts">
  import type { SSLInfo, TLSInfo, PhishingResult } from "../../types";
  import TooltipIcon from "../TooltipIcon.svelte";
  export let sslInfo: SSLInfo | undefined;
  export let tlsInfo: TLSInfo | undefined;
  export let phishing: PhishingResult | undefined;
</script>

{#if sslInfo || tlsInfo || phishing}
  <section
    id="section-security"
    class="bg-gray-900/80 border border-gray-800 rounded-lg p-5 shadow-md hover:shadow-lg hover:scale-[1.01] transition-all scroll-mt-20"
  >
    <div class="flex items-center justify-between mb-4">
      <h3 class="text-base font-semibold text-white">Security & Encryption</h3>
      <div class="flex items-center gap-2">
        {#if phishing}
          <span class="text-[10px] text-gray-400 uppercase tracking-wide px-2 py-0.5 bg-gray-800 rounded">
            Threat Intel
          </span>
        {/if}
        {#if sslInfo || tlsInfo}
          <span class="text-[10px] text-gray-400 uppercase tracking-wide px-2 py-0.5 bg-gray-800 rounded">
            SSL / TLS
          </span>
        {/if}
      </div>
    </div>

    <!-- ── PhishTank block (top) ──────────────────────────────────────── -->
    {#if phishing}
      <div class="mb-5 rounded-lg border {phishing.in_database && phishing.valid ? 'border-red-700/60 bg-red-950/20' : 'border-gray-700 bg-gray-800/30'} p-4">

        <!-- Block header -->
        <div class="flex items-center justify-between mb-3">
          <div class="flex items-center gap-2">
            <img src="https://phishtank.com/favicon.ico" alt="" class="w-4 h-4" />
            <span class="text-sm font-semibold text-white">PhishTank Lookup</span>
            <TooltipIcon text="Community-verified phishing database check. Runs on every scan." />
          </div>
          <div class="flex items-center gap-2">
            <!-- Source badge -->
            <span class="text-[10px] px-2 py-0.5 rounded bg-gray-700 text-gray-400 border border-gray-600 uppercase tracking-wide">
              {phishing.source}
            </span>
            <!-- Cache / Live badge -->
            {#if phishing.from_cache}
              <span class="text-[10px] px-2 py-0.5 rounded bg-gray-700 text-gray-300 border border-gray-600">
                Cached
              </span>
            {:else}
              <span class="text-[10px] px-2 py-0.5 rounded bg-emerald-900/40 text-emerald-300 border border-emerald-700">
                Live
              </span>
            {/if}
          </div>
        </div>

        {#if !phishing.in_database}
          <!-- Not in database at all — cleanest result -->
          <p class="text-sm text-green-400 font-medium flex items-center gap-2">
            ✅ Not found in PhishTank database.
          </p>
        {:else if phishing.verified && !phishing.valid}
          <!-- In database, reviewed, and confirmed NOT phishing -->
          <p class="text-sm text-green-400 font-medium flex items-center gap-2">
            ✅ Reviewed by PhishTank community — confirmed not phishing.
          </p>
          {#if phishing.phish_id}
            <p class="text-xs text-gray-500 mt-1">
              Report
              {#if phishing.phish_detail_page}
                <a href={phishing.phish_detail_page} target="_blank" rel="noopener noreferrer"
                  class="text-gray-400 hover:text-gray-300 underline">#{phishing.phish_id}</a>
              {:else}
                #{phishing.phish_id}
              {/if}
            </p>
          {/if}
        {:else}
          <!-- In database and valid=true (phishing) or not yet reviewed -->
          <div class="space-y-0 divide-y divide-gray-700/50 text-sm text-gray-200">

            <!-- in_database -->
            <div class="flex flex-col md:grid md:grid-cols-[220px,1fr] md:items-center gap-1 md:gap-4 py-2 first:pt-0">
              <div class="flex items-center gap-1 text-gray-400">
                <span>In Database:</span>
                <TooltipIcon text="This URL has been submitted to or reported in PhishTank." />
              </div>
              <span class="text-yellow-400 font-medium">Yes — on record</span>
            </div>

            <!-- valid — the actual phishing signal -->
            <div class="flex flex-col md:grid md:grid-cols-[220px,1fr] md:items-center gap-1 md:gap-4 py-2">
              <div class="flex items-center gap-1 text-gray-400">
                <span>Is Phishing:</span>
                <TooltipIcon text="Whether this URL is confirmed as a phishing site. true = phishing, false = not phishing." />
              </div>
              {#if phishing.valid}
                <span class="font-semibold text-red-400">Yes — phishing confirmed</span>
              {:else}
                <span class="font-semibold text-gray-400">No</span>
              {/if}
            </div>

            <!-- verified — reviewed status -->
            <div class="flex flex-col md:grid md:grid-cols-[220px,1fr] md:items-center gap-1 md:gap-4 py-2">
              <div class="flex items-center gap-1 text-gray-400">
                <span>Community Reviewed:</span>
                <TooltipIcon text="Whether the PhishTank community has reviewed this report (true = reviewed, false = still pending)." />
              </div>
              {#if phishing.verified}
                <span class="text-white font-medium">Yes</span>
              {:else}
                <span class="text-yellow-400 font-medium">No — awaiting review</span>
              {/if}
            </div>

            <!-- verified_at -->
            {#if phishing.verified_at}
              <div class="flex flex-col md:grid md:grid-cols-[220px,1fr] md:items-center gap-1 md:gap-4 py-2">
                <div class="flex items-center gap-1 text-gray-400">
                  <span>Reviewed At:</span>
                  <TooltipIcon text="When the PhishTank community reviewed this report." />
                </div>
                <span class="font-medium text-white">{phishing.verified_at}</span>
              </div>
            {/if}

            <!-- target -->
            {#if phishing.target}
              <div class="flex flex-col md:grid md:grid-cols-[220px,1fr] md:items-center gap-1 md:gap-4 py-2">
                <div class="flex items-center gap-1 text-gray-400">
                  <span>Impersonation Target:</span>
                  <TooltipIcon text="The brand or service this phishing URL is impersonating." />
                </div>
                <span class="font-medium text-white">{phishing.target}</span>
              </div>
            {/if}

            <!-- phish_id + phish_detail_page -->
            {#if phishing.phish_id}
              <div class="flex flex-col md:grid md:grid-cols-[220px,1fr] md:items-center gap-1 md:gap-4 py-2 last:pb-0">
                <div class="flex items-center gap-1 text-gray-400">
                  <span>PhishTank Report:</span>
                  <TooltipIcon text="Unique PhishTank ID. Click to view the full report page." />
                </div>
                {#if phishing.phish_detail_page}
                  <a href={phishing.phish_detail_page} target="_blank" rel="noopener noreferrer"
                    class="font-mono text-sm text-blue-400 hover:text-blue-300 underline"
                  >#{phishing.phish_id}</a>
                {:else}
                  <span class="font-mono text-sm text-gray-300">#{phishing.phish_id}</span>
                {/if}
              </div>
            {/if}

          </div>
        {/if}
      </div>

      {#if sslInfo || tlsInfo}
        <div class="border-t border-gray-800 mb-4"></div>
      {/if}
    {/if}

    <!-- ── SSL / TLS block ───────────────────────────────────────────── -->
    <div class="space-y-0 divide-y divide-gray-800 text-sm text-gray-200 max-w-4xl w-full mx-auto">
      {#if sslInfo}
        <div class="flex flex-col md:grid md:grid-cols-[350px,1fr] md:items-center gap-2 md:gap-4 py-2 first:pt-0">
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
          <div class="flex flex-col md:grid md:grid-cols-[350px,1fr] md:items-center gap-2 md:gap-4 py-2">
            <div class="flex items-center gap-1 text-gray-400">
              <span>Certificate Chain:</span>
              <TooltipIcon text="Verifies if the SSL certificate is issued by a trusted authority and the full chain is valid." />
            </div>
            {#if sslInfo.ChainValid}
              <span class="text-green-400 font-medium flex items-center gap-1">✅ Valid</span>
            {:else}
              <span class="text-red-400 font-medium flex items-center gap-1">❌ Invalid / Self-signed</span>
            {/if}
          </div>

          <div class="flex flex-col md:grid md:grid-cols-[350px,1fr] md:items-center gap-2 md:gap-4 py-2">
            <div class="flex items-center gap-1 text-gray-400">
              <span>Certificate Issuer:</span>
              <TooltipIcon text="The organization that issued the SSL certificate." />
            </div>
            <span class="font-medium text-white">{sslInfo.Issuer || "-"}</span>
          </div>

          <div class="flex flex-col md:grid md:grid-cols-[350px,1fr] md:items-center gap-2 md:gap-4 py-2">
            <div class="flex items-center gap-1 text-gray-400">
              <span>Certificate Age:</span>
              <TooltipIcon text="How many days ago the certificate was issued. Recently issued certificates on new domains can be suspicious." />
            </div>
            <span class="font-medium text-white">{sslInfo.AgeDays} days</span>
          </div>

          <div class="flex flex-col md:grid md:grid-cols-[350px,1fr] md:items-center gap-2 md:gap-4 py-2">
            <div class="flex items-center gap-1 text-gray-400">
              <span>Valid From:</span>
              <TooltipIcon text="The date this certificate first became active." />
            </div>
            <span class="font-medium text-white">{sslInfo.NotBefore}</span>
          </div>

          <div class="flex flex-col md:grid md:grid-cols-[350px,1fr] md:items-center gap-2 md:gap-4 py-2">
            <div class="flex items-center gap-1 text-gray-400">
              <span>Expiry Date:</span>
              <TooltipIcon text="When the current SSL certificate will expire." />
            </div>
            <span class="font-medium text-white">{sslInfo.NotAfter}</span>
          </div>

          <div class="flex flex-col md:grid md:grid-cols-[350px,1fr] md:items-center gap-2 md:gap-4 py-2">
            <div class="flex items-center gap-1 text-gray-400">
              <span>Certificate Risk Level:</span>
              <TooltipIcon text="Overall assessment of the certificate's technical integrity." />
            </div>
            {#if !sslInfo.IsSuspicious}
              <span class="text-green-400 font-medium flex items-center gap-1">✅ Low Risk</span>
            {:else}
              <span class="text-yellow-400 font-medium flex items-center gap-1">⚠️ Suspicious</span>
            {/if}
          </div>

          {#if sslInfo.Reasons && sslInfo.Reasons.length > 0}
            <div class="flex flex-col md:grid md:grid-cols-[350px,1fr] gap-2 md:gap-4 py-2">
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

          <div class="flex flex-col md:grid md:grid-cols-[350px,1fr] md:items-center gap-2 md:gap-4 py-2">
            <div class="flex items-center gap-1 text-gray-400">
              <span>Certificate Fingerprint:</span>
              <TooltipIcon text="A unique identifier (SHA-256 hash) for this specific certificate." />
            </div>
            <span class="font-mono text-[12px] text-gray-300 break-all">{sslInfo.Fingerprint}</span>
          </div>
        {/if}
      {/if}

      {#if tlsInfo}
        <div class="flex flex-col md:grid md:grid-cols-[350px,1fr] md:items-center gap-2 md:gap-4 py-2">
          <div class="flex items-center gap-1 text-gray-400">
            <span>TLS Issuer (Connection):</span>
            <TooltipIcon text="The certificate issuer detected during the live connection." />
          </div>
          <span class="font-medium text-white">{tlsInfo.Issuer || "-"}</span>
        </div>

        <div class="flex flex-col md:grid md:grid-cols-[350px,1fr] md:items-center gap-2 md:gap-4 py-2 last:pb-0">
          <div class="flex items-center gap-1 text-gray-400">
            <span>Hostname Match:</span>
            <TooltipIcon text="Ensures the certificate is actually issued for the domain you are visiting." />
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
