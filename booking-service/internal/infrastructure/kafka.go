package infrastructure

import (
	"booking-service/internal/config"
	"booking-service/internal/dto"
	"booking-service/internal/models"
	"context"
	"encoding/json"
	"fmt"
	"net"
	"os"
	"strconv"
	"time"

	"github.com/segmentio/kafka-go"
)

const (
	kafkaTopic = "bookings"
)

func getKafkaBroker() string {
	broker := os.Getenv("KAFKA_BROKER")
	if broker == "" {
		config.GetLogger().Warn("KAFKA_BROKER not set, using default localhost:9092")
		return "localhost:9092"
	}
	return broker
}

// kafkaWriter — объект для отправки сообщений в Kafka
// Объявляем глобально, чтобы использовать во всех функциях
var kafkaWriter *kafka.Writer

// createTopic создаёт топик в Kafka, если он ещё не существует
// Библиотека kafka-go не создаёт топики автоматически, поэтому делаем это вручную
func createTopic() error {
	kafkaBroker := getKafkaBroker()
	// Устанавливаем соединение с Kafka
	conn, err := kafka.Dial("tcp", kafkaBroker)
	if err != nil {
		config.GetLogger().Error("Failed to connect to Kafka", "error", err, "broker", kafkaBroker)
		return err
	}
	defer conn.Close()

	// Получаем информацию о контроллере кластера
	controller, err := conn.Controller()
	if err != nil {
		config.GetLogger().Error("Failed to get Kafka controller", "error", err)
		return err
	}

	// Подключаемся к контроллеру для создания топика
	controllerConn, err := kafka.Dial("tcp", net.JoinHostPort(controller.Host, strconv.Itoa(controller.Port)))
	if err != nil {
		config.GetLogger().Error("Failed to connect to Kafka controller", "error", err, "host", controller.Host, "port", controller.Port)
		return err
	}
	defer controllerConn.Close()

	// Проверяем, существует ли топик
	partitions, err := conn.ReadPartitions()
	if err == nil {
		for _, p := range partitions {
			if p.Topic == kafkaTopic {
				config.GetLogger().Info("Kafka topic already exists", "topic", kafkaTopic)
				return nil
			}
		}
	}

	// Конфигурация топика
	topicConfig := kafka.TopicConfig{
		Topic:             kafkaTopic,
		NumPartitions:     1,
		ReplicationFactor: 1,
	}

	// Создаём топик
	err = controllerConn.CreateTopics(topicConfig)
	if err != nil {
		// Игнорируем ошибку, если топик уже существует
		if err.Error() == "topic already exists" || 
		   err.Error() == "TopicExistsException" ||
		   err.Error() == "TopicExistsException: Topic 'bookings' already exists" {
			config.GetLogger().Info("Kafka topic already exists", "topic", kafkaTopic)
			return nil
		}
		config.GetLogger().Error("Failed to create Kafka topic", "topic", kafkaTopic, "error", err)
		return err
	}

	config.GetLogger().Info("Kafka topic created successfully", "topic", kafkaTopic)
	return nil
}

// initKafkaWriter создаёт и настраивает Kafka Writer
func InitKafkaWriter() {
	kafkaBroker := getKafkaBroker()
	
	// Сначала создаём топик, если его нет
	if err := createTopic(); err != nil {
		config.GetLogger().Error("Failed to create Kafka topic, continuing anyway", "error", err)
		// Не создаём writer, если не удалось создать топик
		return
	}

	kafkaWriter = &kafka.Writer{
		Addr:     kafka.TCP(kafkaBroker),
		Topic:    kafkaTopic,
		Balancer: &kafka.LeastBytes{},
		WriteTimeout: 10 * time.Second,
		RequiredAcks: 1,
	}
	config.GetLogger().Info("Kafka writer initialized successfully", "topic", kafkaTopic, "broker", kafkaBroker)
}

// publishOrderCreated отправляет событие о создании заказа в Kafka
func PublishOrderCreated(booking models.Booking) error {
	if kafkaWriter == nil {
		config.GetLogger().Error("Kafka writer is not initialized")
		return fmt.Errorf("kafka writer is not initialized")
	}

	event := dto.BookingConfirmResponse{
		SessionID:     booking.SessionID,
		UserID:        booking.UserID,
		BookingStatus: &booking.BookingStatus,
	}

	eventJSON, err := json.Marshal(event)
	if err != nil {
		config.GetLogger().Error("Failed to marshal event", "error", err)
		return err
	}

	msg := kafka.Message{
		Key:   []byte(fmt.Sprintf("booking-%d", booking.ID)),
		Value: eventJSON,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	err = kafkaWriter.WriteMessages(ctx, msg)
	if err != nil {
		config.GetLogger().Error("Failed to publish event to Kafka", "error", err, "booking_id", booking.ID)
		return err
	}

	config.GetLogger().Info("Event published to Kafka", "booking_id", booking.ID, "topic", kafkaTopic)
	return nil
}
