package storage

import (
	"context"
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

type PostgresStore struct {
	db *sql.DB
}

func NewPostgresStore(connStr string) (Store, error) {
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	if err := db.PingContext(context.Background()); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(5)

	return &PostgresStore{db: db}, nil
}

func (ps *PostgresStore) UpsertScan(ctx context.Context, scan Scan) error {
	query := `
		INSERT INTO scans (ip, port, service, last_scanned_at, response_text)
		VALUES ($1, $2, $3, $4, $5)
		ON CONFLICT (ip, port, service) DO UPDATE SET
			last_scanned_at = EXCLUDED.last_scanned_at,
			response_text = EXCLUDED.response_text,
			updated_at = CURRENT_TIMESTAMP
		WHERE scans.last_scanned_at < EXCLUDED.last_scanned_at
	`

	_, err := ps.db.ExecContext(ctx, query,
		scan.IP,
		scan.Port,
		scan.Service,
		scan.LastScannedAt,
		scan.ResponseText,
	)

	if err != nil {
		return fmt.Errorf("failed to upsert scan: %w", err)
	}

	return nil
}

func (ps *PostgresStore) GetScan(ctx context.Context, ip string, port uint32, service string) (Scan, bool, error) {
	query := `
		SELECT ip, port, service, last_scanned_at, response_text
		FROM scans
		WHERE ip = $1 AND port = $2 AND service = $3
	`

	var scan Scan
	err := ps.db.QueryRowContext(ctx, query, ip, port, service).Scan(
		&scan.IP,
		&scan.Port,
		&scan.Service,
		&scan.LastScannedAt,
		&scan.ResponseText,
	)

	if err == sql.ErrNoRows {
		return Scan{}, false, nil
	}
	if err != nil {
		return Scan{}, false, fmt.Errorf("failed to get scan: %w", err)
	}

	return scan, true, nil
}

func (ps *PostgresStore) Close() error {
	return ps.db.Close()
}
