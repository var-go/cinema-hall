package main

import (
	"cinema-service/internal/config"
	"cinema-service/internal/models"
	"cinema-service/internal/repository"
	"cinema-service/internal/services"
	"cinema-service/internal/transport"
	"log/slog"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/gofiber/fiber/v2/log"
)

func main() {
	db := config.Connect()
	logger := config.InitLogger()

	if err := db.AutoMigrate(
		&models.Hall{},
	); err != nil {
		log.Error("failed to migrate database", "error", err)
		os.Exit(1)
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8081"
	}

	r := gin.Default()

	hallRepo := repository.NewHallRepository(db, logger)

	hallService := services.NewHallService(hallRepo, logger)

	transport.RegisterRoutes(r, logger, hallService)

	if err := r.Run(":" + port); err != nil {
		log.Error("не удалось запустить сервер", slog.Any("error", err))
		os.Exit(1)
	}

}
