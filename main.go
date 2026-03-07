package main

import (
	"fmt"
	"log"

	"cbs-simulator/api/routes"
	"cbs-simulator/config"
	"cbs-simulator/database"
	"cbs-simulator/utils"

	"github.com/gin-gonic/gin"
)

func main() {
	// Load configuration
	log.Println("Loading configuration...")
	config.LoadConfig()

	// Initialize logger
	log.Println("Initializing logger...")
	utils.InitLogger()

	// Initialize database
	log.Println("Initializing database...")
	if err := database.InitDB(); err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	defer database.CloseDB()

	// Set Gin mode
	if config.AppConfig.Environment == "production" {
		gin.SetMode(gin.ReleaseMode)
	} else {
		gin.SetMode(gin.DebugMode)
	}

	// Create Gin router
	router := gin.Default()

	// Setup routes
	routes.SetupRoutes(router)

	// Start server
	log.Println("")
	addr := fmt.Sprintf("0.0.0.0:%s", config.AppConfig.ServerPort)
	log.Printf("CBS Simulator starting on %s", addr)
	log.Printf("Environment: %s", config.AppConfig.Environment)
	log.Printf("Database: %s", config.AppConfig.DatabasePath)
	log.Println("Server is ready to accept requests")
	log.Println("")
	log.Println("Available endpoints:")
	log.Println("  - Health Check: GET http://localhost:8080/health")
	log.Println("  - Login: POST http://localhost:8080/api/v1/auth/login")
	log.Println("  - API Documentation: See docs/API.md")
	log.Println("")

	if err := router.Run(addr); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
