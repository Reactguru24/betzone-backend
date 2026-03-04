package routes

import (
	"github.com/betzone/backend/handlers"
	"github.com/betzone/backend/services"
	"github.com/gin-gonic/gin"
)

// RegisterBetRoutes registers bet-related routes (protected)
func RegisterBetRoutes(router *gin.RouterGroup, betkraftService *services.BetkraftService, dbService *services.DatabaseService) {
	router.POST("/bets", handlers.CreateBetHandler)
	router.GET("/bets", func(c *gin.Context) {
		handlers.GetBetsHandler(c, dbService)
	})
	router.GET("/bets/:id", handlers.GetBetByIDHandler)
	router.GET("/bets/status/:game_uuid", func(c *gin.Context) {
		handlers.GetBetStatusHandler(c, betkraftService)
	})
}
