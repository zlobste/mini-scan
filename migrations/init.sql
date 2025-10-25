CREATE TABLE IF NOT EXISTS scans (
    ip TEXT NOT NULL,
    port INTEGER NOT NULL,
    service TEXT NOT NULL,
    last_scanned_at BIGINT NOT NULL,
    response_text TEXT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (ip, port, service)
);

CREATE INDEX IF NOT EXISTS idx_scans_timestamp ON scans(last_scanned_at DESC);

CREATE INDEX IF NOT EXISTS idx_scans_service_timestamp ON scans(service, last_scanned_at DESC);

GRANT ALL PRIVILEGES ON TABLE scans TO postgres;
GRANT ALL PRIVILEGES ON ALL SEQUENCES IN SCHEMA public TO postgres;
