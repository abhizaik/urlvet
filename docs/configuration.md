# Configuration

All configuration is via environment variables. Copy `server/.env.example` → `server/.env` and fill in your values.

## Backend (`server/.env`)

| Variable | Required | Description |
|---|---|---|
| `CACHE_ADDR` | Yes | Valkey/Redis address — `safesurf-valkey:6379` |
| `CACHE_PASSWORD` | Yes (prod) | Cache auth password — must match `CACHE_PASSWORD` set in `docker-compose.prod.yml` |
| `CACHE_DB` | No | Redis DB index (default `0`) |
| `CACHE_POOL_SIZE` | No | Connection pool size (default `50`) |
| `CACHE_MIN_IDLE_CONNS` | No | Min idle connections (default `10`) |
| `ENV` | No | Set to `DEV` to enable the `/api/v1/test` endpoint |
| `CORS_ALLOWED_ORIGINS` | No | Comma-separated allowed origins (default `*`) — e.g. `https://safesurf.xorwave.com,https://url.vet` |
| `PORT` | No | HTTP port (default `8080`) |
| `ADMIN_PASSWORD_HASH` | Yes | Argon2id hash of the admin password — see [security.md](security.md#setup) |
| `ADMIN_JWT_SECRET` | Yes | Signing secret for session tokens — `openssl rand -hex 32` |

## Docker Compose (dev: `docker/dev/.env`, prod: `docker/prod/.env`)

Create the appropriate `.env` file with secrets used by docker-compose variable substitution:

| Variable | Required | Description |
|---|---|---|
| `CACHE_PASSWORD` | Yes | Password set on the Valkey container (`--requirepass`) and passed to the backend as `CACHE_PASSWORD`. Generate with `openssl rand -hex 32`. |

## Frontend (`web/website/.env`)

| Variable | Required | Description |
|---|---|---|
| `PUBLIC_BASE_URL` | No | Go API base URL (default `http://localhost:8080/api/v1`) |
