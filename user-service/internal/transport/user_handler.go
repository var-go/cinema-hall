package transport

import (
	"net/http"
	"strconv"
	"user-service/internal/dto"
	"user-service/internal/models"
	"user-service/internal/services"

	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	service services.UserService
}

func NewUserHandler(service services.UserService) *UserHandler {
	return &UserHandler{service: service}
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
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	user, err := h.service.Create(req)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(201, toUserResponse(user))
}

func (h *UserHandler) Get(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	user, err := h.service.Get(uint(id))
	if err != nil {
		c.JSON(404, gin.H{"error": "not found"})
		return
	}
	c.JSON(200, toUserResponse(user))
}

func (h *UserHandler) List(c *gin.Context) {
	users, _ := h.service.List()
	resp := make([]dto.UserResponse, 0)
	for _, u := range users {
		resp = append(resp, toUserResponse(&u))
	}
	c.JSON(200, resp)
}

func (h *UserHandler) Update(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	var req dto.UpdateUserRequest

	c.ShouldBindJSON(&req)
	user, err := h.service.Update(uint(id), req)
	if err != nil {
		c.JSON(404, gin.H{"error": "not found"})
		return
	}

	c.JSON(200, toUserResponse(user))
}

func (h *UserHandler) Delete(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	if err := h.service.Delete(uint(id)); err != nil {
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
	// 1. Берём user_id из JWT middleware
	userID := c.GetUint("user_id")

	// 2. Получаем пользователя
	user, err := h.service.Get(userID)
	if err != nil {
		c.JSON(404, gin.H{"error": "user not found"})
		return
	}

	// 3. Отдаём пользователя
	c.JSON(200, toUserResponse(user))
}

func (h *UserHandler) MyBookings(c *gin.Context) {
	userID := c.GetUint("user_id")

	url := "http://localhost:8081/bookings/user/" + strconv.Itoa(int(userID))

	resp, err := http.Get(url)
	if err != nil {
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

// admin := r.Group("/users")
// admin.Use(AuthMiddleware(), AdminOnly())
