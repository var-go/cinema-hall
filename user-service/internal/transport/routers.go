package transport

import (
	"user-service/internal/middleware"

	"github.com/gin-gonic/gin"
)

func RegisterRouters(r *gin.Engine,
	auth *AuthHandler,
	users *UserHandler,
) {
	{
		authGroup := r.Group("/auth")

		authGroup.POST("/register", auth.Register)
		authGroup.POST("/login", auth.Login)
	}

	admin := r.Group("/admin")
	admin.Use(middleware.JWTMiddleware())
	admin.Use(middleware.AdminMiddleware())
	{
		admin.DELETE("/:id", users.Delete)
		admin.PUT("/:id", users.Update)
	}

	user := r.Group("/users")
	user.Use(middleware.JWTMiddleware())

	{
		user.GET("", users.List)
		user.GET("/:id", users.Get)
		user.DELETE("/:id", users.Delete)
		user.PUT("/:id", users.Update)
	}

	protected := r.Group("")
	protected.Use(middleware.JWTMiddleware())

	{
		protected.GET("/me", users.Me)
		protected.GET("/me/bookings", users.MyBookings)
	}

}
