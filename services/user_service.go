package services

import (
	"own-paynet/models"
	"own-paynet/repository"
)

type UserService struct {
	userRepo *repository.UserRepository
}

func NewUserService(userRepo *repository.UserRepository) *UserService {
	return &UserService{
		userRepo: userRepo,
	}
}

// GetUserByID retrieves a user by their ID
func (s *UserService) GetUserByID(id uint) (*models.User, error) {
	return s.userRepo.FindByID(id)
}
