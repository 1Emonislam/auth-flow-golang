package models

import (
	"time"

	"gorm.io/gorm"
)

type Payment struct {
	gorm.Model
	PaymentID      string    `json:"payment_id" gorm:"unique"`
	UserID         uint      `json:"user_id"`
	Amount         float64   `json:"amount"`
	Currency       string    `json:"currency"`
	Status         string    `json:"status"`
	PaymentURL     string    `json:"payment_url"`
	BitcoinAddress string    `json:"bitcoin_address"`
	MerchantWallet string    `json:"merchant_wallet"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}
