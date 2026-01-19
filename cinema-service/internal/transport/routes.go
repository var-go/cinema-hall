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

) {

	hallHandler := NewHallHandler(hallService, logger)
	seatHandler := NewSeatHandler(seatService, hallService, logger)

	hallHandler.RegisterRoutes(router)
	seatHandler.RegisterRoutes(router)
}
