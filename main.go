package main

import (
	"log"
	"own-paynet/api/routes"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	// Setup router
	router := gin.Default()

	// Add CORS middleware
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"}, // Allow all origins; adjust as needed
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * 60 * 60, // 12 hours
	}))

	// Setup all routes
	routes.SetupRoutes(router)

	// Start server
	log.Printf("Server running on port %s", "8080")
	err := router.Run(":8080")
	if err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
