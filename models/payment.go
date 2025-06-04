package models

import (
	"gorm.io/gorm"
)

type Payment struct {
	gorm.Model
	PaymentID      string  `json:"payment_id" gorm:"unique"`
	UserID         uint    `json:"user_id"`
	User           User    `json:"user" gorm:"foreignKey:UserID"`
	Amount         float64 `json:"amount"`
	Currency       string  `json:"currency"`
	Status         string  `json:"status"`
	PaymentURL     string  `json:"payment_url"`
	BitcoinAddress string  `json:"bitcoin_address"`
	MerchantWallet string  `json:"merchant_wallet"`
	TransactionID  string  `json:"transaction_id" gorm:"index"`
	Confirmations  int64   `json:"confirmations"`
}
