package routes

import (
	"own-paynet/api/handlers"
	"own-paynet/api/middleware"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(router *gin.Engine, authHandler *handlers.AuthHandler, paymentHandler *handlers.PaymentHandler) {
	api := router.Group("/api/v1")
	{
		api.POST("/signup", authHandler.Signup)
		api.POST("/signin", authHandler.Signin)

		// Password reset flow
		api.POST("/request-password-reset", authHandler.RequestPasswordReset)
		api.POST("/verify-reset-token", authHandler.VerifyResetToken)
		api.POST("/reset-password-with-token", authHandler.ResetPasswordWithToken)
		// Add these routes to the existing routes
		api.POST("/verify-email", authHandler.VerifyEmail)
		api.POST("/resend-verification-email", authHandler.ResendVerificationEmail)

		// Payment-related routes
		api.POST("/webhook", paymentHandler.HandleWebhook)

		protected := api.Group("/")
		protected.Use(middleware.AuthMiddleware())
		{
			protected.POST("/logout", authHandler.Logout)
			protected.POST("/payments", paymentHandler.CreatePayment)
		}
	}

}
