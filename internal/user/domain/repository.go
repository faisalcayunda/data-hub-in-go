package domain

import (
	"context"
)

// Repository defines the interface for user data operations
type Repository interface {
	// GetUserByID retrieves a user by ID
	GetUserByID(ctx context.Context, id string) (*User, error)

	// ListUsers retrieves a list of users with filters and pagination
	ListUsers(ctx context.Context, filter *UserFilter, limit, offset int, sortBy, sortOrder string) ([]*User, int, error)

	// UpdateUser updates an existing user
	UpdateUser(ctx context.Context, user *User) error

	// DeleteUser soft deletes a user by ID
	DeleteUser(ctx context.Context, id string) error

	// UpdateStatus updates user status
	UpdateStatus(ctx context.Context, id string, status string) error
}
