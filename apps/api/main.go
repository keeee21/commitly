package main

import (
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/keeee21/commitly/api/db"
	"github.com/keeee21/commitly/api/router"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

// @title           Commitly API
// @version         1.0
// @description     API for tracking GitHub commit activity and comparing with rivals
// @host            localhost:8080
// @BasePath        /
// @securityDefinitions.apikey  GitHubUserID
// @in                          header
// @name                        X-GitHub-User-ID
func main() {
	// Load .env file
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}

	// Connect to database
	database, err := db.NewDatabase(os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// Auto-migrate models
	if err := db.AutoMigrate(database); err != nil {
		log.Fatalf("Failed to migrate database: %v", err)
	}

	// Initialize Echo
	e := echo.New()

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORS())

	// Setup routes
	router.SetupRoutes(e, database)

	// Start server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	e.Logger.Fatal(e.Start(":" + port))
}
