package usecase

import (
	"context"

	"portal-data-backend/internal/auth/domain"
)

// Usecase defines the interface for auth business logic
type Usecase interface {
	// Login authenticates a user and returns tokens
	Login(ctx context.Context, req *domain.LoginRequest) (*domain.AuthResponse, error)

	// Register creates a new user account
	Register(ctx context.Context, req *domain.RegisterRequest) (*domain.AuthResponse, error)

	// Logout logs out a user by revoking their tokens
	Logout(ctx context.Context, accessToken, refreshToken string) error

	// RefreshToken refreshes an access token using a refresh token
	RefreshToken(ctx context.Context, refreshToken string) (*domain.AuthResponse, error)

	// RevokeAllTokens revokes all tokens for a user
	RevokeAllTokens(ctx context.Context, userID string) error

	// ValidateToken validates a token and returns the claims
	ValidateToken(ctx context.Context, token string) (*domain.TokenClaims, error)

	// GetCurrentUser retrieves the current user by ID
	GetCurrentUser(ctx context.Context, userID string) (*domain.UserInfo, error)
}
