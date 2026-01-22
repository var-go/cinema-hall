package kafka

import (
	"context"
	"log/slog"

	kafka "github.com/segmentio/kafka-go"
)

func StartUserCreatedConsumer(ctx context.Context, broker string, log *slog.Logger) {
	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers: []string{broker},
		Topic:   "user.created",
		GroupID: "user-service-debug",
	})

	go func() {
		defer r.Close()

		log.Info("kafka consumer started")

		for {
			msg, err := r.ReadMessage(ctx)
			if err != nil {
				if ctx.Err() != nil {
					log.Info("kafka consumer stopped")
					return
				}
				log.Error("failed to read kafka message", "err:", err)
				continue
			}

			log.Info("user created event received", "value", string(msg.Value))
		}
	}()
}
