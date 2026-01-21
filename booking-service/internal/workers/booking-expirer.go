package workers

import (
	"booking-service/internal/config"
	"booking-service/internal/services"
	"time"
)

func StartExpiredBookingsWorker(bookingService services.BookingService) {
	ticker := time.NewTicker(30 * time.Second)

	logger := config.GetLogger()
	logger.Info("Expired bookings worker started", "interval", "30 second")

	if err := bookingService.ExpireOldBookings(); err != nil {
		logger.Error("Failed to expire old bookings on startup", "error", err)
	}

	for range ticker.C {
		if err := bookingService.ExpireOldBookings(); err != nil {
			logger.Error("Failed to expire old bookings", "error", err)
		}
	}
}

func StartEndedSessionsWorker(bookingService services.BookingService) {
	ticker := time.NewTicker(30 * time.Second)

	logger := config.GetLogger()
	logger.Info("Ended sessions worker started", "interval", "30 second")

	if err := bookingService.FreeSeatsForEndedSessions(); err != nil {
		logger.Error("Failed to free seats for ended sessions on startup", "error", err)
	}

	for range ticker.C {
		if err := bookingService.FreeSeatsForEndedSessions(); err != nil {
			logger.Error("Failed to free seats for ended sessions", "error", err)
		}
	}
}
