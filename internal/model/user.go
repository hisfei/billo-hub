package model

import "time"

// User represents a user in the system.
type User struct {
	ID             int       `gorm:"primaryKey" json:"id"`
	Username       string    `gorm:"unique;not null" json:"username"`
	HashedPassword string    `gorm:"not null" json:"-"`
	CreatedAt      time.Time `json:"created_at,omitempty"`
	UpdatedAt      time.Time `json:"updated_at,omitempty"`
}

// ResetPasswordRequest defines the structure for a password reset request.
type ResetPasswordRequest struct {
	Username    string `json:"username" binding:"required"`
	OldPassword string `json:"old_password" binding:"required"`
	NewPassword string `json:"new_password" binding:"required"`
}
