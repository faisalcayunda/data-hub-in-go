package usecase

import (
	"context"
	"fmt"
	"math"
	"time"

	"portal-data-backend/internal/user/domain"
)

// userUsecase implements the Usecase interface
type userUsecase struct {
	userRepo domain.Repository
}

// NewUserUsecase creates a new user usecase
func NewUserUsecase(userRepo domain.Repository) Usecase {
	return &userUsecase{
		userRepo: userRepo,
	}
}

// GetUserByID retrieves a user by ID
func (u *userUsecase) GetUserByID(ctx context.Context, id string) (*domain.UserInfo, error) {
	user, err := u.userRepo.GetUserByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	return u.toUserInfo(user), nil
}

// ListUsers retrieves a paginated list of users
func (u *userUsecase) ListUsers(ctx context.Context, req *domain.ListUsersRequest) (*domain.UserListResponse, error) {
	// Set default values
	if req.Page < 1 {
		req.Page = 1
	}
	if req.Limit < 1 {
		req.Limit = 20
	}
	if req.Limit > 100 {
		req.Limit = 100
	}

	offset := (req.Page - 1) * req.Limit

	// Build filter
	filter := &domain.UserFilter{
		OrganizationID: req.OrganizationID,
		RoleID:         req.RoleID,
		Status:         req.Status,
		Search:         req.Search,
	}

	sortBy := req.SortBy
	if sortBy == "" {
		sortBy = "created_at"
	}
	sortOrder := req.SortOrder
	if sortOrder == "" {
		sortOrder = "DESC"
	}

	// Get users
	users, total, err := u.userRepo.ListUsers(ctx, filter, req.Limit, offset, sortBy, sortOrder)
	if err != nil {
		return nil, fmt.Errorf("failed to list users: %w", err)
	}

	// Convert to response
	userInfos := make([]domain.UserInfo, len(users))
	for i, user := range users {
		userInfos[i] = *u.toUserInfo(user)
	}

	totalPage := int(math.Ceil(float64(total) / float64(req.Limit)))

	return &domain.UserListResponse{
		Users: userInfos,
		Meta: domain.ListMeta{
			Page:      req.Page,
			Limit:     req.Limit,
			Total:     total,
			TotalPage: totalPage,
		},
	}, nil
}

// UpdateUser updates an existing user
func (u *userUsecase) UpdateUser(ctx context.Context, id string, req *domain.UpdateUserRequest) (*domain.UserInfo, error) {
	// Get existing user
	user, err := u.userRepo.GetUserByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	// Update fields
	user.Name = req.Name
	user.UpdatedAt = time.Now()

	if req.Position != "" {
		user.Position = &req.Position
	}
	if req.Address != "" {
		user.Address = &req.Address
	}
	if req.Phone != "" {
		user.Phone = &req.Phone
	}
	if req.Bio != "" {
		user.Bio = &req.Bio
	}

	// Save
	if err := u.userRepo.UpdateUser(ctx, user); err != nil {
		return nil, fmt.Errorf("failed to update user: %w", err)
	}

	return u.toUserInfo(user), nil
}

// DeleteUser soft deletes a user
func (u *userUsecase) DeleteUser(ctx context.Context, id string) error {
	if err := u.userRepo.DeleteUser(ctx, id); err != nil {
		return fmt.Errorf("failed to delete user: %w", err)
	}
	return nil
}

// UpdateUserStatus updates user status
func (u *userUsecase) UpdateUserStatus(ctx context.Context, id string, status string) error {
	if err := u.userRepo.UpdateStatus(ctx, id, status); err != nil {
		return fmt.Errorf("failed to update user status: %w", err)
	}
	return nil
}

// toUserInfo converts User to UserInfo
func (u *userUsecase) toUserInfo(user *domain.User) *domain.UserInfo {
	return &domain.UserInfo{
		ID:             user.ID,
		OrganizationID: user.OrganizationID,
		RoleID:         user.RoleID,
		Name:           user.Name,
		Username:       user.Username,
		Email:          user.Email,
		Position:       user.Position,
		Thumbnail:      user.Thumbnail,
		Status:         string(user.Status),
		CreatedAt:      user.CreatedAt,
		UpdatedAt:      user.UpdatedAt,
	}
}
