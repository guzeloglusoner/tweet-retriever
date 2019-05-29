package main

import (
	"context"
	"time"

	"github.com/segmentio/kafka-go"
)

// ConsumeKafkaTopic is used for consuming messages of the kafka topic
func ConsumeKafkaTopic(hub *Hub) {
	config := kafka.ReaderConfig{
		Brokers:         kafkaBrokerURL,
		GroupID:         kafkaClientID,
		Topic:           kafkaTopic,
		MinBytes:        10e3,
		MaxBytes:        10e6,
		MaxWait:         1 * time.Second,
		ReadLagInterval: 1, // updating the reader lag. if its setted as negative value, then the reporting will be disableds
	}

	reader := kafka.NewReader(config)
	defer reader.Close()

	for {
		m, err := reader.ReadMessage(context.Background())
		if err != nil {
			logger.Printf("something has happened %v", err)
		}

		logger.Printf("topic / partition / offset / value : %v / %v / %v : %s \n", m.Topic, m.Partition, m.Offset, string(m.Value)) // converting byte array to string - deserialization
		hub.inbound <- m.Value
	}
}
