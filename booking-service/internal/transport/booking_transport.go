package transport

import (
	"booking-service/internal/config"
	"booking-service/internal/constants"
	"booking-service/internal/dto"
	"booking-service/internal/infrastructure"
	"booking-service/internal/services"
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type bookingTransport struct {
	service services.BookingService
}

func NewBookingHandler(service services.BookingService) *bookingTransport {
	return &bookingTransport{
		service: service,
	}
}

func (h *bookingTransport) BookingRoutes(ctx *gin.Engine) {
	api := ctx.Group("/bookings")
	{
		api.POST("", h.Create)
		api.GET("", h.List)
		api.GET("/:id", h.GetByID)
		api.PATCH("/:id", h.Update)
		api.DELETE("/:id", h.Delete)
		api.POST("/:id/confirm", h.ConfirmBooking)
		api.POST("/:id/cancel", h.CancelBooking)
	}
}

func (h *bookingTransport) Create(ctx *gin.Context) {
	var req dto.BookingCreateRequest

	if err := ctx.ShouldBindJSON(&req); err != nil {
		config.GetLogger().Warn("Invalid JSON in booking request", "error", err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid JSON"})
		return
	}

	config.GetLogger().Info("Creating booking", "session_id", req.SessionID, "user_id", req.UserID, "seats", req.SeatsID)

	booking, err := h.service.Create(req)
	if err != nil {
		config.GetLogger().Error("Failed to create booking", "error", err, "session_id", req.SessionID, "user_id", req.UserID, "seats", req.SeatsID)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	config.GetLogger().Info("Booking created successfully", "booking_id", booking.ID, "session_id", booking.SessionID, "user_id", booking.UserID)

	ctx.JSON(http.StatusOK, booking)
}

func (h *bookingTransport) List(ctx *gin.Context) {
	list, err := h.service.List()
	if err != nil {
		config.GetLogger().Error("Failed to list bookings", "error", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, list)
}

func (h *bookingTransport) GetByID(ctx *gin.Context) {
	id, err := parseID(ctx.Param("id"))
	if err != nil {
		config.GetLogger().Warn("Invalid booking ID in request", "error", err, "id_param", ctx.Param("id"))
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	booking, err := h.service.GetByID(uint(id))
	if err != nil {
		if errors.Is(err, constants.ErrBookingNotFound) {
			ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		config.GetLogger().Error("Failed to get booking by ID", "error", err, "booking_id", id)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, booking)
}

func (h *bookingTransport) Update(ctx *gin.Context) {
	id, err := parseID(ctx.Param("id"))
	if err != nil {
		config.GetLogger().Warn("Invalid booking ID in update request", "error", err, "id_param", ctx.Param("id"))
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	var req dto.BookingUpdateRequest

	if err := ctx.ShouldBindJSON(&req); err != nil {
		config.GetLogger().Warn("Invalid JSON in update request", "error", err, "booking_id", id)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid JSON"})
		return
	}

	booking, err := h.service.Update(uint(id), req)
	if err != nil {
		if errors.Is(err, constants.ErrBookingNotFound) {
			ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		config.GetLogger().Error("Failed to update booking", "error", err, "booking_id", id)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, booking)
}

func (h *bookingTransport) Delete(ctx *gin.Context) {
	id, err := parseID(ctx.Param("id"))
	if err != nil {
		config.GetLogger().Warn("Invalid booking ID in delete request", "error", err, "id_param", ctx.Param("id"))
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	if err := h.service.Delete(uint(id)); err != nil {
		if errors.Is(err, constants.ErrBookingNotFound) {
			ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		config.GetLogger().Error("Failed to delete booking", "error", err, "booking_id", id)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	config.GetLogger().Info("Booking deleted successfully", "booking_id", id)
	ctx.JSON(http.StatusOK, gin.H{"message": "booking deleted"})
}

func (h *bookingTransport) ConfirmBooking(ctx *gin.Context) {
	id, err := parseID(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	confirmed, err := h.service.ConfirmBooking(uint(id))
	if err != nil {
		switch {

		case errors.Is(err, constants.ErrBookingAlreadyCancelled),
			errors.Is(err, constants.ErrBookingAlreadyConfirmed),
			errors.Is(err, constants.ErrBookingExpired):
			ctx.JSON(http.StatusConflict, gin.H{"error": err.Error()})
			return

		case errors.Is(err, constants.ErrBookingNotFound):
			ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return

		default:
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
	}

	if err := infrastructure.PublishOrderCreated(*confirmed); err != nil {
		config.GetLogger().Error("Failed to publish event to Kafka",
			"error", err,
			"booking_id", confirmed.ID,
			"session_id", confirmed.SessionID)
	} else {
		config.GetLogger().Info("Successfully published booking confirm event to Kafka",
			"booking_id", confirmed.ID)
	}

	ctx.JSON(http.StatusOK, confirmed)
}

func (h *bookingTransport) CancelBooking(ctx *gin.Context) {
	id, err := parseID(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	cancelled, err := h.service.CancelBooking(uint(id))
	if err != nil {
		switch {

		case errors.Is(err, constants.ErrBookingAlreadyCancelled),
			errors.Is(err, constants.ErrBookingExpired):
			ctx.JSON(http.StatusConflict, gin.H{"error": err.Error()})
			return

		case errors.Is(err, constants.ErrBookingNotFound):
			ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return

		default:
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
	}

	if err := infrastructure.PublishOrderCreated(*cancelled); err != nil {
		config.GetLogger().Error("Failed to publish cancel event to Kafka",
			"error", err,
			"booking_id", cancelled.ID,
			"session_id", cancelled.SessionID)
	} else {
		config.GetLogger().Info("Successfully published booking cancel event to Kafka",
			"booking_id", cancelled.ID,
			"session_id", cancelled.SessionID,
			"user_id", cancelled.UserID,
			"status", cancelled.BookingStatus)
	}

	ctx.JSON(http.StatusOK, cancelled)
}

func parseID(idStr string) (uint, error) {
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		return 0, constants.ErrInvalidID
	}
	return uint(id), nil
}
