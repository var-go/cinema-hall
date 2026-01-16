package main

import (
	"log/slog"
	"movie-service/internal/config"
	"movie-service/internal/models"
	"movie-service/internal/repository"
	"movie-service/internal/services"
	"movie-service/internal/transport"
	"os"

	"github.com/gin-gonic/gin"
)

func main() {
	logger := config.InitLogger()

	r := gin.Default()

	db, err := config.SetUpDatabaseConnection(logger)
	if err != nil {
		logger.Error("failed to set up database", slog.Any("error", err))
		os.Exit(1)
	}

	if err := db.AutoMigrate(&models.Movie{}, &models.Genre{}); err != nil {
		logger.Error("failed to migrate database", "error", err)
		os.Exit(1)
	}

	logger.Info("migrations completed")

	movieRepo := repository.NewMovieRepository(db, logger)
	genreRepo := repository.NewGenreRepository(db, logger)

	movieService := services.NewMovieService(movieRepo, logger)
	genreService := services.NewGenreService(genreRepo, logger)

	transport.RegisterRoutes(r, movieService, genreService, logger)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8083"
	}

	logger.Info("application started successfully")

	if err := r.Run(":" + port); err != nil {
		logger.Error("ошибка запуска сервера", slog.Any("error", err))
	}

}
