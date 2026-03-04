package routes

import (
	"github.com/betzone/backend/handlers"
	"github.com/betzone/backend/services"
	"github.com/gin-gonic/gin"
)

// RegisterGameRoutes registers game-related routes
func RegisterGameRoutes(router *gin.RouterGroup, betkraftService *services.BetkraftService) {
	router.GET("/games", func(c *gin.Context) {
		handlers.GetGamesHandler(c, betkraftService)
	})
	router.GET("/games/:id", handlers.GetGameByIDHandler)
	router.POST("/launch", func(c *gin.Context) {
		handlers.LaunchGameHandler(c, betkraftService)
	})
}
