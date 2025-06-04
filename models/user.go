package models

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Email           string     `json:"email" gorm:"unique"`
	Password        string     `json:"-"`
	EmailVerified   bool       `json:"email_verified" gorm:"default:false"`
	EmailVerifiedAt *time.Time `json:"email_verified_at"`
	CompanyID       uint       `json:"company_id"`
	Company         Company    `json:"company" gorm:"foreignKey:CompanyID"`

	// 2FA fields
	Email2FAEnabled         bool       `json:"email_2fa_enabled" gorm:"default:false"`
	Authenticator2FAEnabled bool       `json:"authenticator_2fa_enabled" gorm:"default:false"`
	TwoFactorSecret         string     `json:"two_factor_secret"`
	BackupCodes             string     `json:"-" gorm:"type:text"`
	LastOTPSentAt           *time.Time `json:"last_otp_sent_at"`
	OTPSecret               string     `json:"-"`

	// OAuth fields
	GoogleID string `json:"google_id" gorm:"unique"`
	Name     string `json:"name"`
	Avatar   string `json:"avatar"`
	Locale   string `json:"locale"`
}
