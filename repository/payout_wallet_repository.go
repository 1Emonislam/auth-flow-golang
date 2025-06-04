package repository

import (
	"own-paynet/models"

	"gorm.io/gorm"
)

type PayoutWalletRepository struct {
	db *gorm.DB
}

func NewPayoutWalletRepository(db *gorm.DB) *PayoutWalletRepository {
	return &PayoutWalletRepository{db: db}
}

// Create a new payout wallet
func (r *PayoutWalletRepository) Create(wallet *models.PayoutWallet) error {
	return r.db.Preload("User").Preload("User.Company").Create(wallet).Error
}

// FindByID retrieves a payout wallet by its ID
func (r *PayoutWalletRepository) FindByID(id uint) (*models.PayoutWallet, error) {
	var wallet models.PayoutWallet
	err := r.db.Preload("User").Preload("User.Company").First(&wallet, id).Error
	if err != nil {
		return nil, err
	}
	wallet.User.Password = ""
	return &wallet, nil
}

// FindByUserID retrieves all payout wallets for a specific user
func (r *PayoutWalletRepository) FindByUserID(userID uint) ([]models.PayoutWallet, error) {
	var wallets []models.PayoutWallet
	err := r.db.Preload("User").Preload("User.Company").Where("user_id = ?", userID).Find(&wallets).Error
	if err != nil {
		return nil, err
	}

	// Clear sensitive data from all wallets
	for i := range wallets {
		wallets[i].User.Password = ""
	}

	return wallets, nil
}

// Update updates a payout wallet
func (r *PayoutWalletRepository) Update(wallet *models.PayoutWallet) error {
	err := r.db.Preload("User").Preload("User.Company").Save(wallet).Error
	if err != nil {
		return err
	}
	wallet.User.Password = ""
	return nil
}

// Delete deletes a payout wallet
func (r *PayoutWalletRepository) Delete(id uint) error {
	return r.db.Delete(&models.PayoutWallet{}, id).Error
}

// UpdateBalance updates the balance of a payout wallet
func (r *PayoutWalletRepository) UpdateBalance(id uint, amount float64) error {
	var wallet models.PayoutWallet
	err := r.db.Model(&models.PayoutWallet{}).Where("id = ?", id).Update("balance", amount).Error
	if err != nil {
		return err
	}

	err = r.db.Preload("User").Preload("User.Company").First(&wallet, id).Error
	if err != nil {
		return err
	}

	wallet.User.Password = ""
	return nil
}

// FindByWalletAddress finds a wallet by its address
func (r *PayoutWalletRepository) FindByWalletAddress(walletAddress string) (*models.PayoutWallet, error) {
	var wallet models.PayoutWallet
	err := r.db.Where("wallet_address = ?", walletAddress).First(&wallet).Error
	if err != nil {
		return nil, err
	}
	return &wallet, nil
}

// FindDefaultWallet finds the default wallet for a user and currency
func (r *PayoutWalletRepository) FindDefaultWallet(userID uint, currency string) (*models.PayoutWallet, error) {
	var wallet models.PayoutWallet
	err := r.db.Where("user_id = ? AND currency = ? AND is_default = ?", userID, currency, true).First(&wallet).Error
	if err != nil {
		return nil, err
	}
	return &wallet, nil
}

// UnsetDefaultWallet unsets the default wallet for a user and currency
func (r *PayoutWalletRepository) UnsetDefaultWallet(userID uint, currency string) error {
	return r.db.Model(&models.PayoutWallet{}).
		Where("user_id = ? AND currency = ? AND is_default = ?", userID, currency, true).
		Update("is_default", false).Error
}

// SetDefaultWallet sets a wallet as the default for its currency
func (r *PayoutWalletRepository) SetDefaultWallet(userID uint, walletID uint, currency string) error {
	// Start a transaction
	tx := r.db.Begin()

	// Unset any existing default wallet
	if err := tx.Model(&models.PayoutWallet{}).
		Where("user_id = ? AND currency = ? AND is_default = ?", userID, currency, true).
		Update("is_default", false).Error; err != nil {
		tx.Rollback()
		return err
	}

	// Set the new default wallet
	if err := tx.Model(&models.PayoutWallet{}).
		Where("id = ? AND user_id = ?", walletID, userID).
		Update("is_default", true).Error; err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit().Error
}
