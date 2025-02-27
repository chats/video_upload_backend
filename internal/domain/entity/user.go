// internal/domain/entity/user.go
package entity

import (
	"time"
)

// UserRole defines the role of a user in the system
type UserRole string

const (
	RoleAdmin  UserRole = "admin"
	RoleUser   UserRole = "user"
	RoleEditor UserRole = "editor"
)

// User represents a user entity in the system
type User struct {
	ID           string    `json:"id"`
	Username     string    `json:"username"`
	Email        string    `json:"email"`
	PasswordHash string    `json:"-"` // Not exposed in JSON
	FullName     string    `json:"fullName,omitempty"`
	AvatarURL    string    `json:"avatarUrl,omitempty"`
	Role         UserRole  `json:"role"`
	CreatedAt    time.Time `json:"createdAt"`
	UpdatedAt    time.Time `json:"updatedAt"`
}

// NewUser creates a new user with default role
func NewUser(id, username, email, passwordHash, fullName string) *User {
	now := time.Now()
	return &User{
		ID:           id,
		Username:     username,
		Email:        email,
		PasswordHash: passwordHash,
		FullName:     fullName,
		Role:         RoleUser, // Default role
		CreatedAt:    now,
		UpdatedAt:    now,
	}
}

// IsAdmin checks if the user has admin role
func (u *User) IsAdmin() bool {
	return u.Role == RoleAdmin
}

// SetPassword updates the user's password hash
func (u *User) SetPassword(passwordHash string) {
	u.PasswordHash = passwordHash
	u.UpdatedAt = time.Now()
}

// UpdateProfile updates the user's profile information
func (u *User) UpdateProfile(fullName, avatarURL string) {
	u.FullName = fullName
	u.AvatarURL = avatarURL
	u.UpdatedAt = time.Now()
}

// SetRole updates the user's role
func (u *User) SetRole(role UserRole) {
	u.Role = role
	u.UpdatedAt = time.Now()
}
