package repository

import (
	"own-paynet/models"

	"gorm.io/gorm"
)

type TransactionRepository struct {
	db *gorm.DB
}

func NewTransactionRepository(db *gorm.DB) *TransactionRepository {
	return &TransactionRepository{db: db}
}

// Create creates a new transaction
func (r *TransactionRepository) Create(transaction *models.Transaction) error {
	return r.db.Create(transaction).Error
}

// FindByID retrieves a transaction by its ID
func (r *TransactionRepository) FindByID(id uint) (*models.Transaction, error) {
	var transaction models.Transaction
	err := r.db.First(&transaction, id).Error
	return &transaction, err
}

// FindByPayoutWalletID retrieves all transactions for a specific payout wallet
func (r *TransactionRepository) FindByPayoutWalletID(walletID uint) ([]models.Transaction, error) {
	var transactions []models.Transaction
	err := r.db.Where("payout_wallet_id = ?", walletID).Find(&transactions).Error
	return transactions, err
}

// FindByUserID retrieves all transactions where the user is either sender or receiver
func (r *TransactionRepository) FindByUserID(userID uint) ([]models.Transaction, error) {
	var transactions []models.Transaction
	err := r.db.Where("sender_id = ?", userID).Find(&transactions).Error
	return transactions, err
}

// Update updates a transaction
func (r *TransactionRepository) Update(transaction *models.Transaction) error {
	return r.db.Save(transaction).Error
}

// UpdateStatus updates the status of a transaction
func (r *TransactionRepository) UpdateStatus(id uint, status string) error {
	return r.db.Model(&models.Transaction{}).Where("id = ?", id).Update("status", status).Error
}

// MarkAsProcessed marks a transaction as processed
func (r *TransactionRepository) MarkAsProcessed(id uint) error {
	return r.db.Model(&models.Transaction{}).Where("id = ?", id).Update("is_processed", true).Error
}
