package handlers

import (
	"net/http"
	response "own-paynet/api/response"
	"own-paynet/services"

	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	authService *services.AuthService
}

func NewAuthHandler(authService *services.AuthService) *AuthHandler {
	return &AuthHandler{authService: authService}
}

type SignupRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
}

func (h *AuthHandler) Signup(c *gin.Context) {
	var req SignupRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	if err := h.authService.Signup(req.Email, req.Password); err != nil {
		// Check for specific error types and return appropriate status codes
		if err.Error() == "an account with this email already exists" {
			response.ErrorResponse(c, http.StatusConflict, "This email address is already registered. Please use a different email or try logging in.")
			return
		}
		// For other errors, use a more specific message with 500 status
		response.ErrorResponse(c, http.StatusInternalServerError, "Unable to create account at this time. Please try again later.")
		return
	}

	response.SuccessResponse(c, http.StatusOK, "User created successfully. Please check your email to verify your account.", nil)
}

type SigninRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

func (h *AuthHandler) Signin(c *gin.Context) {
	var req SigninRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	token, user, err := h.authService.Signin(req.Email, req.Password)
	if err != nil {
		switch err.Error() {
		case "email not verified":
			response.ErrorResponse(c, http.StatusUnauthorized, "Email not verified. Please check your email for verification instructions.")
		case "invalid credentials":
			response.ErrorResponse(c, http.StatusUnauthorized, "The email or password you entered is incorrect. Please try again.")
		case "user not found":
			response.ErrorResponse(c, http.StatusNotFound, "No account found with this email address. Please check your email or sign up for a new account.")
		default:
			response.ErrorResponse(c, http.StatusInternalServerError, "Unable to sign in at this time. Please try again later.")
		}
		return
	}

	response.SuccessResponse(c, http.StatusOK, "Login successful", gin.H{
		"token": token,
		"user": gin.H{
			"id":    user.ID,
			"email": user.Email,
		},
	})
}

type ResetPasswordRequest struct {
	Email       string `json:"email" binding:"required,email"`
	NewPassword string `json:"new_password" binding:"required,min=6"`
}

func (h *AuthHandler) ResetPassword(c *gin.Context) {
	var req ResetPasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	if err := h.authService.ResetPassword(req.Email, req.NewPassword); err != nil {
		response.ErrorResponse(c, http.StatusInternalServerError, "Failed to reset password")
		return
	}

	response.SuccessResponse(c, http.StatusOK, "Password reset successfully", nil)
}

type RequestPasswordResetRequest struct {
	Email string `json:"email" binding:"required,email"`
}

func (h *AuthHandler) RequestPasswordReset(c *gin.Context) {
	var req RequestPasswordResetRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	_, err := h.authService.RequestPasswordReset(req.Email)
	if err != nil {
		response.ErrorResponse(c, http.StatusInternalServerError, "Failed to request password reset")
		return
	}

	// In development, you might want to return the token
	// In production, you would typically not return the token
	response.SuccessResponse(c, http.StatusOK, "Password reset instructions sent to your email", nil)
}

type VerifyResetTokenRequest struct {
	Email string `json:"email" binding:"required,email"`
	Token string `json:"token" binding:"required"`
}

func (h *AuthHandler) VerifyResetToken(c *gin.Context) {
	var req VerifyResetTokenRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	valid, err := h.authService.VerifyResetToken(req.Email, req.Token)
	if err != nil || !valid {
		response.ErrorResponse(c, http.StatusBadRequest, "Invalid or expired token")
		return
	}

	response.SuccessResponse(c, http.StatusOK, "Token is valid", nil)
}

type ResetPasswordWithTokenRequest struct {
	Email       string `json:"email" binding:"required,email"`
	Token       string `json:"token" binding:"required"`
	NewPassword string `json:"new_password" binding:"required,min=6"`
}

func (h *AuthHandler) ResetPasswordWithToken(c *gin.Context) {
	var req ResetPasswordWithTokenRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	err := h.authService.ResetPasswordWithToken(req.Email, req.Token, req.NewPassword)
	if err != nil {
		response.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	response.SuccessResponse(c, http.StatusOK, "Password reset successfully", nil)
}

func (h *AuthHandler) Logout(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		response.ErrorResponse(c, http.StatusUnauthorized, "User not authenticated")
		return
	}

	if err := h.authService.Logout(userID.(uint)); err != nil {
		response.ErrorResponse(c, http.StatusInternalServerError, "Failed to logout")
		return
	}

	response.SuccessResponse(c, http.StatusOK, "Successfully logged out", nil)
}

// VerifyEmailRequest is the request body for email verification
type VerifyEmailRequest struct {
	Email string `json:"email" binding:"required,email"`
	Token string `json:"token" binding:"required"`
}

// VerifyEmail handles email verification
func (h *AuthHandler) VerifyEmail(c *gin.Context) {
	var req VerifyEmailRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	err := h.authService.VerifyEmail(req.Email, req.Token)
	if err != nil {
		response.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	response.SuccessResponse(c, http.StatusOK, "Email verified successfully", nil)
}

// ResendVerificationEmailRequest is the request body for resending verification email
type ResendVerificationEmailRequest struct {
	Email string `json:"email" binding:"required,email"`
}

// ResendVerificationEmail handles resending verification email
func (h *AuthHandler) ResendVerificationEmail(c *gin.Context) {
	var req ResendVerificationEmailRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	err := h.authService.ResendVerificationEmail(req.Email)
	if err != nil {
		response.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	response.SuccessResponse(c, http.StatusOK, "Verification email sent successfully", nil)
}
