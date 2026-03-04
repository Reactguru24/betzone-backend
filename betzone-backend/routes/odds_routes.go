package routes

import (
	"github.com/betzone/backend/handlers"
	"github.com/gin-gonic/gin"
)

// RegisterOddsRoutes registers odds-related routes
func RegisterOddsRoutes(router *gin.RouterGroup) {
	router.GET("/odds/:gameId", handlers.GetOddsHandler)
}
