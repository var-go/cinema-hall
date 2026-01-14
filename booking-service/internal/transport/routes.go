package transport

import (
	"booking-service/internal/services"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(router *gin.Engine, bookingService services.BookingService) {
	bookingHandler := NewBookingHandler(bookingService)

	bookingHandler.BookingRoutes(router)
}
