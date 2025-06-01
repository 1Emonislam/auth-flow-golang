package services

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
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

	// ✅ FIX: Add third argument (network parameters)
	if err := s.bitcoin.MonitorAddress(btcAddress, func(txID, status string) {
		_ = s.repo.UpdateStatus(paymentID, status)
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
