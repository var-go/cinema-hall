package transport

import "github.com/gin-gonic/gin"

func RegisterRouters(r *gin.Engine,
	auth *AuthHandler,
	users *UserHandler,
) {
	{
		authGroup := r.Group("/auth")

		authGroup.POST("/auth/register", auth.Register)
		authGroup.POST("/auth/login", auth.Login)
	}
	user := r.Group("/users")
	{
		user.POST("", users.Create)
		user.GET("", users.List)
		user.GET("/:id", users.Get)
		user.PUT("/:id", users.Update)
		user.DELETE("/:id", users.Delete)
	}
	// protected
	protected := r.Group("/users")
	protected.Use(JWTMiddleware())

	{
		protected.GET("/me", users.Me)
		protected.GET("/me/bookings", users.MyBookings)
	}

}
