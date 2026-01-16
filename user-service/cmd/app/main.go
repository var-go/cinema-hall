package main

import (
	"log"
	"user-service/internal/config"
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

	authService := services.NewAuthService(userRepo)
	userService := services.NewUserService(userRepo)

	authHandler := transport.NewAuthHandler(authService)
	userHandler := transport.NewUserHandler(userService)

	r := gin.Default()
	transport.RegisterRouters(r, authHandler, userHandler)

}
