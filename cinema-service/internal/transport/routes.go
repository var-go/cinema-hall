package transport

import (
	"cinema-service/internal/services"
	"log/slog"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(
	router *gin.Engine,
	logger *slog.Logger,
	service services.HallService,
) {

	hallHandler := NewHallHandler(service, logger)

	hallHandler.RegisterRoutes(router)
}
