package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"log/slog"
	"math/rand"
	"os"
	"time"

	"cloud.google.com/go/pubsub"
	"github.com/zlobste/mini-scan/pkg/scanning"
)

var (
	services = []string{"HTTP", "SSH", "DNS"}
)

func main() {
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	slog.SetDefault(logger)

	projectId := flag.String("project", "test-project", "GCP Project ID")
	topicId := flag.String("topic", "scan-topic", "GCP PubSub Topic ID")
	flag.Parse()

	logger.Info("scanner starting", "project", *projectId, "topic", *topicId)

	ctx := context.Background()

	client, err := pubsub.NewClient(ctx, *projectId)
	if err != nil {
		logger.Error("failed to create pubsub client", "error", err)
		os.Exit(1)
	}
	defer client.Close()
	logger.Info("pubsub client created")

	topic := client.Topic(*topicId)

	logger.Info("publishing scans", "topic", *topicId, "interval", "1s")

	for range time.Tick(time.Second) {
		scan := &scanning.Scan{
			Ip:        fmt.Sprintf("1.1.1.%d", rand.Intn(255)),
			Port:      uint32(rand.Intn(65535)),
			Service:   services[rand.Intn(len(services))],
			Timestamp: time.Now().Unix(),
		}

		serviceResp := fmt.Sprintf("service response: %d", rand.Intn(100))

		if rand.Intn(2) == 0 {
			scan.DataVersion = scanning.V1
			scan.Data = &scanning.V1Data{ResponseBytesUtf8: []byte(serviceResp)}
		} else {
			scan.DataVersion = scanning.V2
			scan.Data = &scanning.V2Data{ResponseStr: serviceResp}
		}

		encoded, err := json.Marshal(scan)
		if err != nil {
			logger.Error("failed to marshal scan", "error", err)
			continue
		}

		_, err = topic.Publish(ctx, &pubsub.Message{Data: encoded}).Get(ctx)
		if err != nil {
			logger.Error("failed to publish message", "error", err)
			continue
		}

		logger.Info("scan published", "ip", scan.Ip, "port", scan.Port, "service", scan.Service)
	}
}
