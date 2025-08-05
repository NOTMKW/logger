package services

import (
	"fmt"
	"time"

	"github.com/notmkw/log/internal/models"
	"github.com/notmkw/log/internal/repositories"

	"github.com/google/uuid"
)

type UserService struct {
	userRepo *repositories.UserRepository
}

func NewUserService(userRepo *repositories.UserRepository) *UserService {
	return &UserService{
		userRepo: userRepo,
	}
}

func (s *UserService) GetUser(userID string) (*models.User, error) {
	if userID == "" {
		return nil, fmt.Errorf("user ID cannot be empty")
	}

	user := s.userRepo.GetByID(userID)
	if user == nil {
		return nil, fmt.Errorf("user not found")
	}

	return user, nil
}

func (s *UserService) CreateUser() (*models.User, error) {
	user := &models.User{
		ID:        uuid.New().String(),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
	}

	if err := s.userRepo.Create(user); err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	return user, nil
}

func (s *UserService) UpdateUser(userID string) (*models.User, error) {
	if userID == "" {
		return nil, fmt.Errorf("user ID cannot be empty")
	}

	user := s.userRepo.GetByID(userID)
	if user == nil {
		return nil, fmt.Errorf("user not found")
	}

	user.UpdatedAt = time.Now().UTC()

	if err := s.userRepo.Update(user); err != nil {
		return nil, fmt.Errorf("failed to update user: %w", err)
	}

	return user, nil
}

func (s *UserService) DeleteUser(userID string) error {
	if userID == "" {
		return fmt.Errorf("user ID cannot be empty")
	}

	user := s.userRepo.GetByID(userID)
	if user == nil {
		return fmt.Errorf("user not found")
	}

	return s.userRepo.Delete(userID)
}
