package repository

import (
	"own-paynet/models"

	"gorm.io/gorm"
)

type APIKeyRepository struct {
	db *gorm.DB
}

func NewAPIKeyRepository(db *gorm.DB) *APIKeyRepository {
	return &APIKeyRepository{db: db}
}

func (r *APIKeyRepository) Create(apiKey *models.APIKey) error {
	return r.db.Create(apiKey).Error
}

func (r *APIKeyRepository) GetByID(id uint) (*models.APIKey, error) {
	var apiKey models.APIKey
	err := r.db.First(&apiKey, id).Error
	if err != nil {
		return nil, err
	}
	return &apiKey, nil
}

func (r *APIKeyRepository) GetByPublicKey(publicKey string) (*models.APIKey, error) {
	var apiKey models.APIKey
	err := r.db.Where("public_key = ?", publicKey).First(&apiKey).Error
	if err != nil {
		return nil, err
	}
	return &apiKey, nil
}

func (r *APIKeyRepository) GetByUserID(userID uint) ([]models.APIKey, error) {
	var apiKeys []models.APIKey
	err := r.db.Where("user_id = ?", userID).Find(&apiKeys).Error
	return apiKeys, err
}

func (r *APIKeyRepository) GetDefaultByUserID(userID uint) (*models.APIKey, error) {
	var apiKey models.APIKey
	err := r.db.Where("user_id = ? AND is_default = ?", userID, true).First(&apiKey).Error
	if err != nil {
		return nil, err
	}
	return &apiKey, nil
}

func (r *APIKeyRepository) Update(apiKey *models.APIKey) error {
	return r.db.Save(apiKey).Error
}

func (r *APIKeyRepository) Delete(id uint) error {
	return r.db.Delete(&models.APIKey{}, id).Error
}

func (r *APIKeyRepository) SetDefault(userID uint, apiKeyID uint) error {
	// Start a transaction
	tx := r.db.Begin()

	// Reset all default keys for the user
	if err := tx.Model(&models.APIKey{}).Where("user_id = ?", userID).Update("is_default", false).Error; err != nil {
		tx.Rollback()
		return err
	}

	// Set the new default key
	if err := tx.Model(&models.APIKey{}).Where("id = ? AND user_id = ?", apiKeyID, userID).Update("is_default", true).Error; err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit().Error
}

func (r *APIKeyRepository) UpdateLastUsed(id uint) error {
	return r.db.Model(&models.APIKey{}).Where("id = ?", id).Update("last_used_at", gorm.Expr("NOW()")).Error
}
