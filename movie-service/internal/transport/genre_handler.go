package transport

import (
	"errors"
	"log/slog"
	"movie-service/internal/dto"
	"movie-service/internal/services"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type GenreHandler struct {
	service services.GenreService
	logger  *slog.Logger
}

func NewGenreHandler(service services.GenreService, logger *slog.Logger) *GenreHandler {
	return &GenreHandler{
		service: service,
		logger:  logger,
	}
}

func (h *GenreHandler) RegisterRoutes(ctx *gin.Engine) {
	api := ctx.Group("/genres")
	{
		api.POST("/", h.Create)
		api.GET("/", h.List)
		api.GET("/:id", h.GetByID)
		api.PUT("/:id", h.Update)
		api.DELETE("/:id", h.Delete)

	}
}

func (h *GenreHandler) Create(ctx *gin.Context) {

	var req dto.GenreCreateRequest

	if err := ctx.ShouldBindJSON(&req); err != nil {
		h.logger.Error("invalid create genre request", slog.Any("error", err))
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	genre, err := h.service.Create(&req)

	if err != nil {
		h.logger.Error("genre create failed", slog.Any("error", err), slog.String("name", req.Name))
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "genre create error"})
		return
	}

	ctx.JSON(http.StatusCreated, genre)
}

func (h *GenreHandler) List(ctx *gin.Context) {

	genres, err := h.service.List()

	if err != nil {
		h.logger.Error("genre list handler failed", slog.Any("error", err))
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "genre list error"})
		return
	}
	ctx.JSON(http.StatusOK, genres)
}

func (h *GenreHandler) GetByID(ctx *gin.Context) {

	id, err := strconv.ParseUint(ctx.Param("id"), 10, 64)

	if err != nil {
		h.logger.Error("invalid genre id param", slog.String("param", ctx.Param("id")), slog.Any("error", err))
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid genre id"})
		return
	}

	genre, err := h.service.GetByID(uint(id))
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			h.logger.Info("genre not found", slog.Any("id", id))
			ctx.JSON(http.StatusNotFound, gin.H{"error": "genre not found"})
			return
		}
		h.logger.Error("failed to get genre", slog.Any("id", id), slog.Any("error", err))
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get genre"})
		return
	}

	ctx.JSON(http.StatusOK, genre)
}

func (h *GenreHandler) Update(ctx *gin.Context) {

	id, err := strconv.ParseUint(ctx.Param("id"), 10, 64)

	if err != nil {
		h.logger.Error("invalid genre id param", slog.String("param", ctx.Param("id")), slog.Any("error", err))
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid genre id"})
		return
	}

	var req dto.GenreUpdateRequest

	if err := ctx.ShouldBindJSON(&req); err != nil {
		h.logger.Error("invalid update genre request", slog.Any("error", err), slog.Any("id", id))
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	genre, err := h.service.Update(uint(id), &req)

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			h.logger.Info("genre not found for update", slog.Any("id", id))
			ctx.JSON(http.StatusNotFound, gin.H{"error": "genre not found"})
			return
		}
		h.logger.Error("genre update failed", slog.Any("id", id), slog.Any("error", err))
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, genre)
}

func (h *GenreHandler) Delete(ctx *gin.Context) {

	id, err := strconv.ParseUint(ctx.Param("id"), 10, 64)

	if err != nil {
		h.logger.Error("invalid genre id param", slog.String("param", ctx.Param("id")), slog.Any("error", err))
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid genre id"})
		return
	}

	if err := h.service.Delete(uint(id)); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			h.logger.Info("genre not found for delete", slog.Any("id", id))
			ctx.JSON(http.StatusNotFound, gin.H{"error": "genre not found"})
			return
		}
		h.logger.Error("failed to delete genre", slog.Any("id", id), slog.Any("error", err))
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to delete genre"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "genre deleted successfully"})
}
