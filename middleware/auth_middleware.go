package middleware

import (
	"net/http"
	"restaurant-app/utils"
	"strings"

	"github.com/gin-gonic/gin"
)

// AuthMiddleware validates JWT tokens from the Authorization header
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")

		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "Authorization header missing or invalid",
			})
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")

		// ✅ Validate the token using utils.ValidateJWT
		claims, err := utils.ValidateJWT(tokenString)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "Invalid or expired token",
			})
			return
		}

		// ✅ Store values in context for future use (e.g., in handlers)
		if userID, ok := claims["user_id"].(string); ok {
			c.Set("user_id", userID)
		}
		if role, ok := claims["role"].(string); ok {
			c.Set("role", role)
		}

		c.Next()
	}
}
