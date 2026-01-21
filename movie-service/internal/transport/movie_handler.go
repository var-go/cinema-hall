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

type MovieHandler struct {
	service services.MovieService
	logger  *slog.Logger
}

func NewMovieHandler(service services.MovieService, logger *slog.Logger) *MovieHandler {
	return &MovieHandler{
		service: service,
		logger:  logger,
	}
}

func (h *MovieHandler) RegisterRoutes(ctx *gin.Engine) {
	api := ctx.Group("/movies")
	{
		api.POST("/", h.Create)
		api.GET("/", h.List)
		api.GET("/now-showing", h.NowShowing)
		api.GET("/coming-soon", h.ComingSoon)
		api.GET("/:id", h.GetByID)
		api.PUT("/:id", h.Update)
		api.DELETE("/:id", h.Delete)

	}
}

func (h *MovieHandler) Create(ctx *gin.Context) {

	var req dto.MovieCreateRequest

	if err := ctx.ShouldBindJSON(&req); err != nil {
		h.logger.Error("invalid create movie request", slog.Any("error", err))
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	movie, err := h.service.Create(&req)

	if err != nil {
		h.logger.Error("movie create failed", slog.Any("error", err), slog.String("title", req.Title))
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "movie create error"})
		return
	}

	ctx.JSON(http.StatusCreated, movie)
}

func (h *MovieHandler) List(ctx *gin.Context) {

	movies, err := h.service.List()

	if err != nil {
		h.logger.Error("movie list handler failed", slog.Any("error", err))
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "movie list error"})
		return
	}
	ctx.JSON(http.StatusOK, movies)
}

func (h *MovieHandler) GetByID(ctx *gin.Context) {

	id, err := strconv.ParseUint(ctx.Param("id"), 10, 64)

	if err != nil {
		h.logger.Error("invalid movie id param", slog.String("param", ctx.Param("id")), slog.Any("error", err))
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid movie id"})
		return
	}

	movie, err := h.service.GetByID(uint(id))

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			h.logger.Info("movie not found", slog.Any("id", id))
			ctx.JSON(http.StatusNotFound, gin.H{"error": "movie not found"})
			return
		}

		h.logger.Error("failed to get movie", slog.Any("id", id), slog.Any("error", err))
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get movie"})
		return
	}

	ctx.JSON(http.StatusOK, movie)
}

func (h *MovieHandler) NowShowing(ctx *gin.Context) {

	movies, err := h.service.GetNowShowing()

	if err != nil {

		h.logger.Error("failed to get now showing movies", slog.Any("error", err))
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get now showing movies"})
		return
	}

	ctx.JSON(http.StatusOK, movies)

}

func (h *MovieHandler) ComingSoon(ctx *gin.Context) {

	movies, err := h.service.GetComingSoon()

	if err != nil {
		h.logger.Error("failed to get coming soon movies", slog.Any("error", err))
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get coming soon movies"})
		return
	}

	ctx.JSON(http.StatusOK, movies)

}

func (h *MovieHandler) Update(ctx *gin.Context) {

	id, err := strconv.ParseUint(ctx.Param("id"), 10, 64)

	if err != nil {
		h.logger.Error("invalid movie id param", slog.String("param", ctx.Param("id")), slog.Any("error", err))
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid movie id"})
		return
	}

	var req dto.MovieUpdateRequest

	if err := ctx.ShouldBindJSON(&req); err != nil {
		h.logger.Error("invalid update movie request", slog.Any("error", err), slog.Any("id", id))
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	movie, err := h.service.Update(uint(id), &req)

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			h.logger.Info("movie not found for update", slog.Any("id", id))
			ctx.JSON(http.StatusNotFound, gin.H{"error": "movie not found"})
			return
		}
		h.logger.Error("movie update failed", slog.Any("id", id), slog.Any("error", err))
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, movie)
}

func (h *MovieHandler) Delete(ctx *gin.Context) {

	id, err := strconv.ParseUint(ctx.Param("id"), 10, 64)

	if err != nil {
		h.logger.Error("invalid movie id param", slog.String("param", ctx.Param("id")), slog.Any("error", err))
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid movie id"})
		return
	}

	if err := h.service.Delete(uint(id)); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			h.logger.Info("movie not found for delete", slog.Any("id", id))
			ctx.JSON(http.StatusNotFound, gin.H{"error": "movie not found"})
			return
		}
		h.logger.Error("failed to delete movie", slog.Any("id", id), slog.Any("error", err))
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to delete movie"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "movie deleted successfully"})
}
