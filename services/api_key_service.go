package services

import (
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"
	"own-paynet/models"
	"own-paynet/repository"
	"strings"

	"github.com/google/uuid"
)

type APIKeyService struct {
	apiKeyRepo *repository.APIKeyRepository
}

func NewAPIKeyService(apiKeyRepo *repository.APIKeyRepository) *APIKeyService {
	return &APIKeyService{
		apiKeyRepo: apiKeyRepo,
	}
}

func (s *APIKeyService) generateKey(length int) (string, error) {
	bytes := make([]byte, length)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(bytes), nil
}

func (s *APIKeyService) generateUniquePublicKey() (string, error) {
	// Generate a UUID v4
	u := uuid.New()

	// Convert UUID to base64 URL-safe string
	key := base64.URLEncoding.EncodeToString(u[:])

	// Remove padding characters
	key = strings.TrimRight(key, "=")

	// Add prefix for identification
	key = "pk_" + key

	return key, nil
}

func (s *APIKeyService) generateUniqueSecretKey() (string, error) {
	// Generate a UUID v4
	u := uuid.New()

	// Convert UUID to base64 URL-safe string
	key := base64.URLEncoding.EncodeToString(u[:])

	// Remove padding characters
	key = strings.TrimRight(key, "=")

	// Add prefix for identification
	key = "sk_" + key

	return key, nil
}

func (s *APIKeyService) ensureUniqueKeys(publicKey, _ string) error {
	// Check if public key exists
	existingPublic, err := s.apiKeyRepo.GetByPublicKey(publicKey)
	if err == nil && existingPublic != nil {
		return errors.New("public key collision detected")
	}

	// Check if secret key exists (you might need to add this method to your repository)
	// For now, we'll rely on the database unique constraint
	return nil
}

func (s *APIKeyService) GenerateAPIKey(userID uint, description string) (*models.APIKey, error) {
	// Generate unique public and secret keys
	publicKey, err := s.generateUniquePublicKey()
	if err != nil {
		return nil, fmt.Errorf("failed to generate public key: %v", err)
	}

	secretKey, err := s.generateUniqueSecretKey()
	if err != nil {
		return nil, fmt.Errorf("failed to generate secret key: %v", err)
	}

	// Ensure keys are unique
	if err := s.ensureUniqueKeys(publicKey, secretKey); err != nil {
		return nil, err
	}

	// Check if this is the first key for the user
	existingKeys, err := s.apiKeyRepo.GetByUserID(userID)
	if err != nil {
		return nil, err
	}

	isDefault := len(existingKeys) == 0

	apiKey := &models.APIKey{
		UserID:      userID,
		PublicKey:   publicKey,
		SecretKey:   secretKey,
		IsDefault:   isDefault,
		Description: description,
	}

	if err := s.apiKeyRepo.Create(apiKey); err != nil {
		return nil, fmt.Errorf("failed to create API key: %v", err)
	}

	return apiKey, nil
}

func (s *APIKeyService) GetUserAPIKeys(userID uint) ([]models.APIKey, error) {
	return s.apiKeyRepo.GetByUserID(userID)
}

func (s *APIKeyService) GetAPIKeyByPublicKey(publicKey string) (*models.APIKey, error) {
	return s.apiKeyRepo.GetByPublicKey(publicKey)
}

func (s *APIKeyService) SetDefaultKey(userID uint, apiKeyID uint) error {
	return s.apiKeyRepo.SetDefault(userID, apiKeyID)
}

func (s *APIKeyService) DeleteAPIKey(userID uint, apiKeyID uint) error {
	apiKey, err := s.apiKeyRepo.GetByID(apiKeyID)
	if err != nil {
		return err
	}

	if apiKey.UserID != userID {
		return errors.New("unauthorized to delete this API key")
	}

	if apiKey.IsDefault {
		return errors.New("cannot delete the default API key")
	}

	return s.apiKeyRepo.Delete(apiKeyID)
}

func (s *APIKeyService) UpdateLastUsed(apiKeyID uint) error {
	return s.apiKeyRepo.UpdateLastUsed(apiKeyID)
}
