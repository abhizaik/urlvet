<script lang="ts">
  import { browser } from "$app/environment";
  import { replaceState } from "$app/navigation";
  import { onMount } from "svelte";
  import { api } from "../lib/api";
  import ResultSection from "../lib/components/ResultSection.svelte";
  import ScanProgress from "../lib/components/ScanProgress.svelte";
  import Shoutouts from "../lib/components/Shoutouts.svelte";
  import type { AnalyzeResult } from "../lib/types";
  import { formatUrl, isValidUrl } from "../lib/utils";

  // Page load data from +page.ts — runs server-side so bots get correct OG meta tags.
  export let data: {
    queryDomain: string;
    queryUrl: string;
    formattedQueryUrl: string;
    verdict: string;
    score: string;
  };

  let input = "";
  let loading = false;
  let error: string | null = null;
  let formError: string | null = null;
  let scanResult: AnalyzeResult | null = null;
  let screenshotUrl: string | null = null;
  let screenshotLoading = false;
  let screenshotFailed = false;
  let scanDone = false;

  type Verdict = "Safe" | "Risky" | "Suspicious";
  const ACCENT_RING: Record<Verdict, string> = {
    Safe: "focus:ring-emerald-600",
    Risky: "focus:ring-red-600",
    Suspicious: "focus:ring-yellow-600",
  };

  function normalizeVerdict(v: string | null | undefined): Verdict {
    if (v === "Safe" || v === "Risky" || v === "Suspicious") return v;
    return "Suspicious";
  }

  $: verdict = normalizeVerdict(scanResult?.result?.verdict);
  $: accentRing = ACCENT_RING[verdict];
  $: isLanding = !scanResult && !loading && !error && !formError;
  $: currentUrl = browser ? window.location.href : "";
  $: shareDomain = scanResult?.domain || data.queryDomain;

  let justPasted = false;

  async function pasteFromClipboard() {
    try {
      const text = await navigator.clipboard.readText();
      if (text) {
        input = text.trim();
        error = null;
        formError = null;
        justPasted = true;
        setTimeout(() => (justPasted = false), 1500);
      }
    } catch {
      /* clipboard access denied */
    }
  }

  async function runAnalyze(q: string) {
    const url = formatUrl(q);
    if (!isValidUrl(url)) {
      formError = "Please enter a valid URL";
      return;
    }

    loading = true;
    scanDone = false;
    error = null;
    formError = null;
    scanResult = null;
    if (screenshotUrl) {
      URL.revokeObjectURL(screenshotUrl);
      screenshotUrl = null;
    }
    screenshotLoading = true;
    screenshotFailed = false;

    try {
      api
        .screenshot(url)
        .then((res) => {
          if (res.data) screenshotUrl = res.data as string;
          else screenshotFailed = true;
        })
        .catch(() => {
          screenshotFailed = true;
        })
        .finally(() => {
          screenshotLoading = false;
        });

      const res = await api.analyze(url);
      if (res.error) {
        error = res.error;
      } else {
        scanResult = res.data as AnalyzeResult;
        const share = new URL(window.location.href);
        share.searchParams.set("q", url);
        if (scanResult.result?.verdict) share.searchParams.set("v", scanResult.result.verdict);
        if (scanResult.result?.final_score !== undefined)
          share.searchParams.set("s", String(scanResult.result.final_score));
        replaceState(share.toString(), {});
      }
    } catch {
      error = "Analyze request failed";
    } finally {
      loading = false;
      scanDone = true;
      setTimeout(() => (scanDone = false), 1200);
    }
  }

  onMount(() => {
    const q = new URLSearchParams(window.location.search).get("q");
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
        (document.getElementById("url-input") as HTMLInputElement | null)?.focus();
      }
      if (e.key === "Escape") input = "";
    };
    window.addEventListener("keydown", onKey);
    return () => window.removeEventListener("keydown", onKey);
  });
</script>

<svelte:head>
  {#if shareDomain}
    {@const ogVerdict = scanResult?.result?.verdict || data.verdict}
    {@const ogScore =
      scanResult?.result?.final_score ?? (data.score ? Number(data.score) : undefined)}
    {@const desc = ogVerdict
      ? `SafeSurf verdict: ${ogVerdict} — see the full breakdown for ${shareDomain}.`
      : `SafeSurf scanned ${shareDomain}. Is it safe to open? See the full phishing detection report.`}
    {@const ogImage = `https://safesurf.xorwave.com/og?domain=${encodeURIComponent(shareDomain)}${ogVerdict ? `&v=${encodeURIComponent(ogVerdict)}` : ""}${ogScore !== undefined ? `&s=${ogScore}` : ""}`}
    <title>SafeSurf — Is {shareDomain} safe?</title>
    <meta name="description" content={desc} />
    <meta property="og:title" content="SafeSurf — Is {shareDomain} safe?" />
    <meta property="og:description" content={desc} />
    <meta property="og:type" content="website" />
    <meta
      property="og:url"
      content={currentUrl || `https://safesurf.xorwave.com/?q=${encodeURIComponent(data.queryUrl)}`}
    />
    <meta property="og:image" content={ogImage} />
    <meta property="og:image:width" content="1200" />
    <meta property="og:image:height" content="630" />
    <meta name="twitter:card" content="summary_large_image" />
    <meta name="twitter:title" content="SafeSurf — Is {shareDomain} safe?" />
    <meta name="twitter:description" content={desc} />
    <meta name="twitter:image" content={ogImage} />
  {:else}
    <title>SafeSurf — Instantly analyze any link.</title>
    <meta
      name="description"
      content="Instantly check whether this URL is safe, suspicious, or risky. Free URL scanner with threat intelligence."
    />
    <meta
      property="og:title"
      content="SafeSurf — Instantly analyze any link for phishing, hidden threats, and suspicious behavior before you open it."
    />
    <meta
      property="og:description"
      content="Instantly check whether this URL is safe, suspicious, or risky. Free URL scanner with threat intelligence."
    />
    <meta property="og:type" content="website" />
    <meta property="og:url" content="https://safesurf.xorwave.com" />
    <meta property="og:image" content="https://safesurf.xorwave.com/safesurf.png" />
    <meta name="twitter:card" content="summary" />
    <meta
      name="twitter:title"
      content="SafeSurf — Instantly analyze any link for phishing, hidden threats, and suspicious behavior before you open it."
    />
    <meta
      name="twitter:description"
      content="Instantly check if a URL is safe, suspicious, or risky."
    />
  {/if}
</svelte:head>

<section>
  <div
    class={`max-w-4xl mx-auto px-6 ${isLanding ? "min-h-[80vh] flex flex-col items-center justify-center text-center pt-10 pb-6" : "py-12"}`}
  >
    <header class="relative mb-14 flex flex-col items-center text-center">
      <div
        class="absolute -top-16 -left-8 w-48 h-48 bg-blue-600/20 rounded-full blur-3xl animate-blob z-0"
      ></div>
      <div
        class="absolute -top-8 right-0 w-36 h-36 bg-emerald-500/15 rounded-full blur-3xl animate-blob animation-delay-2000 z-0"
      ></div>

      <h1 class="relative text-5xl md:text-6xl font-extrabold tracking-tight text-white z-10">
        <a
          href="/"
          on:click={() => (location.href = "/")}
          class="bg-clip-text text-transparent bg-gradient-to-r from-blue-400 via-indigo-500 to-purple-500 hover:from-purple-500 hover:to-pink-500 transition-all"
        >
          SafeSurf
        </a>
      </h1>

      <p
        class="relative mt-3 text-gray-300 text-base md:text-lg font-normal leading-relaxed tracking-wide max-w-md z-10 animate-fadeIn"
      >
        Is this link safe? Find out in seconds.
      </p>
    </header>

    <form
      class={isLanding ? "w-full max-w-2xl" : "w-full"}
      on:submit|preventDefault={() => runAnalyze(input)}
    >
      <div class="flex flex-col sm:flex-row gap-3">
        <div class="relative flex-1">
          <input
            id="url-input"
            type="text"
            class={`w-full rounded-xl bg-gray-900 border px-4 py-3.5 pr-24 text-sm placeholder-gray-500 text-gray-200 focus:outline-none focus:ring-2 focus:ring-offset-0 transition-all duration-200 ${formError ? "border-red-600/70 focus:ring-red-600" : `border-gray-700/80 ${accentRing}`}`}
            placeholder="Enter a link (e.g. example.com)"
            bind:value={input}
            on:input={() => {
              if (formError) formError = null;
            }}
            autocomplete="url"
            inputmode="url"
            aria-invalid={formError ? "true" : undefined}
            aria-describedby={formError ? "url-error" : undefined}
            required
          />
          <div class="absolute right-0 top-0 bottom-0 flex items-center">
            <button
              type="button"
              on:click={pasteFromClipboard}
              class="h-full flex items-center gap-1.5 px-3.5 rounded-r-xl border-l text-[11px] font-medium transition-all duration-200 {justPasted
                ? 'border-white/5 bg-emerald-500/10 text-emerald-400'
                : 'border-white/5 bg-transparent hover:bg-white/5 text-gray-500 hover:text-gray-200'}"
              aria-label="Paste from clipboard"
              title="Paste from clipboard"
            >
              {#if justPasted}
                <svg
                  class="w-3.5 h-3.5"
                  fill="none"
                  stroke="currentColor"
                  stroke-width="2.5"
                  viewBox="0 0 24 24"
                  aria-hidden="true"
                >
                  <path stroke-linecap="round" stroke-linejoin="round" d="M5 13l4 4L19 7" />
                </svg>
                Pasted
              {:else}
                <svg
                  class="w-3.5 h-3.5"
                  fill="none"
                  stroke="currentColor"
                  stroke-width="2"
                  viewBox="0 0 24 24"
                  aria-hidden="true"
                >
                  <path
                    stroke-linecap="round"
                    stroke-linejoin="round"
                    d="M9 5H7a2 2 0 00-2 2v12a2 2 0 002 2h10a2 2 0 002-2V7a2 2 0 00-2-2h-2M9 5a2 2 0 002 2h2a2 2 0 002-2M9 5a2 2 0 012-2h2a2 2 0 012 2"
                  />
                </svg>
                Paste
              {/if}
            </button>
          </div>
        </div>

        <button
          type="submit"
          class="w-full sm:w-auto inline-flex items-center justify-center gap-2 px-6 py-3.5 rounded-xl bg-gradient-to-r from-blue-600 to-indigo-600 hover:from-blue-500 hover:to-indigo-500 text-white text-sm font-semibold shadow-lg shadow-blue-900/30 transition-all duration-200 disabled:opacity-50 focus:outline-none focus:ring-2 focus:ring-offset-0 focus:ring-blue-500 active:scale-95"
          disabled={loading}
          aria-busy={loading}
          aria-label={loading ? "Scanning URL, please wait" : "Scan Now"}
        >
          {#if loading}
            <svg
              class="w-4 h-4 animate-spin"
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
            Scanning..
          {:else}
            <svg
              class="w-4 h-4"
              fill="none"
              stroke="currentColor"
              stroke-width="2"
              viewBox="0 0 24 24"
              aria-hidden="true"
            >
              <path
                stroke-linecap="round"
                stroke-linejoin="round"
                d="M21 21l-4.35-4.35M17 11A6 6 0 1 1 5 11a6 6 0 0 1 12 0z"
              />
            </svg>
            Scan Now
          {/if}
        </button>
      </div>

      {#if formError}
        <p id="url-error" class="text-red-400 text-xs mt-2 text-left" role="alert">{formError}</p>
      {/if}
    </form>

    {#if isLanding}
      <div class="mt-5 flex flex-wrap justify-center items-center gap-2">
        <span class="text-[11px] text-gray-500 mr-0.5">Try:</span>
        {#each [{ label: "google.com", url: "google.com", dot: "bg-emerald-500", hint: "Safe" }, { label: "аррӏе.com", url: "аррӏе.com", dot: "bg-red-500", hint: "Risky" }, { label: "wikipedia.org", url: "wikipedia.org", dot: "bg-emerald-500", hint: "Safe" }, { label: "pаypal.com", url: "pаypal.com", dot: "bg-red-500", hint: "Risky" }] as example}
          <button
            type="button"
            on:click={() => {
              input = example.url;
              runAnalyze(example.url);
            }}
            title="Expected: {example.hint}"
            class="inline-flex items-center gap-1.5 px-3 py-1.5 rounded-full bg-gray-900 border border-gray-800 hover:border-gray-600 text-gray-400 hover:text-gray-200 text-xs transition-all"
          >
            <span class={`w-1.5 h-1.5 rounded-full flex-shrink-0 ${example.dot}`}></span>
            {example.label}
          </button>
        {/each}
      </div>

      <div class="mt-5 flex flex-wrap justify-center gap-x-5 gap-y-2">
        {#each ["Free & open source", "No signup required", "Explains every verdict", "Live page preview"] as pill}
          <span class="flex items-center gap-1.5 text-[11px] text-gray-500">
            <svg
              class="w-3 h-3 text-emerald-500/70 flex-shrink-0"
              fill="none"
              stroke="currentColor"
              stroke-width="2.5"
              viewBox="0 0 24 24"
            >
              <path stroke-linecap="round" stroke-linejoin="round" d="M5 13l4 4L19 7" />
            </svg>
            {pill}
          </span>
        {/each}
      </div>
    {/if}

    <div class={`w-full ${isLanding ? "mt-12" : "mt-8"}`} aria-live="polite">
      {#if loading || scanDone}
        <ScanProgress {loading} done={scanDone} />
      {:else}
        <ResultSection
          data={scanResult}
          {loading}
          {error}
          {screenshotUrl}
          {screenshotLoading}
          {screenshotFailed}
        />
      {/if}
    </div>

    {#if isLanding}
      <div class="mt-8 w-full">
        <Shoutouts />
      </div>
    {/if}
  </div>
</section>
