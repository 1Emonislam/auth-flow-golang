package handlers

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"net/http"

	response "own-paynet/api/response"
	"own-paynet/config"
	"own-paynet/services"

	"github.com/gin-gonic/gin"
)

type PaymentHandler struct {
	paymentService *services.PaymentService
	webhookSecret  string
}

func NewPaymentHandler(paymentService *services.PaymentService, cfg *config.Config) *PaymentHandler {
	return &PaymentHandler{
		paymentService: paymentService,
		webhookSecret:  cfg.WebhookSecret,
	}
}

type CreatePaymentRequest struct {
	Amount         float64 `json:"amount" binding:"required,gt=0"`
	MerchantWallet string  `json:"merchant_wallet" binding:"required"`
	Currency       string  `json:"currency" binding:"required"`
}

func (h *PaymentHandler) CreatePayment(c *gin.Context) {
	var req CreatePaymentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	userID, exists := c.Get("user_id")
	if !exists {
		response.ErrorResponse(c, http.StatusUnauthorized, "User not authenticated")
		return
	}

	payment, err := h.paymentService.CreatePayment(userID.(uint), req.Amount, req.MerchantWallet, req.Currency)
	if err != nil {
		response.ErrorResponse(c, http.StatusInternalServerError, "Failed to create payment")
		return
	}

	paymentData := gin.H{
		"payment_id":      payment.PaymentID,
		"payment_url":     payment.PaymentURL,
		"bitcoin_address": payment.BitcoinAddress,
		"status":          payment.Status,
	}

	response.SuccessResponse(c, http.StatusOK, "Payment created successfully", paymentData)
}

type WebhookRequest struct {
	PaymentID string `json:"payment_id"`
	Status    string `json:"status"`
	Address   string `json:"address"`
}

func (h *PaymentHandler) HandleWebhook(c *gin.Context) {
	signature := c.GetHeader("X-Webhook-Signature")
	payload, err := c.GetRawData()
	if err != nil {
		response.ErrorResponse(c, http.StatusBadRequest, "Invalid payload")
		return
	}

	mac := hmac.New(sha256.New, []byte(h.webhookSecret))
	mac.Write(payload)
	expectedSignature := hex.EncodeToString(mac.Sum(nil))
	if signature != expectedSignature {
		response.ErrorResponse(c, http.StatusUnauthorized, "Invalid webhook signature")
		return
	}

	var webhook WebhookRequest
	if err := json.Unmarshal(payload, &webhook); err != nil {
		response.ErrorResponse(c, http.StatusBadRequest, "Invalid webhook payload")
		return
	}

	if err := h.paymentService.UpdatePaymentStatus(webhook.PaymentID, webhook.Status); err != nil {
		response.ErrorResponse(c, http.StatusInternalServerError, "Failed to update payment status")
		return
	}

	response.SuccessResponse(c, http.StatusOK, "Webhook processed successfully", nil)
}
