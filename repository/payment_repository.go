package repository

import (
	"own-paynet/models"

	"gorm.io/gorm"
)

type PaymentRepository struct {
	db *gorm.DB
}

func NewPaymentRepository(db *gorm.DB) *PaymentRepository {
	return &PaymentRepository{db: db}
}

func (r *PaymentRepository) Create(payment *models.Payment) error {
	return r.db.Create(payment).Error
}

func (r *PaymentRepository) UpdateStatus(paymentID, status string) error {
	return r.db.Model(&models.Payment{}).Where("payment_id = ?", paymentID).Update("status", status).Error
}

func (r *PaymentRepository) UpdateTransaction(paymentID, txID string, confirmations int64) error {
	return r.db.Model(&models.Payment{}).Where("payment_id = ?", paymentID).Updates(map[string]interface{}{
		"transaction_id": txID,
		"confirmations":  confirmations,
	}).Error
}

func (r *PaymentRepository) FindByID(paymentID string) (*models.Payment, error) {
	var payment models.Payment
	err := r.db.Where("payment_id = ?", paymentID).First(&payment).Error
	return &payment, err
}

func (r *PaymentRepository) FindByTransactionID(txID string) (*models.Payment, error) {
	var payment models.Payment
	err := r.db.Where("transaction_id = ?", txID).First(&payment).Error
	return &payment, err
}
