package config

import (
	"log/slog"
	"os"
	"strings"
)

var (
	Logger *slog.Logger
)

func InitLogger() *slog.Logger {
	level := slog.LevelInfo

	logLevel := strings.ToLower(os.Getenv("LOG_LEVEL"))
	switch logLevel {
	case "debug":
		level = slog.LevelDebug
	case "warn":
		level = slog.LevelWarn
	case "info":
		level = slog.LevelInfo
	case "error":
		level = slog.LevelError
	default:
		level = slog.LevelInfo
	}

	handler := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level:     level,
		AddSource: false,
	})

	logger := slog.New(handler)

	logger = logger.With(
		"service", "booking-service",
	)

	Logger = logger

	return logger
}

func GetLogger() *slog.Logger {
	if Logger == nil {
		return InitLogger()
	}
	return Logger
}

