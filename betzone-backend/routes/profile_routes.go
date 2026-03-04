package routes

import (
	"github.com/betzone/backend/handlers"
	"github.com/betzone/backend/services"
	"github.com/gin-gonic/gin"
)

// RegisterProfileRoutes registers user profile routes (protected)
func RegisterProfileRoutes(router *gin.RouterGroup, authService *services.AuthService) {
	router.GET("/auth/profile", func(c *gin.Context) {
		handlers.GetProfileHandler(c, authService)
	})
}
