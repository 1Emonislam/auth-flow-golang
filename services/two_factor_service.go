package services

import (
	"crypto/rand"
	"encoding/base32"
	"fmt"
	"math/big"
	"own-paynet/models"
	"own-paynet/repository"
	"strings"
	"time"

	"github.com/pquerna/otp/totp"
)

type TwoFactorService struct {
	userRepo     *repository.UserRepository
	emailService *EmailService
}

func NewTwoFactorService(userRepo *repository.UserRepository, emailService *EmailService) *TwoFactorService {
	return &TwoFactorService{
		userRepo:     userRepo,
		emailService: emailService,
	}
}

// GenerateOTP generates a new OTP for email verification
func (s *TwoFactorService) GenerateOTP() (string, error) {
	// Generate a 6-digit OTP
	otp := ""
	for i := 0; i < 6; i++ {
		num, err := rand.Int(rand.Reader, big.NewInt(10))
		if err != nil {
			return "", err
		}
		otp += fmt.Sprintf("%d", num)
	}
	return otp, nil
}

// GenerateBackupCodes generates 8 backup codes
func (s *TwoFactorService) GenerateBackupCodes() ([]string, error) {
	codes := make([]string, 8)
	for i := 0; i < 8; i++ {
		// Generate 10 random bytes
		bytes := make([]byte, 10)
		if _, err := rand.Read(bytes); err != nil {
			return nil, err
		}
		// Convert to base32 and take first 8 characters
		code := base32.StdEncoding.EncodeToString(bytes)[:8]
		codes[i] = code
	}
	return codes, nil
}

// EnableEmail2FA generates and sends OTP for email 2FA verification
func (s *TwoFactorService) EnableEmail2FA(userID uint) (*models.User, error) {
	user, err := s.userRepo.FindByID(userID)
	if err != nil {
		return nil, err
	}
	if user.Email2FAEnabled {
		return nil, fmt.Errorf("email 2FA is already enabled")
	}
	otp, err := s.GenerateOTP()
	if err != nil {
		return nil, err
	}
	now := time.Now()
	user.OTPSecret = otp
	user.LastOTPSentAt = &now
	err = s.userRepo.Update(user)
	if err != nil {
		return nil, err
	}
	err = s.emailService.SendOTPEmail(user.Email, otp)
	if err != nil {
		return nil, fmt.Errorf("failed to send initial OTP: %w", err)
	}
	return user, nil
}

// VerifyAndEnableEmail2FA verifies OTP and enables email 2FA
func (s *TwoFactorService) VerifyAndEnableEmail2FA(userID uint, otp string) (*models.User, error) {
	user, err := s.userRepo.FindByID(userID)
	if err != nil {
		return nil, err
	}
	if user.Email2FAEnabled {
		return nil, fmt.Errorf("email 2FA is already enabled")
	}
	if user.OTPSecret != otp {
		return nil, fmt.Errorf("invalid OTP")
	}
	if user.LastOTPSentAt != nil && time.Since(*user.LastOTPSentAt) > 5*time.Minute {
		return nil, fmt.Errorf("OTP expired")
	}
	user.Email2FAEnabled = true
	backupCodes, err := s.GenerateBackupCodes()
	if err != nil {
		return nil, err
	}
	user.BackupCodes = strings.Join(backupCodes, ",")
	err = s.userRepo.Update(user)
	if err != nil {
		return nil, err
	}
	return user, nil
}

// EnableAuthenticator2FA generates and sends OTP for authenticator 2FA verification
func (s *TwoFactorService) EnableAuthenticator2FA(userID uint) (*models.User, error) {
	user, err := s.userRepo.FindByID(userID)
	if err != nil {
		return nil, err
	}
	if user.Authenticator2FAEnabled {
		return nil, fmt.Errorf("authenticator 2FA is already enabled")
	}
	secret, err := totp.Generate(totp.GenerateOpts{
		Issuer:      "OwnPaynet",
		AccountName: user.Email,
	})
	if err != nil {
		return nil, err
	}
	user.TwoFactorSecret = secret.Secret()
	otp, err := s.GenerateOTP()
	if err != nil {
		return nil, err
	}
	now := time.Now()
	user.OTPSecret = otp
	user.LastOTPSentAt = &now
	err = s.userRepo.Update(user)
	if err != nil {
		return nil, err
	}
	err = s.emailService.SendOTPEmail(user.Email, otp)
	if err != nil {
		return nil, fmt.Errorf("failed to send initial OTP: %w", err)
	}
	return user, nil
}

// VerifyAndEnableAuthenticator2FA verifies OTP and enables authenticator 2FA
func (s *TwoFactorService) VerifyAndEnableAuthenticator2FA(userID uint, otp string) (*models.User, error) {
	user, err := s.userRepo.FindByID(userID)
	if err != nil {
		return nil, err
	}
	if user.Authenticator2FAEnabled {
		return nil, fmt.Errorf("authenticator 2FA is already enabled")
	}
	if user.OTPSecret != otp {
		return nil, fmt.Errorf("invalid OTP")
	}
	if user.LastOTPSentAt != nil && time.Since(*user.LastOTPSentAt) > 5*time.Minute {
		return nil, fmt.Errorf("OTP expired")
	}
	user.Authenticator2FAEnabled = true
	backupCodes, err := s.GenerateBackupCodes()
	if err != nil {
		return nil, err
	}
	user.BackupCodes = strings.Join(backupCodes, ",")
	err = s.userRepo.Update(user)
	if err != nil {
		return nil, err
	}
	return user, nil
}

// DisableEmail2FA generates and sends OTP for email 2FA verification before disabling
func (s *TwoFactorService) DisableEmail2FA(userID uint) (*models.User, error) {
	user, err := s.userRepo.FindByID(userID)
	if err != nil {
		return nil, err
	}
	if !user.Email2FAEnabled {
		return nil, fmt.Errorf("email 2FA is not enabled")
	}
	otp, err := s.GenerateOTP()
	if err != nil {
		return nil, err
	}
	now := time.Now()
	user.OTPSecret = otp
	user.LastOTPSentAt = &now
	err = s.userRepo.Update(user)
	if err != nil {
		return nil, err
	}
	err = s.emailService.SendOTPEmail(user.Email, otp)
	if err != nil {
		return nil, fmt.Errorf("failed to send initial OTP: %w", err)
	}
	return user, nil
}

// VerifyAndDisableEmail2FA verifies OTP and disables email 2FA
func (s *TwoFactorService) VerifyAndDisableEmail2FA(userID uint, otp string) error {
	user, err := s.userRepo.FindByID(userID)
	if err != nil {
		return err
	}
	if !user.Email2FAEnabled {
		return fmt.Errorf("email 2FA is not enabled")
	}
	if user.OTPSecret != otp {
		return fmt.Errorf("invalid OTP")
	}
	if user.LastOTPSentAt != nil && time.Since(*user.LastOTPSentAt) > 5*time.Minute {
		return fmt.Errorf("OTP expired")
	}
	user.Email2FAEnabled = false
	user.OTPSecret = ""
	user.LastOTPSentAt = nil
	if !user.Authenticator2FAEnabled {
		user.BackupCodes = ""
	}
	return s.userRepo.Update(user)
}

// DisableAuthenticator2FA generates and sends OTP for authenticator 2FA verification before disabling
func (s *TwoFactorService) DisableAuthenticator2FA(userID uint) (*models.User, error) {
	user, err := s.userRepo.FindByID(userID)
	if err != nil {
		return nil, err
	}
	if !user.Authenticator2FAEnabled {
		return nil, fmt.Errorf("authenticator 2FA is not enabled")
	}
	otp, err := s.GenerateOTP()
	if err != nil {
		return nil, err
	}
	now := time.Now()
	user.OTPSecret = otp
	user.LastOTPSentAt = &now
	err = s.userRepo.Update(user)
	if err != nil {
		return nil, err
	}
	err = s.emailService.SendOTPEmail(user.Email, otp)
	if err != nil {
		return nil, fmt.Errorf("failed to send initial OTP: %w", err)
	}
	return user, nil
}

// VerifyAndDisableAuthenticator2FA verifies OTP and disables authenticator 2FA
func (s *TwoFactorService) VerifyAndDisableAuthenticator2FA(userID uint, otp string) error {
	user, err := s.userRepo.FindByID(userID)
	if err != nil {
		return err
	}
	if !user.Authenticator2FAEnabled {
		return fmt.Errorf("authenticator 2FA is not enabled")
	}
	if user.OTPSecret != otp {
		return fmt.Errorf("invalid OTP")
	}
	if user.LastOTPSentAt != nil && time.Since(*user.LastOTPSentAt) > 5*time.Minute {
		return fmt.Errorf("OTP expired")
	}
	user.Authenticator2FAEnabled = false
	user.TwoFactorSecret = ""
	if !user.Email2FAEnabled {
		user.BackupCodes = ""
	}
	return s.userRepo.Update(user)
}

// Verify2FA verifies the 2FA code (checks both methods if both enabled)
func (s *TwoFactorService) Verify2FA(userID uint, code string) (bool, error) {
	user, err := s.userRepo.FindByID(userID)
	if err != nil {
		return false, err
	}
	if !user.Email2FAEnabled && !user.Authenticator2FAEnabled {
		return false, fmt.Errorf("2FA is not enabled")
	}
	// Check if it's a backup code
	backupCodes := strings.Split(user.BackupCodes, ",")
	for i, backupCode := range backupCodes {
		if backupCode == code {
			// Remove used backup code
			backupCodes = append(backupCodes[:i], backupCodes[i+1:]...)
			user.BackupCodes = strings.Join(backupCodes, ",")
			err = s.userRepo.Update(user)
			if err != nil {
				return false, err
			}
			return true, nil
		}
	}
	// Check authenticator
	if user.Authenticator2FAEnabled && totp.Validate(code, user.TwoFactorSecret) {
		return true, nil
	}
	// Check email OTP
	if user.Email2FAEnabled && user.OTPSecret == code {
		if user.LastOTPSentAt != nil && time.Since(*user.LastOTPSentAt) > 5*time.Minute {
			return false, fmt.Errorf("OTP expired")
		}
		return true, nil
	}
	return false, nil
}

// SendOTP sends a new OTP to the user's email
func (s *TwoFactorService) SendOTP(userID uint) error {
	user, err := s.userRepo.FindByID(userID)
	if err != nil {
		return err
	}
	if !user.Email2FAEnabled {
		return fmt.Errorf("email 2FA is not enabled")
	}
	if user.LastOTPSentAt != nil && time.Since(*user.LastOTPSentAt) < 1*time.Minute {
		return fmt.Errorf("please wait before requesting a new OTP")
	}
	otp, err := s.GenerateOTP()
	if err != nil {
		return err
	}
	now := time.Now()
	user.OTPSecret = otp
	user.LastOTPSentAt = &now
	err = s.userRepo.Update(user)
	if err != nil {
		return err
	}
	return s.emailService.SendOTPEmail(user.Email, otp)
}

// GetAuthenticatorQRCode returns the QR code data for authenticator app setup
func (s *TwoFactorService) GetAuthenticatorQRCode(userID uint) (string, error) {
	user, err := s.userRepo.FindByID(userID)
	if err != nil {
		return "", err
	}
	if !user.Authenticator2FAEnabled {
		return "", fmt.Errorf("authenticator 2FA is not enabled")
	}
	secret, err := totp.Generate(totp.GenerateOpts{
		Issuer:      "OwnPaynet",
		AccountName: user.Email,
		Secret:      []byte(user.TwoFactorSecret),
	})
	if err != nil {
		return "", err
	}
	return secret.Secret(), nil
}
