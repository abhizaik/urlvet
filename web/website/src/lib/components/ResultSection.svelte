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
  export let onScanAnother: (() => void) | undefined = undefined;

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

  function computeExpanded(d: AnalyzeResult | null): Record<string, boolean> {
    if (!d)
      return {
        domain: false,
        analysis: false,
        threatintel: false,
        security: false,
        content: false,
        features: false,
        infrastructure: false,
        performance: false,
      };
    return {
      domain: d.domain_info?.age_days !== undefined && d.domain_info.age_days < 365,
      analysis:
        !!d.analysis?.redirection_result?.has_domain_jump ||
        (d.analysis?.redirection_result?.chain_length ?? 0) > 3,
      threatintel: !!d.phishing?.valid,
      security: !!(
        d.ssl_info?.IsSuspicious ||
        d.ssl_info?.KnownBadChain ||
        (d.ssl_info && !d.ssl_info.HasTLS) ||
        (d.ssl_info?.HasTLS && !d.ssl_info.ChainValid) ||
        d.tls_info?.HostnameMismatch
      ),
      content: !!(
        d.content_data?.brand_check?.is_mismatch ||
        d.content_data?.has_hidden_iframe ||
        d.content_data?.forms?.some((f) => f.is_external) ||
        d.content_data?.has_login_form ||
        d.content_data?.has_payment_form
      ),
      features: !!(
        d.features?.url?.has_homoglyph ||
        d.features?.url?.contains_punycode ||
        d.features?.url?.uses_ip ||
        d.features?.url?.url_shortener ||
        d.typosquat_result?.is_suspicious ||
        (d.domain_randomness?.entropy ?? 0) > 3.8 ||
        d.features?.tld?.is_risky_tld ||
        (d.features?.tld && !d.features?.tld?.is_icann && !d.features?.tld?.is_hosting_platform)
      ),
      infrastructure: !!(
        d.infrastructure &&
        !d.infrastructure.nameservers_valid &&
        !d.features?.tld?.is_hosting_platform
      ),
      performance: false,
    };
  }

  $: if (data) sectionExpanded = computeExpanded(data);

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
    const url = window.location.href;
    if (navigator.share) {
      try {
        await navigator.share({ title: document.title, url });
        return;
      } catch {
        // user cancelled or share failed — fall through to clipboard
      }
    }
    try {
      await navigator.clipboard.writeText(url);
      shareCopied = true;
      setTimeout(() => (shareCopied = false), 2000);
    } catch (err) {
      console.error("Clipboard copy failed:", err);
    }
  }

  function scrollToTop() {
    window.scrollTo({ top: 0, behavior: "smooth" });
    setTimeout(() => {
      (document.getElementById("url-input") as HTMLInputElement | null)?.focus();
    }, 400);
  }

  const CHEVRON_PATH =
    "M5.23 7.21a.75.75 0 011.06.02L10 11.188l3.71-3.958a.75.75 0 111.08 1.04l-4.25 4.53a.75.75 0 01-1.08 0l-4.25-4.53a.75.75 0 01.02-1.06z";
</script>

{#if error}
  <div class="max-w-xl mx-auto flex flex-col items-center gap-4 py-10 text-center">
    <div class="w-12 h-12 rounded-full bg-red-500/10 flex items-center justify-center">
      <svg
        class="w-6 h-6 text-red-400"
        fill="none"
        stroke="currentColor"
        stroke-width="1.5"
        viewBox="0 0 24 24"
      >
        <path
          stroke-linecap="round"
          stroke-linejoin="round"
          d="M12 9v3.75m-9.303 3.376c-.866 1.5.217 3.374 1.948 3.374h14.71c1.73 0 2.813-1.874 1.948-3.374L13.949 3.378c-.866-1.5-3.032-1.5-3.898 0L2.697 16.126zM12 15.75h.007v.008H12v-.008z"
        />
      </svg>
    </div>
    <div>
      <p class="text-gray-900 dark:text-white font-semibold mb-1">Scan Failed</p>
      <p class="text-gray-600 dark:text-gray-400 text-sm max-w-sm">{error}</p>
      <p class="text-gray-600 text-xs mt-2">
        The site may be unreachable, blocking automated access, or the URL may be invalid.
      </p>
    </div>
  </div>
{:else if loading}
  <div class="max-w-3xl mx-auto space-y-4">
    <div
      class="animate-pulse rounded-xl border border-gray-300 dark:border-gray-800 bg-gray-100/60 dark:bg-gray-950/60 p-6 h-32"
    ></div>
    <div
      class="animate-pulse rounded-xl border border-gray-300 dark:border-gray-800 bg-gray-100/60 dark:bg-gray-950/60 p-6 h-24"
    ></div>
  </div>
{:else if data}
  <section class="max-w-4xl mx-auto space-y-8 px-4">
    <!-- Header & Share Button -->
    <div
      class="flex flex-col md:flex-row items-start md:items-center justify-between gap-3 animate-fadeIn"
    >
      <div class="flex flex-col">
        <h2 class="text-2xl font-semibold text-gray-900 dark:text-white" id="analysis-summary">
          Analysis Summary
        </h2>
        <div class="flex items-center gap-2 mt-1.5 flex-wrap">
          <span class="text-gray-600 dark:text-gray-400 text-sm">Security profile for</span>
          <span
            class="inline-flex items-center gap-1.5 px-2.5 py-0.5 rounded-full bg-gray-100 dark:bg-gray-800 border border-gray-300 dark:border-gray-700 text-sm font-medium text-gray-800 dark:text-gray-200"
          >
            <svg
              class="w-3.5 h-3.5 text-blue-500 flex-shrink-0"
              fill="none"
              stroke="currentColor"
              stroke-width="2"
              viewBox="0 0 24 24"
            >
              <path
                stroke-linecap="round"
                stroke-linejoin="round"
                d="M12 2a10 10 0 100 20A10 10 0 0012 2zm0 0c2.21 0 4 4.477 4 10s-1.79 10-4 10-4-4.477-4-10 1.79-10 4-10zm-10 10h20"
              />
            </svg>
            {data?.domain}
          </span>
        </div>
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
    <div class="flex flex-col md:flex-row gap-4 items-stretch animate-fadeIn delay-100">
      <div class="flex-1 min-w-0">
        <VerdictCard
          verdict={primary?.verdict}
          finalScore={primary?.final_score}
          unreachable={httpStatusCode === 0 || httpStatusCode === null}
        />
      </div>
      {#if screenshotLoading || screenshotUrl}
        <div
          class="md:w-56 md:flex-shrink-0 rounded-xl border border-gray-300 dark:border-gray-800 bg-gray-50 dark:bg-gray-900/70 p-3 shadow-md flex flex-col gap-2"
        >
          <p class="text-xs font-semibold text-gray-600 dark:text-gray-400">Real-time Preview</p>
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
    <div class="animate-fadeIn delay-200"><FlagsGrid reasons={primary?.reasons} /></div>

    <!-- Detailed Sections -->
    {#if data.domain_info || data.analysis || data.phishing || data.ssl_info || data.tls_info || data.content_data || data.features || data.infrastructure}
      <div class="space-y-2 animate-fadeIn delay-300">
        {#if data.domain_info}
          <div
            id="section-domain"
            class="rounded-xl border border-gray-300 dark:border-gray-800 overflow-hidden scroll-mt-20"
          >
            <button
              type="button"
              class="w-full flex items-center justify-between px-5 py-4 bg-white dark:bg-gray-900/50 hover:bg-gray-50 dark:hover:bg-gray-900/80 transition-colors text-left focus:outline-none focus-visible:ring-2 focus-visible:ring-inset focus-visible:ring-emerald-400"
              on:click={() => toggleSection("domain")}
              aria-expanded={sectionExpanded.domain}
            >
              <div class="flex items-center gap-3">
                <span class="text-base leading-none">🏷️</span>
                <span class="text-sm font-semibold text-gray-800 dark:text-gray-100"
                  >Domain Info</span
                >
                {#if data?.domain_info?.age_days !== undefined && data.domain_info.age_days < 365}
                  <span class="w-2 h-2 rounded-full bg-yellow-400 flex-shrink-0" title="New domain"
                  ></span>
                {/if}
              </div>
              <svg
                class="w-4 h-4 text-gray-600 dark:text-gray-500 transition-transform duration-200 {sectionExpanded.domain
                  ? 'rotate-180'
                  : ''}"
                viewBox="0 0 20 20"
                fill="currentColor"
                ><path fill-rule="evenodd" d={CHEVRON_PATH} clip-rule="evenodd" /></svg
              >
            </button>
            <div
              class="transition-all duration-300 ease-in-out {sectionExpanded.domain
                ? 'max-h-[5000px] opacity-100'
                : 'max-h-0 opacity-0 overflow-hidden'}"
            >
              <div class="acc-body border-t border-gray-100 dark:border-gray-800">
                <DomainInfoSection domainInfo={data.domain_info} rank={data.features?.rank} />
              </div>
            </div>
          </div>
        {/if}

        {#if data.analysis}
          <div
            id="section-analysis"
            class="rounded-xl border border-gray-300 dark:border-gray-800 overflow-hidden scroll-mt-20"
          >
            <button
              type="button"
              class="w-full flex items-center justify-between px-5 py-4 bg-white dark:bg-gray-900/50 hover:bg-gray-50 dark:hover:bg-gray-900/80 transition-colors text-left focus:outline-none focus-visible:ring-2 focus-visible:ring-inset focus-visible:ring-emerald-400"
              on:click={() => toggleSection("analysis")}
              aria-expanded={sectionExpanded.analysis}
            >
              <div class="flex items-center gap-3">
                <span class="text-base leading-none">🔀</span>
                <span class="text-sm font-semibold text-gray-800 dark:text-gray-100"
                  >Redirection</span
                >
                {#if data?.analysis?.redirection_result?.has_domain_jump}
                  <span
                    class="w-2 h-2 rounded-full bg-red-500 flex-shrink-0"
                    title="Domain jump detected"
                  ></span>
                {:else if (data?.analysis?.redirection_result?.chain_length ?? 0) > 3}
                  <span
                    class="w-2 h-2 rounded-full bg-yellow-400 flex-shrink-0"
                    title="Long redirect chain"
                  ></span>
                {/if}
              </div>
              <svg
                class="w-4 h-4 text-gray-600 dark:text-gray-500 transition-transform duration-200 {sectionExpanded.analysis
                  ? 'rotate-180'
                  : ''}"
                viewBox="0 0 20 20"
                fill="currentColor"
                ><path fill-rule="evenodd" d={CHEVRON_PATH} clip-rule="evenodd" /></svg
              >
            </button>
            <div
              class="transition-all duration-300 ease-in-out {sectionExpanded.analysis
                ? 'max-h-[5000px] opacity-100'
                : 'max-h-0 opacity-0 overflow-hidden'}"
            >
              <div class="acc-body border-t border-gray-100 dark:border-gray-800">
                <RedirectionSection analysis={data.analysis} domain={data.domain} />
              </div>
            </div>
          </div>
        {/if}

        {#if data.phishing}
          <div
            id="section-threatintel"
            class="rounded-xl border border-gray-300 dark:border-gray-800 overflow-hidden scroll-mt-20"
          >
            <button
              type="button"
              class="w-full flex items-center justify-between px-5 py-4 bg-white dark:bg-gray-900/50 hover:bg-gray-50 dark:hover:bg-gray-900/80 transition-colors text-left focus:outline-none focus-visible:ring-2 focus-visible:ring-inset focus-visible:ring-emerald-400"
              on:click={() => toggleSection("threatintel")}
              aria-expanded={sectionExpanded.threatintel}
            >
              <div class="flex items-center gap-3">
                <span class="text-base leading-none">🛡️</span>
                <span class="text-sm font-semibold text-gray-800 dark:text-gray-100"
                  >Threat Intel</span
                >
                {#if data?.phishing?.valid}
                  <span
                    class="w-2 h-2 rounded-full bg-red-500 flex-shrink-0"
                    title="Confirmed phishing"
                  ></span>
                {/if}
              </div>
              <svg
                class="w-4 h-4 text-gray-600 dark:text-gray-500 transition-transform duration-200 {sectionExpanded.threatintel
                  ? 'rotate-180'
                  : ''}"
                viewBox="0 0 20 20"
                fill="currentColor"
                ><path fill-rule="evenodd" d={CHEVRON_PATH} clip-rule="evenodd" /></svg
              >
            </button>
            <div
              class="transition-all duration-300 ease-in-out {sectionExpanded.threatintel
                ? 'max-h-[5000px] opacity-100'
                : 'max-h-0 opacity-0 overflow-hidden'}"
            >
              <div class="acc-body border-t border-gray-100 dark:border-gray-800">
                <ThreatIntelSection phishing={data.phishing} />
              </div>
            </div>
          </div>
        {/if}

        {#if data.ssl_info || data.tls_info}
          <div
            id="section-security"
            class="rounded-xl border border-gray-300 dark:border-gray-800 overflow-hidden scroll-mt-20"
          >
            <button
              type="button"
              class="w-full flex items-center justify-between px-5 py-4 bg-white dark:bg-gray-900/50 hover:bg-gray-50 dark:hover:bg-gray-900/80 transition-colors text-left focus:outline-none focus-visible:ring-2 focus-visible:ring-inset focus-visible:ring-emerald-400"
              on:click={() => toggleSection("security")}
              aria-expanded={sectionExpanded.security}
            >
              <div class="flex items-center gap-3">
                <span class="text-base leading-none">🔒</span>
                <span class="text-sm font-semibold text-gray-800 dark:text-gray-100"
                  >Security & SSL</span
                >
                {#if data?.ssl_info?.IsSuspicious || data?.ssl_info?.KnownBadChain || (data?.ssl_info && !data?.ssl_info?.HasTLS) || (data?.ssl_info?.HasTLS && !data?.ssl_info?.ChainValid) || data?.tls_info?.HostnameMismatch}
                  <span
                    class="w-2 h-2 rounded-full bg-red-500 flex-shrink-0"
                    title="SSL/TLS issues detected"
                  ></span>
                {/if}
              </div>
              <svg
                class="w-4 h-4 text-gray-600 dark:text-gray-500 transition-transform duration-200 {sectionExpanded.security
                  ? 'rotate-180'
                  : ''}"
                viewBox="0 0 20 20"
                fill="currentColor"
                ><path fill-rule="evenodd" d={CHEVRON_PATH} clip-rule="evenodd" /></svg
              >
            </button>
            <div
              class="transition-all duration-300 ease-in-out {sectionExpanded.security
                ? 'max-h-[5000px] opacity-100'
                : 'max-h-0 opacity-0 overflow-hidden'}"
            >
              <div class="acc-body border-t border-gray-100 dark:border-gray-800">
                <SecuritySection sslInfo={data.ssl_info} tlsInfo={data.tls_info} />
              </div>
            </div>
          </div>
        {/if}

        {#if data.content_data}
          <div
            id="section-content"
            class="rounded-xl border border-gray-300 dark:border-gray-800 overflow-hidden scroll-mt-20"
          >
            <button
              type="button"
              class="w-full flex items-center justify-between px-5 py-4 bg-white dark:bg-gray-900/50 hover:bg-gray-50 dark:hover:bg-gray-900/80 transition-colors text-left focus:outline-none focus-visible:ring-2 focus-visible:ring-inset focus-visible:ring-emerald-400"
              on:click={() => toggleSection("content")}
              aria-expanded={sectionExpanded.content}
            >
              <div class="flex items-center gap-3">
                <span class="text-base leading-none">📄</span>
                <span class="text-sm font-semibold text-gray-800 dark:text-gray-100"
                  >Page Content</span
                >
                {#if data?.content_data?.brand_check?.is_mismatch || data?.content_data?.has_hidden_iframe || data?.content_data?.forms?.some((f) => f.is_external)}
                  <span
                    class="w-2 h-2 rounded-full bg-red-500 flex-shrink-0"
                    title="Content issues detected"
                  ></span>
                {:else if data?.content_data?.has_login_form || data?.content_data?.has_payment_form}
                  <span
                    class="w-2 h-2 rounded-full bg-yellow-400 flex-shrink-0"
                    title="Sensitive forms present"
                  ></span>
                {/if}
              </div>
              <svg
                class="w-4 h-4 text-gray-600 dark:text-gray-500 transition-transform duration-200 {sectionExpanded.content
                  ? 'rotate-180'
                  : ''}"
                viewBox="0 0 20 20"
                fill="currentColor"
                ><path fill-rule="evenodd" d={CHEVRON_PATH} clip-rule="evenodd" /></svg
              >
            </button>
            <div
              class="transition-all duration-300 ease-in-out {sectionExpanded.content
                ? 'max-h-[5000px] opacity-100'
                : 'max-h-0 opacity-0 overflow-hidden'}"
            >
              <div class="acc-body border-t border-gray-100 dark:border-gray-800">
                <ContentSection contentData={data.content_data} />
              </div>
            </div>
          </div>
        {/if}

        {#if data.features}
          <div
            id="section-features"
            class="rounded-xl border border-gray-300 dark:border-gray-800 overflow-hidden scroll-mt-20"
          >
            <button
              type="button"
              class="w-full flex items-center justify-between px-5 py-4 bg-white dark:bg-gray-900/50 hover:bg-gray-50 dark:hover:bg-gray-900/80 transition-colors text-left focus:outline-none focus-visible:ring-2 focus-visible:ring-inset focus-visible:ring-emerald-400"
              on:click={() => toggleSection("features")}
              aria-expanded={sectionExpanded.features}
            >
              <div class="flex items-center gap-3">
                <span class="text-base leading-none">📡</span>
                <span class="text-sm font-semibold text-gray-800 dark:text-gray-100"
                  >URL Signals</span
                >
                {#if data?.features?.url?.has_homoglyph || data?.features?.url?.contains_punycode || data?.features?.url?.uses_ip || data?.features?.url?.url_shortener || data?.typosquat_result?.is_suspicious || (data?.domain_randomness?.entropy ?? 0) > 3.8}
                  <span
                    class="w-2 h-2 rounded-full bg-red-500 flex-shrink-0"
                    title="URL signal issues detected"
                  ></span>
                {:else if data?.features?.tld?.is_risky_tld || (data?.features?.tld && !data?.features?.tld?.is_icann)}
                  <span class="w-2 h-2 rounded-full bg-yellow-400 flex-shrink-0" title="Risky TLD"
                  ></span>
                {/if}
              </div>
              <svg
                class="w-4 h-4 text-gray-600 dark:text-gray-500 transition-transform duration-200 {sectionExpanded.features
                  ? 'rotate-180'
                  : ''}"
                viewBox="0 0 20 20"
                fill="currentColor"
                ><path fill-rule="evenodd" d={CHEVRON_PATH} clip-rule="evenodd" /></svg
              >
            </button>
            <div
              class="transition-all duration-300 ease-in-out {sectionExpanded.features
                ? 'max-h-[5000px] opacity-100'
                : 'max-h-0 opacity-0 overflow-hidden'}"
            >
              <div class="acc-body border-t border-gray-100 dark:border-gray-800">
                <URLSignalsSection
                  features={data.features}
                  domainRandomness={data.domain_randomness}
                  typosquatResult={data.typosquat_result}
                />
              </div>
            </div>
          </div>
        {/if}

        {#if data.infrastructure}
          <div
            id="section-infrastructure"
            class="rounded-xl border border-gray-300 dark:border-gray-800 overflow-hidden scroll-mt-20"
          >
            <button
              type="button"
              class="w-full flex items-center justify-between px-5 py-4 bg-white dark:bg-gray-900/50 hover:bg-gray-50 dark:hover:bg-gray-900/80 transition-colors text-left focus:outline-none focus-visible:ring-2 focus-visible:ring-inset focus-visible:ring-emerald-400"
              on:click={() => toggleSection("infrastructure")}
              aria-expanded={sectionExpanded.infrastructure}
            >
              <div class="flex items-center gap-3">
                <span class="text-base leading-none">🖥️</span>
                <span class="text-sm font-semibold text-gray-800 dark:text-gray-100"
                  >Hosting & Server</span
                >
                {#if data?.infrastructure && !data?.infrastructure?.nameservers_valid && !data?.features?.tld?.is_hosting_platform}
                  <span
                    class="w-2 h-2 rounded-full bg-red-500 flex-shrink-0"
                    title="DNS issue detected"
                  ></span>
                {/if}
              </div>
              <svg
                class="w-4 h-4 text-gray-600 dark:text-gray-500 transition-transform duration-200 {sectionExpanded.infrastructure
                  ? 'rotate-180'
                  : ''}"
                viewBox="0 0 20 20"
                fill="currentColor"
                ><path fill-rule="evenodd" d={CHEVRON_PATH} clip-rule="evenodd" /></svg
              >
            </button>
            <div
              class="transition-all duration-300 ease-in-out {sectionExpanded.infrastructure
                ? 'max-h-[5000px] opacity-100'
                : 'max-h-0 opacity-0 overflow-hidden'}"
            >
              <div class="acc-body border-t border-gray-100 dark:border-gray-800">
                <InfrastructureSection
                  infrastructure={data.infrastructure}
                  isHostingPlatform={!!data.features?.tld?.is_hosting_platform}
                />
              </div>
            </div>
          </div>
        {/if}

        <!-- Back to top -->
        <div class="pt-2 flex justify-center">
          <button
            class="inline-flex items-center gap-1.5 px-4 py-2 rounded-full bg-gray-100 dark:bg-gray-800/50 hover:bg-gray-200 dark:hover:bg-gray-700 text-gray-500 hover:text-gray-700 dark:hover:text-gray-200 text-xs font-medium focus:outline-none focus-visible:ring-2 focus-visible:ring-emerald-400 transition-colors duration-150"
            on:click={scrollToTop}
            aria-label="Scroll to top"
          >
            <svg class="w-3.5 h-3.5" viewBox="0 0 20 20" fill="currentColor"
              ><path
                fill-rule="evenodd"
                d="M5.23 12.79a.75.75 0 001.06.02L10 8.812l3.71 3.998a.75.75 0 101.08-1.04l-4.25-4.53a.75.75 0 00-1.08 0l-4.25 4.53a.75.75 0 00.02 1.06z"
                clip-rule="evenodd"
              /></svg
            >
            Back to top
          </button>
        </div>
      </div>
    {/if}

    <!-- Scan another CTA -->
    <div
      class="flex flex-col items-center gap-3 pt-4 pb-2 border-t border-gray-300 dark:border-gray-800/50"
    >
      <p class="text-sm text-gray-600 dark:text-gray-500">Want to scan another link?</p>
      <button
        type="button"
        class="inline-flex items-center gap-2 px-6 py-3 rounded-xl bg-gradient-to-r from-blue-600 to-indigo-600 hover:from-blue-500 hover:to-indigo-500 text-white text-sm font-semibold shadow-lg shadow-blue-900/30 transition-all duration-200 active:scale-95 focus:outline-none focus:ring-2 focus:ring-blue-500"
        on:click={() => {
          onScanAnother ? onScanAnother() : scrollToTop();
        }}
      >
        <svg class="w-4 h-4" fill="none" stroke="currentColor" stroke-width="2" viewBox="0 0 24 24">
          <path
            stroke-linecap="round"
            stroke-linejoin="round"
            d="M21 21l-4.35-4.35M17 11A6 6 0 1 1 5 11a6 6 0 0 1 12 0z"
          />
        </svg>
        Scan Another Link
      </button>
    </div>

    <!-- Performance timings — dev/power user only -->
    {#if data.performance}
      <div class="flex flex-col items-center gap-2 pt-2">
        <button
          type="button"
          class="inline-flex items-center gap-1.5 text-[11px] text-gray-600 hover:text-gray-400 transition-colors focus:outline-none"
          on:click={() => toggleSection("performance")}
        >
          <svg
            class="w-3 h-3 transition-transform duration-200 {sectionExpanded.performance
              ? 'rotate-180'
              : ''}"
            viewBox="0 0 20 20"
            fill="currentColor"
            ><path fill-rule="evenodd" d={CHEVRON_PATH} clip-rule="evenodd" /></svg
          >
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
  button[aria-expanded="true"] {
    box-shadow: inset 3px 0 0 #3b82f6;
  }
  button[aria-expanded="false"] {
    box-shadow: inset 3px 0 0 transparent;
    transition: box-shadow 0.2s ease;
  }
  button[aria-expanded="true"] {
    transition: box-shadow 0.15s ease;
  }

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
