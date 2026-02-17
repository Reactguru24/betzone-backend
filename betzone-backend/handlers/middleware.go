package handlers

import (
	"net/http"
	"strings"

	"github.com/betzone/backend/models"
	"github.com/betzone/backend/services"
	"github.com/gin-gonic/gin"
)

// AuthMiddleware validates JWT token in Authorization header
func AuthMiddleware(authService *services.AuthService) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, models.ErrorResponse{
				Success: false,
				Message: "Missing authorization header",
			})
			c.Abort()
			return
		}

		// Extract token from "Bearer <token>"
		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.JSON(http.StatusUnauthorized, models.ErrorResponse{
				Success: false,
				Message: "Invalid authorization header format",
			})
			c.Abort()
			return
		}

		token := parts[1]

		// Verify token
		claims, err := authService.VerifyToken(token)
		if err != nil {
			c.JSON(http.StatusUnauthorized, models.ErrorResponse{
				Success: false,
				Message: "Invalid or expired token",
				Error:   err.Error(),
			})
			c.Abort()
			return
		}

		// Store claims in context for use in handlers
		c.Set("user_id", claims.UserID)
		c.Set("phone", claims.Phone)

		c.Next()
	}
}
