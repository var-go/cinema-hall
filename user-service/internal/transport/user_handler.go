package transport

import (
	"log/slog"
	"net/http"
	"strconv"
	"user-service/internal/config"
	"user-service/internal/dto"
	"user-service/internal/models"
	"user-service/internal/services"

	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	service services.UserService
	log     *slog.Logger
}

func NewUserHandler(service services.UserService, log *slog.Logger) *UserHandler {
	return &UserHandler{service: service, log: log}
}

func (h *UserHandler) RegisterUserRoutes(r *gin.Engine) {
	user := r.Group("/users")
	{
		user.POST("", h.Create)
		user.GET("", h.List)
		user.GET("/:id", h.Get)
		user.PUT("/:id", h.Update)
		user.DELETE("/:id", h.Delete)
	}
}

func (h *UserHandler) Create(c *gin.Context) {
	var req dto.CreateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.log.Warn("invalid create user request", "err", err)
		c.JSON(400, gin.H{"error": "invalid request body"})
		return
	}

	user, err := h.service.Create(req)
	if err != nil {
		h.log.Error("failed to create user", "email", req.Email, "err", err)
		c.JSON(500, gin.H{"error": "internal error"})
		return
	}

	c.JSON(201, toUserResponse(user))
}

func (h *UserHandler) Get(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		h.log.Warn("invalid user id", "id", c.Param("id"))
		c.JSON(400, gin.H{"error": "invalid user id"})
		return
	}
	user, err := h.service.Get(uint(id))
	if err != nil {
		h.log.Warn("user not found", "id", id)
		c.JSON(404, gin.H{"error": "not found"})
		return
	}
	c.JSON(200, toUserResponse(user))
}

func (h *UserHandler) List(c *gin.Context) {
	users, err := h.service.List()
	if err != nil {
		h.log.Error("failed to list users", "err", err)
		c.JSON(500, gin.H{"error": "internal error"})
		return
	}
	resp := make([]dto.UserResponse, 0)
	for _, u := range users {
		resp = append(resp, toUserResponse(&u))
	}
	c.JSON(200, resp)
}

func (h *UserHandler) Update(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		h.log.Warn("invalid user id for update", "id", c.Param("id"))
		c.JSON(400, gin.H{"error": "invalid user id"})
		return
	}

	var req dto.UpdateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.log.Warn("invalid update user request", "id", id, "err", err)
		c.JSON(400, gin.H{"error": "invalid request body"})
		return
	}

	user, err := h.service.Update(uint(id), req)
	if err != nil {
		h.log.Warn("user not found for update", "id", id)
		c.JSON(404, gin.H{"error": "not found"})
		return
	}

	c.JSON(200, toUserResponse(user))
}

func (h *UserHandler) Delete(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		h.log.Warn("invalid user id for delete", "id", c.Param("id"))
		c.JSON(400, gin.H{"error": "invalid user id"})
		return
	}
	if err := h.service.Delete(uint(id)); err != nil {
		h.log.Warn("user not found for delete", "id", id)
		c.JSON(404, gin.H{"error": "not found"})
		return
	}
	c.Status(204)
}

func toUserResponse(u *models.User) dto.UserResponse {
	return dto.UserResponse{
		ID:    u.ID,
		Email: u.Email,
		Name:  u.Name,
		Role:  u.Role,
	}
}

func (h *UserHandler) Me(c *gin.Context) {
	userID := c.GetUint("user_id")

	user, err := h.service.Get(userID)
	if err != nil {
		h.log.Warn("me: user not found", "user_id", userID)
		c.JSON(404, gin.H{"error": "user not found"})
		return
	}

	c.JSON(200, toUserResponse(user))
}

func (h *UserHandler) MyBookings(c *gin.Context) {
	userID := c.GetUint("user_id")

	url := config.BookingServiceURL() + "/bookings/user/" + strconv.Itoa(int(userID))

	resp, err := http.Get(url)
	if err != nil {
		h.log.Error("booking service unavailable", "url", url, "err", err)
		c.JSON(500, gin.H{"error": "booking service unavailable"})
		return
	}
	defer resp.Body.Close()

	c.DataFromReader(
		resp.StatusCode,
		resp.ContentLength,
		resp.Header.Get("Content-Type"),
		resp.Body,
		nil,
	)
}
