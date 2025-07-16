package routes

import (
	"restaurant-app/controllers"
	"restaurant-app/middleware"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(router *gin.Engine) {
	api := router.Group("/api/v1")
	{
		api.POST("/register", controllers.Register)
		api.POST("/login", controllers.Login)
	}

	// ✅ Protected routes (require JWT auth)
	protected := api.Group("/")
	protected.Use(middleware.AuthMiddleware())
	{
		protected.GET("/profile", controllers.GetProfile)
		protected.PUT("/profile", controllers.UpdateProfile)
		protected.PUT("/change-password", controllers.ChangePassword)
	}
}
