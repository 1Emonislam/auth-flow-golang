package routes

import (
	"log"
	"own-paynet/api/handlers"
	"own-paynet/api/middleware"
	"own-paynet/config"
	"own-paynet/database"
	"own-paynet/repository"
	"own-paynet/services"
	"own-paynet/services/bitcoin"
	"own-paynet/utils/email"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(router *gin.Engine) {
	// Load config
	cfg := config.LoadConfig()

	// Initialize database
	db := database.InitDB(cfg)

	// Initialize Redis
	_ = database.InitRedis(cfg)

	// Initialize Bitcoin service
	bitcoinService, err := bitcoin.NewBitcoinService(cfg)
	if err != nil {
		log.Fatal("failed to initialize Bitcoin service:", err)
	}

	// Initialize email service
	emailService := email.NewEmailService(cfg)

	// Initialize repositories
	userRepo := repository.NewUserRepository(db)
	apiKeyRepo := repository.NewAPIKeyRepository(db)
	apiKeyService := services.NewAPIKeyService(apiKeyRepo)
	payoutWalletRepo := repository.NewPayoutWalletRepository(db)
	payoutWalletService := services.NewPayoutWalletService(payoutWalletRepo)
	authService := services.NewAuthService(userRepo, emailService, apiKeyService, payoutWalletService)
	authHandler := handlers.NewAuthHandler(authService)

	paymentRepo := repository.NewPaymentRepository(db)
	paymentService := services.NewPaymentService(paymentRepo, bitcoinService, cfg.BaseURL, cfg.BitcoinNetwork)
	paymentHandler := handlers.NewPaymentHandler(paymentService, cfg)

	// Initialize company service and handler
	companyRepo := repository.NewCompanyRepository(db)
	companyService := services.NewCompanyService(companyRepo)
	companyHandler := handlers.NewCompanyHandler(companyService)

	// Initialize payout wallet repository, service, and handler
	payoutWalletRepo = repository.NewPayoutWalletRepository(db)
	payoutWalletService = services.NewPayoutWalletService(payoutWalletRepo)
	payoutWalletHandler := handlers.NewPayoutWalletHandler(payoutWalletService)

	// Initialize transaction repository, service, and handler
	transactionRepo := repository.NewTransactionRepository(db)
	transactionService := services.NewTransactionService(transactionRepo, payoutWalletService, bitcoinService)
	userService := services.NewUserService(userRepo)
	transactionHandler := handlers.NewTransactionHandler(transactionService, payoutWalletService, userService)

	// Initialize API key handler
	apiKeyHandler := handlers.NewAPIKeyHandler(apiKeyService)

	// Initialize 2FA service and handler
	twoFactorService := services.NewTwoFactorService(userRepo, services.NewEmailService(cfg))
	twoFactorHandler := handlers.NewTwoFactorHandler(twoFactorService)

	// Setup routes
	api := router.Group("/api/v1")
	{
		api.POST("/signup", authHandler.Signup)
		api.POST("/signin", authHandler.Signin)

		// Google OAuth routes
		api.GET("/auth/google", authHandler.GoogleLogin)
		api.GET("/auth/google/callback", authHandler.GoogleCallback)

		// Password reset flow
		api.POST("/request-password-reset", authHandler.RequestPasswordReset)
		api.POST("/verify-reset-token", authHandler.VerifyResetToken)
		api.POST("/reset-password-with-token", authHandler.ResetPasswordWithToken)

		// Email verification routes
		api.POST("/verify-email", authHandler.VerifyEmail)
		api.POST("/resend-verification-email", authHandler.ResendVerificationEmail)

		// Payment-related routes
		api.POST("/webhook", paymentHandler.HandleWebhook)

		protected := api.Group("/")
		protected.Use(middleware.AuthMiddleware())
		{
			protected.POST("/logout", authHandler.Logout)
			// 2FA routes
			// Enable 2FA flow
			protected.POST("/2fa/enable/email", twoFactorHandler.EnableEmail2FA)
			protected.POST("/2fa/enable/email/verify", twoFactorHandler.VerifyAndEnableEmail2FA)
			protected.POST("/2fa/enable/authenticator", twoFactorHandler.EnableAuthenticator2FA)
			protected.POST("/2fa/enable/authenticator/verify", twoFactorHandler.VerifyAndEnableAuthenticator2FA)

			// Disable 2FA flow
			protected.POST("/2fa/disable/email", twoFactorHandler.DisableEmail2FA)
			protected.POST("/2fa/disable/email/verify", twoFactorHandler.VerifyAndDisableEmail2FA)
			protected.POST("/2fa/disable/authenticator", twoFactorHandler.DisableAuthenticator2FA)
			protected.POST("/2fa/disable/authenticator/verify", twoFactorHandler.VerifyAndDisableAuthenticator2FA)

			// Other 2FA operations
			protected.POST("/2fa/verify", twoFactorHandler.Verify2FA)
			protected.POST("/2fa/send-otp", twoFactorHandler.SendOTP)
			protected.GET("/2fa/authenticator-qr", twoFactorHandler.GetAuthenticatorQRCode)

			protected.POST("/payments", paymentHandler.CreatePayment)
			protected.PUT("/company/:id", companyHandler.UpdateCompany)

			// Payout wallet routes
			protected.POST("/payout-wallets", payoutWalletHandler.CreatePayoutWallet)
			protected.GET("/payout-wallets", payoutWalletHandler.GetUserPayoutWallets)
			protected.GET("/payout-wallets/:id", payoutWalletHandler.GetPayoutWallet)
			protected.PUT("/payout-wallets/:id", payoutWalletHandler.UpdatePayoutWallet)
			protected.DELETE("/payout-wallets/:id", payoutWalletHandler.DeletePayoutWallet)

			// Transaction routes
			protected.POST("/transactions", transactionHandler.CreateTransaction)
			protected.GET("/transactions/:id", transactionHandler.GetTransaction)
			protected.GET("/wallets/:wallet_id/transactions", transactionHandler.GetWalletTransactions)
			protected.GET("/transactions", transactionHandler.GetUserTransactions)

			// API Key routes
			protected.POST("/api-keys", apiKeyHandler.GenerateAPIKey)
			protected.GET("/api-keys", apiKeyHandler.GetUserAPIKeys)
			protected.PUT("/api-keys/:id/default", apiKeyHandler.SetDefaultKey)
			protected.DELETE("/api-keys/:id", apiKeyHandler.DeleteAPIKey)

		}
	}
}
