package config

import (
	"log/slog"
	"os"
	"strings"
)

var (
	// Logger глобальный экземпляр логгера для использования во всем сервисе
	Logger *slog.Logger
)

// InitLogger инициализирует структурированный логгер для всего сервиса
func InitLogger() *slog.Logger {
	// Определяем уровень логирования из переменной окружения
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
		// Если не указан, используем info по умолчанию
		level = slog.LevelInfo
	}

	// JSON формат для структурированного логирования
	handler := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level:     level,
		AddSource: false, // Можно включить для отладки
	})

	logger := slog.New(handler)

	// Устанавливаем контекст по умолчанию для всех логов
	logger = logger.With(
		"service", "booking-service",
	)

	// Сохраняем в глобальную переменную для удобства доступа
	Logger = logger

	return logger
}

// GetLogger возвращает глобальный экземпляр логгера
// Если логгер не инициализирован, создаёт новый
func GetLogger() *slog.Logger {
	if Logger == nil {
		return InitLogger()
	}
	return Logger
}

