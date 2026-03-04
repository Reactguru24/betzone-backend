package routes

import (
	"github.com/betzone/backend/handlers"
	"github.com/betzone/backend/services"
	"github.com/gin-gonic/gin"
)

// RegisterCallbackRoutes registers Betkraft callback routes
func RegisterCallbackRoutes(router *gin.RouterGroup, authService *services.AuthService, dbService *services.DatabaseService) {
	callbacks := router.Group("/callbacks")
	{
		callbacks.POST("/player_info", func(c *gin.Context) {
			handlers.PlayerInfoCallback(c, authService, dbService)
		})
		callbacks.POST("/bet", func(c *gin.Context) {
			handlers.BetCallback(c, authService, dbService)
		})
		callbacks.POST("/win", func(c *gin.Context) {
			handlers.WinCallback(c, authService, dbService)
		})
		callbacks.POST("/rollback", func(c *gin.Context) {
			handlers.RollbackCallback(c, authService, dbService)
		})
	}
}
