package transport

import (
	"cinema-service/internal/dto"
	"cinema-service/internal/services"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type SeatHandler struct {
	seatService services.SeatService
	hallService services.HallService
	logger      *slog.Logger
}

func NewSeatHandler(seatService services.SeatService, hallService services.HallService, logger *slog.Logger) *SeatHandler {
	return &SeatHandler{
		seatService: seatService,
		hallService: hallService,
		logger:      logger,
	}
}

func (h *SeatHandler) RegisterRoutes(r *gin.Engine) {
	seats := r.Group("/")
	{
		seats.POST("/halls/:id/seats", h.Create)
		seats.GET("/seats", h.GetAllSeats)
		seats.PATCH("/seats/:id", h.Patch)
		seats.DELETE("/seats/:id", h.RemoveSeat)
	}
}

func (h *SeatHandler) Create(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	var req dto.CreateSeatRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Error("handler: failed to bind JSON", "err", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	seat, err := h.seatService.Create(uint(id), req)
	if err != nil {
		h.logger.Error("failed to fetch seats")
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, seat)
}

func (h *SeatHandler) GetAllSeats(c *gin.Context) {
	seats, err := h.seatService.List()
	if err != nil {
		h.logger.Error("failed to fetch seats", "err", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "failed to fetch seats"})
		return
	}
	c.JSON(http.StatusOK, seats)
}

func (h *SeatHandler) Patch(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	var req dto.UpdateSeatRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Error("handler: failed to bind JSON", "err", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	seat, err := h.seatService.UpdateSeat(uint(id), req)
	if err != nil {
		h.logger.Error("failed to fetch halls")
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, seat)
}

func (h *SeatHandler) RemoveSeat(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		h.logger.Error("handler: invalid hall id", "id", id)
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	if err := h.seatService.Delete(uint(id)); err != nil {
		h.logger.Error("handler: failed to delete seat", "id", id)
		c.JSON(http.StatusBadRequest, gin.H{"error": "failed to delete seat"})
		return
	}
	h.logger.Info("handler: seat deleted successfully", "id", id)
	c.JSON(http.StatusOK, gin.H{"message": "seat deleted successfully"})
}
