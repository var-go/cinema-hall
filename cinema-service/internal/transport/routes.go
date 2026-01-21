package transport

import (
	"cinema-service/internal/services"
	"log/slog"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(
	router *gin.Engine,
	logger *slog.Logger,
	hallService services.HallService,
	seatService services.SeatService,
	sessionsService services.SessionService,

) {

	hallHandler := NewHallHandler(hallService, logger)
	seatHandler := NewSeatHandler(seatService, logger)
	sessionHandler := NewSessionHandler(sessionsService, logger)

	hallHandler.RegisterRoutes(router)
	seatHandler.RegisterRoutes(router)
	sessionHandler.RegisterRoutes(router)
}
