package services

import (
	"context"
	"errors"
	"own-paynet/database"
	"own-paynet/models"
	"own-paynet/repository"
	"own-paynet/utils"
	"own-paynet/utils/email"
	"time"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type AuthService struct {
	repo                *repository.UserRepository
	emailService        *email.EmailService
	apiKeyService       *APIKeyService
	payoutWalletService *PayoutWalletService
}

func NewAuthService(repo *repository.UserRepository, emailService *email.EmailService, apiKeyService *APIKeyService, payoutWalletService *PayoutWalletService) *AuthService {
	return &AuthService{
		repo:                repo,
		emailService:        emailService,
		apiKeyService:       apiKeyService,
		payoutWalletService: payoutWalletService,
	}
}

func (s *AuthService) Signup(email, password string) error {
	// Check if user already exists
	existingUser, err := s.repo.FindByEmail(email)
	if err == nil && existingUser != nil {
		return errors.New("an account with this email already exists")
	}

	// Create company first
	company := &models.Company{}
	err = database.DB.Create(company).Error
	if err != nil {
		return errors.New("unable to create company at this time, please try again later")
	}

	// Create user with company relationship
	user := &models.User{
		Email:     email,
		Password:  password,
		CompanyID: company.ID,
	}
	err = s.repo.Create(user)
	if err != nil {
		// Handle database-specific errors with user-friendly messages
		if err.Error() == "duplicate key value violates unique constraint" {
			return errors.New("an account with this email already exists")
		}
		return errors.New("unable to create account at this time, please try again later")
	}

	// Generate default API key for the new user
	_, err = s.apiKeyService.GenerateAPIKey(user.ID, "Default API Key")
	if err != nil {
		// Log the error but don't fail the signup process
		// The user can generate a new API key later
		// TODO: Add proper logging here
	}

	// Create default BTC wallet for the new user
	_, err = s.payoutWalletService.CreateDefaultBTCWallet(user.ID)
	if err != nil {
		// Log the error but don't fail the signup process
		// The user can create a new wallet later
		// TODO: Add proper logging here
	}

	// Generate verification token
	token := utils.GenerateRandomToken(32)

	// Store token in Redis with 24 hours expiry
	ctx := context.Background()
	err = database.StoreEmailVerificationToken(ctx, email, token, 24*time.Hour)
	if err != nil {
		return errors.New("unable to generate verification token, please try again later")
	}

	// Send verification email
	err = s.emailService.SendVerificationEmail(email, token)
	if err != nil {
		return errors.New("account created but unable to send verification email, please request a new verification email")
	}

	return nil
}
func (s *AuthService) Signin(email, password string) (string, *models.User, error) {
	// Check if user exists
	user, err := s.repo.FindByEmail(email)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return "", nil, errors.New("user not found")
		}
		return "", nil, errors.New("failed to find user")
	}

	// Verify password
	if errCompare := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); errCompare != nil {
		return "", nil, errors.New("invalid credentials")
	}

	// Check if email is verified
	verified, err := s.repo.IsEmailVerified(email)
	if err != nil {
		return "", nil, errors.New("failed to check email verification status")
	}

	if !verified {
		// Generate new verification token
		token := utils.GenerateRandomToken(32)

		// Store token in Redis with 24 hours expiry
		ctx := context.Background()
		err = database.StoreEmailVerificationToken(ctx, email, token, 24*time.Hour)
		if err != nil {
			return "", nil, errors.New("failed to generate verification token")
		}

		// Send verification email
		err = s.emailService.SendVerificationEmail(email, token)
		if err != nil {
			return "", nil, errors.New("failed to send verification email")
		}

		return "", nil, errors.New("email not verified - please check your inbox for a new verification email and follow the instructions to verify your account")
	}

	// Generate JWT token
	token, err := utils.GenerateJWT(user.ID)
	if err != nil {
		return "", nil, errors.New("failed to generate authentication token")
	}

	return token, user, nil
}

func (s *AuthService) ResetPassword(email, newPassword string) error {
	return s.repo.UpdatePassword(email, newPassword)
}

// Logout invalidates a user's token
func (s *AuthService) Logout(userID uint) error {
	return utils.InvalidateToken(userID)
}

// RequestPasswordReset generates a reset token and stores it in Redis
func (s *AuthService) RequestPasswordReset(email string) (string, error) {
	// Check if user exists
	_, err := s.repo.FindByEmail(email)
	if err != nil {
		return "", errors.New("user not found")
	}

	// Generate a random token
	token := utils.GenerateRandomToken(32)

	// Store token in Redis with 15 minutes expiry
	ctx := context.Background()
	err = database.StorePasswordResetToken(ctx, email, token, 15*time.Minute)
	if err != nil {
		return "", err
	}

	// Send password reset email
	err = s.emailService.SendPasswordResetEmail(email, token)
	if err != nil {
		return "", err
	}

	return token, nil
}

// VerifyResetToken verifies a password reset token
func (s *AuthService) VerifyResetToken(email, token string) (bool, error) {
	ctx := context.Background()
	return database.ValidatePasswordResetToken(ctx, email, token)
}

// ResetPasswordWithToken resets a password after verifying the token
func (s *AuthService) ResetPasswordWithToken(email, token, newPassword string) error {
	// Verify token
	valid, err := s.VerifyResetToken(email, token)
	if err != nil {
		return err
	}

	if !valid {
		return errors.New("invalid or expired token")
	}

	// Reset password
	err = s.repo.UpdatePassword(email, newPassword)
	if err != nil {
		return err
	}

	// Delete token after successful password reset
	ctx := context.Background()
	return database.DeletePasswordResetToken(ctx, email)
}

// VerifyEmail verifies a user's email address
func (s *AuthService) VerifyEmail(email, token string) error {
	// Verify token
	ctx := context.Background()
	valid, err := database.ValidateEmailVerificationToken(ctx, email, token)
	if err != nil {
		return err
	}

	if !valid {
		return errors.New("invalid or expired token")
	}

	// Mark email as verified
	err = s.repo.VerifyEmail(email)
	if err != nil {
		return err
	}

	// Delete token after successful verification
	err = database.DeleteEmailVerificationToken(ctx, email)
	if err != nil {
		return err
	}

	// Send welcome email
	return s.emailService.SendWelcomeEmail(email)
}

// ResendVerificationEmail resends the verification email
func (s *AuthService) ResendVerificationEmail(email string) error {
	// Check if user exists
	_, err := s.repo.FindByEmail(email)
	if err != nil {
		return errors.New("user not found")
	}

	// Check if email is already verified
	verified, err := s.repo.IsEmailVerified(email)
	if err != nil {
		return err
	}

	if verified {
		return errors.New("email already verified")
	}

	// Generate new verification token
	token := utils.GenerateRandomToken(32)

	// Store token in Redis with 24 hours expiry
	ctx := context.Background()
	err = database.StoreEmailVerificationToken(ctx, email, token, 24*time.Hour)
	if err != nil {
		return err
	}

	// Send verification email
	return s.emailService.SendVerificationEmail(email, token)
}

// HandleGoogleUser handles authentication for Google OAuth users
func (s *AuthService) HandleGoogleUser(email, name, googleID, avatar, locale string) (string, *models.User, error) {
	// Check if user exists
	user, err := s.repo.FindByEmail(email)
	if err != nil {
		// If user doesn't exist, create new user
		user = &models.User{
			Email:         email,
			Name:          name,
			GoogleID:      googleID,
			Avatar:        avatar,
			Locale:        locale,
			EmailVerified: true, // Google emails are pre-verified
		}
		if err := s.repo.Create(user); err != nil {
			return "", nil, err
		}
	} else {
		// Update Google info if not set
		if user.GoogleID == "" {
			user.GoogleID = googleID
			user.Name = name
			user.Avatar = avatar
			user.Locale = locale
			if err := s.repo.Update(user); err != nil {
				return "", nil, err
			}
		}
	}

	// Generate JWT token
	token, err := utils.GenerateJWT(user.ID)
	if err != nil {
		return "", nil, err
	}

	return token, user, nil
}
