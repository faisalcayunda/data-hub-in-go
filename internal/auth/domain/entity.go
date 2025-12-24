package domain

import (
	"time"
)

// User represents a user entity
type User struct {
	ID             string    `db:"id" json:"id"`
	OrganizationID string    `db:"organization_id" json:"organization_id"`
	RoleID         string    `db:"role_id" json:"role_id"`
	Name           string    `db:"name" json:"name"`
	Username       string    `db:"username" json:"username"`
	EmployeeID     *string   `db:"employee_id" json:"employee_id,omitempty"`
	Position       *string   `db:"position" json:"position,omitempty"`
	Email          string    `db:"email" json:"email"`
	PasswordHash   string    `db:"password_hash" json:"-"`
	Address        *string   `db:"address" json:"address,omitempty"`
	Phone          *string   `db:"phone" json:"phone,omitempty"`
	Thumbnail      *string   `db:"thumbnail" json:"thumbnail,omitempty"`
	Status         UserStatus `db:"status" json:"status"`
	CreatedAt      time.Time `db:"created_at" json:"created_at"`
	UpdatedAt      time.Time `db:"updated_at" json:"updated_at"`
}

// UserStatus represents the status of a user
type UserStatus string

const (
	UserStatusActive   UserStatus = "active"
	UserStatusInactive UserStatus = "inactive"
	UserStatusSuspended UserStatus = "suspended"
	UserStatusDeleted  UserStatus = "deleted"
)

// IsActive checks if user is active
func (u *User) IsActive() bool {
	return u.Status == UserStatusActive
}

// Token represents a refresh token entity
type Token struct {
	ID           string    `db:"id" json:"id"`
	UserID       string    `db:"user_id" json:"user_id"`
	AccessToken  string    `db:"access_token" json:"access_token"`
	RefreshToken string    `db:"refresh_token" json:"refresh_token"`
	ExpiresAt    time.Time `db:"expires_at" json:"expires_at"`
	Revoked      bool      `db:"revoked" json:"revoked"`
	CreatedAt    time.Time `db:"created_at" json:"created_at"`
}

// IsExpired checks if token is expired
func (t *Token) IsExpired() bool {
	return time.Now().After(t.ExpiresAt)
}

// IsValid checks if token is valid (not revoked and not expired)
func (t *Token) IsValid() bool {
	return !t.Revoked && !t.IsExpired()
}

// LoginRequest represents login input
type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8"`
}

// RegisterRequest represents registration input
type RegisterRequest struct {
	OrganizationID string `json:"organization_id" validate:"required"`
	RoleID         string `json:"role_id" validate:"required"`
	Name           string `json:"name" validate:"required,min=2"`
	Username       string `json:"username" validate:"required,min=3,alphanum"`
	EmployeeID     string `json:"employee_id,omitempty"`
	Position       string `json:"position,omitempty"`
	Email          string `json:"email" validate:"required,email"`
	Password       string `json:"password" validate:"required,min=8"`
	Address        string `json:"address,omitempty"`
	Phone          string `json:"phone,omitempty"`
}

// RefreshTokenRequest represents refresh token input
type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token" validate:"required"`
}

// AuthResponse represents authentication response
type AuthResponse struct {
	User         UserInfo    `json:"user"`
	AccessToken  string      `json:"access_token"`
	RefreshToken string      `json:"refresh_token"`
	ExpiresIn    int64       `json:"expires_in"`
	TokenType    string      `json:"token_type"`
}

// UserInfo represents user information in auth response
type UserInfo struct {
	ID             string `json:"id"`
	OrganizationID string `json:"organization_id"`
	RoleID         string `json:"role_id"`
	Name           string `json:"name"`
	Username       string `json:"username"`
	Email          string `json:"email"`
	Thumbnail      string `json:"thumbnail,omitempty"`
}

// TokenClaims represents JWT token claims
type TokenClaims struct {
	UserID         string
	OrganizationID string
	RoleID         string
	Name           string
	Username       string
	Email          string
}

// ToUserInfo converts User to UserInfo
func (u *User) ToUserInfo() UserInfo {
	info := UserInfo{
		ID:             u.ID,
		OrganizationID: u.OrganizationID,
		RoleID:         u.RoleID,
		Name:           u.Name,
		Username:       u.Username,
		Email:          u.Email,
	}
	if u.Thumbnail != nil {
		info.Thumbnail = *u.Thumbnail
	}
	return info
}
