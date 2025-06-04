package handlers

import (
	"net/http"

	response "own-paynet/api/response"
	"own-paynet/services"

	"github.com/gin-gonic/gin"
)

type TwoFactorHandler struct {
	twoFactorService *services.TwoFactorService
}

func NewTwoFactorHandler(twoFactorService *services.TwoFactorService) *TwoFactorHandler {
	return &TwoFactorHandler{twoFactorService: twoFactorService}
}

type Verify2FARequest struct {
	Code string `json:"code" binding:"required"`
}

type VerifyAndEnable2FARequest struct {
	Code string `json:"code" binding:"required"`
}

// Verify2FA handles verifying a 2FA code
func (h *TwoFactorHandler) Verify2FA(c *gin.Context) {
	var req Verify2FARequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	userID, exists := c.Get("user_id")
	if !exists {
		response.ErrorResponse(c, http.StatusUnauthorized, "User not authenticated")
		return
	}

	valid, err := h.twoFactorService.Verify2FA(userID.(uint), req.Code)
	if err != nil {
		response.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	if !valid {
		response.ErrorResponse(c, http.StatusBadRequest, "Invalid 2FA code")
		return
	}

	response.SuccessResponse(c, http.StatusOK, "2FA code verified successfully", nil)
}

// SendOTP handles sending a new OTP to the user's email
func (h *TwoFactorHandler) SendOTP(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		response.ErrorResponse(c, http.StatusUnauthorized, "User not authenticated")
		return
	}

	err := h.twoFactorService.SendOTP(userID.(uint))
	if err != nil {
		response.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	response.SuccessResponse(c, http.StatusOK, "OTP sent successfully", nil)
}

// GetAuthenticatorQRCode handles getting the QR code for authenticator app setup
func (h *TwoFactorHandler) GetAuthenticatorQRCode(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		response.ErrorResponse(c, http.StatusUnauthorized, "User not authenticated")
		return
	}

	qrCode, err := h.twoFactorService.GetAuthenticatorQRCode(userID.(uint))
	if err != nil {
		response.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	response.SuccessResponse(c, http.StatusOK, "QR code generated successfully", gin.H{
		"qr_code": qrCode,
	})
}

// EnableEmail2FA handles enabling email 2FA for a user
func (h *TwoFactorHandler) EnableEmail2FA(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		response.ErrorResponse(c, http.StatusUnauthorized, "User not authenticated")
		return
	}
	_, err := h.twoFactorService.EnableEmail2FA(userID.(uint))
	if err != nil {
		response.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.SuccessResponse(c, http.StatusOK, "We've sent a 2FA verification code to your email. Please check your inbox and enter the code to complete the verification process.", nil)
}

// EnableAuthenticator2FA handles enabling authenticator 2FA for a user
func (h *TwoFactorHandler) EnableAuthenticator2FA(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		response.ErrorResponse(c, http.StatusUnauthorized, "User not authenticated")
		return
	}
	_, err := h.twoFactorService.EnableAuthenticator2FA(userID.(uint))
	if err != nil {
		response.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.SuccessResponse(c, http.StatusOK, "We've sent a 2FA verification code to your email. Please check your inbox and enter the code to complete the verification process.", nil)
}

// DisableEmail2FA handles initiating email 2FA disable process
func (h *TwoFactorHandler) DisableEmail2FA(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		response.ErrorResponse(c, http.StatusUnauthorized, "User not authenticated")
		return
	}
	_, err := h.twoFactorService.DisableEmail2FA(userID.(uint))
	if err != nil {
		response.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.SuccessResponse(c, http.StatusOK, "We've sent a verification code to your email. Please enter the code to disable 2FA.", nil)
}

// DisableAuthenticator2FA handles initiating authenticator 2FA disable process
func (h *TwoFactorHandler) DisableAuthenticator2FA(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		response.ErrorResponse(c, http.StatusUnauthorized, "User not authenticated")
		return
	}
	_, err := h.twoFactorService.DisableAuthenticator2FA(userID.(uint))
	if err != nil {
		response.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.SuccessResponse(c, http.StatusOK, "We've sent a verification code to your email. Please enter the code to disable 2FA.", nil)
}

// VerifyAndEnableEmail2FA handles verifying OTP and enabling email 2FA
func (h *TwoFactorHandler) VerifyAndEnableEmail2FA(c *gin.Context) {
	var req VerifyAndEnable2FARequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	userID, exists := c.Get("user_id")
	if !exists {
		response.ErrorResponse(c, http.StatusUnauthorized, "User not authenticated")
		return
	}

	_, err := h.twoFactorService.VerifyAndEnableEmail2FA(userID.(uint), req.Code)
	if err != nil {
		response.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	response.SuccessResponse(c, http.StatusOK, "Email 2FA enabled successfully", nil)
}

// VerifyAndEnableAuthenticator2FA handles verifying OTP and enabling authenticator 2FA
func (h *TwoFactorHandler) VerifyAndEnableAuthenticator2FA(c *gin.Context) {
	var req VerifyAndEnable2FARequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	userID, exists := c.Get("user_id")
	if !exists {
		response.ErrorResponse(c, http.StatusUnauthorized, "User not authenticated")
		return
	}

	_, err := h.twoFactorService.VerifyAndEnableAuthenticator2FA(userID.(uint), req.Code)
	if err != nil {
		response.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	response.SuccessResponse(c, http.StatusOK, "Please open your authenticator app and enter the verification code to complete the two-factor authentication setup.", nil)
}

// VerifyAndDisableEmail2FA handles verifying OTP and disabling email 2FA
func (h *TwoFactorHandler) VerifyAndDisableEmail2FA(c *gin.Context) {
	var req VerifyAndEnable2FARequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	userID, exists := c.Get("user_id")
	if !exists {
		response.ErrorResponse(c, http.StatusUnauthorized, "User not authenticated")
		return
	}

	err := h.twoFactorService.VerifyAndDisableEmail2FA(userID.(uint), req.Code)
	if err != nil {
		response.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	response.SuccessResponse(c, http.StatusOK, "Email 2FA disabled successfully", nil)
}

// VerifyAndDisableAuthenticator2FA handles verifying OTP and disabling authenticator 2FA
func (h *TwoFactorHandler) VerifyAndDisableAuthenticator2FA(c *gin.Context) {
	var req VerifyAndEnable2FARequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	userID, exists := c.Get("user_id")
	if !exists {
		response.ErrorResponse(c, http.StatusUnauthorized, "User not authenticated")
		return
	}

	err := h.twoFactorService.VerifyAndDisableAuthenticator2FA(userID.(uint), req.Code)
	if err != nil {
		response.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	response.SuccessResponse(c, http.StatusOK, "Authenticator 2FA disabled successfully", nil)
}
