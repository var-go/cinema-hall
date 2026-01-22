package main

import (
	"log"
	"os"
	"user-service/internal/config"
	"user-service/internal/kafka"
	"user-service/internal/models"
	"user-service/internal/repository"
	"user-service/internal/services"
	"user-service/internal/transport"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {

	_ = godotenv.Load()

	db := config.SetupDatabase()

	if err := db.AutoMigrate(&models.User{}); err != nil {
		log.Fatal(err)
	}

	userRepo := repository.NewUserRepository(db)

	broker := os.Getenv("KAFKA_BROKER")
	if broker == "" {
		log.Fatal("KAFKA_BROKER is not set")
	}

	producer := kafka.NewProducer(broker)

	authService := services.NewAuthService(userRepo, producer)

	userService := services.NewUserService(userRepo)

	authHandler := transport.NewAuthHandler(authService)
	userHandler := transport.NewUserHandler(userService)

	r := gin.Default()
	transport.RegisterRouters(r, authHandler, userHandler)

	if err := r.Run(":8080"); err != nil {
		log.Fatal(err)
	}

}
