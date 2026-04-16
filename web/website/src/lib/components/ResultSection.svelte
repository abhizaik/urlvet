<script lang="ts">
  import { onDestroy } from "svelte";
  import { browser } from "$app/environment";
  import type { AnalyzeResult } from "../types";
  import VerdictCard from "./sections/VerdictCard.svelte";
  import FlagsGrid from "./sections/FlagsGrid.svelte";
  import ScreenshotViewer from "./sections/ScreenshotViewer.svelte";
  import DomainInfoSection from "./sections/DomainInfoSection.svelte";
  import RedirectionSection from "./sections/RedirectionSection.svelte";
  import ThreatIntelSection from "./sections/ThreatIntelSection.svelte";
  import SecuritySection from "./sections/SecuritySection.svelte";
  import ContentSection from "./sections/ContentSection.svelte";
  import URLSignalsSection from "./sections/URLSignalsSection.svelte";
  import InfrastructureSection from "./sections/InfrastructureSection.svelte";
  import PerformanceSection from "./sections/PerformanceSection.svelte";

  export let data: AnalyzeResult | null = null;
  export let screenshotUrl: string | null = null;
  export let screenshotLoading = false;
  export let loading = false;
  export let error: string | null = null;

  let showAdvanced = false;
  let shareCopied = false;

  $: primary = data?.result;

  onDestroy(() => {
    if (screenshotUrl) {
      URL.revokeObjectURL(screenshotUrl);
    }
  });

  $: availableTabs = [
    { id: "domain", label: "Domain Info", icon: "🏷️", condition: data?.domain_info },
    { id: "analysis", label: "Redirection", icon: "🔀", condition: data?.analysis },
    { id: "threatintel", label: "Threat Intel", icon: "🛡️", condition: data?.phishing },
    {
      id: "security",
      label: "Security/SSL",
      icon: "🔒",
      condition: data?.ssl_info || data?.tls_info,
    },
    { id: "content", label: "Page Content", icon: "📄", condition: data?.content_data },
    { id: "features", label: "URL Signals", icon: "📡", condition: data?.features },
    {
      id: "infrastructure",
      label: "Hosting & Server",
      icon: "🖥️",
      condition: data?.infrastructure,
    },
  ].filter((tab) => tab.condition);

  $: gridColumns = availableTabs.length > 0 ? `repeat(${availableTabs.length}, 1fr)` : "1fr";

  function scrollToSection(sectionId: string) {
    const element = document.getElementById(`section-${sectionId}`);
    if (element) {
      const offset = 80;
      const top = element.getBoundingClientRect().top + window.scrollY - offset;
      window.scrollTo({ top, behavior: "smooth" });
    }
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
    const target = document.getElementById("analysis-summary");
    if (target) {
      const top = target.getBoundingClientRect().top + window.scrollY - 50;
      window.scrollTo({ top, behavior: "smooth" });
    }
  }
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
        class="group inline-flex items-center justify-center gap-2 px-5 py-3 rounded-lg bg-gradient-to-r from-blue-600 to-blue-500 hover:from-blue-500 hover:to-blue-400 text-white text-sm font-semibold shadow-md hover:shadow-lg transition-all duration-200 hover:scale-105 active:scale-95 focus:outline-none focus:ring-2 focus:ring-offset-0 focus:ring-blue-400 focus:ring-offset-gray-950 disabled:opacity-50 {shareCopied
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

    <!-- Verdict & Trust Score -->
    <VerdictCard verdict={primary?.verdict} finalScore={primary?.final_score} />

    <!-- Red / Green Flags -->
    <FlagsGrid reasons={primary?.reasons} />

    <!-- Screenshot -->
    <ScreenshotViewer {screenshotUrl} loading={screenshotLoading} />

    <!-- Advanced Panel Toggle -->
    <div class="mt-6 flex justify-center">
      <button
        id="full-report-button"
        class="inline-flex items-center justify-center gap-2 px-5 py-3 rounded-lg {showAdvanced
          ? 'bg-gray-800 hover:bg-gray-700 text-white'
          : 'bg-blue-600 hover:bg-blue-500 text-white'} text-sm font-medium disabled:opacity-50 focus:outline-none focus:ring-2 focus:ring-offset-0 focus:ring-blue-600 transition-all duration-200"
        on:click={() => (showAdvanced = !showAdvanced)}
        aria-expanded={showAdvanced}
        aria-controls="advanced-panel"
      >
        <svg class="w-4 h-4" viewBox="0 0 20 20" fill="currentColor">
          {#if showAdvanced}
            <path
              fill-rule="evenodd"
              d="M14.77 12.79a.75.75 0 01-1.06-.02L10 8.812l-3.71 3.958a.75.75 0 11-1.08-1.04l4.25-4.53a.75.75 0 011.08 0l4.25 4.53a.75.75 0 01-.02 1.06z"
              clip-rule="evenodd"
            />
          {:else}
            <path
              fill-rule="evenodd"
              d="M5.23 7.21a.75.75 0 011.06.02L10 11.188l3.71-3.958a.75.75 0 111.08 1.04l-4.25 4.53a.75.75 0 01-1.08 0l-4.25-4.53a.75.75 0 01.02-1.06z"
              clip-rule="evenodd"
            />
          {/if}
        </svg>
        <span>{showAdvanced ? "Hide Full Report" : "View Full Report"}</span>
      </button>
    </div>

    <!-- Advanced Panel -->
    <div
      id="advanced-panel"
      class="transition-all duration-500 ease-in-out {showAdvanced
        ? 'max-h-[50000px] opacity-100 mt-4 overflow-visible'
        : 'max-h-0 opacity-0 overflow-hidden'}"
    >
      <div class="rounded-xl border border-gray-800 bg-gray-950 shadow-md overflow-visible">
        <!-- Tab Nav -->
        {#if availableTabs.length > 0}
          <div class="border-b border-gray-800 bg-gray-900/50">
            <div
              class="flex overflow-x-auto scrollbar-hide md:grid"
              style="grid-template-columns: {gridColumns};"
            >
              {#each availableTabs as tab}
                <button
                  class="flex items-center justify-center gap-1.5 md:gap-2 px-3 md:px-4 py-3 md:py-4 text-xs md:text-sm font-medium transition-all duration-200 whitespace-nowrap border-b-2 border-transparent text-gray-400 hover:text-gray-300 hover:bg-gray-900/40 focus:outline-none focus:ring-2 focus:ring-inset focus:ring-emerald-400 flex-shrink-0 min-w-fit"
                  on:click={() => scrollToSection(tab.id)}
                  role="tab"
                >
                  <span class="text-sm md:text-base">{tab.icon}</span>
                  <span>{tab.label}</span>
                </button>
              {/each}
            </div>
          </div>
        {/if}

        <!-- Section Content -->
        <div class="p-6 space-y-6">
          <DomainInfoSection domainInfo={data.domain_info} rank={data.features?.rank} />
          <RedirectionSection analysis={data.analysis} domain={data.domain} />
          <ThreatIntelSection phishing={data.phishing} />
          <SecuritySection sslInfo={data.ssl_info} tlsInfo={data.tls_info} />
          <ContentSection contentData={data.content_data} />
          <URLSignalsSection features={data.features} domainRandomness={data.domain_randomness} typosquatResult={data.typosquat_result} />
          <InfrastructureSection infrastructure={data.infrastructure} />
          <PerformanceSection performance={data.performance} />

          <!-- Scroll to Top -->
          <div class="mt-6 flex justify-center">
            <button
              class="inline-flex items-center justify-center gap-2 px-5 py-3 rounded-full bg-gray-800 hover:bg-gray-700 text-white text-sm font-medium focus:outline-none focus-visible:ring-2 focus-visible:ring-emerald-400 focus-visible:ring-offset-0 transition-colors duration-150"
              on:click={scrollToTop}
              aria-label="Scroll to top"
            >
              <svg class="w-4 h-4" viewBox="0 0 20 20" fill="currentColor">
                <path
                  fill-rule="evenodd"
                  d="M5.23 12.79a.75.75 0 001.06.02L10 8.812l3.71 3.998a.75.75 0 101.08-1.04l-4.25-4.53a.75.75 0 00-1.08 0l-4.25 4.53a.75.75 0 00.02 1.06z"
                  clip-rule="evenodd"
                />
              </svg>
            </button>
          </div>
        </div>
      </div>
    </div>
  </section>
{/if}

<style>
  /* Hide scrollbar for tabs on mobile */
  .scrollbar-hide {
    -ms-overflow-style: none;
    scrollbar-width: none;
  }
  .scrollbar-hide::-webkit-scrollbar {
    display: none;
  }
</style>
