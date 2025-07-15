package routes

import (
	"restaurant-app/controllers"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(router *gin.Engine) {
	api := router.Group("/api/v1")
	{
		api.POST("/register", controllers.Register)
		api.POST("/login", controllers.Login)
	}
}
