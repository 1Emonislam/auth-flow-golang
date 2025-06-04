package services

import (
	"errors"
	"own-paynet/models"
	"own-paynet/repository"

	"github.com/btcsuite/btcd/btcec/v2"
	"github.com/btcsuite/btcd/btcutil"
	"github.com/btcsuite/btcd/chaincfg"
)

type PayoutWalletService struct {
	repo *repository.PayoutWalletRepository
}

func NewPayoutWalletService(repo *repository.PayoutWalletRepository) *PayoutWalletService {
	return &PayoutWalletService{repo: repo}
}

func (s *PayoutWalletService) generateBTCAddress() (string, error) {
	// Generate a new random private key
	privKey, err := btcec.NewPrivateKey()
	if err != nil {
		return "", err
	}

	// Get the public key
	pubKey := privKey.PubKey()

	// Create a P2PKH address
	addr, err := btcutil.NewAddressPubKey(pubKey.SerializeCompressed(), &chaincfg.MainNetParams)
	if err != nil {
		return "", err
	}

	return addr.EncodeAddress(), nil
}

// CreatePayoutWallet creates a new payout wallet for a user
func (s *PayoutWalletService) CreatePayoutWallet(userID uint, currency, walletAddress string, isDefault bool) (*models.PayoutWallet, error) {
	// Validate inputs
	if currency == "" {
		return nil, errors.New("currency is required")
	}

	// If no wallet address provided and currency is BTC, generate one
	if walletAddress == "" && currency == "BTC" {
		var err error
		walletAddress, err = s.generateBTCAddress()
		if err != nil {
			return nil, errors.New("failed to generate BTC address")
		}
	} else if walletAddress == "" {
		return nil, errors.New("wallet address is required for non-BTC currencies")
	}

	// If this is going to be the default wallet, unset any existing default
	if isDefault {
		err := s.repo.UnsetDefaultWallet(userID, currency)
		if err != nil {
			return nil, err
		}
	}

	// Create the wallet
	wallet := &models.PayoutWallet{
		UserID:        userID,
		Currency:      currency,
		WalletAddress: walletAddress,
		Balance:       0, // Initial balance is zero
		IsDefault:     isDefault,
	}

	err := s.repo.Create(wallet)
	if err != nil {
		return nil, err
	}

	wallet.User.Password = ""
	return wallet, nil
}

// CreateDefaultBTCWallet creates a default BTC wallet for a new user
func (s *PayoutWalletService) CreateDefaultBTCWallet(userID uint) (*models.PayoutWallet, error) {
	return s.CreatePayoutWallet(userID, "BTC", "", true)
}

// GetPayoutWallet retrieves a payout wallet by ID
func (s *PayoutWalletService) GetPayoutWallet(id uint) (*models.PayoutWallet, error) {
	return s.repo.FindByID(id)
}

// GetUserPayoutWallets retrieves all payout wallets for a user
func (s *PayoutWalletService) GetUserPayoutWallets(userID uint) ([]models.PayoutWallet, error) {
	return s.repo.FindByUserID(userID)
}

// GetDefaultWallet retrieves the default wallet for a specific currency
func (s *PayoutWalletService) GetDefaultWallet(userID uint, currency string) (*models.PayoutWallet, error) {
	return s.repo.FindDefaultWallet(userID, currency)
}

// UpdatePayoutWallet updates a payout wallet's address
func (s *PayoutWalletService) UpdatePayoutWallet(id uint, currency, walletAddress string, isDefault bool) (*models.PayoutWallet, error) {
	wallet, err := s.repo.FindByID(id)
	if err != nil {
		return nil, err
	}

	// If setting as default, unset any existing default
	if isDefault && !wallet.IsDefault {
		err := s.repo.UnsetDefaultWallet(wallet.UserID, wallet.Currency)
		if err != nil {
			return nil, err
		}
	}

	// Update fields
	if walletAddress != "" {
		wallet.WalletAddress = walletAddress
	}
	if currency != "" {
		wallet.Currency = currency
	}
	wallet.IsDefault = isDefault

	err = s.repo.Update(wallet)
	if err != nil {
		return nil, err
	}

	return wallet, nil
}

// SetDefaultWallet sets a wallet as the default for its currency
func (s *PayoutWalletService) SetDefaultWallet(id uint) error {
	wallet, err := s.repo.FindByID(id)
	if err != nil {
		return err
	}

	return s.repo.SetDefaultWallet(wallet.UserID, id, wallet.Currency)
}

// DeletePayoutWallet deletes a payout wallet
func (s *PayoutWalletService) DeletePayoutWallet(id uint) error {
	wallet, err := s.repo.FindByID(id)
	if err != nil {
		return err
	}

	if wallet.IsDefault {
		return errors.New("cannot delete the default wallet")
	}

	return s.repo.Delete(id)
}

// UpdateBalance updates the balance of a payout wallet
func (s *PayoutWalletService) UpdateBalance(id uint, newBalance float64) error {
	if newBalance < 0 {
		return errors.New("balance cannot be negative")
	}

	return s.repo.UpdateBalance(id, newBalance)
}

// AddFunds adds funds to a payout wallet
func (s *PayoutWalletService) AddFunds(id uint, amount float64) error {
	if amount <= 0 {
		return errors.New("amount must be positive")
	}

	wallet, err := s.repo.FindByID(id)
	if err != nil {
		return err
	}

	newBalance := wallet.Balance + amount
	return s.repo.UpdateBalance(id, newBalance)
}

// WithdrawFunds withdraws funds from a payout wallet
func (s *PayoutWalletService) WithdrawFunds(id uint, amount float64) error {
	if amount <= 0 {
		return errors.New("amount must be positive")
	}

	wallet, err := s.repo.FindByID(id)
	if err != nil {
		return err
	}

	if wallet.Balance < amount {
		return errors.New("insufficient funds")
	}

	newBalance := wallet.Balance - amount
	return s.repo.UpdateBalance(id, newBalance)
}

// GetPayoutWalletByAddress retrieves a payout wallet by its address
func (s *PayoutWalletService) GetPayoutWalletByAddress(walletAddress string) (*models.PayoutWallet, error) {
	return s.repo.FindByWalletAddress(walletAddress)
}
