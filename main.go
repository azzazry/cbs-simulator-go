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
	config.LoadConfig()

	// Initialize logger
	utils.InitLogger()

	// Initialize database
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

	// Create router
	router := gin.Default()

	// Setup routes
	routes.SetupRoutes(router)

	// Server address
	addr := fmt.Sprintf("0.0.0.0:%s", config.AppConfig.ServerPort)

	// Show startup screen
	utils.StartupScreen(addr)

	// Start server
	if err := router.Run(addr); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
