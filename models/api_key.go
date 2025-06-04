package models

import (
	"time"

	"gorm.io/gorm"
)

type APIKey struct {
	gorm.Model
	ID          uint      `json:"id" gorm:"primaryKey"`
	UserID      uint      `json:"user_id" gorm:"not null"`
	PublicKey   string    `json:"public_key" gorm:"uniqueIndex;not null"`
	SecretKey   string    `json:"secret_key" gorm:"uniqueIndex;not null"`
	IsDefault   bool      `json:"is_default" gorm:"default:false"`
	LastUsedAt  time.Time `json:"last_used_at"`
	Description string    `json:"description"`
}
