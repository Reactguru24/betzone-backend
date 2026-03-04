package routes

import (
	"github.com/betzone/backend/handlers"
	"github.com/betzone/backend/services"
	"github.com/gin-gonic/gin"
)

// RegisterRoutes registers all API routes
func RegisterRoutes(router *gin.Engine, betkraftService *services.BetkraftService, authService *services.AuthService, dbService *services.DatabaseService) {
	// Health check
	router.GET("/health", handlers.HealthHandler)

	// API v1 routes
	v1 := router.Group("/api/v1")
	{
		// Authentication routes (no auth required)
		RegisterAuthRoutes(v1, authService)

		// Games routes
		RegisterGameRoutes(v1, betkraftService)

		// Odds routes
		RegisterOddsRoutes(v1)

		// Protected routes (require authentication)
		protected := v1.Group("")
		protected.Use(handlers.AuthMiddleware(authService))
		{
			// User profile
			RegisterProfileRoutes(protected, authService)

			// Bets
			RegisterBetRoutes(protected, betkraftService, dbService)
		}

		// Callback routes (from Betkraft provider)
		RegisterCallbackRoutes(v1, authService, dbService)
	}
}
