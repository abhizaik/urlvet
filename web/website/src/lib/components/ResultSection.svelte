<script lang="ts">
  import { browser } from "$app/environment";
  import { onDestroy } from "svelte";
  import type { AnalyzeResult } from "../types";
  import ContentSection from "./sections/ContentSection.svelte";
  import DomainInfoSection from "./sections/DomainInfoSection.svelte";
  import FlagsGrid from "./sections/FlagsGrid.svelte";
  import InfrastructureSection from "./sections/InfrastructureSection.svelte";
  import PerformanceSection from "./sections/PerformanceSection.svelte";
  import RedirectionSection from "./sections/RedirectionSection.svelte";
  import ScreenshotViewer from "./sections/ScreenshotViewer.svelte";
  import SecuritySection from "./sections/SecuritySection.svelte";
  import ThreatIntelSection from "./sections/ThreatIntelSection.svelte";
  import URLSignalsSection from "./sections/URLSignalsSection.svelte";
  import VerdictCard from "./sections/VerdictCard.svelte";

  export let data: AnalyzeResult | null = null;
  export let screenshotUrl: string | null = null;
  export let screenshotLoading = false;
  export let screenshotFailed = false;
  export let loading = false;
  export let error: string | null = null;

  let shareCopied = false;
  let sectionExpanded: Record<string, boolean> = {
    domain: false,
    analysis: false,
    threatintel: false,
    security: false,
    content: false,
    features: false,
    infrastructure: false,
    performance: false,
  };

  $: primary = data?.result;
  $: httpStatusCode = data?.analysis?.http_status?.code ?? null;
  $: screenshotUnavailableReason = (() => {
    if (!screenshotFailed) return null;
    if (httpStatusCode === 0 || httpStatusCode === null)
      return "No web content detected. Preview unavailable.";
    return "Screenshot unavailable — the site may be blocking automated access or failed to respond.";
  })();

  onDestroy(() => {
    if (screenshotUrl) {
      URL.revokeObjectURL(screenshotUrl);
    }
  });

  function toggleSection(id: string) {
    sectionExpanded = { ...sectionExpanded, [id]: !sectionExpanded[id] };
  }

  async function shareLink() {
    if (!browser) return;
    try {
      await navigator.clipboard.writeText(window.location.href);
      shareCopied = true;
      setTimeout(() => (shareCopied = false), 2000);
    } catch (err) {
      console.error("Clipboard copy failed:", err);
    }
  }

  function scrollToTop() {
    window.scrollTo({ top: 0, behavior: "smooth" });
  }

  const CHEVRON_PATH =
    "M5.23 7.21a.75.75 0 011.06.02L10 11.188l3.71-3.958a.75.75 0 111.08 1.04l-4.25 4.53a.75.75 0 01-1.08 0l-4.25-4.53a.75.75 0 01.02-1.06z";
</script>

{#if error}
  <div class="max-w-3xl mx-auto p-4 bg-red-900/30 border border-red-700 text-red-200 rounded-md">
    {error}
  </div>
{:else if loading}
  <div class="max-w-3xl mx-auto space-y-4">
    <div class="animate-pulse rounded-xl border border-gray-800 bg-gray-950/60 p-6 h-32"></div>
    <div class="animate-pulse rounded-xl border border-gray-800 bg-gray-950/60 p-6 h-24"></div>
  </div>
{:else if data}
  <section class="max-w-4xl mx-auto space-y-8 px-4">
    <!-- Header & Share Button -->
    <div class="flex flex-col md:flex-row items-start md:items-center justify-between gap-3">
      <div class="flex flex-col">
        <h2 class="text-2xl font-semibold text-white" id="analysis-summary">Analysis Summary</h2>
        <p class="text-gray-400 text-sm mt-1">Here's the security profile for {data?.domain}</p>
      </div>

      <button
        class="self-end md:self-auto group inline-flex items-center justify-center gap-2 px-5 py-3 rounded-lg bg-gradient-to-r from-blue-600 to-blue-500 hover:from-blue-500 hover:to-blue-400 text-white text-sm font-semibold shadow-md hover:shadow-lg transition-all duration-200 hover:scale-105 active:scale-95 focus:outline-none focus:ring-2 focus:ring-offset-0 focus:ring-blue-400 focus:ring-offset-gray-950 disabled:opacity-50 {shareCopied
          ? 'from-emerald-600 to-emerald-500 hover:from-emerald-500 hover:to-emerald-400 animate-pulse'
          : ''}"
        on:click={shareLink}
        disabled={shareCopied}
      >
        {#if shareCopied}
          <svg
            class="w-4 h-4 text-white animate-in fade-in"
            fill="none"
            stroke="currentColor"
            stroke-width="2.5"
            viewBox="0 0 24 24"
          >
            <path stroke-linecap="round" stroke-linejoin="round" d="M5 13l4 4L19 7" />
          </svg>
          <span class="text-white font-medium">Copied!</span>
        {:else}
          <svg
            class="w-4 h-4 text-white group-hover:rotate-12 transition-transform duration-200"
            fill="none"
            stroke="currentColor"
            stroke-width="2"
            viewBox="0 0 24 24"
          >
            <path
              stroke-linecap="round"
              stroke-linejoin="round"
              d="M7.217 10.907a2.25 2.25 0 100 2.186m0-2.186c.18.324.283.696.283 1.093s-.103.77-.283 1.093m0-2.186l9.566-5.314m-9.566 7.5l9.566 5.314m0 0a2.25 2.25 0 103.935 2.186 2.25 2.25 0 00-3.935-2.186zm0-12.814a2.25 2.25 0 103.935-2.186 2.25 2.25 0 00-3.935 2.186z"
            />
          </svg>
          <span class="font-medium">Share</span>
        {/if}
      </button>
    </div>

    <!-- Verdict + Screenshot side-by-side on desktop, stacked on mobile -->
    <div class="flex flex-col md:flex-row gap-4 items-stretch">
      <div class="flex-1 min-w-0">
        <VerdictCard verdict={primary?.verdict} finalScore={primary?.final_score} />
      </div>
      {#if screenshotLoading || screenshotUrl}
        <div
          class="md:w-56 md:flex-shrink-0 rounded-xl border border-gray-800 bg-gray-900/70 p-3 shadow-md flex flex-col gap-2"
        >
          <p class="text-xs font-semibold text-gray-400">Preview</p>
          <ScreenshotViewer
            {screenshotUrl}
            loading={screenshotLoading}
            failed={screenshotFailed}
            unavailableReason={screenshotUnavailableReason}
            compact={true}
          />
        </div>
      {/if}
    </div>

    <!-- Red / Green Flags -->
    <FlagsGrid reasons={primary?.reasons} />

    <!-- Detailed Sections -->
    {#if data.domain_info || data.analysis || data.phishing || data.ssl_info || data.tls_info || data.content_data || data.features || data.infrastructure}
      <div class="space-y-2">

        {#if data.domain_info}
          <div id="section-domain" class="rounded-xl border border-gray-800 overflow-hidden scroll-mt-20">
            <button type="button" class="w-full flex items-center justify-between px-5 py-4 bg-gray-900/50 hover:bg-gray-900/80 transition-colors text-left focus:outline-none focus-visible:ring-2 focus-visible:ring-inset focus-visible:ring-emerald-400" on:click={() => toggleSection("domain")} aria-expanded={sectionExpanded.domain}>
              <div class="flex items-center gap-3">
                <span class="text-base leading-none">🏷️</span>
                <span class="text-sm font-semibold text-gray-100">Domain Info</span>
                {#if data?.domain_info?.age_days !== undefined && data.domain_info.age_days < 365}
                  <span class="w-2 h-2 rounded-full bg-yellow-400 flex-shrink-0" title="New domain"></span>
                {/if}
              </div>
              <svg class="w-4 h-4 text-gray-500 transition-transform duration-200 {sectionExpanded.domain ? 'rotate-180' : ''}" viewBox="0 0 20 20" fill="currentColor"><path fill-rule="evenodd" d={CHEVRON_PATH} clip-rule="evenodd" /></svg>
            </button>
            <div class="transition-all duration-300 ease-in-out {sectionExpanded.domain ? 'max-h-[5000px] opacity-100' : 'max-h-0 opacity-0 overflow-hidden'}">
              <div class="acc-body border-t border-gray-800"><DomainInfoSection domainInfo={data.domain_info} rank={data.features?.rank} /></div>
            </div>
          </div>
        {/if}

        {#if data.analysis}
          <div id="section-analysis" class="rounded-xl border border-gray-800 overflow-hidden scroll-mt-20">
            <button type="button" class="w-full flex items-center justify-between px-5 py-4 bg-gray-900/50 hover:bg-gray-900/80 transition-colors text-left focus:outline-none focus-visible:ring-2 focus-visible:ring-inset focus-visible:ring-emerald-400" on:click={() => toggleSection("analysis")} aria-expanded={sectionExpanded.analysis}>
              <div class="flex items-center gap-3">
                <span class="text-base leading-none">🔀</span>
                <span class="text-sm font-semibold text-gray-100">Redirection</span>
                {#if data?.analysis?.redirection_result?.has_domain_jump}
                  <span class="w-2 h-2 rounded-full bg-red-500 flex-shrink-0" title="Domain jump detected"></span>
                {:else if (data?.analysis?.redirection_result?.chain_length ?? 0) > 3}
                  <span class="w-2 h-2 rounded-full bg-yellow-400 flex-shrink-0" title="Long redirect chain"></span>
                {/if}
              </div>
              <svg class="w-4 h-4 text-gray-500 transition-transform duration-200 {sectionExpanded.analysis ? 'rotate-180' : ''}" viewBox="0 0 20 20" fill="currentColor"><path fill-rule="evenodd" d={CHEVRON_PATH} clip-rule="evenodd" /></svg>
            </button>
            <div class="transition-all duration-300 ease-in-out {sectionExpanded.analysis ? 'max-h-[5000px] opacity-100' : 'max-h-0 opacity-0 overflow-hidden'}">
              <div class="acc-body border-t border-gray-800"><RedirectionSection analysis={data.analysis} domain={data.domain} /></div>
            </div>
          </div>
        {/if}

        {#if data.phishing}
          <div id="section-threatintel" class="rounded-xl border border-gray-800 overflow-hidden scroll-mt-20">
            <button type="button" class="w-full flex items-center justify-between px-5 py-4 bg-gray-900/50 hover:bg-gray-900/80 transition-colors text-left focus:outline-none focus-visible:ring-2 focus-visible:ring-inset focus-visible:ring-emerald-400" on:click={() => toggleSection("threatintel")} aria-expanded={sectionExpanded.threatintel}>
              <div class="flex items-center gap-3">
                <span class="text-base leading-none">🛡️</span>
                <span class="text-sm font-semibold text-gray-100">Threat Intel</span>
                {#if data?.phishing?.valid}
                  <span class="w-2 h-2 rounded-full bg-red-500 flex-shrink-0" title="Confirmed phishing"></span>
                {/if}
              </div>
              <svg class="w-4 h-4 text-gray-500 transition-transform duration-200 {sectionExpanded.threatintel ? 'rotate-180' : ''}" viewBox="0 0 20 20" fill="currentColor"><path fill-rule="evenodd" d={CHEVRON_PATH} clip-rule="evenodd" /></svg>
            </button>
            <div class="transition-all duration-300 ease-in-out {sectionExpanded.threatintel ? 'max-h-[5000px] opacity-100' : 'max-h-0 opacity-0 overflow-hidden'}">
              <div class="acc-body border-t border-gray-800"><ThreatIntelSection phishing={data.phishing} /></div>
            </div>
          </div>
        {/if}

        {#if data.ssl_info || data.tls_info}
          <div id="section-security" class="rounded-xl border border-gray-800 overflow-hidden scroll-mt-20">
            <button type="button" class="w-full flex items-center justify-between px-5 py-4 bg-gray-900/50 hover:bg-gray-900/80 transition-colors text-left focus:outline-none focus-visible:ring-2 focus-visible:ring-inset focus-visible:ring-emerald-400" on:click={() => toggleSection("security")} aria-expanded={sectionExpanded.security}>
              <div class="flex items-center gap-3">
                <span class="text-base leading-none">🔒</span>
                <span class="text-sm font-semibold text-gray-100">Security & SSL</span>
                {#if data?.ssl_info?.IsSuspicious || data?.ssl_info?.KnownBadChain || (data?.ssl_info && !data?.ssl_info?.HasTLS) || (data?.ssl_info?.HasTLS && !data?.ssl_info?.ChainValid) || data?.tls_info?.HostnameMismatch}
                  <span class="w-2 h-2 rounded-full bg-red-500 flex-shrink-0" title="SSL/TLS issues detected"></span>
                {/if}
              </div>
              <svg class="w-4 h-4 text-gray-500 transition-transform duration-200 {sectionExpanded.security ? 'rotate-180' : ''}" viewBox="0 0 20 20" fill="currentColor"><path fill-rule="evenodd" d={CHEVRON_PATH} clip-rule="evenodd" /></svg>
            </button>
            <div class="transition-all duration-300 ease-in-out {sectionExpanded.security ? 'max-h-[5000px] opacity-100' : 'max-h-0 opacity-0 overflow-hidden'}">
              <div class="acc-body border-t border-gray-800"><SecuritySection sslInfo={data.ssl_info} tlsInfo={data.tls_info} /></div>
            </div>
          </div>
        {/if}

        {#if data.content_data}
          <div id="section-content" class="rounded-xl border border-gray-800 overflow-hidden scroll-mt-20">
            <button type="button" class="w-full flex items-center justify-between px-5 py-4 bg-gray-900/50 hover:bg-gray-900/80 transition-colors text-left focus:outline-none focus-visible:ring-2 focus-visible:ring-inset focus-visible:ring-emerald-400" on:click={() => toggleSection("content")} aria-expanded={sectionExpanded.content}>
              <div class="flex items-center gap-3">
                <span class="text-base leading-none">📄</span>
                <span class="text-sm font-semibold text-gray-100">Page Content</span>
                {#if data?.content_data?.brand_check?.is_mismatch || data?.content_data?.has_hidden_iframe || data?.content_data?.forms?.some(f => f.is_external)}
                  <span class="w-2 h-2 rounded-full bg-red-500 flex-shrink-0" title="Content issues detected"></span>
                {:else if data?.content_data?.has_login_form || data?.content_data?.has_payment_form}
                  <span class="w-2 h-2 rounded-full bg-yellow-400 flex-shrink-0" title="Sensitive forms present"></span>
                {/if}
              </div>
              <svg class="w-4 h-4 text-gray-500 transition-transform duration-200 {sectionExpanded.content ? 'rotate-180' : ''}" viewBox="0 0 20 20" fill="currentColor"><path fill-rule="evenodd" d={CHEVRON_PATH} clip-rule="evenodd" /></svg>
            </button>
            <div class="transition-all duration-300 ease-in-out {sectionExpanded.content ? 'max-h-[5000px] opacity-100' : 'max-h-0 opacity-0 overflow-hidden'}">
              <div class="acc-body border-t border-gray-800"><ContentSection contentData={data.content_data} /></div>
            </div>
          </div>
        {/if}

        {#if data.features}
          <div id="section-features" class="rounded-xl border border-gray-800 overflow-hidden scroll-mt-20">
            <button type="button" class="w-full flex items-center justify-between px-5 py-4 bg-gray-900/50 hover:bg-gray-900/80 transition-colors text-left focus:outline-none focus-visible:ring-2 focus-visible:ring-inset focus-visible:ring-emerald-400" on:click={() => toggleSection("features")} aria-expanded={sectionExpanded.features}>
              <div class="flex items-center gap-3">
                <span class="text-base leading-none">📡</span>
                <span class="text-sm font-semibold text-gray-100">URL Signals</span>
                {#if data?.features?.url?.has_homoglyph || data?.features?.url?.contains_punycode || data?.features?.url?.uses_ip || data?.features?.url?.url_shortener || data?.typosquat_result?.is_suspicious || (data?.domain_randomness?.entropy ?? 0) > 3.8}
                  <span class="w-2 h-2 rounded-full bg-red-500 flex-shrink-0" title="URL signal issues detected"></span>
                {:else if data?.features?.tld?.is_risky_tld || (data?.features?.tld && !data?.features?.tld?.is_icann)}
                  <span class="w-2 h-2 rounded-full bg-yellow-400 flex-shrink-0" title="Risky TLD"></span>
                {/if}
              </div>
              <svg class="w-4 h-4 text-gray-500 transition-transform duration-200 {sectionExpanded.features ? 'rotate-180' : ''}" viewBox="0 0 20 20" fill="currentColor"><path fill-rule="evenodd" d={CHEVRON_PATH} clip-rule="evenodd" /></svg>
            </button>
            <div class="transition-all duration-300 ease-in-out {sectionExpanded.features ? 'max-h-[5000px] opacity-100' : 'max-h-0 opacity-0 overflow-hidden'}">
              <div class="acc-body border-t border-gray-800">
                <URLSignalsSection features={data.features} domainRandomness={data.domain_randomness} typosquatResult={data.typosquat_result} />
              </div>
            </div>
          </div>
        {/if}

        {#if data.infrastructure}
          <div id="section-infrastructure" class="rounded-xl border border-gray-800 overflow-hidden scroll-mt-20">
            <button type="button" class="w-full flex items-center justify-between px-5 py-4 bg-gray-900/50 hover:bg-gray-900/80 transition-colors text-left focus:outline-none focus-visible:ring-2 focus-visible:ring-inset focus-visible:ring-emerald-400" on:click={() => toggleSection("infrastructure")} aria-expanded={sectionExpanded.infrastructure}>
              <div class="flex items-center gap-3">
                <span class="text-base leading-none">🖥️</span>
                <span class="text-sm font-semibold text-gray-100">Hosting & Server</span>
                {#if data?.infrastructure && !data?.infrastructure?.nameservers_valid}
                  <span class="w-2 h-2 rounded-full bg-red-500 flex-shrink-0" title="DNS issue detected"></span>
                {/if}
              </div>
              <svg class="w-4 h-4 text-gray-500 transition-transform duration-200 {sectionExpanded.infrastructure ? 'rotate-180' : ''}" viewBox="0 0 20 20" fill="currentColor"><path fill-rule="evenodd" d={CHEVRON_PATH} clip-rule="evenodd" /></svg>
            </button>
            <div class="transition-all duration-300 ease-in-out {sectionExpanded.infrastructure ? 'max-h-[5000px] opacity-100' : 'max-h-0 opacity-0 overflow-hidden'}">
              <div class="acc-body border-t border-gray-800"><InfrastructureSection infrastructure={data.infrastructure} /></div>
            </div>
          </div>
        {/if}

        <!-- Back to top -->
        <div class="pt-2 flex justify-center">
          <button class="inline-flex items-center gap-1.5 px-4 py-2 rounded-full bg-gray-800/50 hover:bg-gray-700 text-gray-500 hover:text-gray-200 text-xs font-medium focus:outline-none focus-visible:ring-2 focus-visible:ring-emerald-400 transition-colors duration-150" on:click={scrollToTop} aria-label="Scroll to top">
            <svg class="w-3.5 h-3.5" viewBox="0 0 20 20" fill="currentColor"><path fill-rule="evenodd" d="M5.23 12.79a.75.75 0 001.06.02L10 8.812l3.71 3.998a.75.75 0 101.08-1.04l-4.25-4.53a.75.75 0 00-1.08 0l-4.25 4.53a.75.75 0 00.02 1.06z" clip-rule="evenodd" /></svg>
            Back to top
          </button>
        </div>
      </div>
    {/if}

    <!-- Performance timings — dev/power user only -->
    {#if data.performance}
      <div class="flex flex-col items-center gap-2 pt-2">
        <button
          type="button"
          class="inline-flex items-center gap-1.5 text-[11px] text-gray-600 hover:text-gray-400 transition-colors focus:outline-none"
          on:click={() => toggleSection("performance")}
        >
          <svg class="w-3 h-3 transition-transform duration-200 {sectionExpanded.performance ? 'rotate-180' : ''}" viewBox="0 0 20 20" fill="currentColor"><path fill-rule="evenodd" d={CHEVRON_PATH} clip-rule="evenodd" /></svg>
          Performance timings
        </button>
        {#if sectionExpanded.performance}
          <div class="w-full acc-body">
            <PerformanceSection performance={data.performance} />
          </div>
        {/if}
      </div>
    {/if}
  </section>
{/if}

<style>
  /* Strip section components' own card wrapper inside accordion — the accordion card provides the container */
  .acc-body > :global(section),
  .acc-body > :global(div) {
    border: none !important;
    border-radius: 0 !important;
    box-shadow: none !important;
    background: transparent !important;
    margin: 0 !important;
  }
  .acc-body > :global(section):hover,
  .acc-body > :global(div):hover {
    transform: none !important;
    box-shadow: none !important;
  }
</style>
