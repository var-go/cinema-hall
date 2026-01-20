package transport

import (
	"cinema-service/internal/dto"
	"cinema-service/internal/services"
	"errors"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type SessionHandler struct {
	sessionService services.SessionService
	logger         *slog.Logger
}

func NewSessionHandler(sessionService services.SessionService, logger *slog.Logger) *SessionHandler {
	return &SessionHandler{
		sessionService: sessionService,
		logger:         logger,
	}
}

func (h *SessionHandler) RegisterRoutes(r *gin.Engine) {
	sessions := r.Group("/")
	{
		sessions.GET("/sessions", h.List)
		sessions.GET("/sessions/:id", h.GetById)
		sessions.POST("/sessions", h.Create)
		sessions.PATCH("/sessions/:id", h.Update)
		sessions.DELETE("/sessions/:id", h.Delete)
		sessions.GET("/movies/:id/sessions", h.ListByMovieID)
	}
}

func (h *SessionHandler) Create(c *gin.Context) {

	var req dto.CreateSessionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Error("handler: failed to bind JSON", "err", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	session, err := h.sessionService.Create(req)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "hall not found"})
			return
		}

		h.logger.Error("failed to create session", "err", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, session)
}

func (h *SessionHandler) List(c *gin.Context) {
	sessions, err := h.sessionService.List()
	if err != nil {
		h.logger.Error("failed to list sessions", "err", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "failed to list sessions"})
		return
	}

	c.JSON(http.StatusOK, sessions)
}

func (h *SessionHandler) GetById(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	session, err := h.sessionService.GetById(uint(id))
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "session not found"})
			return
		}

		h.logger.Error("failed to fetch session", "id", id, "err", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch session"})
		return
	}

	c.JSON(http.StatusOK, session)
}

func (h *SessionHandler) Update(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	var req dto.UpdateSessionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Error("handler: failed to bind JSON", "err", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	session, err := h.sessionService.Update(uint(id), req)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "session not found"})
			return
		}

		h.logger.Error("failed to update session", "id", id, "err", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, session)
}

func (h *SessionHandler) Delete(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	if err := h.sessionService.Delete(uint(id)); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "session not found"})
			return
		}

		h.logger.Error("failed to delete session", "id", id, "err", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "session deleted successfully"})
}

func (h *SessionHandler) ListByMovieID(c *gin.Context) {
	idStr := c.Param("id")
	movieID, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid movie id"})
		return
	}

	sessions, err := h.sessionService.ListByMovieID(uint(movieID))
	if err != nil {
		h.logger.Error(
			"failed to list sessions by movie id",
			"movie_id", movieID,
			"err", err,
		)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, sessions)
}
