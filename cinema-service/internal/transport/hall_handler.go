package transport

import (
	"cinema-service/internal/dto"
	"cinema-service/internal/services"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type HallHandler struct {
	hallService services.HallService
	logger      *slog.Logger
}

func NewHallHandler(hallService services.HallService, logger *slog.Logger) *HallHandler {
	return &HallHandler{
		hallService: hallService,
		logger:      logger,
	}
}

func (h *HallHandler) RegisterRoutes(r *gin.Engine) {
	halls := r.Group("/halls")
	{
		halls.POST("", h.Create)
		halls.DELETE("/:id", h.RemoveHall)
		halls.PATCH("/:id", h.Patch)
		halls.GET("/:id", h.GetById)
		halls.GET("", h.GetAllHalls)
	}
}

func (h *HallHandler) Create(c *gin.Context) {
	var req dto.CreateHallRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Error("handler: failed to bind JSON", "err", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	hall, err := h.hallService.CreateHall(req)
	if err != nil {
		h.logger.Error("handler: failed to create hall", "err", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, hall)
}

func (h *HallHandler) Patch(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	var req dto.UpdateHallRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Error("handler: failed to bind JSON", "err", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	hall, err := h.hallService.UpdateHall(uint(id), req)
	if err != nil {
		h.logger.Error("failed to fetch halls")
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, hall)
}

func (h *HallHandler) GetById(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	hall, err := h.hallService.GetHallByID(uint(id))
	if err != nil {
		h.logger.Error("failed to fetch halls")
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, hall)
}

func (h *HallHandler) GetAllHalls(c *gin.Context) {
	halls, err := h.hallService.ListHall()
	if err != nil {
		h.logger.Error("failed to fetch halls", "err", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "failed to fetch halls"})
		return
	}
	c.JSON(http.StatusOK, halls)
}

func (h *HallHandler) RemoveHall(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		h.logger.Error("handler: invalid hall id", "id", id)
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	if err := h.hallService.DeleteHall(uint(id)); err != nil {
		h.logger.Error("handler: failed to delete hall", "id", id)
		c.JSON(http.StatusBadRequest, gin.H{"error": "failed to delete hall"})
		return
	}
	h.logger.Info("handler: hall deleted successfully", "id", id)
	c.JSON(http.StatusOK, gin.H{"message": "hall deleted successfully"})
}
