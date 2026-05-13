# Maintenance

Routine operational tasks for keeping a url.vet deployment healthy.

---

## Cache management

### Flush the entire cache

Requires a valid admin Bearer token. Obtain one via `POST /api/v1/admin/login` with your admin password.

```bash
TOKEN=$(curl -s -X POST http://localhost:8080/api/v1/admin/login \
  -H "Content-Type: application/json" \
  -d '{"password":"your-password"}' | jq -r '.token')

curl -X DELETE http://localhost:8080/api/v1/cache \
  -H "Authorization: Bearer $TOKEN"
```

### Delete a specific cache entry

```bash
curl -X DELETE "http://localhost:8080/api/v1/cache/full_result:https://example.com" \
  -H "Authorization: Bearer $TOKEN"
```

### Inspect cached keys directly in Valkey

```bash
docker exec -it urlvet-valkey redis-cli -a "$CACHE_PASSWORD"

# List all keys
KEYS *

# Check TTL of a key
TTL "full_result:https://example.com"

# Delete a key
DEL "full_result:https://example.com"

# Memory usage
INFO memory
```

### Cache TTLs

| Cache type | TTL | Key prefix |
|---|---|---|
| Full analysis result | 24 hours | `full_result:` |
| Content analysis | configurable | `content_check:` |
| PhishTank lookup | 3 hours | `phishtank:` |
| HTTP combined result | configurable | `http_combined:` |

---

## Log management

Logs are written to `server/logs/YYYY-MM/DD.log` (when `LOG_DIR=logs`) and always to stderr.

### View today's log

```bash
cat server/logs/$(date +%Y-%m)/$(date +%d).log

# Follow live
tail -f server/logs/$(date +%Y-%m)/$(date +%d).log
```

### View container logs

```bash
make logs          # production
make dev-logs      # development
```

### Purge old logs

Log files rotate automatically at midnight — the logger opens a new file per calendar day. Clean up old files periodically:

```bash
# Delete logs older than 30 days
find server/logs -name "*.log" -mtime +30 -delete

# Remove empty monthly directories
find server/logs -type d -empty -delete
```

Add to a host cron job (`crontab -e`):

```cron
0 3 * * * find /path/to/urlvet/server/logs -name "*.log" -mtime +30 -delete
```

---

## Updating data files

### Domain rank list

The top-1M domain rank list is bundled at build time from `server/assets/`. To refresh it with the latest Tranco list:

```bash
cd server/assets
curl -L "https://tranco-list.eu/top-1m.csv.zip" -o top-1m.csv.zip
unzip -o top-1m.csv.zip
rm top-1m.csv.zip
```

Rebuild and restart the backend to pick up the new list.

### Typosquatting brand list

The brand list lives in `server/assets/`. Edit it directly and restart the backend.

---

## Container updates

```bash
# Pull latest code and rebuild all images
git pull
make start
```

Update only the backend:

```bash
docker compose -f docker/prod/docker-compose.prod.yml up -d --build backend
```

Update the headless Chrome image:

```bash
docker pull chromedp/headless-shell:latest
docker compose -f docker/prod/docker-compose.prod.yml up -d chrome
```

---

## Valkey backup and restore

Valkey snapshots to the `valkey_data` Docker volume automatically. To back up manually:

```bash
docker run --rm \
  -v urlvet_valkey_data:/data \
  -v $(pwd):/backup \
  alpine tar czf /backup/valkey-$(date +%Y%m%d).tar.gz /data
```

To restore:

```bash
docker run --rm \
  -v urlvet_valkey_data:/data \
  -v $(pwd):/backup \
  alpine tar xzf /backup/valkey-YYYYMMDD.tar.gz -C /
```

---

## Health checks

```bash
# Backend liveness
curl http://localhost:8080/health

# Prometheus metrics snapshot
curl -s http://localhost:8080/metrics | grep urlvet_

# Container status and restart counts
docker compose -f docker/prod/docker-compose.prod.yml ps
docker stats --no-stream
```
