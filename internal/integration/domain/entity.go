package domain

import "time"

// Integration represents an external service integration
type Integration struct {
	ID             string     `db:"id" json:"id"`
	Name           string     `db:"name" json:"name"`
	Type           string     `db:"type" json:"type"` // api, webhook, database, etc.
	Description    *string    `db:"description" json:"description,omitempty"`
	Config         string     `db:"config" json:"config"` // JSON config
	Endpoint       *string    `db:"endpoint" json:"endpoint,omitempty"`
	APIKey         *string    `db:"api_key" json:"api_key,omitempty"`
	Status         string     `db:"status" json:"status"`
	LastSyncAt     *time.Time `db:"last_sync_at" json:"last_sync_at,omitempty"`
	OrganizationID *string    `db:"organization_id" json:"organization_id,omitempty"`
	CreatedBy      string     `db:"created_by" json:"created_by"`
	CreatedAt      time.Time  `db:"created_at" json:"created_at"`
	UpdatedAt      time.Time  `db:"updated_at" json:"updated_at"`
	DeletedAt      *time.Time `db:"deleted_at" json:"deleted_at,omitempty"`
}

// IntegrationType represents integration type
type IntegrationType string

const (
	IntegrationTypeAPI     IntegrationType = "api"
	IntegrationTypeWebhook IntegrationType = "webhook"
	IntegrationTypeDatabase IntegrationType = "database"
	IntegrationTypeCustom  IntegrationType = "custom"
)

// IntegrationStatus represents integration status
type IntegrationStatus string

const (
	IntegrationStatusActive   IntegrationStatus = "active"
	IntegrationStatusInactive IntegrationStatus = "inactive"
	IntegrationStatusError    IntegrationStatus = "error"
)

// ListIntegrationsRequest represents list integrations input
type ListIntegrationsRequest struct {
	Page           int     `json:"page" validate:"min=1"`
	Limit          int     `json:"limit" validate:"min=1,max=100"`
	OrganizationID *string `json:"organization_id,omitempty"`
	Type           *string `json:"type,omitempty"`
	Status         *string `json:"status,omitempty"`
	Search         string  `json:"search,omitempty"`
}

// CreateIntegrationRequest represents create integration input
type CreateIntegrationRequest struct {
	Name           string  `json:"name" validate:"required,min=2,max=100"`
	Type           string  `json:"type" validate:"required"`
	Description    *string `json:"description,omitempty"`
	Config         string  `json:"config" validate:"required"`
	Endpoint       *string `json:"endpoint,omitempty"`
	APIKey         *string `json:"api_key,omitempty"`
	OrganizationID *string `json:"organization_id,omitempty"`
}

// UpdateIntegrationRequest represents update integration input
type UpdateIntegrationRequest struct {
	Name           *string `json:"name" validate:"omitempty,min=2,max=100"`
	Description    *string `json:"description,omitempty"`
	Config         *string `json:"config,omitempty"`
	Endpoint       *string `json:"endpoint,omitempty"`
	APIKey         *string `json:"api_key,omitempty"`
	Status         *string `json:"status,omitempty"`
}

// IntegrationInfo represents integration information for API responses
type IntegrationInfo struct {
	ID             string     `json:"id"`
	Name           string     `json:"name"`
	Type           string     `json:"type"`
	Description    *string    `json:"description,omitempty"`
	Endpoint       *string    `json:"endpoint,omitempty"`
	Status         string     `json:"status"`
	LastSyncAt     *time.Time `json:"last_sync_at,omitempty"`
	OrganizationID *string    `json:"organization_id,omitempty"`
	CreatedBy      string     `json:"created_by"`
	CreatedAt      time.Time  `json:"created_at"`
	UpdatedAt      time.Time  `json:"updated_at"`
}

// IntegrationListResponse represents paginated integration list
type IntegrationListResponse struct {
	Integrations []IntegrationInfo `json:"integrations"`
	Meta         ListMeta          `json:"meta"`
}

// ListMeta represents pagination metadata
type ListMeta struct {
	Page      int `json:"page"`
	Limit     int `json:"limit"`
	Total     int `json:"total"`
	TotalPage int `json:"total_page"`
}
