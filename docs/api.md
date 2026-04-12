# API Reference

## Interactive Docs (Swagger UI)

When the server is running, the full OpenAPI spec is browsable at:

```
http://localhost:8080/swagger/index.html
```

The raw spec files are also available:

| Format | URL |
|---|---|
| JSON | `http://localhost:8080/swagger/doc.json` |
| YAML | `server/internal/docs/swagger.yaml` (generated, committed) |

To regenerate the spec after editing handler annotations:

```bash
cd server
swag init -g cmd/safesurf/safesurf.go -o internal/docs
```

Install the CLI once with: `go install github.com/swaggo/swag/cmd/swag@v1.16.4`

---

## Endpoints

All endpoints are under `GET /api/v1/` and accept a `url` query parameter (max 2048 chars).

### Analysis

| Method | Path | Description |
|---|---|---|
| `GET` | `/api/v1/analyze` | Full URL analysis — runs all checks in parallel, returns scored report |

### URL Structure

| Method | Path | Description |
|---|---|---|
| `GET` | `/api/v1/length` | Check if URL exceeds safe length |
| `GET` | `/api/v1/depth` | Check if URL path is suspiciously deep |
| `GET` | `/api/v1/punycode` | Detect IDN/punycode characters |
| `GET` | `/api/v1/ip/check` | Detect raw IP address usage |
| `GET` | `/api/v1/url-shortener` | Detect known URL shortener services |
| `GET` | `/api/v1/trusted-tld` | Check for trusted TLD (gov/edu) |
| `GET` | `/api/v1/risky-tld` | Check for high-risk TLD |

### DNS / Infrastructure

| Method | Path | Description |
|---|---|---|
| `GET` | `/api/v1/ip/resolve` | Resolve domain to IP addresses |
| `GET` | `/api/v1/rank` | Global popularity rank of the domain |
| `GET` | `/api/v1/domain-info` | WHOIS / RDAP registration data |

### Security

| Method | Path | Description |
|---|---|---|
| `GET` | `/api/v1/hsts` | Check if host enforces HSTS |
| `GET` | `/api/v1/redirects` | Follow and report redirect chain |
| `GET` | `/api/v1/status-code` | Fetch HTTP status code |

### Utility

| Method | Path | Description |
|---|---|---|
| `GET` | `/health` | Service liveness check |
| `GET` | `/api/v1/health` | Same, versioned |
| `GET` | `/api/v1/screenshot` | Headless screenshot of the URL |
| `DELETE` | `/api/v1/cache` | Flush the Valkey cache |
| `GET` | `/metrics` | Prometheus metrics scrape endpoint |
| `GET` | `/swagger/*` | Swagger UI and spec |

---

## Rate Limiting

20 requests per minute per IP. Headers returned on every response:

```
X-RateLimit-Limit: 20
X-RateLimit-Remaining: 19
X-RateLimit-Reset: <unix timestamp>
```

Exceeding the limit returns `429 Too Many Requests`.

---

## Error Format

All 4xx/5xx responses return JSON:

```json
{ "error": "description of the problem" }
```

---

## Example

```bash
curl "http://localhost:8080/api/v1/analyze?url=https://example.com"
```

See the main [README](../README.md) for a full example response.
