package main

import (
	"log"
	"own-paynet/api/handlers"
	"own-paynet/api/routes"
	"own-paynet/config"
	"own-paynet/database"
	"own-paynet/repository"
	"own-paynet/services"
	"own-paynet/services/bitcoin"
	"own-paynet/utils/email"

	"github.com/gin-gonic/gin"
)

func main() {
	// Load config
	cfg := config.LoadConfig()

	// Initialize database
	db := database.InitDB(cfg)

	// Initialize Redis
	_ = database.InitRedis(cfg)

	// Initialize Bitcoin service
	bitcoinService, err := bitcoin.NewBitcoinService(cfg)
	if err != nil {
		log.Fatal("Failed to initialize Bitcoin service:", err)
	}

	// Initialize email service
	emailService := email.NewEmailService(cfg)

	// Initialize repositories
	userRepo := repository.NewUserRepository(db)
	paymentRepo := repository.NewPaymentRepository(db)

	// Initialize services
	authService := services.NewAuthService(userRepo, emailService)

	// Pass bitcoinNetwork to NewPaymentService
	paymentService := services.NewPaymentService(paymentRepo, bitcoinService, cfg.BaseURL, cfg.BitcoinNetwork)

	// Initialize handlers
	authHandler := handlers.NewAuthHandler(authService)
	paymentHandler := handlers.NewPaymentHandler(paymentService, cfg)

	// Setup router
	router := gin.Default()
	routes.SetupRoutes(router, authHandler, paymentHandler)

	// Start server
	log.Printf("Server running on port %s", cfg.ServerPort)
	err = router.Run(":" + cfg.ServerPort)
	if err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
