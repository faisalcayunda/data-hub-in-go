package http

import (
	"portal-data-backend/internal/auth/domain"
)

// LoginRequest represents HTTP request for login
type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8"`
}

// ToDomain converts HTTP request to domain
func (r *LoginRequest) ToDomain() *domain.LoginRequest {
	return &domain.LoginRequest{
		Email:    r.Email,
		Password: r.Password,
	}
}

// RegisterRequest represents HTTP request for registration
type RegisterRequest struct {
	OrganizationID string `json:"organization_id" validate:"required"`
	RoleID         string `json:"role_id" validate:"required"`
	Name           string `json:"name" validate:"required,min=2"`
	Username       string `json:"username" validate:"required,min=3,alphanum"`
	EmployeeID     string `json:"employee_id,omitempty" validate:"omitempty"`
	Position       string `json:"position,omitempty" validate:"omitempty"`
	Email          string `json:"email" validate:"required,email"`
	Password       string `json:"password" validate:"required,min=8"`
	Address        string `json:"address,omitempty" validate:"omitempty"`
	Phone          string `json:"phone,omitempty" validate:"omitempty"`
}

// ToDomain converts HTTP request to domain
func (r *RegisterRequest) ToDomain() *domain.RegisterRequest {
	return &domain.RegisterRequest{
		OrganizationID: r.OrganizationID,
		RoleID:         r.RoleID,
		Name:           r.Name,
		Username:       r.Username,
		EmployeeID:     r.EmployeeID,
		Position:       r.Position,
		Email:          r.Email,
		Password:       r.Password,
		Address:        r.Address,
		Phone:          r.Phone,
	}
}

// RefreshTokenRequest represents HTTP request for token refresh
type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token" validate:"required"`
}

// ToDomain converts HTTP request to domain
func (r *RefreshTokenRequest) ToDomain() *domain.RefreshTokenRequest {
	return &domain.RefreshTokenRequest{
		RefreshToken: r.RefreshToken,
	}
}

// LogoutRequest represents HTTP request for logout
type LogoutRequest struct {
	RefreshToken string `json:"refresh_token" validate:"required"`
}
