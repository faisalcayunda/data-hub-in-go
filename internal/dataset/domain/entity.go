package domain

import (
	"time"
)

// Dataset represents a dataset entity
type Dataset struct {
	ID                string        `db:"id" json:"id"`
	Name              string        `db:"name" json:"name"`
	Slug              string        `db:"slug" json:"slug"`
	Description       *string       `db:"description" json:"description,omitempty"`
	Period            *string       `db:"period" json:"period,omitempty"`
	UnitID            *string       `db:"unit_id" json:"unit_id,omitempty"`
	BusinessFieldID   *string       `db:"business_field_id" json:"business_field_id,omitempty"`
	Image             *string       `db:"image" json:"image,omitempty"`
	TopicID           *string       `db:"topic_id" json:"topic_id,omitempty"`
	OrganizationID    string        `db:"organization_id" json:"organization_id"`
	ReferenceID       *string       `db:"reference_id" json:"reference_id,omitempty"`
	Classification    string        `db:"classification" json:"classification"`
	Category          string        `db:"category" json:"category"`
	DataFixed         bool          `db:"data_fixed" json:"data_fixed"`
	ValidationStatus  ValidationStatus `db:"validation_status" json:"validation_status"`
	Metadata          *string       `db:"metadatas" json:"metadatas,omitempty"`
	CreatedBy         string        `db:"created_by" json:"created_by"`
	UpdatedBy         *string       `db:"updated_by" json:"updated_by,omitempty"`
	CreatedAt         time.Time     `db:"created_at" json:"created_at"`
	UpdatedAt         time.Time     `db:"updated_at" json:"updated_at"`
	IsHighlight       bool          `db:"is_highlight" json:"is_highlight"`
	Status            DatasetStatus `db:"status" json:"status"`

	// Relations
	Tags              []Tag         `json:"tags,omitempty"`
	Unit              *Unit         `json:"unit,omitempty"`
	BusinessField     *BusinessField `json:"business_field,omitempty"`
	Topic             *Topic        `json:"topic,omitempty"`
	Organization      *OrganizationSummary `json:"organization,omitempty"`
}

// ValidationStatus represents dataset validation status
type ValidationStatus string

const (
	ValidationStatusValid     ValidationStatus = "valid"
	ValidationStatusInvalid   ValidationStatus = "invalid"
	ValidationStatusPending   ValidationStatus = "pending"
)

// DatasetStatus represents dataset status
type DatasetStatus string

const (
	DatasetStatusDraft     DatasetStatus = "draft"
	DatasetStatusPublished DatasetStatus = "published"
	DatasetStatusArchived  DatasetStatus = "archived"
)

// Tag represents a tag entity
type Tag struct {
	ID        string    `db:"id" json:"id"`
	Name      string    `db:"name" json:"name"`
	Slug      string    `db:"slug" json:"slug"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
}

// Unit represents a unit of measurement
type Unit struct {
	ID          string    `db:"id" json:"id"`
	Name        string    `db:"name" json:"name"`
	Symbol      string    `db:"symbol" json:"symbol"`
	CreatedAt   time.Time `db:"created_at" json:"created_at"`
}

// BusinessField represents a business field
type BusinessField struct {
	ID        string    `db:"id" json:"id"`
	Name      string    `db:"name" json:"name"`
	Slug      string    `db:"slug" json:"slug"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
}

// Topic represents a topic
type Topic struct {
	ID        string    `db:"id" json:"id"`
	Name      string    `db:"name" json:"name"`
	Slug      string    `db:"slug" json:"slug"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
}

// OrganizationSummary represents a summary of organization
type OrganizationSummary struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	Slug string `json:"slug"`
}

// CreateDatasetRequest represents dataset creation input
type CreateDatasetRequest struct {
	Name            string   `json:"name" validate:"required,min=2"`
	Description     string   `json:"description,omitempty"`
	Period          string   `json:"period,omitempty"`
	UnitID          string   `json:"unit_id,omitempty"`
	BusinessFieldID string   `json:"business_field_id,omitempty"`
	Image           string   `json:"image,omitempty"`
	TopicID         string   `json:"topic_id,omitempty"`
	ReferenceID     string   `json:"reference_id,omitempty"`
	Classification  string   `json:"classification" validate:"required"`
	Category        string   `json:"category" validate:"required"`
	DataFixed       bool     `json:"data_fixed"`
	ValidationStatus string  `json:"validation_status,omitempty"`
	Metadata        string   `json:"metadatas,omitempty"`
	TagIDs          []string `json:"tag_ids,omitempty"`
	IsHighlight     bool     `json:"is_highlight"`
}

// UpdateDatasetRequest represents dataset update input
type UpdateDatasetRequest struct {
	Name            string   `json:"name" validate:"required,min=2"`
	Description     string   `json:"description,omitempty"`
	Period          string   `json:"period,omitempty"`
	UnitID          string   `json:"unit_id,omitempty"`
	BusinessFieldID string   `json:"business_field_id,omitempty"`
	Image           string   `json:"image,omitempty"`
	TopicID         string   `json:"topic_id,omitempty"`
	ReferenceID     string   `json:"reference_id,omitempty"`
	Classification  string   `json:"classification" validate:"required"`
	Category        string   `json:"category" validate:"required"`
	DataFixed       bool     `json:"data_fixed"`
	ValidationStatus string  `json:"validation_status,omitempty"`
	Metadata        string   `json:"metadatas,omitempty"`
	TagIDs          []string `json:"tag_ids,omitempty"`
	IsHighlight     bool     `json:"is_highlight"`
}

// ListDatasetsRequest represents list datasets input
type ListDatasetsRequest struct {
	Page            int    `json:"page" validate:"min=1"`
	Limit           int    `json:"limit" validate:"min=1,max=100"`
	OrganizationID  string `json:"organization_id,omitempty"`
	TopicID         string `json:"topic_id,omitempty"`
	BusinessFieldID string `json:"business_field_id,omitempty"`
	TagID           string `json:"tag_id,omitempty"`
	Status          string `json:"status,omitempty"`
	ValidationStatus string `json:"validation_status,omitempty"`
	Classification  string `json:"classification,omitempty"`
	Search          string `json:"search,omitempty"`
	SortBy          string `json:"sort_by,omitempty"`
	SortOrder       string `json:"sort_order,omitempty"`
}

// DatasetResponse represents dataset response
type DatasetResponse struct {
	ID               string              `json:"id"`
	Name             string              `json:"name"`
	Slug             string              `json:"slug"`
	Description      *string             `json:"description,omitempty"`
	Period           *string             `json:"period,omitempty"`
	Unit             *Unit               `json:"unit,omitempty"`
	BusinessField    *BusinessField      `json:"business_field,omitempty"`
	Image            *string             `json:"image,omitempty"`
	Topic            *Topic              `json:"topic,omitempty"`
	OrganizationID   string              `json:"organization_id"`
	ReferenceID      *string             `json:"reference_id,omitempty"`
	Classification   string              `json:"classification"`
	Category         string              `json:"category"`
	DataFixed        bool                `json:"data_fixed"`
	ValidationStatus string              `json:"validation_status"`
	Metadata         *string             `json:"metadatas,omitempty"`
	CreatedBy        string              `json:"created_by"`
	UpdatedBy        *string             `json:"updated_by,omitempty"`
	CreatedAt        time.Time           `json:"created_at"`
	UpdatedAt        time.Time           `json:"updated_at"`
	IsHighlight      bool                `json:"is_highlight"`
	Status           string              `json:"status"`
	Tags             []Tag               `json:"tags,omitempty"`
}

// DatasetListResponse represents paginated dataset list
type DatasetListResponse struct {
	Datasets []DatasetResponse `json:"datasets"`
	Meta     ListMeta          `json:"meta"`
}

// ListMeta represents pagination metadata
type ListMeta struct {
	Page      int `json:"page"`
	Limit     int `json:"limit"`
	Total     int `json:"total"`
	TotalPage int `json:"total_page"`
}
