package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"math/rand"
	"time"

	"cloud.google.com/go/pubsub"
	"github.com/censys/scan-takehome/pkg/scanning"
)

var (
	services = []string{"HTTP", "SSH", "DNS"}
)

func main() {
	projectId := flag.String("project", "test-project", "GCP Project ID")
	topicId := flag.String("topic", "scan-topic", "GCP PubSub Topic ID")

	ctx := context.Background()

	client, err := pubsub.NewClient(ctx, *projectId)
	if err != nil {
		panic(err)
	}

	topic := client.Topic(*topicId)

	for range time.Tick(time.Second) {

		var data []byte
		var dataFormat int
		if rand.Intn(2) == 0 {
			dataFormat = scanning.BinaryFormat
			data = []byte(fmt.Sprintf("this is binary data: %d", rand.Intn(100)))
		} else {
			dataFormat = scanning.JsonFormat
			data = []byte(fmt.Sprintf(`{"key": "value-%d"}`, rand.Intn(100)))
		}

		scan := &scanning.Scan{
			Ip:         fmt.Sprintf("1.1.1.%d", rand.Intn(255)),
			Port:       uint32(rand.Intn(65535)),
			Service:    services[rand.Intn(len(services))],
			Timestamp:  time.Now().Unix(),
			DataFormat: dataFormat,
			Data:       data,
		}

		encoded, err := json.Marshal(scan)
		if err != nil {
			panic(err)
		}

		_, err = topic.Publish(ctx, &pubsub.Message{Data: encoded}).Get(ctx)
		if err != nil {
			panic(err)
		}
	}
}
