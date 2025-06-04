package models

import (
	"gorm.io/gorm"
)

type TransactionType string

const (
	TransactionTypeDebit  TransactionType = "debit"
	TransactionTypeCredit TransactionType = "credit"
)

type Transaction struct {
	gorm.Model
	PayoutWalletID uint            `json:"payout_wallet_id"`
	PayoutWallet   PayoutWallet    `json:"payout_wallet" gorm:"foreignKey:PayoutWalletID"`
	Type           TransactionType `json:"type"`
	Amount         float64         `json:"amount"`
	PriceCurrency  string          `json:"price_currency"`
	PayCurrency    string          `json:"pay_currency"`
	Comment        string          `json:"comment"`
	Status         string          `json:"status"` // pending, completed, failed
	SenderID       uint            `json:"sender_id"`
	Sender         User            `json:"sender" gorm:"foreignKey:SenderID"`
	ReceiverWallet string          `json:"receiver_wallet"` // Store receiver's wallet address
	TransactionID  string          `json:"transaction_id" gorm:"index"`
	IsProcessed    bool            `json:"is_processed" gorm:"default:false"` // Track if transaction has been processed
}
