package domain

import (
	"time"
)

// Feedback represents user feedback
type Feedback struct {
	ID          string       `db:"id" json:"id"`
	UserID      string       `db:"user_id" json:"user_id"`
	DatasetID   *string      `db:"dataset_id" json:"dataset_id,omitempty"`
	Rating      int          `db:"rating" json:"rating" validate:"min=1,max=5"`
	Comment     string       `db:"comment" json:"comment"`
	Category    FeedbackCategory `db:"category" json:"category"`
	Status      FeedbackStatus `db:"status" json:"status"`
	CreatedAt   time.Time    `db:"created_at" json:"created_at"`
	UpdatedAt   time.Time    `db:"updated_at" json:"updated_at"`
}

// FeedbackCategory represents feedback category
type FeedbackCategory string

const (
	FeedbackCategoryDataQuality    FeedbackCategory = "data_quality"
	FeedbackCategoryUsability      FeedbackCategory = "usability"
	FeedbackCategoryFeatureRequest FeedbackCategory = "feature_request"
	FeedbackCategoryBug            FeedbackCategory = "bug"
	FeedbackCategoryOther          FeedbackCategory = "other"
)

// FeedbackStatus represents feedback status
type FeedbackStatus string

const (
	FeedbackStatusPending  FeedbackStatus = "pending"
	FeedbackStatusReview   FeedbackStatus = "in_review"
	FeedbackStatusResolved FeedbackStatus = "resolved"
	FeedbackStatusClosed   FeedbackStatus = "closed"
)

// CreateFeedbackRequest represents feedback creation input
type CreateFeedbackRequest struct {
	DatasetID *string            `json:"dataset_id,omitempty"`
	Rating    int                `json:"rating" validate:"required,min=1,max=5"`
	Comment   string             `json:"comment" validate:"required,min=10,max=1000"`
	Category  FeedbackCategory   `json:"category" validate:"required"`
}

// UpdateFeedbackStatusRequest represents feedback status update
type UpdateFeedbackStatusRequest struct {
	Status FeedbackStatus `json:"status" validate:"required"`
}

// ListFeedbacksRequest represents list feedbacks input
type ListFeedbacksRequest struct {
	Page       int                `json:"page" validate:"min=1"`
	Limit      int                `json:"limit" validate:"min=1,max=100"`
	DatasetID  *string             `json:"dataset_id,omitempty"`
	Category   *string             `json:"category,omitempty"`
	Status     *string             `json:"status,omitempty"`
	UserID     *string             `json:"user_id,omitempty"`
	Search     string             `json:"search,omitempty"`
	SortBy     string             `json:"sort_by,omitempty"`
	SortOrder  string             `json:"sort_order,omitempty"`
}

// FeedbackResponse represents feedback response
type FeedbackResponse struct {
	ID        string            `json:"id"`
	UserID    string            `json:"user_id"`
	DatasetID *string           `json:"dataset_id,omitempty"`
	Rating    int               `json:"rating"`
	Comment   string            `json:"comment"`
	Category  string            `json:"category"`
	Status    string            `json:"status"`
	CreatedAt time.Time         `json:"created_at"`
	UpdatedAt time.Time         `json:"updated_at"`
}

// FeedbackListResponse represents paginated feedback list
type FeedbackListResponse struct {
	Feedbacks []FeedbackResponse `json:"feedbacks"`
	Meta      ListMeta           `json:"meta"`
}

// ListMeta represents pagination metadata
type ListMeta struct {
	Page      int `json:"page"`
	Limit     int `json:"limit"`
	Total     int `json:"total"`
	TotalPage int `json:"total_page"`
}
