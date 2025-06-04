package handlers

import (
	"net/http"
	"strconv"

	response "own-paynet/api/response"
	"own-paynet/models"
	"own-paynet/services"

	"github.com/gin-gonic/gin"
)

type TransactionHandler struct {
	transactionService *services.TransactionService
	walletService      *services.PayoutWalletService
	userService        *services.UserService
}

func NewTransactionHandler(transactionService *services.TransactionService, walletService *services.PayoutWalletService, userService *services.UserService) *TransactionHandler {
	return &TransactionHandler{
		transactionService: transactionService,
		walletService:      walletService,
		userService:        userService,
	}
}

type CreateTransactionRequest struct {
	WalletID       uint    `json:"wallet_id" binding:"required"`
	Amount         float64 `json:"amount" binding:"required"`
	PriceCurrency  string  `json:"price_currency" binding:"required"`
	PayCurrency    string  `json:"pay_currency" binding:"required"`
	Comment        string  `json:"comment"`
	ReceiverWallet string  `json:"receiver_wallet" binding:"required"`
}

// CreateTransaction handles the creation of a new transaction
func (h *TransactionHandler) CreateTransaction(c *gin.Context) {
	var req CreateTransactionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	// Get the authenticated user's ID
	userID, exists := c.Get("user_id")
	if !exists {
		response.ErrorResponse(c, http.StatusUnauthorized, "User not authenticated")
		return
	}
	senderID := userID.(uint)

	// Verify the wallet belongs to the authenticated user
	wallet, err := h.walletService.GetPayoutWallet(req.WalletID)
	if err != nil {
		response.ErrorResponse(c, http.StatusNotFound, "Payout wallet not found")
		return
	}

	if wallet.UserID != senderID {
		response.ErrorResponse(c, http.StatusForbidden, "You don't have permission to perform transactions on this wallet")
		return
	}

	// Create the transaction
	transaction, err := h.transactionService.CreateTransaction(
		req.WalletID,
		models.TransactionTypeDebit, // Always debit for sender
		req.Amount,
		req.PriceCurrency,
		req.PayCurrency,
		req.Comment,
		senderID,
		req.ReceiverWallet,
	)
	if err != nil {
		response.ErrorResponse(c, http.StatusInternalServerError, "Failed to create transaction: "+err.Error())
		return
	}

	response.SuccessResponse(c, http.StatusCreated, "Transaction created successfully", transaction)
}

// GetTransaction handles retrieving a single transaction
func (h *TransactionHandler) GetTransaction(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		response.ErrorResponse(c, http.StatusBadRequest, "Invalid transaction ID")
		return
	}

	transaction, err := h.transactionService.GetTransaction(uint(id))
	if err != nil {
		response.ErrorResponse(c, http.StatusNotFound, "Transaction not found")
		return
	}

	// Verify the transaction belongs to the authenticated user
	userID, exists := c.Get("user_id")
	if !exists {
		response.ErrorResponse(c, http.StatusUnauthorized, "User not authenticated")
		return
	}
	authUserID := userID.(uint)

	// Check if user is the sender
	if transaction.SenderID != authUserID {
		// If not sender, check if user owns the receiver wallet
		receiverWallet, err := h.walletService.GetPayoutWalletByAddress(transaction.ReceiverWallet)
		if err != nil || receiverWallet.UserID != authUserID {
			response.ErrorResponse(c, http.StatusForbidden, "You don't have permission to access this transaction")
			return
		}
	}

	response.SuccessResponse(c, http.StatusOK, "Transaction retrieved successfully", transaction)
}

// GetWalletTransactions handles retrieving all transactions for a wallet
func (h *TransactionHandler) GetWalletTransactions(c *gin.Context) {
	walletIDStr := c.Param("wallet_id")
	walletID, err := strconv.ParseUint(walletIDStr, 10, 64)
	if err != nil {
		response.ErrorResponse(c, http.StatusBadRequest, "Invalid wallet ID")
		return
	}

	// Verify the wallet belongs to the authenticated user
	wallet, err := h.walletService.GetPayoutWallet(uint(walletID))
	if err != nil {
		response.ErrorResponse(c, http.StatusNotFound, "Payout wallet not found")
		return
	}

	userID, exists := c.Get("user_id")
	if !exists || wallet.UserID != userID.(uint) {
		response.ErrorResponse(c, http.StatusForbidden, "You don't have permission to access these transactions")
		return
	}

	transactions, err := h.transactionService.GetWalletTransactions(uint(walletID))
	if err != nil {
		response.ErrorResponse(c, http.StatusInternalServerError, "Failed to retrieve transactions")
		return
	}

	response.SuccessResponse(c, http.StatusOK, "Transactions retrieved successfully", transactions)
}

// GetUserTransactions handles retrieving all transactions for the authenticated user
func (h *TransactionHandler) GetUserTransactions(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		response.ErrorResponse(c, http.StatusUnauthorized, "User not authenticated")
		return
	}

	transactions, err := h.transactionService.GetUserTransactions(userID.(uint))
	if err != nil {
		response.ErrorResponse(c, http.StatusInternalServerError, "Failed to retrieve transactions")
		return
	}

	response.SuccessResponse(c, http.StatusOK, "Transactions retrieved successfully", transactions)
}
