package main

import (
	"log"
	"ticert/config"
	"ticert/routes"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Println("Warning: .env file not found:", err)
	}

	// Initialize database
	db := config.InitDatabase()

	// Auto migrate database
	config.AutoMigrate(db)

	// Initialize Redis
	config.InitRedis(config.GetConfig())

	// Setup Gin router
	r := gin.Default()

	// Setup routes
	routes.SetupRoutes(db, r)

	// Start server
	cfg := config.GetConfig()
	log.Printf("Server starting on port %s", cfg.Port)
	if err := r.Run(":" + cfg.Port); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
