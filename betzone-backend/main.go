package main

import (
	"log"
	"net/http"

	"github.com/betzone/backend/config"
	_ "github.com/betzone/backend/docs"
	"github.com/betzone/backend/routes"
	"github.com/betzone/backend/services"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
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

	// Initialize Database Service
	log.Println("Connecting to MySQL database...")
	dbService, err := services.NewDatabaseService(cfg.GetDatabaseDSN())
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	defer dbService.Close()

	// Run database migrations
	if err := dbService.Migrate(); err != nil {
		log.Fatalf("Failed to run migrations: %v", err)
	}

	// Initialize AuthService with database and config
	authService := services.NewAuthServiceWithConfig(cfg, dbService.DB)

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
	routes.RegisterRoutes(router, betkraftService, authService, dbService)

	// Setup Swagger
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Start server
	log.Printf("Starting server on port %s", cfg.Port)
	server := &http.Server{
		Addr:    ":" + cfg.Port,
		Handler: router,
	}

	if err := server.ListenAndServe(); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

func middlewareLogger() gin.HandlerFunc {
	return gin.Logger()
}

func corsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", c.Request.Header.Get("Origin"))
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With, x-api-key, x-timestamp, x-signature-key, Upgrade, Connection")
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
