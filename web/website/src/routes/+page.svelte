<script lang="ts">
  import { browser } from "$app/environment";
  import { replaceState } from "$app/navigation";
  import { onMount } from "svelte";
  import { api } from "../lib/api";
  import ResultSection from "../lib/components/ResultSection.svelte";
  import type { AnalyzeResult } from "../lib/types";
  import { formatUrl, formatUrlForShare, isValidUrl } from "../lib/utils";

  // Page load data from +page.ts — runs server-side so bots get correct OG meta tags.
  export let data: { queryDomain: string; queryUrl: string; formattedQueryUrl: string };

  let input = "";
  let loading = false;
  let error: string | null = null;
  let scanResult: AnalyzeResult | null = null;
  let screenshotUrl: string | null = null;
  let screenshotLoading = false;

  type Verdict = "Safe" | "Risky" | "Unclear" | "Suspicious";
  const ACCENTS: Record<Verdict, { ring: string; glow: string; badge: string }> = {
    Safe: {
      ring: "focus:ring-emerald-600",
      glow: "from-emerald-600/20",
      badge: "bg-emerald-600/20 text-emerald-300 border-emerald-700",
    },
    Risky: {
      ring: "focus:ring-red-600",
      glow: "from-red-600/20",
      badge: "bg-red-600/20 text-red-300 border-red-700",
    },
    Unclear: {
      ring: "focus:ring-yellow-600",
      glow: "from-yellow-500/20",
      badge: "bg-yellow-600/20 text-yellow-300 border-yellow-700",
    },
    Suspicious: {
      ring: "focus:ring-yellow-600",
      glow: "from-yellow-500/20",
      badge: "bg-yellow-600/20 text-yellow-300 border-yellow-700",
    },
  };

  function normalizeVerdict(v: string | null | undefined): Verdict {
    switch (v) {
      case "Safe":
      case "Risky":
      case "Unclear":
      case "Suspicious":
        return v;
      default:
        return "Unclear";
    }
  }

  $: verdict = normalizeVerdict(scanResult?.result?.verdict);
  $: accent = ACCENTS[verdict];
  $: isLanding = !scanResult && !loading && !error;

  // currentUrl is only available in the browser; SSR og:url falls back to the canonical base.
  $: currentUrl = browser ? window.location.href : "";

  // shareDomain: scan result domain after a scan, or the SSR-provided query domain for bots.
  $: shareDomain = scanResult?.domain || data.queryDomain;
  // formattedInput: defanged URL for display in meta description.
  $: formattedInput = scanResult?.url
    ? formatUrlForShare(scanResult.url)
    : data.formattedQueryUrl;

  async function runAnalyze(q: string) {
    const url = formatUrl(q);
    if (!isValidUrl(url)) {
      error = "Please enter a valid URL";
      return;
    }

    loading = true;
    error = null;
    scanResult = null;
    if (screenshotUrl) {
      URL.revokeObjectURL(screenshotUrl);
      screenshotUrl = null;
    }
    screenshotLoading = true;

    try {
      // Kick off screenshot in parallel — don't block on it.
      api
        .screenshot(url)
        .then((res) => {
          if (res.data) screenshotUrl = res.data as string;
        })
        .catch(() => console.warn("Screenshot request failed"))
        .finally(() => { screenshotLoading = false; });

      const res = await api.analyze(url);

      if (res.error) {
        error = res.error;
      } else {
        scanResult = res.data as AnalyzeResult;
        const share = new URL(window.location.href);
        share.searchParams.set("q", url);
        replaceState(share.toString(), {});
      }
    } catch (err) {
      error = "Analyze request failed";
    } finally {
      loading = false;
    }
  }

  function onSubmit(e: Event) {
    e.preventDefault();
    runAnalyze(input);
  }

  onMount(() => {
    const params = new URLSearchParams(window.location.search);
    const q = params.get("q");
    if (q) {
      input = q;
      runAnalyze(q);
    }

    const onKey = (e: KeyboardEvent) => {
      if (
        e.key === "/" &&
        !(e.target instanceof HTMLInputElement || e.target instanceof HTMLTextAreaElement)
      ) {
        e.preventDefault();
        const el = document.getElementById("url-input") as HTMLInputElement | null;
        el?.focus();
      }
      if (e.key === "Escape") {
        input = "";
      }
    };
    window.addEventListener("keydown", onKey);
    return () => window.removeEventListener("keydown", onKey);
  });
</script>

<svelte:head>
  {#if shareDomain}
    {@const desc = formattedInput
      ? `Security scan report for ${formattedInput}. Check if this URL is safe, phishing, or suspicious.`
      : `Security scan report for ${shareDomain}. Check if this URL is safe, phishing, or suspicious.`}
    <title>SafeSurf — Is {shareDomain} safe?</title>
    <meta name="description" content={desc} />
    <meta property="og:title" content="SafeSurf — Is {shareDomain} safe?" />
    <meta property="og:description" content={desc} />
    <meta property="og:type" content="website" />
    <meta property="og:url" content={currentUrl || `https://safesurf.vercel.app/?q=${encodeURIComponent(data.queryUrl)}`} />
    <meta property="og:image" content="https://safesurf.vercel.app/safesurf.png" />
    <meta name="twitter:card" content="summary" />
    <meta name="twitter:title" content="SafeSurf — Is {shareDomain} safe?" />
    <meta name="twitter:description" content={desc} />
  {:else}
    <title>SafeSurf — Check if a link is safe</title>
    <meta name="description" content="Instantly check if a URL is safe, phishing, or suspicious. Free URL scanner with threat intelligence." />
    <meta property="og:title" content="SafeSurf — Check if a link is safe" />
    <meta property="og:description" content="Instantly check if a URL is safe, phishing, or suspicious. Free URL scanner with threat intelligence." />
    <meta property="og:type" content="website" />
    <meta property="og:url" content="https://safesurf.vercel.app" />
    <meta property="og:image" content="https://safesurf.vercel.app/safesurf.png" />
    <meta name="twitter:card" content="summary" />
    <meta name="twitter:title" content="SafeSurf — Check if a link is safe" />
    <meta name="twitter:description" content="Instantly check if a URL is safe, phishing, or suspicious." />
  {/if}
</svelte:head>

<section>
  <div
    class={`max-w-4xl mx-auto px-6 ${isLanding ? "min-h-[70vh] flex flex-col justify-center" : "py-12"}`}
  >
    <header
      class="relative mb-14 flex flex-col items-center md:items-start text-center md:text-left"
    >
      <!-- Background accent -->
      <div
        class="absolute -top-10 -left-10 w-40 h-40 bg-blue-600/30 rounded-full blur-3xl animate-blob z-0"
      ></div>
      <div
        class="absolute top-0 right-0 w-32 h-32 bg-emerald-500/20 rounded-full blur-3xl animate-blob animation-delay-2000 z-0"
      ></div>

      <!-- Heading -->
      <h1 class="relative text-6xl md:text-6xl font-extrabold tracking-tight text-white z-10">
        <a
          href="/"
          on:click={() => (location.href = "/")}
          class="bg-clip-text text-transparent bg-gradient-to-r from-blue-400 via-indigo-500 to-purple-500 hover:from-purple-500 hover:to-pink-500 transition-all"
        >
          SafeSurf
        </a>
      </h1>

      <!-- Subheading -->
      <p
        class="relative mt-4 text-gray-300 md:text-lg text-center md:text-left max-w-xl z-10 animate-fadeIn"
      >
        Check if a link is safe.
      </p>
    </header>

    <style>
      /* Blob animation */
      @keyframes blob {
        0%,
        100% {
          transform: translate(0px, 0px) scale(1);
        }
        33% {
          transform: translate(20px, -10px) scale(1.1);
        }
        66% {
          transform: translate(-15px, 15px) scale(0.95);
        }
      }
      .animate-blob {
        animation: blob 8s infinite;
      }
      .animation-delay-2000 {
        animation-delay: 2s;
      }

      /* Fade-in animation for subheading */
      @keyframes fadeIn {
        from {
          opacity: 0;
          transform: translateY(10px);
        }
        to {
          opacity: 1;
          transform: translateY(0);
        }
      }
      .animate-fadeIn {
        animation: fadeIn 1s ease-out forwards;
      }

      /* Autofill fix for dark input added by browsers */
      input:-webkit-autofill,
      input:-webkit-autofill:hover,
      input:-webkit-autofill:focus,
      input:-webkit-autofill:active {
        -webkit-text-fill-color: #e5e7eb; /* Tailwind text-gray-200 */
        transition: background-color 5000s ease-in-out 0s; /* prevent yellow flash */
        box-shadow: 0 0 0px 1000px #1f2937 inset; /* Tailwind bg-gray-900 */
        -webkit-box-shadow: 0 0 0px 1000px #1f2937 inset;
      }
    </style>

    <form
      class="relative bg-gray-950 rounded-xl border border-gray-800 p-6 md:p-8 overflow-hidden"
      on:submit|preventDefault={onSubmit}
    >
      <!-- Background accent -->
      <div
        class="absolute -top-10 -left-10 w-32 h-32 bg-blue-600/20 rounded-full blur-3xl animate-blob z-0"
      ></div>
      <div
        class="absolute -bottom-10 -right-10 w-28 h-28 bg-purple-500/20 rounded-full blur-3xl animate-blob animation-delay-3000 z-0"
      ></div>

      <label for="url-input" class="sr-only">URL to analyze</label>
      <div class="relative flex flex-col gap-2 z-10">
        <div class="flex flex-col md:flex-row gap-3">
          <!-- Input -->
          <input
            id="url-input"
            type="text"
            class={`flex-1 rounded-lg bg-gray-900 border px-4 py-3 text-sm placeholder-gray-500 text-gray-200 focus:outline-none focus:ring-2 focus:ring-offset-0 transition-all duration-200 focus:shadow-lg ${error ? 'border-red-600 focus:ring-red-600' : `border-gray-800 ${accent.ring}`}`}
            placeholder="Enter a URL (e.g. example.com)"
            bind:value={input}
            on:input={() => { if (error) error = null; }}
            autocomplete="url"
            inputmode="url"
            aria-invalid={error ? "true" : undefined}
            aria-describedby={error ? "url-error" : undefined}
            required
          />

          <!-- Button -->
          <button
            type="submit"
            class="inline-flex items-center justify-center gap-2 px-5 py-3 rounded-lg bg-blue-600 hover:bg-blue-500 text-white text-sm font-medium disabled:opacity-50 focus:outline-none focus:ring-2 focus:ring-offset-0 focus:ring-blue-600"
            disabled={loading}
            aria-busy={loading}
            aria-label={loading ? "Scanning URL, please wait" : "Scan Now"}
          >
            {#if loading}
              <svg
                class="w-4 h-4 animate-spin text-white"
                fill="none"
                stroke="currentColor"
                stroke-width="2"
                viewBox="0 0 24 24"
                aria-hidden="true"
              >
                <path
                  stroke-linecap="round"
                  stroke-linejoin="round"
                  d="M12 4v4m0 8v4m8-8h-4M4 12H0"
                />
              </svg>
              Scanning ..
            {:else}
              Scan Now
            {/if}
          </button>
        </div>

        {#if error}
          <p id="url-error" class="text-red-400 text-xs mt-0.5" role="alert">{error}</p>
        {/if}
      </div>
    </form>

    <style>
      /* Blob animation */
      @keyframes blob {
        0%,
        100% {
          transform: translate(0px, 0px) scale(1);
        }
        33% {
          transform: translate(20px, -10px) scale(1.1);
        }
        66% {
          transform: translate(-15px, 15px) scale(0.95);
        }
      }
      .animate-blob {
        animation: blob 8s infinite;
      }
      .animation-delay-3000 {
        animation-delay: 3s;
      }
    </style>

    <div class="mt-8" aria-live="polite">
      <ResultSection data={scanResult} {loading} {error} {screenshotUrl} {screenshotLoading} />
    </div>
  </div>
</section>
