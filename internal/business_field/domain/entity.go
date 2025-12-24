package domain

import (
	"time"
)

// BusinessField represents a business field/industry entity
type BusinessField struct {
	ID        string    `db:"id" json:"id"`
	Name      string    `db:"name" json:"name"`
	Slug      string    `db:"slug" json:"slug"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
}

// CreateBusinessFieldRequest represents business field creation input
type CreateBusinessFieldRequest struct {
	Name string `json:"name" validate:"required,min=2"`
}

// UpdateBusinessFieldRequest represents business field update input
type UpdateBusinessFieldRequest struct {
	Name string `json:"name" validate:"required,min=2"`
}

// ListBusinessFieldsRequest represents list business fields input
type ListBusinessFieldsRequest struct {
	Page  int    `json:"page" validate:"min=1"`
	Limit int    `json:"limit" validate:"min=1,max=100"`
	Search string `json:"search,omitempty"`
}

// BusinessFieldResponse represents business field response
type BusinessFieldResponse struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	Slug      string    `json:"slug"`
	CreatedAt time.Time `json:"created_at"`
}

// BusinessFieldListResponse represents paginated business field list
type BusinessFieldListResponse struct {
	BusinessFields []BusinessFieldResponse `json:"business_fields"`
	Meta            ListMeta               `json:"meta"`
}

// ListMeta represents pagination metadata
type ListMeta struct {
	Page      int `json:"page"`
	Limit     int `json:"limit"`
	Total     int `json:"total"`
	TotalPage int `json:"total_page"`
}
