# Performance

## Scan latency

A full analysis (`/api/v1/analyze`) runs 18 checks concurrently. Wall time is dominated by the slowest check in the batch.

| Scenario | Typical latency |
|---|---|
| Cache hit (Valkey) | < 5 ms |
| Cache miss — structural checks only | 100–300 ms |
| Cache miss — with DNS + TLS + WHOIS | 400–900 ms |
| Cache miss — with content fetch + screenshot | 2–15 s |
| Cache miss — with slow WHOIS server | up to 30 s |

Screenshot is the single biggest latency contributor. It runs headless Chrome, navigates to the page, and waits for network idle — this can take 5–30s on slow or complex pages.

---

## Concurrency model

18 goroutines launch simultaneously via `sync.WaitGroup`. Each goroutine is independent — a slow or hung check does not block others. Panics are recovered per-goroutine so one failing check does not abort the request.

Go goroutines are cheap (~2KB initial stack, grows on demand). 18 concurrent goroutines per request adds negligible scheduling overhead.

---

## Caching impact

The Valkey cache with 24h TTL dramatically reduces load for repeated URLs. Cache hit rate grows quickly for any deployment that sees repeat traffic (e.g. the same phishing URL submitted multiple times during a campaign).

At the 256MB Valkey memory limit with `allkeys-lru` eviction, the cache holds roughly 5,000–20,000 full results depending on response size.

---

## Rate limiting

20 requests per minute per IP, enforced in the `RateLimiter` middleware backed by Valkey. Applies to all `/api/v1/` endpoints. Returns `429 Too Many Requests` with `Retry-After` and `X-RateLimit-*` headers.

---

## Resource usage (steady state, production)

| Container | Typical RAM | CPU (idle) |
|---|---|---|
| `safesurf-backend` | 50–150 MB | < 1% |
| `safesurf-web` | 20–50 MB (Nginx) | < 1% |
| `safesurf-chrome` | 300–600 MB | < 1% |
| `safesurf-valkey` | 10–256 MB (bounded by `maxmemory`) | < 1% |

Chrome dominates memory. If screenshots are not needed, the Chrome container can be removed and the screenshot task disabled to save 300–600MB.

---

## Prometheus metrics

All per-route latency and per-task duration are tracked. Scrape `http://localhost:8080/metrics`:

```
safesurf_task_duration_seconds{task="..."}         histogram per analyzer task
safesurf_http_request_duration_seconds{path="..."}  histogram per route
safesurf_risk_score                                  distribution of risk scores
safesurf_trust_score                                 distribution of trust scores
```

---

## Tuning knobs

| Knob | Where | Effect |
|---|---|---|
| `CACHE_POOL_SIZE` | `server/.env` | Valkey connection pool — increase under high concurrency |
| `CACHE_MIN_IDLE_CONNS` | `server/.env` | Minimum warm connections — reduces connection latency spikes |
| Valkey `maxmemory` | `docker/prod/valkey.conf` | Cap memory; `allkeys-lru` evicts cold entries automatically |
| Rate limit | `internal/handler/router.go` — `RateLimiter(20, time.Minute)` | Adjust requests-per-minute per IP |

---

## Profiling

To profile the Go backend in development:

```bash
# Add to safesurf.go (dev only — do not expose in production):
import _ "net/http/pprof"
go http.ListenAndServe(":6060", nil)
```

```bash
go tool pprof http://localhost:6060/debug/pprof/profile?seconds=30
go tool pprof http://localhost:6060/debug/pprof/heap
```
