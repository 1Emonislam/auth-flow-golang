package handlers

import (
	"net/http"
	"strconv"

	response "own-paynet/api/response"
	"own-paynet/services"

	"github.com/gin-gonic/gin"
)

type PayoutWalletHandler struct {
	payoutWalletService *services.PayoutWalletService
}

func NewPayoutWalletHandler(payoutWalletService *services.PayoutWalletService) *PayoutWalletHandler {
	return &PayoutWalletHandler{payoutWalletService: payoutWalletService}
}

type CreatePayoutWalletRequest struct {
	Currency      string `json:"currency" binding:"required"`
	WalletAddress string `json:"wallet_address" binding:"required"`
	IsDefault     bool   `json:"is_default"`
}

// CreatePayoutWallet handles the creation of a new payout wallet
func (h *PayoutWalletHandler) CreatePayoutWallet(c *gin.Context) {
	var req CreatePayoutWalletRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	userID, exists := c.Get("user_id")
	if !exists {
		response.ErrorResponse(c, http.StatusUnauthorized, "User not authenticated")
		return
	}

	// Check if wallet with same address already exists
	existingWallet, err := h.payoutWalletService.GetPayoutWalletByAddress(req.WalletAddress)
	if err == nil && existingWallet != nil {
		response.ErrorResponse(c, http.StatusBadRequest, "Wallet with this address already exists")
		return
	}

	// Check if user already has a wallet with this currency
	defaultWallet, err := h.payoutWalletService.GetDefaultWallet(userID.(uint), req.Currency)
	if err == nil && defaultWallet != nil {
		response.ErrorResponse(c, http.StatusBadRequest, "You already have a wallet for this currency")
		return
	}

	wallet, err := h.payoutWalletService.CreatePayoutWallet(userID.(uint), req.Currency, req.WalletAddress, req.IsDefault)
	if err != nil {
		response.ErrorResponse(c, http.StatusInternalServerError, "Failed to create payout wallet: "+err.Error())
		return
	}

	response.SuccessResponse(c, http.StatusCreated, "Payout wallet created successfully", wallet)
}

// GetPayoutWallet handles retrieving a single payout wallet
func (h *PayoutWalletHandler) GetPayoutWallet(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		response.ErrorResponse(c, http.StatusBadRequest, "Invalid wallet ID")
		return
	}

	wallet, err := h.payoutWalletService.GetPayoutWallet(uint(id))
	if err != nil {
		response.ErrorResponse(c, http.StatusNotFound, "Payout wallet not found")
		return
	}

	// Verify the wallet belongs to the authenticated user
	userID, exists := c.Get("user_id")
	if !exists || wallet.UserID != userID.(uint) {
		response.ErrorResponse(c, http.StatusForbidden, "You don't have permission to access this wallet")
		return
	}

	response.SuccessResponse(c, http.StatusOK, "Payout wallet retrieved successfully", wallet)
}

// GetUserPayoutWallets handles retrieving all payout wallets for a user
func (h *PayoutWalletHandler) GetUserPayoutWallets(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		response.ErrorResponse(c, http.StatusUnauthorized, "User not authenticated")
		return
	}

	wallets, err := h.payoutWalletService.GetUserPayoutWallets(userID.(uint))
	if err != nil {
		response.ErrorResponse(c, http.StatusInternalServerError, "Failed to retrieve payout wallets")
		return
	}

	response.SuccessResponse(c, http.StatusOK, "Payout wallets retrieved successfully", wallets)
}

type UpdatePayoutWalletRequest struct {
	Currency      string `json:"currency"`
	WalletAddress string `json:"wallet_address"`
	IsDefault     bool   `json:"is_default"`
}

// UpdatePayoutWallet handles updating a payout wallet
func (h *PayoutWalletHandler) UpdatePayoutWallet(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		response.ErrorResponse(c, http.StatusBadRequest, "Invalid wallet ID")
		return
	}

	// Verify the wallet belongs to the authenticated user
	wallet, err := h.payoutWalletService.GetPayoutWallet(uint(id))
	if err != nil {
		response.ErrorResponse(c, http.StatusNotFound, "Payout wallet not found")
		return
	}

	userID, exists := c.Get("user_id")
	if !exists || wallet.UserID != userID.(uint) {
		response.ErrorResponse(c, http.StatusForbidden, "You don't have permission to update this wallet")
		return
	}

	var req UpdatePayoutWalletRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	updatedWallet, err := h.payoutWalletService.UpdatePayoutWallet(uint(id), req.Currency, req.WalletAddress, req.IsDefault)
	if err != nil {
		response.ErrorResponse(c, http.StatusInternalServerError, "Failed to update payout wallet")
		return
	}

	response.SuccessResponse(c, http.StatusOK, "Payout wallet updated successfully", updatedWallet)
}

// DeletePayoutWallet handles deleting a payout wallet
func (h *PayoutWalletHandler) DeletePayoutWallet(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		response.ErrorResponse(c, http.StatusBadRequest, "Invalid wallet ID")
		return
	}

	// Verify the wallet belongs to the authenticated user
	wallet, err := h.payoutWalletService.GetPayoutWallet(uint(id))
	if err != nil {
		response.ErrorResponse(c, http.StatusNotFound, "Payout wallet not found")
		return
	}

	userID, exists := c.Get("user_id")
	if !exists || wallet.UserID != userID.(uint) {
		response.ErrorResponse(c, http.StatusForbidden, "You don't have permission to delete this wallet")
		return
	}

	err = h.payoutWalletService.DeletePayoutWallet(uint(id))
	if err != nil {
		response.ErrorResponse(c, http.StatusInternalServerError, "Failed to delete payout wallet")
		return
	}

	response.SuccessResponse(c, http.StatusOK, "Payout wallet deleted successfully", nil)
}
