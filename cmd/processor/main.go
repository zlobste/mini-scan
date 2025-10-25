package main

import (
	"context"
	"log/slog"
	"os"

	"cloud.google.com/go/pubsub"
	"github.com/zlobste/mini-scan/pkg/config"
	"github.com/zlobste/mini-scan/pkg/processor"
	"github.com/zlobste/mini-scan/pkg/scanning"
	"github.com/zlobste/mini-scan/pkg/storage"
)

func main() {
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	slog.SetDefault(logger)

	logger.Info("starting processor")

	cfg, err := config.Load()
	if err != nil {
		logger.Error("failed to load config", "error", err)
		os.Exit(1)
	}
	logger.Info("config loaded", "host", cfg.Database.Host, "pubsub_project", cfg.PubSub.ProjectID, "subscription", cfg.PubSub.Subscription)

	store, err := storage.NewPostgresStore(cfg.Database.ConnectionString())
	if err != nil {
		logger.Error("failed to connect to database", "error", err)
		os.Exit(1)
	}
	defer store.Close()
	logger.Info("connected to database")

	ctx := context.Background()

	client, err := pubsub.NewClient(ctx, cfg.PubSub.ProjectID)
	if err != nil {
		logger.Error("failed to create pubsub client", "error", err)
		os.Exit(1)
	}
	defer client.Close()
	logger.Info("pubsub client created")

	sub := client.Subscription(cfg.PubSub.Subscription)
	proc := processor.New(store)

	logger.Info("listening for messages", "subscription", cfg.PubSub.Subscription)

	if err = sub.Receive(ctx, func(ctx context.Context, msg *pubsub.Message) {
		decoded, err := scanning.DecodeMessage(msg.Data)
		if err != nil {
			logger.Error("failed to decode message", "error", err)
			msg.Nack()
			return
		}

		logger.Info("decoded scan", "ip", decoded.IP, "port", decoded.Port, "service", decoded.Service, "timestamp", decoded.Timestamp)

		if err := proc.ProcessScan(ctx, decoded); err != nil {
			logger.Error("failed to process scan", "ip", decoded.IP, "port", decoded.Port, "service", decoded.Service, "error", err)
			msg.Nack()
			return
		}

		logger.Debug("scan processed successfully", "ip", decoded.IP, "port", decoded.Port, "service", decoded.Service)
		msg.Ack()
	}); err != nil {
		logger.Error("receive error", "error", err)
		os.Exit(1)
	}
}
