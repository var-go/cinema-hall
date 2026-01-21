package infrastructure

import (
	"booking-service/internal/config"
	"booking-service/internal/dto"
	"booking-service/internal/models"
	"context"
	"encoding/json"
	"fmt"
	"net"
	"strconv"

	"github.com/segmentio/kafka-go"
)

const (
	kafkaBroker = "localhost:9092"
	kafkaTopic  = "orders"
)

// kafkaWriter — объект для отправки сообщений в Kafka
// Объявляем глобально, чтобы использовать во всех функциях
var kafkaWriter *kafka.Writer

// createTopic создаёт топик в Kafka, если он ещё не существует
// Библиотека kafka-go не создаёт топики автоматически, поэтому делаем это вручную
func createTopic() {
	// Устанавливаем соединение с Kafka
	conn, err := kafka.Dial("tcp", kafkaBroker)
	if err != nil {
		config.GetLogger().Error("Failed to connect to Kafka", "error", err, "broker", kafkaBroker)
		return
	}
	defer conn.Close()

	// Получаем информацию о контроллере кластера
	controller, err := conn.Controller()
	if err != nil {
		config.GetLogger().Error("Failed to get Kafka controller", "error", err)
		return
	}

	// Подключаемся к контроллеру для создания топика
	controllerConn, err := kafka.Dial("tcp", net.JoinHostPort(controller.Host, strconv.Itoa(controller.Port)))
	if err != nil {
		config.GetLogger().Error("Failed to connect to Kafka controller", "error", err, "host", controller.Host, "port", controller.Port)
		return
	}
	defer controllerConn.Close()

	// Конфигурация топика
	topicConfig := kafka.TopicConfig{
		Topic:             kafkaTopic,
		NumPartitions:     1,
		ReplicationFactor: 1,
	}

	// Создаём топик (если уже существует — ошибка будет проигнорирована)
	err = controllerConn.CreateTopics(topicConfig)
	if err != nil {
		config.GetLogger().Warn("Kafka topic already exists or creation failed", "topic", kafkaTopic, "error", err)
	} else {
		config.GetLogger().Info("Kafka topic created successfully", "topic", kafkaTopic)
	}
}

// initKafkaWriter создаёт и настраивает Kafka Writer
func InitKafkaWriter() {
	// Сначала создаём топик, если его нет
	createTopic()

	kafkaWriter = &kafka.Writer{
		// Адрес Kafka брокера (localhost, порт 9092)
		Addr: kafka.TCP(kafkaBroker),
		// Имя топика, в который будем отправлять сообщения
		Topic: kafkaTopic,
		// Балансировщик определяет, в какую партицию отправить сообщение
		Balancer: &kafka.LeastBytes{},
	}
	config.GetLogger().Info("Kafka writer initialized", "topic", kafkaTopic, "broker", kafkaBroker)
}

// publishOrderCreated отправляет событие о создании заказа в Kafka
func PublishOrderCreated(booking models.Booking) error {
	// Создаём событие с нужными полями
	// Не отправляем весь заказ — только то, что нужно для уведомления
	event := dto.BookingCreateRequest{
		SessionID: booking.SessionID,
		UserID:    booking.UserID,
	}

	// Преобразуем структуру в JSON (массив байтов)
	eventJSON, err := json.Marshal(event)
	if err != nil {
		return err // Если не удалось сериализовать — возвращаем ошибку
	}

	// Создаём сообщение для Kafka
	msg := kafka.Message{
		// Key — ключ сообщения, используется для группировки
		// Сообщения с одинаковым ключом попадают в одну партицию
		Key: []byte(fmt.Sprintf("booking-%d", booking.ID)),
		// Value — само содержимое сообщения (наш JSON)
		Value: eventJSON,
	}

	// Отправляем сообщение в Kafka
	// context.Background() — пустой контекст без таймаута
	err = kafkaWriter.WriteMessages(context.Background(), msg)
	if err != nil {
		return err // Если не удалось отправить — возвращаем ошибку
	}

	config.GetLogger().Info("Event published to Kafka", "booking_id", booking.ID, "topic", kafkaTopic)
	return nil // Всё прошло успешно
}
