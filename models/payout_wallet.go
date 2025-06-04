package models

import (
	"gorm.io/gorm"
)

type PayoutWallet struct {
	gorm.Model
	UserID        uint    `json:"user_id"`
	User          User    `json:"user" gorm:"foreignKey:UserID"`
	Currency      string  `json:"currency"`
	WalletAddress string  `json:"wallet_address" gorm:"uniqueIndex"`
	Balance       float64 `json:"balance"`
	IsDefault     bool    `json:"is_default" gorm:"default:false"`
}
