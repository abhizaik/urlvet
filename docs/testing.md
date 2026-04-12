# Testing

## Running Tests

All tests live under `server/`. Run them from that directory.

```bash
cd server

# All tests
go test ./...

# Specific package
go test ./internal/analyzer/
go test ./internal/handler/
go test ./internal/service/checks/

# Verbose output
go test -v ./internal/analyzer/ -run TestGenerateResult

# With race detector
go test -race ./...
```

## What's Covered

### `internal/analyzer` — Scorer (`result_test.go`)

Table-driven tests for `GenerateResult()`. Every scoring branch is exercised:

| Test group | What it checks |
|---|---|
| `TestGenerateResult_Verdict` | 28 cases — rank, TLD, HSTS, IP, punycode, homoglyph, phishtank, redirects, domain age, DNSSEC, brand mismatch, login/payment/hidden-iframe signals |
| `TestGenerateResult_ScoreClamping` | Risk + final scores stay within 0–100 under max-risk pile-on |
| `TestGenerateResult_VerdictBoundaries` | Exact formula boundary: `finalScore = trust - risk×0.2` produces correct verdict bucket |

### `internal/handler` — HTTP smoke tests (`handler_test.go`)

Uses `net/http/httptest` against the full router (no mocks). Covers:

| Handler | Cases |
|---|---|
| `GET /health`, `GET /api/v1/health` | 200 + `status: ok` |
| `GET /` | 200 + service metadata |
| `GET /api/v1/analyze` | Missing URL → 400; empty-host URL → 400; URL > 2048 chars → 400 |
| `GET /api/v1/length`, `/depth` | Missing → 400; invalid → 400; valid → 200 |
| `GET /api/v1/punycode` | Missing → 400; invalid → 400; valid → 200 |
| `GET /api/v1/url-shortener` | Missing → 400; invalid → 400; valid → 200 |
| `GET /api/v1/trusted-tld`, `/risky-tld` | Missing → 400; valid → 200 |
| `GET /api/v1/ip/check` | Missing → 400; valid → 200 |
| `GET /metrics` | 200 + non-empty `Content-Type` |

### `internal/service/checks` — URL utilities (`checks_test.go`)

Basic unit test for `TooLongUrl`.

### `internal/service/rank` — Rank loader (`load_test.go`, `lookup_test.go`)

Rank file loading and lookup logic.

## URL Validation Behaviour

`IsValidURL` is permissive by design: bare hostnames (`example.com`) are accepted by prepending `https://`. Only truly malformed inputs (empty host, `http://`) are rejected. Handler tests reflect this.

## Coverage Report

```bash
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out
```

## Adding Tests

- Scorer tests: `server/internal/analyzer/result_test.go` — extend `TestGenerateResult_Verdict` table.
- Handler tests: `server/internal/handler/handler_test.go` — add rows to existing test tables or new `Test*` functions.
- Unit tests: place `_test.go` alongside the file under test.
