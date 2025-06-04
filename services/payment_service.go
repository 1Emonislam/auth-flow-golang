package services

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"log"
	"own-paynet/models"
	"own-paynet/repository"
	"own-paynet/services/bitcoin"

	"github.com/btcsuite/btcd/chaincfg"
)

type PaymentService struct {
	repo      *repository.PaymentRepository
	bitcoin   *bitcoin.BitcoinService
	baseURL   string
	netParams *chaincfg.Params
}

// Add network params in the constructor
func NewPaymentService(repo *repository.PaymentRepository, bitcoin *bitcoin.BitcoinService, baseURL string, network string) *PaymentService {
	return &PaymentService{
		repo:      repo,
		bitcoin:   bitcoin,
		baseURL:   baseURL,
		netParams: getNetParams(network),
	}
}

// Utility to select correct Bitcoin network
func getNetParams(network string) *chaincfg.Params {
	switch network {
	case "mainnet":
		return &chaincfg.MainNetParams
	case "testnet":
		return &chaincfg.TestNet3Params
	case "regtest":
		return &chaincfg.RegressionNetParams
	default:
		return &chaincfg.TestNet3Params
	}
}

func (s *PaymentService) CreatePayment(userID uint, amount float64, merchantWallet, currency string) (*models.Payment, error) {
	bytes := make([]byte, 16)
	if _, err := rand.Read(bytes); err != nil {
		return nil, err
	}
	paymentID := hex.EncodeToString(bytes)

	btcAddress, err := s.bitcoin.GenerateAddress()
	if err != nil {
		return nil, err
	}

	paymentURL := fmt.Sprintf("%s/pay/%s", s.baseURL, paymentID)

	// Monitor the address for transactions with enhanced confirmation handling
	if err := s.bitcoin.MonitorAddress(btcAddress, func(txID, status string, confirmations int64) {
		// Update payment status with confirmation count
		_, err := s.repo.FindByID(paymentID)
		if err != nil {
			log.Printf("Failed to find payment %s: %v", paymentID, err)
			return
		}

		// Update status based on confirmations
		if status == "confirmed" {
			_ = s.repo.UpdateStatus(paymentID, "confirmed")
			// Here you could trigger additional business logic for confirmed payments
		} else if status == "pending_confirmation" {
			_ = s.repo.UpdateStatus(paymentID, fmt.Sprintf("pending_confirmation (%d/6)", confirmations))
		} else {
			_ = s.repo.UpdateStatus(paymentID, status)
		}
	}, s.netParams); err != nil {
		return nil, err
	}

	payment := &models.Payment{
		PaymentID:      paymentID,
		UserID:         userID,
		Amount:         amount,
		Currency:       currency,
		Status:         "waiting",
		PaymentURL:     paymentURL,
		BitcoinAddress: btcAddress,
		MerchantWallet: merchantWallet,
	}

	if err := s.repo.Create(payment); err != nil {
		return nil, err
	}

	return payment, nil
}

func (s *PaymentService) UpdatePaymentStatus(paymentID, status string) error {
	return s.repo.UpdateStatus(paymentID, status)
}

// GetPaymentStatus returns the current status and confirmation count of a payment
func (s *PaymentService) GetPaymentStatus(paymentID string) (string, int64, error) {
	payment, err := s.repo.FindByID(paymentID)
	if err != nil {
		return "", 0, err
	}

	// If the payment has a transaction ID, get its confirmations
	if payment.Status == "pending_confirmation" || payment.Status == "confirmed" {
		confirmations, err := s.bitcoin.GetTransactionConfirmations(paymentID)
		if err != nil {
			return payment.Status, 0, err
		}
		return payment.Status, confirmations, nil
	}

	return payment.Status, 0, nil
}
