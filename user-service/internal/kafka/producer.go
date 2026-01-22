package kafka

import (
	"context"
	"encoding/json"

	"github.com/segmentio/kafka-go"
)

type UserCreatedEvent struct {
	ID    uint   `json:"id"`
	Email string `json:"email"`
}

type Producer struct {
	writer *kafka.Writer
}

func NewProducer(broker string) *Producer {
	return &Producer{
		writer: &kafka.Writer{
			Addr:  kafka.TCP(broker),
			Topic: "user.created",
		},
	}
}

func (p *Producer) SendUserCreated(event UserCreatedEvent) error {
	data, err := json.Marshal(event)
	if err != nil {
		return err
	}

	return p.writer.WriteMessages(
		context.Background(),
		kafka.Message{
			Value: data,
		},
	)
}
