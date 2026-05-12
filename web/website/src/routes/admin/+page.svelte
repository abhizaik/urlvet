<script lang="ts">
  import { env } from "$env/dynamic/public";
  import { onMount } from "svelte";

  const BASE = env.PUBLIC_BASE_URL || "http://localhost:8080/api/v1";

  // ── Types ──────────────────────────────────────────────────────────────────
  interface CacheEntry {
    key: string;
    prefix: string;
    ttl_seconds: number;
    value: unknown;
  }

  interface ScanRecord {
    url: string;
    domain: string;
    verdict: string;
    score: number;
    duration: string;
    time: string;
    cached: boolean;
  }

  interface ErrorRecord {
    task: string;
    error: string;
    url: string;
    time: string;
  }

  interface DomainCount {
    domain: string;
    count: number;
  }

  interface Stats {
    total_scans_today: number;
    total_scans_all: number;
    cache_hits: number;
    cache_misses: number;
    cache_hit_rate: number;
    avg_duration_ms: number;
    top_domains: DomainCount[];
    verdict_counts: Record<string, number>;
  }

  // ── State ──────────────────────────────────────────────────────────────────
  type Tab = "overview" | "recent" | "errors" | "cache";
  let activeTab: Tab = "overview";

  let stats: Stats | null = null;
  let recentScans: ScanRecord[] = [];
  let recentErrors: ErrorRecord[] = [];
  let cacheEntries: CacheEntry[] = [];

  let loadingStats = false;
  let loadingRecent = false;
  let loadingErrors = false;
  let loadingCache = false;

  let globalError: string | null = null;

  // ── Auth ───────────────────────────────────────────────────────────────────
  let token: string | null = null;
  let passwordInput = "";
  let loginError: string | null = null;
  let loggingIn = false;

  async function authFetch(url: string, options: RequestInit = {}): Promise<Response> {
    return fetch(url, {
      ...options,
      headers: { ...options.headers, Authorization: `Bearer ${token}` },
    });
  }

  function clearAuth(message?: string) {
    token = null;
    sessionStorage.removeItem("admin_token");
    loginError = message ?? null;
  }

  function logout() {
    clearAuth();
    stats = null;
    recentScans = [];
    recentErrors = [];
    cacheEntries = [];
  }

  async function login() {
    const password = passwordInput.trim();
    if (!password || loggingIn) return;
    loggingIn = true;
    loginError = null;
    try {
      const res = await fetch(`${BASE}/admin/login`, {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({ password }),
      });
      if (res.status === 401) {
        loginError = "Incorrect password.";
        return;
      }
      if (res.status === 503) {
        loginError = "Admin access is disabled on this server.";
        return;
      }
      if (!res.ok) {
        loginError = "Login failed. Please try again.";
        return;
      }
      const data = await res.json();
      token = data.token;
      sessionStorage.setItem("admin_token", token!);
      passwordInput = "";
      fetchStats();
    } catch {
      loginError = "Unable to reach the server.";
    } finally {
      loggingIn = false;
    }
  }

  // cache UI state
  let flushConfirm = false;
  let deletingKey: string | null = null;
  let expandedKeys = new Set<string>();
  let filterText = "";
  let filterPrefix = "";

  $: prefixes = [...new Set(cacheEntries.map((e) => e.prefix))].sort();
  $: filteredCache = cacheEntries.filter((e) => {
    const matchPrefix = !filterPrefix || e.prefix === filterPrefix;
    const matchText = !filterText || e.key.toLowerCase().includes(filterText.toLowerCase());
    return matchPrefix && matchText;
  });

  // ── Fetch helpers ──────────────────────────────────────────────────────────
  async function fetchStats() {
    loadingStats = true;
    try {
      const res = await authFetch(`${BASE}/admin/stats`);
      if (res.status === 401) {
        clearAuth("Invalid token. Please sign in again.");
        return;
      }
      if (!res.ok) throw new Error(`HTTP ${res.status}`);
      stats = await res.json();
    } catch (e) {
      globalError = e instanceof Error ? e.message : "Failed to load stats";
    } finally {
      loadingStats = false;
    }
  }

  async function fetchRecent() {
    loadingRecent = true;
    try {
      const res = await authFetch(`${BASE}/admin/recent`);
      if (res.status === 401) {
        clearAuth("Invalid token. Please sign in again.");
        return;
      }
      if (!res.ok) throw new Error(`HTTP ${res.status}`);
      const data = await res.json();
      recentScans = data.scans ?? [];
    } catch (e) {
      globalError = e instanceof Error ? e.message : "Failed to load recent scans";
    } finally {
      loadingRecent = false;
    }
  }

  async function fetchErrors() {
    loadingErrors = true;
    try {
      const res = await authFetch(`${BASE}/admin/errors`);
      if (res.status === 401) {
        clearAuth("Invalid token. Please sign in again.");
        return;
      }
      if (!res.ok) throw new Error(`HTTP ${res.status}`);
      const data = await res.json();
      recentErrors = data.errors ?? [];
    } catch (e) {
      globalError = e instanceof Error ? e.message : "Failed to load errors";
    } finally {
      loadingErrors = false;
    }
  }

  async function fetchCache() {
    loadingCache = true;
    try {
      const res = await authFetch(`${BASE}/admin/cache`);
      if (res.status === 401) {
        clearAuth("Invalid token. Please sign in again.");
        return;
      }
      if (!res.ok) throw new Error(`HTTP ${res.status}`);
      const data = await res.json();
      cacheEntries = data.keys ?? [];
    } catch (e) {
      globalError = e instanceof Error ? e.message : "Failed to load cache";
    } finally {
      loadingCache = false;
    }
  }

  async function deleteKey(key: string) {
    deletingKey = key;
    try {
      const res = await authFetch(`${BASE}/admin/cache/${encodeURIComponent(key)}`, {
        method: "DELETE",
      });
      if (res.status === 401) {
        clearAuth("Invalid token. Please sign in again.");
        return;
      }
      if (!res.ok) throw new Error(`HTTP ${res.status}`);
      cacheEntries = cacheEntries.filter((e) => e.key !== key);
      expandedKeys.delete(key);
      expandedKeys = expandedKeys;
    } catch (e) {
      globalError = e instanceof Error ? e.message : "Delete failed";
    } finally {
      deletingKey = null;
    }
  }

  async function flushAll() {
    flushConfirm = false;
    loadingCache = true;
    try {
      const res = await authFetch(`${BASE}/admin/cache`, { method: "DELETE" });
      if (res.status === 401) {
        clearAuth("Invalid token. Please sign in again.");
        return;
      }
      if (!res.ok) throw new Error(`HTTP ${res.status}`);
      cacheEntries = [];
      expandedKeys = new Set();
    } catch (e) {
      globalError = e instanceof Error ? e.message : "Flush failed";
    } finally {
      loadingCache = false;
    }
  }

  function switchTab(tab: Tab) {
    activeTab = tab;
    globalError = null;
    if (tab === "overview") fetchStats();
    else if (tab === "recent") fetchRecent();
    else if (tab === "errors") fetchErrors();
    else if (tab === "cache") fetchCache();
  }

  function refresh() {
    globalError = null;
    switchTab(activeTab);
  }

  function toggleExpand(key: string) {
    if (expandedKeys.has(key)) expandedKeys.delete(key);
    else expandedKeys.add(key);
    expandedKeys = expandedKeys;
  }

  // ── Formatters ─────────────────────────────────────────────────────────────
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
    return "text-emerald-400";
  }

  function verdictColor(v: string): string {
    if (v === "Safe") return "text-emerald-400 bg-emerald-500/10 border-emerald-500/20";
    if (v === "Risky") return "text-red-400 bg-red-500/10 border-red-500/20";
    return "text-yellow-400 bg-yellow-500/10 border-yellow-500/20";
  }

  function verdictDot(v: string): string {
    if (v === "Safe") return "bg-emerald-400";
    if (v === "Risky") return "bg-red-400";
    return "bg-yellow-400";
  }

  function relativeTime(iso: string): string {
    const diff = Date.now() - new Date(iso).getTime();
    const s = Math.floor(diff / 1000);
    if (s < 60) return `${s}s ago`;
    const m = Math.floor(s / 60);
    if (m < 60) return `${m}m ago`;
    const h = Math.floor(m / 60);
    if (h < 24) return `${h}h ago`;
    return `${Math.floor(h / 24)}d ago`;
  }

  function formatJSON(val: unknown): string {
    try {
      return JSON.stringify(val, null, 2);
    } catch {
      return String(val);
    }
  }

  onMount(() => {
    token = sessionStorage.getItem("admin_token");
    if (token) fetchStats();
  });
</script>

<svelte:head>
  <title>Admin — SafeSurf</title>
  <meta name="robots" content="noindex, nofollow, noarchive" />
</svelte:head>

{#if !token}
  <div class="min-h-screen bg-gray-950 flex items-center justify-center px-4">
    <div class="w-full max-w-sm">
      <div class="text-center mb-8">
        <h1 class="text-xl font-semibold text-white">SafeSurf Admin</h1>
        <p class="text-sm text-gray-500 mt-1">Enter your admin password to continue</p>
      </div>
      <form
        on:submit|preventDefault={login}
        class="bg-gray-900 border border-gray-800 rounded-2xl p-6 space-y-4"
      >
        {#if loginError}
          <p class="text-sm text-red-400 bg-red-900/20 border border-red-800 rounded-lg px-3 py-2">
            {loginError}
          </p>
        {/if}
        <div>
          <label for="admin-password" class="block text-xs font-medium text-gray-400 mb-1.5"
            >Password</label
          >
          <input
            id="admin-password"
            type="password"
            bind:value={passwordInput}
            placeholder="Admin password"
            autocomplete="current-password"
            disabled={loggingIn}
            class="w-full bg-gray-950 border border-gray-800 rounded-lg px-3 py-2.5 text-sm text-gray-200 placeholder-gray-600 focus:outline-none focus:border-gray-600 disabled:opacity-50"
          />
        </div>
        <button
          type="submit"
          disabled={!passwordInput.trim() || loggingIn}
          class="w-full py-2.5 rounded-lg bg-white text-gray-900 text-sm font-semibold hover:bg-gray-100 transition-colors disabled:opacity-40"
        >
          {loggingIn ? "Signing in…" : "Sign In"}
        </button>
      </form>
    </div>
  </div>
{:else}
  <div class="min-h-screen bg-gray-950 text-gray-200">
    <!-- Top bar -->
    <div class="border-b border-gray-800 bg-gray-900/50 backdrop-blur-sm sticky top-0 z-10">
      <div class="max-w-7xl mx-auto px-6 flex items-center justify-between h-14">
        <div class="flex items-center gap-3">
          <a
            href="/"
            aria-label="Back to SafeSurf"
            class="text-gray-500 hover:text-gray-300 transition-colors"
          >
            <svg
              class="w-4 h-4"
              fill="none"
              viewBox="0 0 24 24"
              stroke="currentColor"
              stroke-width="2"
            >
              <path
                stroke-linecap="round"
                stroke-linejoin="round"
                d="M10 19l-7-7m0 0l7-7m-7 7h18"
              />
            </svg>
          </a>
          <span class="text-gray-700">|</span>
          <span class="text-sm font-semibold text-white">SafeSurf Admin</span>
        </div>
        <div class="flex items-center gap-2">
          <button
            on:click={refresh}
            class="flex items-center gap-1.5 px-3 py-1.5 rounded-lg bg-gray-800 hover:bg-gray-700 border border-gray-700 text-xs font-medium text-gray-300 transition-colors"
          >
            <svg
              class="w-3.5 h-3.5"
              fill="none"
              viewBox="0 0 24 24"
              stroke="currentColor"
              stroke-width="2"
            >
              <path
                stroke-linecap="round"
                stroke-linejoin="round"
                d="M4 4v5h.582m15.356 2A8.001 8.001 0 004.582 9m0 0H9m11 11v-5h-.581m0 0a8.003 8.003 0 01-15.357-2m15.357 2H15"
              />
            </svg>
            Refresh
          </button>
          <button
            on:click={logout}
            class="flex items-center gap-1.5 px-3 py-1.5 rounded-lg bg-gray-800 hover:bg-gray-700 border border-gray-700 text-xs font-medium text-gray-500 hover:text-gray-300 transition-colors"
          >
            Sign Out
          </button>
        </div>
      </div>
    </div>

    <div class="max-w-7xl mx-auto px-6 py-8">
      <!-- Tabs -->
      <nav class="flex gap-1 mb-8 bg-gray-900 rounded-xl p-1 w-fit border border-gray-800">
        {#each [{ id: "overview", label: "Overview", icon: "M9 19v-6a2 2 0 00-2-2H5a2 2 0 00-2 2v6a2 2 0 002 2h2a2 2 0 002-2zm0 0V9a2 2 0 012-2h2a2 2 0 012 2v10m-6 0a2 2 0 002 2h2a2 2 0 002-2m0 0V5a2 2 0 012-2h2a2 2 0 012 2v14a2 2 0 01-2 2h-2a2 2 0 01-2-2z" }, { id: "recent", label: "Recent Scans", icon: "M12 8v4l3 3m6-3a9 9 0 11-18 0 9 9 0 0118 0z" }, { id: "errors", label: "Errors", icon: "M12 9v2m0 4h.01m-6.938 4h13.856c1.54 0 2.502-1.667 1.732-3L13.732 4c-.77-1.333-2.694-1.333-3.464 0L3.34 16c-.77 1.333.192 3 1.732 3z" }, { id: "cache", label: "Cache", icon: "M5 8h14M5 8a2 2 0 110-4h14a2 2 0 110 4M5 8v10a2 2 0 002 2h10a2 2 0 002-2V8m-9 4h4" }] as tab}
          <button
            on:click={() => switchTab(tab.id as Tab)}
            class="flex items-center gap-2 px-4 py-2 rounded-lg text-sm font-medium transition-all duration-150 {activeTab ===
            tab.id
              ? 'bg-gray-800 text-white shadow-sm'
              : 'text-gray-500 hover:text-gray-300'}"
          >
            <svg
              class="w-4 h-4 flex-shrink-0"
              fill="none"
              viewBox="0 0 24 24"
              stroke="currentColor"
              stroke-width="2"
            >
              <path stroke-linecap="round" stroke-linejoin="round" d={tab.icon} />
            </svg>
            {tab.label}
            {#if tab.id === "errors" && recentErrors.length > 0}
              <span
                class="ml-0.5 px-1.5 py-0.5 rounded-full bg-red-500/20 text-red-400 text-[10px] font-semibold"
                >{recentErrors.length}</span
              >
            {/if}
          </button>
        {/each}
      </nav>

      <!-- Global error -->
      {#if globalError}
        <div
          class="flex items-center gap-3 bg-red-900/20 border border-red-800 rounded-xl px-4 py-3 text-sm text-red-300 mb-6"
        >
          <svg class="w-4 h-4 flex-shrink-0" fill="currentColor" viewBox="0 0 20 20">
            <path
              fill-rule="evenodd"
              d="M18 10a8 8 0 11-16 0 8 8 0 0116 0zm-7 4a1 1 0 11-2 0 1 1 0 012 0zm-1-9a1 1 0 00-1 1v4a1 1 0 102 0V6a1 1 0 00-1-1z"
              clip-rule="evenodd"
            />
          </svg>
          {globalError}
          <button
            on:click={() => (globalError = null)}
            class="ml-auto text-red-500 hover:text-red-300">✕</button
          >
        </div>
      {/if}

      <!-- ── OVERVIEW TAB ──────────────────────────────────────────────────── -->
      {#if activeTab === "overview"}
        {#if loadingStats}
          <div class="grid grid-cols-2 md:grid-cols-4 gap-4 mb-8">
            {#each Array(4) as _}
              <div class="h-28 bg-gray-900 rounded-xl border border-gray-800 animate-pulse"></div>
            {/each}
          </div>
        {:else if stats}
          <!-- Stat cards -->
          <div class="grid grid-cols-2 md:grid-cols-4 gap-4 mb-8">
            <div class="bg-gray-900 border border-gray-800 rounded-xl p-5">
              <p class="text-xs text-gray-500 font-medium uppercase tracking-wide mb-1">
                Scans Today
              </p>
              <p class="text-3xl font-bold text-white">{stats.total_scans_today}</p>
            </div>
            <div class="bg-gray-900 border border-gray-800 rounded-xl p-5">
              <p class="text-xs text-gray-500 font-medium uppercase tracking-wide mb-1">
                Total Scans
              </p>
              <p class="text-3xl font-bold text-white">{stats.total_scans_all}</p>
              <p class="text-xs text-gray-600 mt-1">in memory</p>
            </div>
            <div class="bg-gray-900 border border-gray-800 rounded-xl p-5">
              <p class="text-xs text-gray-500 font-medium uppercase tracking-wide mb-1">
                Cache Hit Rate
              </p>
              <p class="text-3xl font-bold text-white">
                {stats.cache_hit_rate.toFixed(1)}<span class="text-lg text-gray-500">%</span>
              </p>
              <p class="text-xs text-gray-600 mt-1">
                {stats.cache_hits} hits · {stats.cache_misses} misses
              </p>
            </div>
            <div class="bg-gray-900 border border-gray-800 rounded-xl p-5">
              <p class="text-xs text-gray-500 font-medium uppercase tracking-wide mb-1">
                Avg Duration
              </p>
              <p class="text-3xl font-bold text-white">
                {stats.avg_duration_ms.toFixed(0)}<span class="text-lg text-gray-500">ms</span>
              </p>
            </div>
          </div>

          <div class="grid grid-cols-1 md:grid-cols-2 gap-6">
            <!-- Verdict breakdown -->
            <div class="bg-gray-900 border border-gray-800 rounded-xl p-5">
              <h2 class="text-sm font-semibold text-gray-300 mb-4">Verdict Breakdown</h2>
              {#if Object.keys(stats.verdict_counts).length === 0}
                <p class="text-sm text-gray-600">No data yet</p>
              {:else}
                {@const total = Object.values(stats.verdict_counts).reduce((a, b) => a + b, 0)}
                <div class="space-y-3">
                  {#each Object.entries(stats.verdict_counts).sort((a, b) => b[1] - a[1]) as [verdict, count]}
                    {@const pct = total > 0 ? (count / total) * 100 : 0}
                    <div>
                      <div class="flex justify-between items-center mb-1">
                        <span
                          class="text-sm font-medium {verdict === 'Safe'
                            ? 'text-emerald-400'
                            : verdict === 'Risky'
                              ? 'text-red-400'
                              : 'text-yellow-400'}">{verdict}</span
                        >
                        <span class="text-xs text-gray-500"
                          >{count} <span class="text-gray-600">({pct.toFixed(0)}%)</span></span
                        >
                      </div>
                      <div class="h-1.5 rounded-full bg-gray-800 overflow-hidden">
                        <div
                          class="h-full rounded-full transition-all duration-500 {verdict === 'Safe'
                            ? 'bg-emerald-500'
                            : verdict === 'Risky'
                              ? 'bg-red-500'
                              : 'bg-yellow-500'}"
                          style="width: {pct}%"
                        ></div>
                      </div>
                    </div>
                  {/each}
                </div>
              {/if}
            </div>

            <!-- Top domains -->
            <div class="bg-gray-900 border border-gray-800 rounded-xl p-5">
              <h2 class="text-sm font-semibold text-gray-300 mb-4">Top Scanned Domains</h2>
              {#if !stats.top_domains || stats.top_domains.length === 0}
                <p class="text-sm text-gray-600">No data yet</p>
              {:else}
                {@const maxCount = stats.top_domains[0]?.count ?? 1}
                <div class="space-y-2.5">
                  {#each stats.top_domains as item, i}
                    <div class="flex items-center gap-3">
                      <span class="text-xs text-gray-600 w-4 text-right flex-shrink-0">{i + 1}</span
                      >
                      <div class="flex-1 min-w-0">
                        <div class="flex items-center justify-between mb-0.5">
                          <span class="text-sm text-gray-300 font-mono truncate">{item.domain}</span
                          >
                          <span class="text-xs text-gray-500 ml-2 flex-shrink-0">{item.count}</span>
                        </div>
                        <div class="h-1 rounded-full bg-gray-800 overflow-hidden">
                          <div
                            class="h-full rounded-full bg-blue-500/60"
                            style="width: {(item.count / maxCount) * 100}%"
                          ></div>
                        </div>
                      </div>
                    </div>
                  {/each}
                </div>
              {/if}
            </div>
          </div>
        {:else}
          <div class="text-center py-20 text-gray-600">
            <p class="text-sm">No stats available — start scanning URLs.</p>
          </div>
        {/if}

        <!-- ── RECENT SCANS TAB ──────────────────────────────────────────────── -->
      {:else if activeTab === "recent"}
        {#if loadingRecent}
          <div class="space-y-2">
            {#each Array(6) as _}
              <div class="h-16 bg-gray-900 rounded-xl border border-gray-800 animate-pulse"></div>
            {/each}
          </div>
        {:else if recentScans.length === 0}
          <div class="text-center py-24 text-gray-600">
            <svg
              class="w-10 h-10 mx-auto mb-3 opacity-30"
              fill="none"
              viewBox="0 0 24 24"
              stroke="currentColor"
            >
              <path
                stroke-linecap="round"
                stroke-linejoin="round"
                stroke-width="1.5"
                d="M12 8v4l3 3m6-3a9 9 0 11-18 0 9 9 0 0118 0z"
              />
            </svg>
            <p class="text-sm">No scans recorded yet.</p>
          </div>
        {:else}
          <div class="space-y-2">
            {#each recentScans as scan}
              <div
                class="bg-gray-900 border border-gray-800 rounded-xl px-4 py-3 flex items-center gap-4 hover:border-gray-700 transition-colors"
              >
                <!-- Verdict dot -->
                <span class="w-2 h-2 rounded-full flex-shrink-0 {verdictDot(scan.verdict)}"></span>

                <!-- URL / Domain -->
                <div class="flex-1 min-w-0">
                  <p class="text-sm text-gray-200 font-mono truncate">{scan.domain}</p>
                  <p class="text-xs text-gray-600 truncate">{scan.url}</p>
                </div>

                <!-- Verdict badge -->
                <span
                  class="flex-shrink-0 text-[11px] font-semibold px-2 py-0.5 rounded-full border {verdictColor(
                    scan.verdict
                  )}"
                >
                  {scan.verdict}
                </span>

                <!-- Score -->
                <span
                  class="flex-shrink-0 text-sm font-bold tabular-nums {scan.verdict === 'Safe'
                    ? 'text-emerald-400'
                    : scan.verdict === 'Risky'
                      ? 'text-red-400'
                      : 'text-yellow-400'}">{scan.score}</span
                >

                <!-- Duration -->
                <span class="flex-shrink-0 text-xs text-gray-500 tabular-nums w-20 text-right"
                  >{scan.duration}</span
                >

                <!-- Cached badge -->
                {#if scan.cached}
                  <span
                    class="flex-shrink-0 text-[10px] font-semibold px-2 py-0.5 rounded-full bg-blue-500/10 text-blue-400 border border-blue-500/20"
                    >cached</span
                  >
                {/if}

                <!-- Time -->
                <span class="flex-shrink-0 text-xs text-gray-600 w-16 text-right"
                  >{relativeTime(scan.time)}</span
                >
              </div>
            {/each}
          </div>
        {/if}

        <!-- ── ERRORS TAB ────────────────────────────────────────────────────── -->
      {:else if activeTab === "errors"}
        {#if loadingErrors}
          <div class="space-y-2">
            {#each Array(4) as _}
              <div class="h-20 bg-gray-900 rounded-xl border border-gray-800 animate-pulse"></div>
            {/each}
          </div>
        {:else if recentErrors.length === 0}
          <div class="text-center py-24 text-gray-600">
            <svg
              class="w-10 h-10 mx-auto mb-3 opacity-30"
              fill="none"
              viewBox="0 0 24 24"
              stroke="currentColor"
            >
              <path
                stroke-linecap="round"
                stroke-linejoin="round"
                stroke-width="1.5"
                d="M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z"
              />
            </svg>
            <p class="text-sm">No errors recorded. All good!</p>
          </div>
        {:else}
          <div class="space-y-2">
            {#each recentErrors as err}
              <div class="bg-gray-900 border border-red-900/40 rounded-xl px-4 py-3">
                <div class="flex items-start justify-between gap-4">
                  <div class="flex-1 min-w-0">
                    <div class="flex items-center gap-2 mb-1">
                      <span
                        class="text-[11px] font-semibold px-2 py-0.5 rounded-full bg-red-500/10 text-red-400 border border-red-500/20"
                        >{err.task}</span
                      >
                      <span class="text-xs text-gray-600">{relativeTime(err.time)}</span>
                    </div>
                    <p class="text-sm text-red-300 font-mono break-all">{err.error}</p>
                    {#if err.url}
                      <p class="text-xs text-gray-600 mt-1 truncate">{err.url}</p>
                    {/if}
                  </div>
                </div>
              </div>
            {/each}
          </div>
        {/if}

        <!-- ── CACHE TAB ──────────────────────────────────────────────────────── -->
      {:else if activeTab === "cache"}
        <!-- Cache header -->
        <div class="flex flex-col sm:flex-row sm:items-center justify-between gap-4 mb-6">
          <p class="text-sm text-gray-500">
            {cacheEntries.length} key{cacheEntries.length !== 1 ? "s" : ""} in store
            {#if filteredCache.length !== cacheEntries.length}
              <span class="text-gray-600"> · {filteredCache.length} shown</span>
            {/if}
          </p>
          <div class="flex items-center gap-2">
            {#if !flushConfirm}
              <button
                on:click={() => (flushConfirm = true)}
                disabled={loadingCache || cacheEntries.length === 0}
                class="flex items-center gap-2 px-3 py-1.5 rounded-lg bg-red-900/30 hover:bg-red-900/60 border border-red-800 text-red-400 text-xs font-medium transition-colors disabled:opacity-40"
              >
                <svg
                  class="w-3.5 h-3.5"
                  fill="none"
                  viewBox="0 0 24 24"
                  stroke="currentColor"
                  stroke-width="2"
                >
                  <path
                    stroke-linecap="round"
                    stroke-linejoin="round"
                    d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16"
                  />
                </svg>
                Flush All
              </button>
            {:else}
              <div class="flex items-center gap-2">
                <span class="text-xs text-red-400">Delete all {cacheEntries.length} keys?</span>
                <button
                  on:click={flushAll}
                  class="px-3 py-1.5 rounded-lg bg-red-700 hover:bg-red-600 text-white text-xs font-semibold transition-colors"
                  >Confirm</button
                >
                <button
                  on:click={() => (flushConfirm = false)}
                  class="px-3 py-1.5 rounded-lg bg-gray-800 hover:bg-gray-700 text-gray-300 text-xs font-semibold transition-colors"
                  >Cancel</button
                >
              </div>
            {/if}
          </div>
        </div>

        <!-- Filters -->
        {#if cacheEntries.length > 0}
          <div class="flex flex-col sm:flex-row gap-3 mb-4">
            <input
              bind:value={filterText}
              placeholder="Search keys…"
              class="flex-1 bg-gray-900 border border-gray-800 rounded-lg px-3 py-2 text-sm text-gray-200 placeholder-gray-500 focus:outline-none focus:border-gray-600"
            />
            <select
              bind:value={filterPrefix}
              class="bg-gray-900 border border-gray-800 rounded-lg px-3 py-2 text-sm text-gray-200 focus:outline-none focus:border-gray-600"
            >
              <option value="">All prefixes</option>
              {#each prefixes as p}
                <option value={p}>{p}</option>
              {/each}
            </select>
          </div>
        {/if}

        {#if loadingCache && cacheEntries.length === 0}
          <div class="space-y-2">
            {#each Array(5) as _}
              <div class="h-14 bg-gray-900 rounded-xl border border-gray-800 animate-pulse"></div>
            {/each}
          </div>
        {:else if !loadingCache && cacheEntries.length === 0}
          <div class="text-center py-24 text-gray-600">
            <svg
              class="w-10 h-10 mx-auto mb-3 opacity-30"
              fill="none"
              viewBox="0 0 24 24"
              stroke="currentColor"
            >
              <path
                stroke-linecap="round"
                stroke-linejoin="round"
                stroke-width="1.5"
                d="M5 8h14M5 8a2 2 0 110-4h14a2 2 0 110 4M5 8v10a2 2 0 002 2h10a2 2 0 002-2V8m-9 4h4"
              />
            </svg>
            <p class="text-sm">Cache is empty.</p>
          </div>
        {:else if filteredCache.length === 0 && cacheEntries.length > 0}
          <p class="text-center text-sm text-gray-600 py-10">No keys match your filter.</p>
        {:else}
          <div class="space-y-2">
            {#each filteredCache as entry (entry.key)}
              {@const expanded = expandedKeys.has(entry.key)}
              <div
                class="bg-gray-900 border border-gray-800 rounded-xl overflow-hidden hover:border-gray-700 transition-colors"
              >
                <div class="flex items-center gap-3 px-4 py-3">
                  <button
                    on:click={() => toggleExpand(entry.key)}
                    class="flex-shrink-0 text-gray-600 hover:text-gray-300 transition-colors"
                    aria-label={expanded ? "Collapse" : "Expand"}
                  >
                    <svg
                      class="w-4 h-4 transition-transform duration-150 {expanded
                        ? 'rotate-90'
                        : ''}"
                      viewBox="0 0 20 20"
                      fill="currentColor"
                    >
                      <path
                        fill-rule="evenodd"
                        d="M7.293 4.293a1 1 0 011.414 0l5 5a1 1 0 010 1.414l-5 5a1 1 0 01-1.414-1.414L11.586 10 7.293 5.707a1 1 0 010-1.414z"
                        clip-rule="evenodd"
                      />
                    </svg>
                  </button>
                  <span
                    class="flex-shrink-0 text-[10px] font-semibold uppercase tracking-wide px-2 py-0.5 rounded bg-gray-800 text-gray-500 border border-gray-700"
                    >{entry.prefix}</span
                  >
                  <button
                    on:click={() => toggleExpand(entry.key)}
                    class="flex-1 text-left font-mono text-sm text-gray-300 truncate hover:text-white transition-colors"
                    title={entry.key}>{entry.key}</button
                  >
                  <span
                    class="flex-shrink-0 text-xs font-medium {ttlColor(
                      entry.ttl_seconds
                    )} tabular-nums w-20 text-right">{formatTTL(entry.ttl_seconds)}</span
                  >
                  <button
                    on:click={() => deleteKey(entry.key)}
                    disabled={deletingKey === entry.key}
                    class="flex-shrink-0 ml-1 p-1.5 rounded-md text-gray-700 hover:text-red-400 hover:bg-red-900/20 transition-colors disabled:opacity-40"
                    aria-label="Delete key"
                    title="Delete key"
                  >
                    <svg
                      class="w-4 h-4 {deletingKey === entry.key ? 'animate-spin' : ''}"
                      fill="none"
                      viewBox="0 0 24 24"
                      stroke="currentColor"
                      stroke-width="2"
                    >
                      <path
                        stroke-linecap="round"
                        stroke-linejoin="round"
                        d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16"
                      />
                    </svg>
                  </button>
                </div>
                {#if expanded}
                  <div class="border-t border-gray-800 px-4 py-3">
                    <pre
                      class="text-xs text-emerald-300 font-mono bg-gray-950 rounded-lg p-4 overflow-x-auto whitespace-pre-wrap break-all leading-relaxed max-h-80 overflow-y-auto">{formatJSON(
                        entry.value
                      )}</pre>
                  </div>
                {/if}
              </div>
            {/each}
          </div>
        {/if}
      {/if}
    </div>
  </div>
{/if}
