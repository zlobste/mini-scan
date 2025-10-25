package storage

import "context"

type Scan struct {
	IP            string
	Port          uint32
	Service       string
	LastScannedAt int64
	ResponseText  string
}

type Store interface {
	UpsertScan(ctx context.Context, scan Scan) error
	GetScan(ctx context.Context, ip string, port uint32, service string) (Scan, bool, error)
	Close() error
}
