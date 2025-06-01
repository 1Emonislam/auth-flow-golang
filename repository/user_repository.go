package repository

import (
	"fmt"
	"own-paynet/models"
	"time"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type UserRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{db: db}
}
func (r *UserRepository) Create(user *models.User) error {
	// Check if email already exists
	var existingUser models.User
	if err := r.db.Where("email = ?", user.Email).First(&existingUser).Error; err == nil {
		return fmt.Errorf("user with email %s already exists", user.Email)
	} else if err != gorm.ErrRecordNotFound {
		return fmt.Errorf("error checking email existence: %w", err)
	}

	// Hash the password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("failed to hash password: %w", err)
	}
	user.Password = string(hashedPassword)

	// Create the user
	if err := r.db.Create(user).Error; err != nil {
		return fmt.Errorf("failed to create user: %w", err)
	}
	return nil
}

func (r *UserRepository) FindByEmail(email string) (*models.User, error) {
	var user models.User
	err := r.db.Where("email = ?", email).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *UserRepository) UpdatePassword(email, password string) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	return r.db.Model(&models.User{}).Where("email = ?", email).Update("password", string(hashedPassword)).Error
}

// VerifyEmail marks a user's email as verified
func (r *UserRepository) VerifyEmail(email string) error {
	now := time.Now()
	return r.db.Model(&models.User{}).Where("email = ?", email).Updates(map[string]interface{}{
		"email_verified":    true,
		"email_verified_at": now,
	}).Error
}

// IsEmailVerified checks if a user's email is verified
func (r *UserRepository) IsEmailVerified(email string) (bool, error) {
	var user models.User
	err := r.db.Select("email_verified").Where("email = ?", email).First(&user).Error
	if err != nil {
		return false, err
	}
	return user.EmailVerified, nil
}
