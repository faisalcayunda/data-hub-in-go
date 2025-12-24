package domain

import (
	"portal-data-backend/internal/auth/domain"
	"time"
)

// User extends auth domain User with additional profile information
type User struct {
	domain.User
	// Additional profile fields can be added here
	Bio       *string    `db:"bio" json:"bio,omitempty"`
	BirthDate *time.Time `db:"birth_date" json:"birth_date,omitempty"`
}

// UserFilter represents filter options for listing users
type UserFilter struct {
	OrganizationID string
	RoleID         string
	Status         string
	Search         string
}

// UpdateUserRequest represents user update input
type UpdateUserRequest struct {
	Name     string `json:"name" validate:"required,min=2"`
	Position string `json:"position,omitempty"`
	Address  string `json:"address,omitempty"`
	Phone    string `json:"phone,omitempty"`
	Bio      string `json:"bio,omitempty"`
}

// ListUsersRequest represents list users input
type ListUsersRequest struct {
	Page          int    `json:"page" validate:"min=1"`
	Limit         int    `json:"limit" validate:"min=1,max=100"`
	OrganizationID string `json:"organization_id,omitempty"`
	RoleID        string `json:"role_id,omitempty"`
	Status        string `json:"status,omitempty"`
	Search        string `json:"search,omitempty"`
	SortBy        string `json:"sort_by,omitempty"`
	SortOrder     string `json:"sort_order,omitempty"`
}

// UserListResponse represents paginated user list response
type UserListResponse struct {
	Users []UserInfo `json:"users"`
	Meta  ListMeta   `json:"meta"`
}

// UserInfo represents user information in response
type UserInfo struct {
	ID             string     `json:"id"`
	OrganizationID string     `json:"organization_id"`
	RoleID         string     `json:"role_id"`
	Name           string     `json:"name"`
	Username       string     `json:"username"`
	Email          string     `json:"email"`
	Position       *string    `json:"position,omitempty"`
	Thumbnail      *string    `json:"thumbnail,omitempty"`
	Status         string     `json:"status"`
	CreatedAt      time.Time  `json:"created_at"`
	UpdatedAt      time.Time  `json:"updated_at"`
}

// ListMeta represents pagination metadata
type ListMeta struct {
	Page      int `json:"page"`
	Limit     int `json:"limit"`
	Total     int `json:"total"`
	TotalPage int `json:"total_page"`
}
