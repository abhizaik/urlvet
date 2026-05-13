# Design Decisions

Key architectural choices in url.vet and the reasoning behind them.

---

## Go for the backend

**Decision:** The analysis server is written in Go.

**Why:** The core bottleneck is I/O — DNS lookups, HTTP fetches, TLS handshakes, WHOIS queries. Go's goroutines make it trivial to run all 18 checks concurrently with a `sync.WaitGroup`. The result is 300–700ms median scan time with no async/await complexity and clean per-goroutine panic recovery. A Python or Node implementation would need explicit async machinery for the same concurrency and would carry more runtime overhead per request.

---

## No machine learning (intentional)

**Decision:** All detection is rule-based heuristics. No ML model.

**Why:** Explainability is a first-class goal. Every verdict can be traced to specific, named signals with human-readable reasons. An ML model would give accurate verdicts but opaque ones — users and integrators couldn't understand or trust why a URL was flagged. The tradeoff is accepted: url.vet will miss sophisticated phishing that mimics legitimate site structure, but it will never produce a verdict nobody can explain.

This decision may be revisited. If added, ML would be a separate signal alongside existing heuristics, not a replacement.

---

## 18 goroutines, all concurrent

**Decision:** All checks launch simultaneously via `sync.WaitGroup`. No prioritization or staged execution.

**Why:** Simplicity. The longest checks (WHOIS ~1–2s, screenshot ~10–30s) dominate latency regardless of ordering. Running fast structural checks first and gating slow network checks on their results would add coordination complexity for minimal gain — the total wall time is still bounded by the slowest check.

Known tradeoff: under high load, all goroutines (including slow screenshot/WHOIS) consume resources for every request. Making screenshot opt-in is on the roadmap.

---

## Valkey (not PostgreSQL) for caching

**Decision:** Full analysis results are cached in Valkey (Redis-compatible), not a relational database.

**Why:** The primary access pattern is `url → result` — a pure key-value lookup. Valkey handles this in sub-millisecond with built-in LRU eviction and TTL. PostgreSQL would add connection overhead, schema complexity, and index maintenance for a use case that maps directly to a key-value store's strengths.

Tradeoff: no persistent scan history, no queries across results. Scan history (for analytics, training data, abuse detection) requires PostgreSQL and is on the roadmap.

---

## Single HTTP request for redirect + HSTS + status code

**Decision:** `httpCombinedTask` makes one HTTP request shared across three checks (redirects, HSTS header, status code) instead of three separate requests.

**Why:** These three checks all need the same HTTP response. Before the optimization, `CheckRedirects` and `CheckHSTS` each made independent requests to the same URL. The combined task issues a single HEAD request (falling back to GET if needed), extracts all three results, and eliminates two redundant network round-trips.

---

## SvelteKit for the frontend

**Decision:** The web UI is a SvelteKit application served as a static build in production.

**Why:** SvelteKit's server-side rendering lets API calls to the Go backend happen server-side, avoiding CORS issues in the browser. The compiled output is a static Nginx container with no Node.js runtime in production. Svelte's compiler-first approach produces minimal bundle sizes compared to React or Vue equivalents.

---

## Gin as the HTTP framework

**Decision:** The REST API uses the Gin framework.

**Why:** Gin provides routing, middleware chaining, and request binding with minimal overhead. It is the most widely used Go HTTP framework with good community support. The alternative (stdlib `net/http`) would require writing routing and middleware manually; heavier frameworks (Echo, Fiber) offer no meaningful advantage for url.vet's API surface.

---

## Structured logging via slog

**Decision:** All logging uses Go's stdlib `log/slog` wrapped in a centralized logger package.

**Why:** `slog` (Go 1.21+) provides structured, leveled logging with no external dependency. The centralized wrapper (`internal/logger`) lets the format be changed in one place — colored text in DEV, JSON in production — without touching any call site. Full URLs are kept out of logs (scheme+host only) to prevent token leakage in log files.

---

## AGPL-3.0 license

**Decision:** url.vet is licensed under AGPL-3.0 with a separate commercial license option.

**Why:** AGPL requires any modified version run over a network to make its source code available. This prevents organizations from running a proprietary fork as a SaaS without contributing back. The commercial license option lets organizations that cannot comply with AGPL (closed-source SaaS) pay for an exception.

---

## chromedp over a REST screenshot API

**Decision:** Screenshots are taken using chromedp (a Go-native Chrome DevTools Protocol client) against a shared headless Chrome container, not an external screenshot API.

**Why:** No external dependency, no API key, no per-screenshot cost, no data leaving the deployment. The shared browser allocator (`screenshot.GetService()`) reuses a single Chrome instance across requests, avoiding the overhead of launching a new browser per scan. Tradeoff: the Chrome container adds ~300MB memory baseline.
