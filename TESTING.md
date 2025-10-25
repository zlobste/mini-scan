# Testing Guide

## Manual Testing

### 1. Start the Stack

```bash
docker compose up -d
```

Wait for all services to be healthy:
```bash
docker compose ps
```

### 2. Verify Scanner is Publishing

Check that the scanner is publishing scans every second:
```bash
docker compose logs -f scanner | head -20
```

Expected output should show `scan published` logs with IP, port, and service information.

### 3. Verify Processor is Consuming

Check that the processor is receiving and processing messages:
```bash
docker compose logs -f processor | head -20
```

Expected output should show:
- `listening for messages`
- `decoded scan` logs with scan details
- `scan processed successfully` debug logs

### 4. Verify Data in Database

Check how many scans have been stored:
```bash
docker compose exec postgres psql -U postgres -d scans -c "SELECT COUNT(*) FROM scans;"
```

View sample scans:
```bash
docker compose exec postgres psql -U postgres -d scans -c "SELECT ip, port, service, last_scanned_at FROM scans LIMIT 5;"
```

### 5. Verify Out-of-Order Handling

The processor handles out-of-order messages by comparing timestamps. To verify:

**Step 1:** Pick an IP from the logs and note its timestamp:
```bash
docker compose logs scanner | grep "scan published" | head -5
# Example: ip=1.1.1.42
```

**Step 2:** Check the database record for that IP:
```bash
docker compose exec postgres psql -U postgres -d scans -c "SELECT ip, port, service, last_scanned_at FROM scans WHERE ip='1.1.1.42' LIMIT 1;"
```

Note the `last_scanned_at` timestamp.

**Step 3:** Wait and watch the processor logs until the same IP appears again:
```bash
docker compose logs -f processor | grep "1.1.1.42"
```

**Step 4:** Once you see the same IP in the logs again, immediately check the database:
```bash
docker compose exec postgres psql -U postgres -d scans -c "SELECT ip, port, service, last_scanned_at FROM scans WHERE ip='1.1.1.42' LIMIT 1;"
```

If the `last_scanned_at` timestamp increased, the processor correctly processed the newer scan and updated the record.

### 6. Stop the Stack

```bash
docker compose down -v
```
