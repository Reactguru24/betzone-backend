package main

import (
	"log"

	"github.com/betzone/backend/config"
	"github.com/betzone/backend/handlers"
	"github.com/betzone/backend/services"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Println("Warning: .env file not found, using system environment variables")
	}

	// Load configuration
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Debug: Print loaded configuration
	log.Printf("Loaded Configuration:")
	log.Printf("  Port: %s", cfg.Port)
	log.Printf("  Environment: %s", cfg.Environment)
	log.Printf("  BetkraftBaseURL: %s", cfg.BetkraftBaseURL)
	log.Printf("  BetkraftAPIKey: %s...", maskSecret(cfg.BetkraftAPIKey))
	log.Printf("  BetkraftAppKey: %s...", maskSecret(cfg.BetkraftAppKey))

	// Initialize BetkraftService
	betkraftService := services.NewBetkraftService(cfg)

	// Set Gin mode
	if cfg.Environment == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	// Create router
	router := gin.Default()

	// Register middleware
	router.Use(middlewareLogger())
	router.Use(corsMiddleware())

	// Register routes
	registerRoutes(router, betkraftService)

	// Start server
	log.Printf("Starting server on port %s", cfg.Port)
	if err := router.Run(":" + cfg.Port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

func registerRoutes(router *gin.Engine, betkraftService *services.BetkraftService) {
	// Health check
	router.GET("/health", handlers.HealthHandler)

	// API v1 routes
	v1 := router.Group("/api/v1")
	{
		// Games
		v1.GET("/games", func(c *gin.Context) {
			handlers.GetGamesHandler(c, betkraftService)
		})
		v1.GET("/games/:id", handlers.GetGameByIDHandler)
		v1.POST("/launch", func(c *gin.Context) {
			handlers.LaunchGameHandler(c, betkraftService)
		})

		// Bets
		v1.POST("/bets", handlers.CreateBetHandler)
		v1.GET("/bets", handlers.GetBetsHandler)
		v1.GET("/bets/:id", handlers.GetBetByIDHandler)

		// Odds
		v1.GET("/odds/:gameId", handlers.GetOddsHandler)
	}
}

func middlewareLogger() gin.HandlerFunc {
	return gin.Logger()
}

func corsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With, x-api-key, x-timestamp, x-signature-key")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}
func maskSecret(secret string) string {
	if len(secret) <= 8 {
		return "****"
	}
	return secret[:4] + "****" + secret[len(secret)-4:]
}