package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/yourusername/online-library/internal/auth"
	"go.uber.org/zap"
)

// AuthMiddleware validates JWT tokens
func AuthMiddleware(authSvc auth.Service, log *zap.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "authorization header required"})
			c.Abort()
			return
		}

		// Extract token from "Bearer <token>"
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid authorization header format"})
			c.Abort()
			return
		}

		token := parts[1]
		claims, err := authSvc.ValidateToken(c.Request.Context(), token)
		if err != nil {
			log.Warn("invalid token", zap.Error(err))
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid or expired token"})
			c.Abort()
			return
		}

		// Set user info in context
		c.Set("user_id", claims.UserID)
		c.Set("user_role", claims.Role)

		c.Next()
	}
}

// GetUserID extracts user ID from context
func GetUserID(c *gin.Context) string {
	userID, _ := c.Get("user_id")
	return userID.(string)
}

// GetUserRole extracts user role from context
func GetUserRole(c *gin.Context) string {
	role, _ := c.Get("user_role")
	return role.(string)
}
