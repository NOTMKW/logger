package repositories

import (
	"fmt"
	"sync"
	"github.com/notmkw/logger/internal/models"
)

type UserRepository struct {
	users map[string]*models.User
	mutex sync.RWMutex
}

func NewUserRepository() *UserRepository {
	return &UserRepository{
		users: make(map[string]*models.User),
	}
}

func (r *UserRepository) GetByID(userID string) *models.User {
	r.mutex.RLock()
	defer r.mutex.RUnlock()
	
	return r.users[userID]
}

func (r *UserRepository) Create(user *models.User) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()
	
	if _, exists := r.users[user.ID]; exists {
		return fmt.Errorf("user with ID %s already exists", user.ID)
	}
	
	r.users[user.ID] = user
	return nil
}

func (r *UserRepository) Update(user *models.User) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()
	
	if _, exists := r.users[user.ID]; !exists {
		return fmt.Errorf("user with ID %s does not exist", user.ID)
	}
	
	r.users[user.ID] = user
	return nil
}

func (r *UserRepository) Delete(userID string) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()
	
	if _, exists := r.users[userID]; !exists {
		return fmt.Errorf("user with ID %s does not exist", userID)
	}
	
	delete(r.users, userID)
	return nil
}

func (r *UserRepository) GetAll() []*models.User {
	r.mutex.RLock()
	defer r.mutex.RUnlock()
	
	users := make([]*models.User, 0, len(r.users))
	for _, user := range r.users {
		users = append(users, user)
	}
	
	return users
}