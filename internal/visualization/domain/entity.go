package domain

import "time"

// Visualization represents a data visualization entity
type Visualization struct {
	ID              string        `db:"id" json:"id"`
	Title           string        `db:"title" json:"title"`
	Description     *string       `db:"description" json:"description,omitempty"`
	Type            string        `db:"type" json:"type"` // chart, map, table, etc.
	Config          string        `db:"config" json:"config"` // JSON config for visualization
	DatasetID       *string       `db:"dataset_id" json:"dataset_id,omitempty"`
	OrganizationID  *string       `db:"organization_id" json:"organization_id,omitempty"`
	TopicID         *string       `db:"topic_id" json:"topic_id,omitempty"`
	IsHighlight     bool          `db:"is_highlight" json:"is_highlight"`
	Status          string        `db:"status" json:"status"`
	CreatedBy       string        `db:"created_by" json:"created_by"`
	UpdatedBy       string        `db:"updated_by" json:"updated_by"`
	CreatedAt       time.Time     `db:"created_at" json:"created_at"`
	UpdatedAt       time.Time     `db:"updated_at" json:"updated_at"`
	DeletedAt       *time.Time    `db:"deleted_at" json:"deleted_at,omitempty"`
}

// VisualizationStatus represents visualization status
type VisualizationStatus string

const (
	VisualizationStatusDraft     VisualizationStatus = "draft"
	VisualizationStatusPublished VisualizationStatus = "published"
	VisualizationStatusArchived  VisualizationStatus = "archived"
)

// VisualizationType represents visualization type
type VisualizationType string

const (
	VisualizationTypeBarChart   VisualizationType = "bar"
	VisualizationTypeLineChart  VisualizationType = "line"
	VisualizationTypePieChart   VisualizationType = "pie"
	VisualizationTypeMap        VisualizationType = "map"
	VisualizationTypeTable      VisualizationType = "table"
	VisualizationTypeScatter    VisualizationType = "scatter"
	VisualizationTypeArea       VisualizationType = "area"
	VisualizationTypeHistogram  VisualizationType = "histogram"
)

// ListVisualizationsRequest represents list visualizations input
type ListVisualizationsRequest struct {
	Page           int     `json:"page" validate:"min=1"`
	Limit          int     `json:"limit" validate:"min=1,max=100"`
	DatasetID      *string `json:"dataset_id,omitempty"`
	OrganizationID *string `json:"organization_id,omitempty"`
	TopicID        *string `json:"topic_id,omitempty"`
	Type           *string `json:"type,omitempty"`
	Status         *string `json:"status,omitempty"`
	IsHighlight    *bool   `json:"is_highlight,omitempty"`
	Search         string  `json:"search,omitempty"`
}

// CreateVisualizationRequest represents create visualization input
type CreateVisualizationRequest struct {
	Title          string  `json:"title" validate:"required,min=2,max=200"`
	Description    *string `json:"description,omitempty"`
	Type           string  `json:"type" validate:"required"`
	Config         string  `json:"config" validate:"required"`
	DatasetID      *string `json:"dataset_id,omitempty"`
	OrganizationID *string `json:"organization_id,omitempty"`
	TopicID        *string `json:"topic_id,omitempty"`
	IsHighlight    bool    `json:"is_highlight"`
}

// UpdateVisualizationRequest represents update visualization input
type UpdateVisualizationRequest struct {
	Title          *string `json:"title" validate:"omitempty,min=2,max=200"`
	Description    *string `json:"description,omitempty"`
	Type           *string `json:"type,omitempty"`
	Config         *string `json:"config,omitempty"`
	DatasetID      *string `json:"dataset_id,omitempty"`
	OrganizationID *string `json:"organization_id,omitempty"`
	TopicID        *string `json:"topic_id,omitempty"`
	IsHighlight    *bool   `json:"is_highlight,omitempty"`
	Status         *string `json:"status,omitempty"`
}

// VisualizationInfo represents visualization information for API responses
type VisualizationInfo struct {
	ID             string    `json:"id"`
	Title          string    `json:"title"`
	Description    *string   `json:"description,omitempty"`
	Type           string    `json:"type"`
	Config         string    `json:"config"`
	DatasetID      *string   `json:"dataset_id,omitempty"`
	OrganizationID *string   `json:"organization_id,omitempty"`
	TopicID        *string   `json:"topic_id,omitempty"`
	IsHighlight    bool      `json:"is_highlight"`
	Status         string    `json:"status"`
	CreatedBy      string    `json:"created_by"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

// VisualizationListResponse represents paginated visualization list
type VisualizationListResponse struct {
	Visualizations []VisualizationInfo `json:"visualizations"`
	Meta           ListMeta            `json:"meta"`
}

// ListMeta represents pagination metadata
type ListMeta struct {
	Page      int `json:"page"`
	Limit     int `json:"limit"`
	Total     int `json:"total"`
	TotalPage int `json:"total_page"`
}

// VisualizationStats represents visualization statistics
type VisualizationStats struct {
	TotalCount      int64     `json:"total_count"`
	PublishedCount  int64     `json:"published_count"`
	DraftCount      int64     `json:"draft_count"`
	HighlightCount  int64     `json:"highlight_count"`
	LastUpdated     time.Time `json:"last_updated"`
}
