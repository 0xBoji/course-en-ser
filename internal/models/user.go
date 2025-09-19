package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// User represents a user in the system (admin users for authentication)
type User struct {
	ID        uuid.UUID `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()" example:"123e4567-e89b-12d3-a456-426614174000"`
	Username  string    `json:"username" gorm:"not null;size:255;unique" validate:"required,min=1,max=255" example:"admin"`
	Password  string    `json:"-" gorm:"not null;size:255" validate:"required,min=1"` // Password is never returned in JSON
	Role      string    `json:"role" gorm:"not null;size:50;default:admin" validate:"required" example:"admin"`
	CreatedAt time.Time `json:"created_at" gorm:"autoCreateTime" example:"2023-01-01T00:00:00Z"`
	UpdatedAt time.Time `json:"updated_at" gorm:"autoUpdateTime" example:"2023-01-01T00:00:00Z"`
}

// BeforeCreate will set a UUID rather than numeric ID
func (u *User) BeforeCreate(tx *gorm.DB) error {
	if u.ID == uuid.Nil {
		u.ID = uuid.New()
	}
	return nil
}

// TableName returns the table name for User model
func (User) TableName() string {
	return "users"
}

// LoginRequest represents the request payload for user login
type LoginRequest struct {
	Username string `json:"username" validate:"required" example:"admin"`
	Password string `json:"password" validate:"required" example:"admin!dev"`
}

// LoginResponse represents the response payload for successful login
type LoginResponse struct {
	Token string `json:"token" example:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."`
	User  UserResponse `json:"user"`
}

// UserResponse represents the response payload for user operations (without password)
type UserResponse struct {
	ID        uuid.UUID `json:"id" example:"123e4567-e89b-12d3-a456-426614174000"`
	Username  string    `json:"username" example:"admin"`
	Role      string    `json:"role" example:"admin"`
	CreatedAt time.Time `json:"created_at" example:"2023-01-01T00:00:00Z"`
}

// ToResponse converts User model to UserResponse
func (u *User) ToResponse() UserResponse {
	return UserResponse{
		ID:        u.ID,
		Username:  u.Username,
		Role:      u.Role,
		CreatedAt: u.CreatedAt,
	}
}
