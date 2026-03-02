package storage

import (
	"billohub/internal/model"
	"billohub/pkg/helper"
	"fmt"
)

// GetUserByUsername retrieves a user by their username.
func (s *Storage) GetUserByUsername(username string) (*model.User, error) {
	var user model.User
	if err := s.db.Where("username = ?", username).First(&user).Error; err != nil {
		return nil, helper.WrapError(err, fmt.Sprintf("failed to get user by username '%s'", username))
	}
	return &user, nil
}

// UpdateUserPassword updates a user's password.
func (s *Storage) UpdateUserPassword(username, newHashedPassword string) error {
	if err := s.db.Model(&model.User{}).Where("username = ?", username).Update("hashed_password", newHashedPassword).Error; err != nil {
		return helper.WrapError(err, fmt.Sprintf("failed to update password for user '%s'", username))
	}
	return nil
}

// CreateUser creates a new user in the database.
func (s *Storage) CreateUser(user *model.User) error {
	if err := s.db.Create(user).Error; err != nil {
		return helper.WrapError(err, fmt.Sprintf("failed to create user '%s'", user.Username))
	}
	return nil
}
