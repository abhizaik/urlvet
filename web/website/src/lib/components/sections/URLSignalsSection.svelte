<script lang="ts">
  import TooltipIcon from "../TooltipIcon.svelte";
  import type { DomainRandomness } from "../../types";
  export let features: any;
  export let domainRandomness: DomainRandomness | undefined;
</script>

{#if features}
  <section
    id="section-features"
    class="bg-gray-900/80 border border-gray-800 rounded-lg p-5 shadow-md hover:shadow-lg hover:scale-[1.01] transition-all scroll-mt-20"
  >
    <div class="flex items-center justify-between mb-4">
      <h3 class="text-base font-semibold text-white">URL Signals</h3>
      <span
        class="text-[10px] text-gray-400 uppercase tracking-wide px-2 py-0.5 bg-gray-800 rounded"
        >URL / TLD</span
      >
    </div>

    <div class="space-y-0 divide-y divide-gray-800 text-sm text-gray-200 max-w-4xl w-full mx-auto">
      {#if features.tld}
        <div
          class="flex flex-col md:grid md:grid-cols-[350px,1fr] gap-2 md:gap-4 py-2 first:pt-0 last:pb-0"
        >
          <div class="flex items-center gap-1 text-gray-400">
            <span>Domain Ending (TLD):</span>
            <TooltipIcon
              text="The last part of a domain name (like .com, .org, .io). It can hint at the site's trust level or purpose."
            />
          </div>
          <span class="font-medium text-white">.{features.tld.tld}</span>
        </div>

        <div
          class="flex flex-col md:grid md:grid-cols-[350px,1fr] md:items-center gap-2 md:gap-4 py-2 first:pt-0 last:pb-0"
        >
          <div class="flex items-center gap-1 text-gray-400">
            <span>High Trust Domain Ending:</span>
            <TooltipIcon
              text="Indicates whether this domain ending (TLD) is widely recognized and commonly used by highly trusted entities like government and other institutions."
            />
          </div>
          {#if features.tld.is_trusted_tld}
            <span class="text-green-400 font-medium flex items-center gap-1">✅ Yes</span>
          {:else}
            <span class="text-red-400 font-medium flex items-center gap-1">❌ No</span>
          {/if}
        </div>

        <div
          class="flex flex-col md:grid md:grid-cols-[350px,1fr] md:items-center gap-2 md:gap-4 py-2 first:pt-0 last:pb-0"
        >
          <div class="flex items-center gap-1 text-gray-400">
            <span>Is Risky TLD:</span>
            <TooltipIcon
              text="Some TLDs are frequently abused by scammers or malicious sites. 'Yes' suggests caution."
            />
          </div>
          {#if features.tld.is_risky_tld}
            <span class="text-red-400 font-medium flex items-center gap-1">❌ Yes</span>
          {:else}
            <span class="text-green-400 font-medium flex items-center gap-1">✅ No</span>
          {/if}
        </div>

        <div
          class="flex flex-col md:grid md:grid-cols-[350px,1fr] md:items-center gap-2 md:gap-4 py-2 first:pt-0 last:pb-0"
        >
          <div class="flex items-center gap-1 text-gray-400">
            <span>TLD Recognized by ICANN:</span>
            <TooltipIcon
              text="ICANN oversees global domain names. Recognition means this TLD is officially managed and monitored."
            />
          </div>
          {#if features.tld.is_icann}
            <span class="text-green-400 font-medium flex items-center gap-1">✅ Yes</span>
          {:else}
            <span class="text-red-400 font-medium flex items-center gap-1">❌ No</span>
          {/if}
        </div>
      {/if}

      {#if features.url}
        <div
          class="flex flex-col md:grid md:grid-cols-[350px,1fr] md:items-center gap-2 md:gap-4 py-2 first:pt-0 last:pb-0"
        >
          <div class="flex items-center gap-1 text-gray-400">
            <span>Uses Link Shortener:</span>
            <TooltipIcon
              text="Shows if the URL uses a shortening service (like bit.ly). Shortened links can hide a site's real destination."
            />
          </div>
          {#if features.url.url_shortener}
            <span class="text-red-400 font-medium flex items-center gap-1">❌ Yes</span>
          {:else}
            <span class="text-green-400 font-medium flex items-center gap-1">✅ No</span>
          {/if}
        </div>

        <div
          class="flex flex-col md:grid md:grid-cols-[350px,1fr] md:items-center gap-2 md:gap-4 py-2 first:pt-0 last:pb-0"
        >
          <div class="flex items-center gap-1 text-gray-400">
            <span>Uses Direct IP Address:</span>
            <TooltipIcon
              text="Some malicious sites use IP addresses instead of domain names to avoid detection. Legitimate websites rarely do this."
            />
          </div>
          {#if features.url.uses_ip}
            <span class="text-red-400 font-medium flex items-center gap-1">❌ Yes</span>
          {:else}
            <span class="text-green-400 font-medium flex items-center gap-1">✅ No</span>
          {/if}
        </div>

        <div
          class="flex flex-col md:grid md:grid-cols-[350px,1fr] md:items-center gap-2 md:gap-4 py-2 first:pt-0 last:pb-0"
        >
          <div class="flex items-center gap-1 text-gray-400">
            <span>Contains Punycode Characters:</span>
            <TooltipIcon
              text="Punycode allows special or non-Latin characters in domains (like xn--example). Scammers sometimes exploit this for lookalike attacks."
            />
          </div>
          {#if features.url.contains_punycode}
            <span class="text-red-400 font-medium flex items-center gap-1">❌ Yes</span>
          {:else}
            <span class="text-green-400 font-medium flex items-center gap-1">✅ No</span>
          {/if}
        </div>

        <div
          class="flex flex-col md:grid md:grid-cols-[350px,1fr] md:items-center gap-2 md:gap-4 py-2 first:pt-0 last:pb-0"
        >
          <div class="flex items-center gap-1 text-gray-400">
            <span>URL Too Long:</span>
            <TooltipIcon
              text="Very long URLs can be used to obscure malicious content or trick users into trusting a fake site."
            />
          </div>
          {#if features.url.too_long}
            <span class="text-red-400 font-medium flex items-center gap-1">❌ Yes</span>
          {:else}
            <span class="text-green-400 font-medium flex items-center gap-1">✅ No</span>
          {/if}
        </div>

        <div
          class="flex flex-col md:grid md:grid-cols-[350px,1fr] md:items-center gap-2 md:gap-4 py-2 first:pt-0 last:pb-0"
        >
          <div class="flex items-center gap-1 text-gray-400">
            <span>URL Too Deep (Many Slashes):</span>
            <TooltipIcon
              text="A URL with many nested paths may indicate redirections or hidden content, often seen in phishing attempts."
            />
          </div>
          {#if features.url.too_deep}
            <span class="text-red-400 font-medium flex items-center gap-1">❌ Yes</span>
          {:else}
            <span class="text-green-400 font-medium flex items-center gap-1">✅ No</span>
          {/if}
        </div>

        <div
          class="flex flex-col md:grid md:grid-cols-[350px,1fr] md:items-center gap-2 md:gap-4 py-2 first:pt-0 last:pb-0"
        >
          <div class="flex items-center gap-1 text-gray-400">
            <span>Has Look-Alike Letters (Homoglyph):</span>
            <TooltipIcon
              text="Detects if the domain name includes characters that look like others (e.g., 'go0gle.com'). Often used for impersonation."
            />
          </div>
          {#if features.url.has_homoglyph}
            <span class="text-red-400 font-medium flex items-center gap-1">❌ Yes</span>
          {:else}
            <span class="text-green-400 font-medium flex items-center gap-1">✅ No</span>
          {/if}
        </div>

        <div
          class="flex flex-col md:grid md:grid-cols-[350px,1fr] md:items-center gap-2 md:gap-4 py-2 first:pt-0 last:pb-0"
        >
          <div class="flex items-center gap-1 text-gray-400">
            <span>Subdomain Count:</span>
            <TooltipIcon
              text="Shows how many subdomains (like shop.example.com) are used. Too many can hint at suspicious or temporary setups."
            />
          </div>
          <span class="font-medium text-white">{features.url.subdomain_count}</span>
        </div>
      {/if}

      {#if domainRandomness && domainRandomness.entropy !== undefined}
        <div
          class="flex flex-col md:grid md:grid-cols-[350px,1fr] md:items-center gap-2 md:gap-4 py-2 first:pt-0 last:pb-0"
        >
          <div class="flex items-center gap-1 text-gray-400">
            <span>Domain Randomness (Entropy):</span>
            <TooltipIcon
              text="Measures the unpredictability of the domain name. High entropy often indicates randomly generated domains used by malware (DGAs)."
            />
          </div>
          <div class="flex items-center gap-2">
            <span class="font-medium text-white">{domainRandomness.entropy.toFixed(2)}</span>
            {#if domainRandomness.entropy > 3.8}
              <span class="text-xs px-1.5 py-0.5 bg-red-900/30 text-red-400 rounded"
                >High Entropy</span
              >
            {/if}
          </div>
        </div>
      {/if}
    </div>
  </section>
{/if}
