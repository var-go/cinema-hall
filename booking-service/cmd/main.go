package main

import (
	"booking-service/internal/config"
	"booking-service/internal/infrastructure"
	"booking-service/internal/models"
	"booking-service/internal/repository"
	"booking-service/internal/services"
	"booking-service/internal/transport"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/gofiber/fiber/v2/log"
)

func main() {
	db := config.Connect()

	router := gin.Default()

	if db == nil {
		log.Error("database is nil")
		return
	}

	if err := db.AutoMigrate(&models.Booking{}); err != nil {
		log.Error("failed to migrate database", err)
		os.Exit(1)
	}

	infrastructure.InitKafkaWriter()

	bookingRepo := repository.NewBookingRepository(db)

	bookingService := services.NewBookingService(bookingRepo)

	transport.RegisterRoutes(router, bookingService)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8081"
	}

	if err := router.Run(":" + port); err != nil {
		log.Error("ошибка запуска сервера")
	}
}
