<script lang="ts">
  import { browser } from "$app/environment";
  import TooltipIcon from "../TooltipIcon.svelte";
  export let analysis: any;
  export let domain: string;

  function openAnalyzeInNewTab(url: string) {
    if (!browser) return;
    const analyzeUrl = new URL(window.location.origin);
    analyzeUrl.searchParams.set("q", url);
    window.open(analyzeUrl.toString(), "_blank");
  }
</script>

{#if analysis}
  <section
    id="section-analysis"
    class="bg-gray-900/80 border border-gray-800 rounded-lg p-5 shadow-md hover:shadow-lg hover:scale-[1.01] transition-all scroll-mt-20"
  >
    <div class="flex items-center justify-between mb-4">
      <h3 class="text-base font-semibold text-white">Redirection & Resolution</h3>
      <span
        class="text-[10px] text-gray-400 uppercase tracking-wide px-2 py-0.5 bg-gray-800 rounded"
        >HTTP / Redirects</span
      >
    </div>

    <div class="space-y-0 divide-y divide-gray-800 text-sm text-gray-200 max-w-4xl w-full mx-auto">
      {#if analysis.redirection_result}
        <div
          class="flex flex-col md:grid md:grid-cols-[minmax(0,280px),1fr] md:items-center gap-2 md:gap-4 py-2 first:pt-0 last:pb-0"
        >
          <div class="flex items-center gap-1 text-gray-400">
            <span>Redirected:</span>
            <TooltipIcon
              text="Indicates whether visiting given URL automatically forwards you to another URL."
            />
          </div>
          {#if analysis.redirection_result.is_redirected}
            <span class="font-medium text-white break-all">Yes</span>
          {:else}
            <span class="font-medium text-white break-all">No</span>
          {/if}
        </div>

        {#if analysis.redirection_result.final_url}
          <div
            class="flex flex-col md:grid md:grid-cols-[minmax(0,280px),1fr] gap-2 md:gap-4 py-2 first:pt-0 last:pb-0"
          >
            <div class="flex items-center gap-1 text-gray-400">
              <span>Final URL Domain:</span>
              <TooltipIcon
                text="The domain where the visitor finally lands after any redirects. Useful to detect domain changes or phishing redirects."
              />
            </div>
            <span class="font-medium text-white break-all"
              >{analysis.redirection_result.final_url_domain}</span
            >
          </div>
          <div
            class="flex flex-col md:grid md:grid-cols-[minmax(0,280px),1fr] gap-2 md:gap-4 py-2 first:pt-0 last:pb-0"
          >
            <div class="flex items-center gap-1 text-gray-400">
              <span>Final URL:</span>
              <TooltipIcon
                text="The complete URL where the user ends up after all redirections."
              />
            </div>
            <span class="font-medium text-white break-all"
              >{analysis.redirection_result.final_url}</span
            >
          </div>
        {/if}

        <div
          class="flex flex-col md:grid md:grid-cols-[minmax(0,280px),1fr] md:items-center gap-2 md:gap-4 py-2 first:pt-0 last:pb-0"
        >
          <div class="flex items-center gap-1 text-gray-400">
            <span>Domain Jumped to Another Domain:</span>
            <TooltipIcon
              text="Checks if the website redirects to a completely different domain, which can indicate phishing or tracking."
            />
          </div>
          {#if analysis.redirection_result.has_domain_jump}
            <span class="text-red-400 font-medium flex items-center gap-1">❌ Yes</span>
          {:else}
            <span class="text-green-400 font-medium flex items-center gap-1">✅ No</span>
          {/if}
        </div>

        <div
          class="flex flex-col md:grid md:grid-cols-[minmax(0,280px),1fr] md:items-center gap-2 md:gap-4 py-2 first:pt-0 last:pb-0"
        >
          <div class="flex items-center gap-1 text-gray-400">
            <span>Redirection Chain Length:</span>
            <TooltipIcon
              text="Shows how many redirect steps the website takes before reaching the final destination."
            />
          </div>
          <span class="font-medium text-white">{analysis.redirection_result.chain_length}</span>
        </div>

        {#if analysis.redirection_result.chain?.length}
          <div
            class="flex flex-col md:grid md:grid-cols-[minmax(0,280px),1fr] gap-2 md:gap-4 py-2 first:pt-0 last:pb-0"
          >
            <div class="flex items-center gap-1 text-gray-400">
              <span>Redirection Chain:</span>
              <TooltipIcon
                text="A step-by-step list of all URLs in the redirection path. Warning icons highlight jumps to unexpected domains. Click any URL to analyze it in a new tab."
              />
            </div>

            {#if !analysis.redirection_result.has_domain_jump}
              <ul class="text-sm text-gray-100 list-none">
                {#each analysis.redirection_result.chain as url, index}
                  <li class="break-all flex items-center gap-2 mb-1">
                    <span class="text-gray-400">{index + 1}.</span>
                    <span class="font-medium text-white">{url}</span>
                    {#if !url.includes(domain)}
                      <span class="text-red-400 text-xs">⚠️</span>
                    {/if}
                  </li>
                {/each}
              </ul>
            {:else}
              <div class="flex flex-col gap-2">
                <p class="text-sm text-gray-300 italic">
                  Click on the URLs to perform a safe scan on them
                </p>
                <ul class="text-sm text-gray-100 list-none">
                  {#each analysis.redirection_result.chain as url, index}
                    <li class="break-all flex items-center gap-2 mb-1">
                      <span class="text-gray-400">{index + 1}.</span>
                      <button
                        type="button"
                        class="font-medium text-blue-400 hover:text-blue-300 underline cursor-pointer text-left focus:outline-none focus:ring-2 focus:ring-blue-500 focus:ring-offset-2 focus:ring-offset-gray-900 rounded px-1 transition-colors"
                        on:click={() => openAnalyzeInNewTab(url)}
                        title="Click to analyze this URL in a new tab"
                        aria-label={`Analyze ${url} in a new tab`}
                      >
                        {url}
                      </button>
                      {#if !url.includes(domain)}
                        <span class="text-red-400 text-xs">⚠️</span>
                      {/if}
                    </li>
                  {/each}
                </ul>
              </div>
            {/if}
          </div>
        {/if}
      {/if}

      {#if analysis.http_status}
        <div
          class="flex flex-col md:grid md:grid-cols-[minmax(0,280px),1fr] md:items-center gap-2 md:gap-4 py-2 first:pt-0 last:pb-0"
        >
          <div class="flex items-center gap-1 text-gray-400">
            <span>HTTP Status Code:</span>
            <TooltipIcon
              text="The server response code returned when accessing the URL (e.g., 200 = OK, 404 = Not Found)."
            />
          </div>
          <span class="font-medium text-white"
            >{analysis.http_status.code} {analysis.http_status.text}</span
          >
        </div>

        <div
          class="flex flex-col md:grid md:grid-cols-[minmax(0,280px),1fr] md:items-center gap-2 md:gap-4 py-2 first:pt-0 last:pb-0"
        >
          <div class="flex items-center gap-1 text-gray-400">
            <span>Redirection Status Code (3xx):</span>
            <TooltipIcon
              text="Indicates whether the HTTP status is a redirection (3xx) code, which automatically sends visitors to another URL."
            />
          </div>
          {#if analysis.http_status.is_redirect}
            <span class="text-red-400 font-medium flex items-center gap-1">❌ Yes</span>
          {:else}
            <span class="text-green-400 font-medium flex items-center gap-1">✅ No</span>
          {/if}
        </div>
      {/if}

      {#if analysis.is_hsts_supported !== undefined}
        <div
          class="flex flex-col md:grid md:grid-cols-[minmax(0,280px),1fr] md:items-center gap-2 md:gap-4 py-2 first:pt-0 last:pb-0"
        >
          <div class="flex items-center gap-1 text-gray-400">
            <span>HSTS Supported (HTTPS Only):</span>
            <TooltipIcon
              text="Shows if the website enforces HTTPS connections automatically to improve security and prevent attacks."
            />
          </div>
          {#if analysis.is_hsts_supported}
            <span class="text-green-400 font-medium flex items-center gap-1">✅ Yes</span>
          {:else}
            <span class="text-red-400 font-medium flex items-center gap-1">❌ No</span>
          {/if}
        </div>
      {/if}
    </div>
  </section>
{/if}
