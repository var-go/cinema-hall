package main

import (
	"booking-service/internal/config"
	"booking-service/internal/models"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/gofiber/fiber/v2/log"
)

func main() {
	db := config.Connect()

	router := gin.Default()

	if db == nil {
		log.Error("database is nill")
		return
	}

	if err := db.AutoMigrate(&models.Booking{}); err != nil {
		log.Error("failed to migrate database", err)
		os.Exit(1)
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8081"
	}

	if err := router.Run(":" + port); err != nil {
		log.Error("ошибка запуска сервера")
	}
}
