# Architecture

## Overview

url.vet is a real-time URL analysis engine. When a URL is submitted — via the web UI, Chrome extension, or REST API — the Go backend runs **18 concurrent analyzers** across 7 signal categories, aggregates a trust/risk score, assigns a verdict, and returns a fully explainable report. Results are cached in Valkey so repeat lookups are instant.

---

## System Architecture

![url.vet Architecture](../assets/architecture.png)

Four containerized services run on a shared Docker bridge network (`urlvet-net`). The Go backend is the **only** service that makes outbound calls — the frontend, Chrome, and cache are strictly internal.

| Service | Container | Role | Port |
|---|---|---|---|
| `urlvet-web` | SvelteKit UI | Renders the web interface; proxies API calls to the backend | `:3000` prod · `:5173` dev |
| `urlvet-backend` | Go REST API | Validates URLs, runs analyzers, aggregates scores, manages cache | `:8080` |
| `urlvet-chrome` | Headless Chrome | Takes page screenshots and serves content via WebSocket (chromedp) | `:9222` |
| `urlvet-valkey` | Valkey (Redis-compatible) | LRU result cache, volume-persisted across restarts | `:6379` |

**External services** (reached only by the backend over HTTPS/TCP):
- **PhishTank** — confirmed phishing database lookups
- **DNS resolvers** — NS, MX, IP resolution checks
- **WHOIS servers** — domain age and registration data

Clients (browser, Chrome extension, API consumers) communicate directly with the Go backend. The SvelteKit frontend forwards all `/api/v1/` requests server-side to avoid CORS complexity.

---

## Request Lifecycle

![url.vet Analyzer Pipeline](../assets/pipeline.png)

```
Client
  │
  │  GET /api/v1/analyze?url=...
  ▼
Go Backend
  ├─ 1. Validate & normalize URL (add scheme if missing, reject private IPs)
  ├─ 2. Check Valkey cache
  │      └─ HIT  → return full cached result immediately (sub-millisecond)
  │      └─ MISS → continue
  ├─ 3. Launch 18 goroutines via sync.WaitGroup
  │      ├─ Each task runs independently; panics are recovered per-task
  │      ├─ Tasks share a read-only Input struct and write to a mutex-guarded Output
  │      └─ All 18 complete (or timeout) before proceeding
  ├─ 4. Aggregate scores → apply formula → assign verdict
  ├─ 5. Store result in Valkey (24 h TTL)
  └─ 6. Return: trust score · verdict · per-signal reasons · redirect chain ·
              screenshot · per-task timings
```

---

## Detection Engine

18 goroutines run across **7 signal categories**, producing **33 individual signals**. Every check emits a labeled reason string — good, bad, or neutral — so the final score is always fully explainable.

### Scoring Formula

```
finalScore = clamp(50 + (trustScore − riskScore) × 0.5, 0, 100)
```

- **50** is the neutral baseline — an unknown URL with no signals scores exactly 50
- Trust signals pull the score up; risk signals pull it down, each weighted at 0.5× so neither dominates
- Both `trustScore` and `riskScore` are individually clamped to 0–100 before the formula runs

| Range | Verdict |
|---|---|
| ≥ 65 | Safe |
| 30 – 64 | Suspicious |
| < 30 | Risky |

### Signal Categories

**URL Signals** — 8 checks, purely structural, no network call

1. Raw IP address as hostname
2. Punycode / IDN encoding (lookalike domain spoofing)
3. URL shortener (hides true destination)
4. Excessive URL length
5. Excessive URL path depth
6. Phishing keywords in URL path (`login`, `verify`, `secure`, `update` …)
7. Excessive subdomain count
8. Non-ASCII Unicode characters in hostname (IDN homograph attack)

**HTTP / Network** — 4 checks, single HTTP request via `httpCombinedTask`

9. Redirect chain hop count
10. Cross-domain redirect (final destination differs from source)
11. HSTS support
12. HTTP status code

**DNS** — 3 checks

13. NS record validity
14. MX record validity
15. IP resolution

**TLS / SSL** — 2 checks, single TLS handshake

16. TLS presence and hostname mismatch
17. Certificate chain — validity, expiry, issuer, CT log status, known-bad fingerprints

**Domain Intelligence** — 6 checks

18. Domain rank (position in top-1M global popularity list)
19. TLD trust / risk / ICANN status
20. Domain age via WHOIS (newly registered = high risk)
21. DNSSEC (cryptographic DNS response integrity)
22. Shannon entropy score (flags algorithmically generated domains)
23. Typosquatting & combo-squatting across 500+ known brands

**Content Analysis** — 8 checks, one HTTP GET to fetch page HTML

24. Login form on unranked or newly registered domain
25. Payment form (credit card, CVV fields)
26. Personal information form
27. Hidden `<iframe>` (credential theft / clickjacking vector)
28. Tracking pixels (1×1 hidden images)
29. Brand name in page content vs. hosting domain
30. Form submitting to an external domain
31. Password field over unencrypted HTTP

**Threat Intelligence** — 2 checks

32. PhishTank confirmed phishing (community-verified)
33. PhishTank reported phishing (awaiting verification, 3 h cache)

---

## Code Layout

```
server/
├── cmd/urlvet/           entry point — init, router setup, graceful shutdown
├── internal/
│   ├── analyzer/
│   │   ├── analyze.go      task registration, cache integration
│   │   ├── runner.go       goroutine runner with panic recovery
│   │   ├── tasks.go        18 task implementations
│   │   └── result.go       score aggregation, verdict assignment
│   ├── handler/
│   │   ├── router.go       Gin router, middleware wiring
│   │   ├── analyze.go      /api/v1/analyze handler
│   │   └── middleware/     rate limiter, auth, Prometheus, request logger
│   ├── service/
│   │   ├── checks/         18 individual analyzer implementations
│   │   ├── screenshot/     headless Chrome integration (chromedp)
│   │   ├── cache/          Valkey client wrapper
│   │   ├── threatfeeds/    PhishTank client
│   │   └── typosquat/      brand similarity engine
│   ├── logger/             centralized slog-based logger (colors in DEV, JSON in prod)
│   └── admintoken/         admin JWT issuance and verification
web/
├── website/                SvelteKit UI
└── chrome-extension/       browser extension
docker/
├── dev/                    dev Compose (hot reload, exposed ports)
└── prod/                   prod Compose (optimized builds, restart policies)
docs/                       API reference, setup guide, architecture, security
```

---

## Deployment

Two fully separated Docker Compose stacks share the same image definitions but differ in configuration:

| | Dev | Prod |
|---|---|---|
| Backend | Air hot-reload, source mounted as volume | Compiled binary in distroless image |
| Frontend | Vite dev server `:5173` | Static build served by Nginx `:3000` |
| Chrome | Same `chromedp/headless-shell` image | Same |
| Valkey | Port exposed (`:6379`) for local inspection | Port not exposed; internal only |
| ENV | `ENV=DEV` — colored logs, debug endpoints | `ENV=PROD` — JSON logs, info level only |

Start everything with one command:

```bash
make start       # production stack
make dev         # development stack (hot reload)
```

See [docs/setup.md](setup.md) for full setup instructions including `.env` configuration.
