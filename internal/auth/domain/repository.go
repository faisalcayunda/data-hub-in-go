package domain

import (
	"context"
)

// UserRepository defines the interface for user data operations
type UserRepository interface {
	// GetUserByID retrieves a user by ID
	GetUserByID(ctx context.Context, id string) (*User, error)

	// GetUserByEmail retrieves a user by email
	GetUserByEmail(ctx context.Context, email string) (*User, error)

	// GetUserByUsername retrieves a user by username
	GetUserByUsername(ctx context.Context, username string) (*User, error)

	// CreateUser creates a new user
	CreateUser(ctx context.Context, user *User) error

	// UpdateUser updates an existing user
	UpdateUser(ctx context.Context, user *User) error

	// DeleteUser soft deletes a user by ID
	DeleteUser(ctx context.Context, id string) error

	// ListUsers retrieves a list of users with pagination
	ListUsers(ctx context.Context, limit, offset int) ([]*User, int, error)

	// IsEmailExists checks if email already exists
	IsEmailExists(ctx context.Context, email string) (bool, error)

	// IsUsernameExists checks if username already exists
	IsUsernameExists(ctx context.Context, username string) (bool, error)
}

// TokenRepository defines the interface for token data operations
type TokenRepository interface {
	// CreateToken creates a new token
	CreateToken(ctx context.Context, token *Token) error

	// GetTokenByRefreshToken retrieves a token by refresh token
	GetTokenByRefreshToken(ctx context.Context, refreshToken string) (*Token, error)

	// GetTokenByAccessToken retrieves a token by access token
	GetTokenByAccessToken(ctx context.Context, accessToken string) (*Token, error)

	// RevokeToken revokes a token by ID
	RevokeToken(ctx context.Context, id string) error

	// RevokeUserTokens revokes all tokens for a user
	RevokeUserTokens(ctx context.Context, userID string) error

	// DeleteToken deletes a token by ID
	DeleteToken(ctx context.Context, id string) error

	// CleanupExpiredTokens deletes expired tokens
	CleanupExpiredTokens(ctx context.Context) error
}
