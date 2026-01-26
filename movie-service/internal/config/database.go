package config

import (
	"fmt"
	"log/slog"
	"os"

	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func SetUpDatabaseConnection(logger *slog.Logger) (*gorm.DB, error) {
	if err := godotenv.Load(); err != nil {
		logger.Info("No .env file found, using environment variables", "error", err)
	}

	dbUrl := os.Getenv("DATABASE_URL")
	if dbUrl == "" {
		logger.Error("DATABASE_URL environment variable is not set")
		return nil, fmt.Errorf("DATABASE_URL is required")
	}

	db, err := gorm.Open(postgres.New(postgres.Config{
		DSN:                  dbUrl,
		PreferSimpleProtocol: true,
	}), &gorm.Config{})

	if err != nil {
		logger.Error("Failed to initialize database", "error", err)
		return nil, err
	}

	logger.Info("Successfully connected to the database")
	return db, nil
}
