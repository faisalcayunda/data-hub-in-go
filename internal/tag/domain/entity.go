package domain

import (
	"time"
)

// Tag represents a tag entity
type Tag struct {
	ID        string    `db:"id" json:"id"`
	Name      string    `db:"name" json:"name"`
	Slug      string    `db:"slug" json:"slug"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
}

// CreateTagRequest represents tag creation input
type CreateTagRequest struct {
	Name string `json:"name" validate:"required,min=2"`
}

// UpdateTagRequest represents tag update input
type UpdateTagRequest struct {
	Name string `json:"name" validate:"required,min=2"`
}

// ListTagsRequest represents list tags input
type ListTagsRequest struct {
	Page  int    `json:"page" validate:"min=1"`
	Limit int    `json:"limit" validate:"min=1,max=100"`
	Search string `json:"search,omitempty"`
}

// TagListResponse represents paginated tag list
type TagListResponse struct {
	Tags []TagResponse `json:"tags"`
	Meta ListMeta      `json:"meta"`
}

// TagResponse represents tag response
type TagResponse struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	Slug      string    `json:"slug"`
	CreatedAt time.Time `json:"created_at"`
}

// ListMeta represents pagination metadata
type ListMeta struct {
	Page      int `json:"page"`
	Limit     int `json:"limit"`
	Total     int `json:"total"`
	TotalPage int `json:"total_page"`
}
