package kafka

import (
	"context"
	"log"

	kafka "github.com/segmentio/kafka-go"
)

func StartUserCreatedConsumer(broker string) {
	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers: []string{broker},
		Topic:   "user.created",
		GroupID: "user-service-debug",
	})

	go func() {
		for {
			msg, err := r.ReadMessage(context.Background())
			if err != nil {
				log.Println("kafka read error:", err)
				continue
			}

			log.Println("📨 Kafka message received:", string(msg.Value))
		}
	}()
}
