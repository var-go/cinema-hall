package transport

import (
	"booking-service/internal/dto"
	"booking-service/internal/infrastructure"
	"booking-service/internal/services"
	"log"
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
	api := ctx.Group("/booking")
	{
		api.POST("", h.Create)
		api.GET("", h.List)
		api.GET("/:id", h.GetByID)
		api.PATCH("/:id", h.Update)
		api.DELETE("/:id", h.Delete)
	}
}

func (h *bookingTransport) Create(ctx *gin.Context) {
	var req dto.BookingCreateRequest

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid JSON"})
		return
	}

	booking, err := h.service.Create(req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}

	if err := infrastructure.PublishOrderCreated(*booking); err != nil {
		// Если не удалось отправить — логируем, но не отменяем заказ
		// Заказ уже создан, клиент получит успешный ответ
		log.Printf("Ошибка отправки в Kafka: %v", err)
	}

	ctx.JSON(http.StatusOK, booking)
}

func (h *bookingTransport) List(ctx *gin.Context) {
	list, err := h.service.List()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}

	ctx.JSON(http.StatusOK, list)
}

func (h *bookingTransport) GetByID(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	booking, err := h.service.GetByID(uint(id))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}

	ctx.JSON(http.StatusOK, booking)
}

func (h *bookingTransport) Update(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	var req dto.BookingUpdateRequest

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid JSON"})
		return
	}

	booking, err := h.service.Update(uint(id), req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}

	ctx.JSON(http.StatusOK, booking)
}

func (h *bookingTransport) Delete(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	if err := h.service.Delete(uint(id)); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "booking deleted"})
}
