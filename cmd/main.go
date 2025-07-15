package main

import (
	"fmt"
	"restaurant-app/config"
	"restaurant-app/middleware"

	"restaurant-app/routes"

	"github.com/gin-gonic/gin"
)

func main() {
	// Load environment variables from .env
	config.LoadEnv()

	// Connect to MongoDB
	config.ConnectDB()

	// Create Gin router with default middleware (logger + recovery)
	r := gin.Default()

	//  Apply middleware
	r.Use(middleware.CORSMiddleware())
	r.Use(middleware.RateLimitMiddleware())

	// Register API routes
	routes.RegisterRoutes(r)

	// Start server on port 5000
	fmt.Println("✅ Server running at port 5000")
	err := r.Run(":5000")
	if err != nil {
		fmt.Println("❌ Failed to start server:", err)
	}
}
