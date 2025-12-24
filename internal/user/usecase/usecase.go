package usecase

import (
	"context"

	"portal-data-backend/internal/user/domain"
)

// Usecase defines the interface for user business logic
type Usecase interface {
	// GetUserByID retrieves a user by ID
	GetUserByID(ctx context.Context, id string) (*domain.UserInfo, error)

	// ListUsers retrieves a paginated list of users
	ListUsers(ctx context.Context, req *domain.ListUsersRequest) (*domain.UserListResponse, error)

	// UpdateUser updates an existing user
	UpdateUser(ctx context.Context, id string, req *domain.UpdateUserRequest) (*domain.UserInfo, error)

	// DeleteUser soft deletes a user
	DeleteUser(ctx context.Context, id string) error

	// UpdateUserStatus updates user status
	UpdateUserStatus(ctx context.Context, id string, status string) error
}
