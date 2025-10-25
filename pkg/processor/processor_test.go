package processor_test

import (
	"context"
	"testing"

	"github.com/zlobste/mini-scan/pkg/processor"
	"github.com/zlobste/mini-scan/pkg/scanning"
	"github.com/zlobste/mini-scan/pkg/storage"
)

type mockStore struct {
	scans map[string]*storage.Scan
}

func newMockStore() *mockStore {
	return &mockStore{
		scans: make(map[string]*storage.Scan),
	}
}

func (m *mockStore) key(ip string, port uint32, service string) string {
	return ip + ":" + string(rune(port)) + ":" + service
}

func (m *mockStore) UpsertScan(ctx context.Context, scan storage.Scan) error {
	m.scans[m.key(scan.IP, scan.Port, scan.Service)] = &scan
	return nil
}

func (m *mockStore) GetScan(ctx context.Context, ip string, port uint32, service string) (storage.Scan, bool, error) {
	scan, exists := m.scans[m.key(ip, port, service)]
	if !exists {
		return storage.Scan{}, false, nil
	}
	return *scan, true, nil
}

func (m *mockStore) Close() error {
	return nil
}

func TestProcessScan_NewScan(t *testing.T) {
	store := newMockStore()
	proc := processor.New(store)

	ctx := context.Background()
	decoded := scanning.DecodedScan{
		IP:           "192.168.1.1",
		Port:         80,
		Service:      "HTTP",
		Timestamp:    1000,
		ResponseText: "test response",
	}

	err := proc.ProcessScan(ctx, decoded)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	stored, found, _ := store.GetScan(ctx, decoded.IP, decoded.Port, decoded.Service)
	if !found {
		t.Fatal("expected scan to be found")
	}
	if stored.ResponseText != decoded.ResponseText {
		t.Errorf("expected %q, got %q", decoded.ResponseText, stored.ResponseText)
	}
}

func TestProcessScan_IgnoreOlderScan(t *testing.T) {
	store := newMockStore()
	proc := processor.New(store)

	ctx := context.Background()

	newer := scanning.DecodedScan{
		IP:           "192.168.1.2",
		Port:         443,
		Service:      "HTTPS",
		Timestamp:    2000,
		ResponseText: "newer",
	}

	proc.ProcessScan(ctx, newer)

	older := scanning.DecodedScan{
		IP:           "192.168.1.2",
		Port:         443,
		Service:      "HTTPS",
		Timestamp:    1000,
		ResponseText: "older",
	}

	err := proc.ProcessScan(ctx, older)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	stored, found, _ := store.GetScan(ctx, newer.IP, newer.Port, newer.Service)
	if !found {
		t.Fatal("expected scan to be found")
	}
	if stored.ResponseText != newer.ResponseText {
		t.Errorf("expected %q (newer), got %q", newer.ResponseText, stored.ResponseText)
	}
}

func TestProcessScan_UpdateWithNewerScan(t *testing.T) {
	store := newMockStore()
	proc := processor.New(store)

	ctx := context.Background()

	older := scanning.DecodedScan{
		IP:           "192.168.1.3",
		Port:         8080,
		Service:      "HTTP",
		Timestamp:    1000,
		ResponseText: "older",
	}

	proc.ProcessScan(ctx, older)

	newer := scanning.DecodedScan{
		IP:           "192.168.1.3",
		Port:         8080,
		Service:      "HTTP",
		Timestamp:    2000,
		ResponseText: "newer",
	}

	err := proc.ProcessScan(ctx, newer)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	stored, found, _ := store.GetScan(ctx, newer.IP, newer.Port, newer.Service)
	if !found {
		t.Fatal("expected scan to be found")
	}
	if stored.ResponseText != newer.ResponseText {
		t.Errorf("expected %q (newer), got %q", newer.ResponseText, stored.ResponseText)
	}
}
