package domain

import "time"

// Setting represents a system or user setting
type Setting struct {
	ID        string     `db:"id" json:"id"`
	Key       string     `db:"key" json:"key"`
	Value     string     `db:"value" json:"value"`
	Type      string     `db:"type" json:"type"` // string, number, boolean, json
	Category  string     `db:"category" json:"category"` // system, user, organization
	UserID    *string    `db:"user_id" json:"user_id,omitempty"`
	IsPublic  bool       `db:"is_public" json:"is_public"`
	CreatedAt time.Time  `db:"created_at" json:"created_at"`
	UpdatedAt time.Time  `db:"updated_at" json:"updated_at"`
	DeletedAt *time.Time `db:"deleted_at" json:"deleted_at,omitempty"`
}

// SettingType represents setting data type
type SettingType string

const (
	SettingTypeString  SettingType = "string"
	SettingTypeNumber  SettingType = "number"
	SettingTypeBoolean SettingType = "boolean"
	SettingTypeJSON    SettingType = "json"
)

// SettingCategory represents setting category
type SettingCategory string

const (
	SettingCategorySystem        SettingCategory = "system"
	SettingCategoryUser          SettingCategory = "user"
	SettingCategoryOrganization  SettingCategory = "organization"
)

// ListSettingsRequest represents list settings input
type ListSettingsRequest struct {
	Page     int     `json:"page" validate:"min=1"`
	Limit    int     `json:"limit" validate:"min=1,max=100"`
	Category *string `json:"category,omitempty"`
	UserID   *string `json:"user_id,omitempty"`
	Type     *string `json:"type,omitempty"`
	Search   string  `json:"search,omitempty"`
}

// CreateSettingRequest represents create setting input
type CreateSettingRequest struct {
	Key      string `json:"key" validate:"required,min=2,max=100"`
	Value    string `json:"value" validate:"required"`
	Type     string `json:"type" validate:"required"`
	Category string `json:"category" validate:"required"`
	UserID   *string `json:"user_id,omitempty"`
	IsPublic bool   `json:"is_public"`
}

// UpdateSettingRequest represents update setting input
type UpdateSettingRequest struct {
	Value    *string `json:"value" validate:"omitempty"`
	Type     *string `json:"type,omitempty"`
	IsPublic *bool   `json:"is_public,omitempty"`
}

// SettingInfo represents setting information for API responses
type SettingInfo struct {
	ID        string    `json:"id"`
	Key       string    `json:"key"`
	Value     string    `json:"value"`
	Type      string    `json:"type"`
	Category  string    `json:"category"`
	UserID    *string   `json:"user_id,omitempty"`
	IsPublic  bool      `json:"is_public"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// SettingListResponse represents paginated setting list
type SettingListResponse struct {
	Settings []SettingInfo `json:"settings"`
	Meta     ListMeta      `json:"meta"`
}

// ListMeta represents pagination metadata
type ListMeta struct {
	Page      int `json:"page"`
	Limit     int `json:"limit"`
	Total     int `json:"total"`
	TotalPage int `json:"total_page"`
}

// GetSettingsByKeysResponse represents response for getting multiple settings by keys
type GetSettingsByKeysResponse struct {
	Settings map[string]string `json:"settings"`
}
