package main

import (
	"log"

	_ "vaxen/api/docs"
	"vaxen/api/internal/config"
	"vaxen/api/internal/database"
	"vaxen/api/internal/middleware"
	"vaxen/api/internal/providers"
	"vaxen/api/internal/routes"
	"vaxen/api/internal/services"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

// @title           Vaxen API
// @version         1.0
// @description     Global treasury and cross-border payments API
// @termsOfService  http://swagger.io/terms/

// @contact.name   API Support
// @contact.url    http://www.vaxen.io/support
// @contact.email  support@vaxen.io

// @license.name  Apache 2.0
// @license.url   http://www.apache.org/licenses/LICENSE-2.0.html

// @host      localhost:8080
// @BasePath  /api/v1

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Type "Bearer" followed by a space and JWT token.

func main() {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using system environment variables")
	}

	// Load configuration
	cfg := config.Load()

	// Connect to database
	if cfg.DatabaseURL != "" {
		if err := database.Connect(cfg); err != nil {
			log.Printf("Database connection failed: %v", err)
			log.Println("Continuing without database...")
		} else {
			// Run migrations
			if err := database.Migrate(); err != nil {
				log.Printf("Migration failed: %v", err)
			} else {
				log.Println("Database migrations completed")
			}
		}
	}

	// Bootstrap providers from config (swap by changing env vars)
	reg := providers.Bootstrap(cfg)
	log.Printf("Providers: KYC=%s MFA=%s Wallet=%s Exchange=%s Payment=%s Email=%s",
		reg.KYC().Name(), reg.MFA().Name(), reg.Wallet().Name(),
		reg.Exchange().Name(), reg.Payment().Name(), reg.Email().Name())

	// Initialize service layer
	svc := services.NewContainer(database.DB, cfg, reg)

	// Set Gin mode
	if cfg.Environment == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	// Initialize router
	router := gin.Default()

	// Apply global middleware
	router.Use(middleware.CORS(cfg))
	router.Use(middleware.Security())
	router.Use(middleware.Logger())
	router.Use(middleware.ErrorHandler())

	// Setup routes
	routes.SetupRoutes(router, cfg, svc)

	// Start server
	port := cfg.Port
	if port == "" {
		port = "8080"
	}

	log.Printf("Vaxen API running on http://localhost:%s\n", port)
	log.Printf("API Documentation: http://localhost:%s/docs/index.html\n", port)
	log.Printf("Health Check: http://localhost:%s/health\n", port)

	if err := router.Run(":" + port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
