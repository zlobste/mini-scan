package processor

import (
	"context"
	"fmt"

	"github.com/zlobste/mini-scan/pkg/scanning"
	"github.com/zlobste/mini-scan/pkg/storage"
)

type Processor interface {
	ProcessScan(ctx context.Context, decoded scanning.DecodedScan) error
}

type processor struct {
	store storage.Store
}

func New(store storage.Store) Processor {
	return &processor{store: store}
}

func (p *processor) ProcessScan(ctx context.Context, decoded scanning.DecodedScan) error {
	existing, found, err := p.store.GetScan(ctx, decoded.IP, decoded.Port, decoded.Service)
	if err != nil {
		return fmt.Errorf("failed to get existing scan: %w", err)
	}

	// Skip update if we already have a newer or equal scan.
	// This handles out-of-order messages: if a 24-hour-old scan arrives after newer data,
	// it's silently ignored. Keeps the latest scan regardless of message arrival order.
	if found && existing.LastScannedAt >= decoded.Timestamp {
		return nil
	}

	scan := storage.Scan{
		IP:            decoded.IP,
		Port:          decoded.Port,
		Service:       decoded.Service,
		LastScannedAt: decoded.Timestamp,
		ResponseText:  decoded.ResponseText,
	}

	if err := p.store.UpsertScan(ctx, scan); err != nil {
		return fmt.Errorf("failed to upsert scan: %w", err)
	}

	return nil
}
