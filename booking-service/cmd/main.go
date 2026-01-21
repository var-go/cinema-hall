package main

import (
	"booking-service/internal/config"
	"booking-service/internal/infrastructure"
	"booking-service/internal/models"
	"booking-service/internal/repository"
	"booking-service/internal/services"
	"booking-service/internal/transport"
	"booking-service/internal/workers"
	"os"

	"github.com/gin-gonic/gin"
)

func main() {
	// Инициализация логгера
	logger := config.InitLogger()
	logger.Info("Starting booking-service")

	// Подключение к базе данных
	db := config.Connect()

	router := gin.Default()

	if db == nil {
		logger.Error("Database connection failed: database is nil")
		return
	}

	logger.Info("Database connected successfully")

	// Автомиграция
	if err := db.AutoMigrate(&models.Booking{}, &models.BookedSeat{}); err != nil {
		logger.Error("Failed to migrate database", "error", err)
		os.Exit(1)
	}

	logger.Info("Database migration completed")

	// Инициализация Kafka
	infrastructure.InitKafkaWriter()
	logger.Info("Kafka writer initialized")

	// Инициализация репозиториев и сервисов
	bookingRepo := repository.NewBookingRepository(db)
	bookingSeatRepo := repository.NewBookingSeatRepository(db)
	bookingService := services.NewBookingService(bookingRepo, bookingSeatRepo, db)

	go workers.StartExpiredBookingsWorker(bookingService)
	go workers.StartEndedSessionsWorker(bookingService)

	transport.RegisterRoutes(router, bookingService)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8081"
	}

	logger.Info("Server starting", "port", port)

	if err := router.Run(":" + port); err != nil {
		logger.Error("Failed to start server", "error", err)
		os.Exit(1)
	}
}
