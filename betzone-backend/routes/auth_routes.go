package routes

import (
	"github.com/betzone/backend/handlers"
	"github.com/betzone/backend/services"
	"github.com/gin-gonic/gin"
)

// RegisterAuthRoutes registers authentication routes
func RegisterAuthRoutes(router *gin.RouterGroup, authService *services.AuthService) {
	auth := router.Group("/auth")
	{
		auth.POST("/signup", func(c *gin.Context) {
			handlers.SignupHandler(c, authService)
		})
		auth.POST("/signin", func(c *gin.Context) {
			handlers.SigninHandler(c, authService)
		})
	}
}
