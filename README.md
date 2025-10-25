# Mini-Scan

A Go application that processes internet scan results from Google Cloud Pub/Sub, handles multiple message formats, and persists data to PostgreSQL.

## Quick Start

```bash
# Set up environment
cp .env.example .env

# Start the stack
docker compose up -d

# Check logs
docker compose logs -f processor
```

## Architecture

The system consists of four components:

1. **Scanner** - Publishes scan results to Pub/Sub topic every second in V1 or V2 format
2. **Pub/Sub** - Distributes messages with at-least-once delivery semantics
3. **Processor** - Consumes and processes messages:
   - Decodes V1 (base64-encoded) and V2 (plain string) formats
   - Deduplicates based on timestamp comparison
   - Persists to database with update-on-conflict logic
4. **PostgreSQL** - Stores unique scans by `(ip, port, service)` composite key

**Scaling:** Multiple processor instances share a single database. Pub/Sub automatically distributes work across them.

## Configuration

Environment variables (see `.env.example` for defaults):

| Variable | Default | Description |
|----------|---------|-------------|
| `POSTGRES_HOST` | `postgres` | PostgreSQL host |
| `POSTGRES_PORT` | `5432` | PostgreSQL port |
| `POSTGRES_USER` | `postgres` | Database user |
| `POSTGRES_PASSWORD` | `postgres` | Database password |
| `POSTGRES_DB` | `scans` | Database name |
| `PUBSUB_EMULATOR_HOST` | `localhost:8085` | Pub/Sub emulator endpoint |
| `PUBSUB_PROJECT` | `test-project` | Project ID |
| `SUBSCRIPTION` | `scan-sub` | Subscription to consume |

## Data Storage

**Table:** `scans`

| Column | Type | Notes |
|--------|------|-------|
| `ip` | TEXT | IP address |
| `port` | INTEGER | Port number |
| `service` | TEXT | Service name |
| `last_scanned_at` | BIGINT | Unix timestamp |
| `response_text` | TEXT | Service response |

**Key:** Composite primary key on `(ip, port, service)`

## Testing

Run all tests:
```bash
go test ./...
```

Run storage tests (requires PostgreSQL):
```bash
TEST_DATABASE_URL="postgres://postgres:postgres@localhost:5432/scans?sslmode=disable" go test ./pkg/storage/... -v
```

See [TESTING.md](TESTING.md) for manual testing instructions.
