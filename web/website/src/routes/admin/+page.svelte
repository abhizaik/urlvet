<script lang="ts">
  import { onMount } from "svelte";
  import { env } from "$env/dynamic/public";

  const BASE = env.PUBLIC_BASE_URL || "http://localhost:8080/api/v1";

  interface CacheEntry {
    key: string;
    prefix: string;
    ttl_seconds: number;
    value: unknown;
  }

  let entries: CacheEntry[] = [];
  let loading = false;
  let error: string | null = null;
  let flushConfirm = false;
  let deletingKey: string | null = null;
  let expandedKeys = new Set<string>();
  let filterPrefix = "";
  let filterText = "";

  $: prefixes = [...new Set(entries.map((e) => e.prefix))].sort();

  $: filtered = entries.filter((e) => {
    const matchPrefix = !filterPrefix || e.prefix === filterPrefix;
    const matchText =
      !filterText ||
      e.key.toLowerCase().includes(filterText.toLowerCase());
    return matchPrefix && matchText;
  });

  async function load() {
    loading = true;
    error = null;
    try {
      const res = await fetch(`${BASE}/admin/cache`);
      if (!res.ok) throw new Error(`HTTP ${res.status}`);
      const data = await res.json();
      entries = data.keys ?? [];
    } catch (e) {
      error = e instanceof Error ? e.message : "Unknown error";
    } finally {
      loading = false;
    }
  }

  async function deleteKey(key: string) {
    deletingKey = key;
    try {
      const res = await fetch(
        `${BASE}/admin/cache/${encodeURIComponent(key)}`,
        { method: "DELETE" }
      );
      if (!res.ok) throw new Error(`HTTP ${res.status}`);
      entries = entries.filter((e) => e.key !== key);
      expandedKeys.delete(key);
      expandedKeys = expandedKeys;
    } catch (e) {
      error = e instanceof Error ? e.message : "Delete failed";
    } finally {
      deletingKey = null;
    }
  }

  async function flushAll() {
    flushConfirm = false;
    loading = true;
    try {
      const res = await fetch(`${BASE}/admin/cache`, { method: "DELETE" });
      if (!res.ok) throw new Error(`HTTP ${res.status}`);
      entries = [];
      expandedKeys = new Set();
    } catch (e) {
      error = e instanceof Error ? e.message : "Flush failed";
    } finally {
      loading = false;
    }
  }

  function toggleExpand(key: string) {
    if (expandedKeys.has(key)) {
      expandedKeys.delete(key);
    } else {
      expandedKeys.add(key);
    }
    expandedKeys = expandedKeys;
  }

  function formatTTL(sec: number): string {
    if (sec === -1) return "No expiry";
    if (sec === -2) return "Gone";
    if (sec < 60) return `${sec}s`;
    if (sec < 3600) return `${Math.floor(sec / 60)}m ${sec % 60}s`;
    return `${Math.floor(sec / 3600)}h ${Math.floor((sec % 3600) / 60)}m`;
  }

  function ttlColor(sec: number): string {
    if (sec < 0) return "text-gray-500";
    if (sec < 60) return "text-red-400";
    if (sec < 300) return "text-yellow-400";
    return "text-green-400";
  }

  function formatJSON(val: unknown): string {
    try {
      return JSON.stringify(val, null, 2);
    } catch {
      return String(val);
    }
  }

  onMount(load);
</script>

<svelte:head>
  <title>Cache Admin — SafeSurf</title>
  <meta name="robots" content="noindex, nofollow, noarchive" />
</svelte:head>

<div class="min-h-screen bg-gray-950 text-gray-200 p-6">
  <div class="max-w-6xl mx-auto space-y-6">

    <!-- Header -->
    <div class="flex flex-col sm:flex-row sm:items-center justify-between gap-4">
      <div>
        <h1 class="text-2xl font-bold text-white">Cache Admin</h1>
        <p class="text-sm text-gray-400 mt-0.5">
          {entries.length} key{entries.length !== 1 ? "s" : ""} in store
          {#if filtered.length !== entries.length}
            <span class="text-gray-500">· {filtered.length} shown</span>
          {/if}
        </p>
      </div>

      <div class="flex items-center gap-3">
        <button
          on:click={load}
          disabled={loading}
          class="flex items-center gap-2 px-4 py-2 rounded-lg bg-gray-800 hover:bg-gray-700 border border-gray-700 text-sm font-medium transition-colors disabled:opacity-50"
        >
          <svg class="w-4 h-4 {loading ? 'animate-spin' : ''}" fill="none" viewBox="0 0 24 24" stroke="currentColor">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
              d="M4 4v5h.582m15.356 2A8.001 8.001 0 004.582 9m0 0H9m11 11v-5h-.581m0 0a8.003 8.003 0 01-15.357-2m15.357 2H15" />
          </svg>
          Refresh
        </button>

        {#if !flushConfirm}
          <button
            on:click={() => (flushConfirm = true)}
            disabled={loading || entries.length === 0}
            class="flex items-center gap-2 px-4 py-2 rounded-lg bg-red-900/40 hover:bg-red-900/70 border border-red-700 text-red-300 text-sm font-medium transition-colors disabled:opacity-40"
          >
            <svg class="w-4 h-4" fill="none" viewBox="0 0 24 24" stroke="currentColor">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16" />
            </svg>
            Flush All
          </button>
        {:else}
          <div class="flex items-center gap-2">
            <span class="text-xs text-red-400">Delete all {entries.length} keys?</span>
            <button
              on:click={flushAll}
              class="px-3 py-1.5 rounded-lg bg-red-700 hover:bg-red-600 text-white text-xs font-semibold transition-colors"
            >Confirm</button>
            <button
              on:click={() => (flushConfirm = false)}
              class="px-3 py-1.5 rounded-lg bg-gray-700 hover:bg-gray-600 text-gray-200 text-xs font-semibold transition-colors"
            >Cancel</button>
          </div>
        {/if}
      </div>
    </div>

    <!-- Error banner -->
    {#if error}
      <div class="flex items-center gap-3 bg-red-900/30 border border-red-700 rounded-lg px-4 py-3 text-sm text-red-300">
        <svg class="w-4 h-4 flex-shrink-0" fill="currentColor" viewBox="0 0 20 20">
          <path fill-rule="evenodd" d="M18 10a8 8 0 11-16 0 8 8 0 0116 0zm-7 4a1 1 0 11-2 0 1 1 0 012 0zm-1-9a1 1 0 00-1 1v4a1 1 0 102 0V6a1 1 0 00-1-1z" clip-rule="evenodd" />
        </svg>
        {error}
        <button on:click={() => (error = null)} class="ml-auto text-red-400 hover:text-red-200">✕</button>
      </div>
    {/if}

    <!-- Filters -->
    {#if entries.length > 0}
      <div class="flex flex-col sm:flex-row gap-3">
        <input
          bind:value={filterText}
          placeholder="Search keys…"
          class="flex-1 bg-gray-900 border border-gray-700 rounded-lg px-3 py-2 text-sm text-gray-200 placeholder-gray-500 focus:outline-none focus:border-gray-500"
        />
        <select
          bind:value={filterPrefix}
          class="bg-gray-900 border border-gray-700 rounded-lg px-3 py-2 text-sm text-gray-200 focus:outline-none focus:border-gray-500"
        >
          <option value="">All prefixes</option>
          {#each prefixes as p}
            <option value={p}>{p}</option>
          {/each}
        </select>
      </div>
    {/if}

    <!-- Empty state -->
    {#if !loading && entries.length === 0 && !error}
      <div class="text-center py-20 text-gray-500">
        <svg class="w-12 h-12 mx-auto mb-4 opacity-30" fill="none" viewBox="0 0 24 24" stroke="currentColor">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.5"
            d="M5 8h14M5 8a2 2 0 110-4h14a2 2 0 110 4M5 8v10a2 2 0 002 2h10a2 2 0 002-2V8m-9 4h4" />
        </svg>
        <p class="text-sm">Cache is empty</p>
      </div>
    {/if}

    <!-- Loading skeleton -->
    {#if loading && entries.length === 0}
      <div class="space-y-2">
        {#each Array(5) as _}
          <div class="h-14 bg-gray-900 rounded-lg animate-pulse border border-gray-800"></div>
        {/each}
      </div>
    {/if}

    <!-- Key list -->
    {#if filtered.length > 0}
      <div class="space-y-2">
        {#each filtered as entry (entry.key)}
          {@const expanded = expandedKeys.has(entry.key)}
          <div class="bg-gray-900 border border-gray-800 rounded-lg overflow-hidden hover:border-gray-700 transition-colors">

            <!-- Row header -->
            <div class="flex items-center gap-3 px-4 py-3">

              <!-- Expand toggle -->
              <button
                on:click={() => toggleExpand(entry.key)}
                class="flex-shrink-0 text-gray-500 hover:text-gray-300 transition-colors"
                aria-label={expanded ? "Collapse" : "Expand"}
              >
                <svg class="w-4 h-4 transition-transform duration-150 {expanded ? 'rotate-90' : ''}"
                  viewBox="0 0 20 20" fill="currentColor">
                  <path fill-rule="evenodd"
                    d="M7.293 4.293a1 1 0 011.414 0l5 5a1 1 0 010 1.414l-5 5a1 1 0 01-1.414-1.414L11.586 10 7.293 5.707a1 1 0 010-1.414z"
                    clip-rule="evenodd" />
                </svg>
              </button>

              <!-- Prefix badge -->
              <span class="flex-shrink-0 text-[10px] font-semibold uppercase tracking-wide px-2 py-0.5 rounded bg-gray-800 text-gray-400 border border-gray-700">
                {entry.prefix}
              </span>

              <!-- Key name -->
              <button
                on:click={() => toggleExpand(entry.key)}
                class="flex-1 text-left font-mono text-sm text-gray-200 truncate hover:text-white transition-colors"
                title={entry.key}
              >
                {entry.key}
              </button>

              <!-- TTL -->
              <span class="flex-shrink-0 text-xs font-medium {ttlColor(entry.ttl_seconds)} tabular-nums w-24 text-right">
                {formatTTL(entry.ttl_seconds)}
              </span>

              <!-- Delete -->
              <button
                on:click={() => deleteKey(entry.key)}
                disabled={deletingKey === entry.key}
                class="flex-shrink-0 ml-2 p-1.5 rounded-md text-gray-600 hover:text-red-400 hover:bg-red-900/20 transition-colors disabled:opacity-40"
                title="Delete key"
              >
                {#if deletingKey === entry.key}
                  <svg class="w-4 h-4 animate-spin" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                      d="M4 4v5h.582m15.356 2A8.001 8.001 0 004.582 9m0 0H9m11 11v-5h-.581m0 0a8.003 8.003 0 01-15.357-2m15.357 2H15" />
                  </svg>
                {:else}
                  <svg class="w-4 h-4" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                      d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16" />
                  </svg>
                {/if}
              </button>
            </div>

            <!-- Expanded value -->
            {#if expanded}
              <div class="border-t border-gray-800 px-4 py-3">
                <pre class="text-xs text-green-300 font-mono bg-gray-950 rounded-lg p-4 overflow-x-auto whitespace-pre-wrap break-all leading-relaxed max-h-96 overflow-y-auto">{formatJSON(entry.value)}</pre>
              </div>
            {/if}

          </div>
        {/each}
      </div>
    {:else if !loading && entries.length > 0}
      <p class="text-center text-sm text-gray-500 py-10">No keys match your filter.</p>
    {/if}

  </div>
</div>
