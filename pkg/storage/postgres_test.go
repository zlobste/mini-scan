package storage_test

import (
	"context"
	"os"
	"testing"

	"github.com/zlobste/mini-scan/pkg/storage"
)

func getTestDB(t *testing.T) storage.Store {
	connStr := os.Getenv("TEST_DATABASE_URL")
	if connStr == "" {
		t.Skip("TEST_DATABASE_URL not set")
	}

	store, err := storage.NewPostgresStore(connStr)
	if err != nil {
		t.Fatalf("failed to connect to test database: %v", err)
	}

	return store
}

func TestUpsertScan_NewScan(t *testing.T) {
	store := getTestDB(t)
	defer store.Close()

	ctx := context.Background()
	scan := storage.Scan{
		IP:            "192.168.1.1",
		Port:          80,
		Service:       "HTTP",
		LastScannedAt: 1000,
		ResponseText:  "test response",
	}

	err := store.UpsertScan(ctx, scan)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	retrieved, found, err := store.GetScan(ctx, scan.IP, scan.Port, scan.Service)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if !found {
		t.Fatal("expected scan to be found")
	}

	if retrieved.ResponseText != scan.ResponseText {
		t.Errorf("expected %q, got %q", scan.ResponseText, retrieved.ResponseText)
	}
}

func TestUpsertScan_UpdateWithNewerTimestamp(t *testing.T) {
	store := getTestDB(t)
	defer store.Close()

	ctx := context.Background()
	scan1 := storage.Scan{
		IP:            "192.168.1.2",
		Port:          443,
		Service:       "HTTPS",
		LastScannedAt: 1000,
		ResponseText:  "first response",
	}

	err := store.UpsertScan(ctx, scan1)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	scan2 := storage.Scan{
		IP:            "192.168.1.2",
		Port:          443,
		Service:       "HTTPS",
		LastScannedAt: 2000,
		ResponseText:  "second response",
	}

	err = store.UpsertScan(ctx, scan2)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	retrieved, found, err := store.GetScan(ctx, scan1.IP, scan1.Port, scan1.Service)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if !found {
		t.Fatal("expected scan to be found")
	}

	if retrieved.ResponseText != scan2.ResponseText {
		t.Errorf("expected %q, got %q", scan2.ResponseText, retrieved.ResponseText)
	}
	if retrieved.LastScannedAt != scan2.LastScannedAt {
		t.Errorf("expected timestamp %d, got %d", scan2.LastScannedAt, retrieved.LastScannedAt)
	}
}

func TestUpsertScan_IgnoreOlderTimestamp(t *testing.T) {
	store := getTestDB(t)
	defer store.Close()

	ctx := context.Background()
	scan1 := storage.Scan{
		IP:            "192.168.1.3",
		Port:          8080,
		Service:       "HTTP",
		LastScannedAt: 2000,
		ResponseText:  "newer response",
	}

	err := store.UpsertScan(ctx, scan1)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	scan2 := storage.Scan{
		IP:            "192.168.1.3",
		Port:          8080,
		Service:       "HTTP",
		LastScannedAt: 1000,
		ResponseText:  "older response",
	}

	err = store.UpsertScan(ctx, scan2)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	retrieved, found, err := store.GetScan(ctx, scan1.IP, scan1.Port, scan1.Service)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if !found {
		t.Fatal("expected scan to be found")
	}

	if retrieved.ResponseText != scan1.ResponseText {
		t.Errorf("expected %q (newer response), got %q", scan1.ResponseText, retrieved.ResponseText)
	}
}

func TestGetScan_NotFound(t *testing.T) {
	store := getTestDB(t)
	defer store.Close()

	ctx := context.Background()
	retrieved, found, err := store.GetScan(ctx, "nonexistent.ip", 99999, "UNKNOWN")

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if found {
		t.Error("expected found to be false for nonexistent scan")
	}

	if retrieved != (storage.Scan{}) {
		t.Errorf("expected empty Scan, got %+v", retrieved)
	}
}
