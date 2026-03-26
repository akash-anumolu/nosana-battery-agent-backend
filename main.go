package main

import (
	"battery-agent/configs"
	"battery-agent/cron"
	"battery-agent/routes"
	"log"

	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	// Load environment variables first
	if err := godotenv.Load(); err != nil {
		log.Println("Warning: .env file not found")
	}

	// Connect to MongoDB
	configs.ConnectDB()

	// Initialize Echo
	e := echo.New()

	// Middleware
	e.Use(middleware.CORS())
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// Register all routes
	routes.RegisterRoutes(e)

	// Start background monitoring
	cron.StartMonitoringCron()

	// Start server
	port := configs.EnvPort()
	log.Printf("🔋 Battery Agent starting on :%s", port)
	e.Logger.Fatal(e.Start(":" + port))
}
