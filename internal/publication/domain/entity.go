package domain

import "time"

// Publication represents a published data entity
type Publication struct {
	ID              string        `db:"id" json:"id"`
	Title           string        `db:"title" json:"title"`
	Description     *string       `db:"description" json:"description,omitempty"`
	Content         string        `db:"content" json:"content"`
	DOI             *string       `db:"doi" json:"doi,omitempty"`
	Publisher       *string       `db:"publisher" json:"publisher,omitempty"`
	PublishedDate   *time.Time    `db:"published_date" json:"published_date,omitempty"`
	DatasetID       *string       `db:"dataset_id" json:"dataset_id,omitempty"`
	OrganizationID  *string       `db:"organization_id" json:"organization_id,omitempty"`
	Authors         *string       `db:"authors" json:"authors,omitempty"` // JSON array
	Tags            *string       `db:"tags" json:"tags,omitempty"` // JSON array
	Status          string        `db:"status" json:"status"`
	IsFeatured      bool          `db:"is_featured" json:"is_featured"`
	ViewCount       int64         `db:"view_count" json:"view_count"`
	DownloadCount   int64         `db:"download_count" json:"download_count"`
	CreatedBy       string        `db:"created_by" json:"created_by"`
	UpdatedBy       string        `db:"updated_by" json:"updated_by"`
	CreatedAt       time.Time     `db:"created_at" json:"created_at"`
	UpdatedAt       time.Time     `db:"updated_at" json:"updated_at"`
	DeletedAt       *time.Time    `db:"deleted_at" json:"deleted_at,omitempty"`
}

// PublicationStatus represents publication status
type PublicationStatus string

const (
	PublicationStatusDraft     PublicationStatus = "draft"
	PublicationStatusPublished PublicationStatus = "published"
	PublicationStatusArchived  PublicationStatus = "archived"
)

// ListPublicationsRequest represents list publications input
type ListPublicationsRequest struct {
	Page           int     `json:"page" validate:"min=1"`
	Limit          int     `json:"limit" validate:"min=1,max=100"`
	DatasetID      *string `json:"dataset_id,omitempty"`
	OrganizationID *string `json:"organization_id,omitempty"`
	Status         *string `json:"status,omitempty"`
	IsFeatured     *bool   `json:"is_featured,omitempty"`
	Search         string  `json:"search,omitempty"`
}

// CreatePublicationRequest represents create publication input
type CreatePublicationRequest struct {
	Title          string     `json:"title" validate:"required,min=2,max=200"`
	Description    *string    `json:"description,omitempty"`
	Content        string     `json:"content" validate:"required"`
	DOI            *string    `json:"doi,omitempty"`
	Publisher      *string    `json:"publisher,omitempty"`
	PublishedDate  *time.Time `json:"published_date,omitempty"`
	DatasetID      *string    `json:"dataset_id,omitempty"`
	OrganizationID *string    `json:"organization_id,omitempty"`
	Authors        *string    `json:"authors,omitempty"`
	Tags           *string    `json:"tags,omitempty"`
	IsFeatured     bool       `json:"is_featured"`
}

// UpdatePublicationRequest represents update publication input
type UpdatePublicationRequest struct {
	Title          *string    `json:"title" validate:"omitempty,min=2,max=200"`
	Description    *string    `json:"description,omitempty"`
	Content        *string    `json:"content,omitempty"`
	DOI            *string    `json:"doi,omitempty"`
	Publisher      *string    `json:"publisher,omitempty"`
	PublishedDate  *time.Time `json:"published_date,omitempty"`
	DatasetID      *string    `json:"dataset_id,omitempty"`
	OrganizationID *string    `json:"organization_id,omitempty"`
	Authors        *string    `json:"authors,omitempty"`
	Tags           *string    `json:"tags,omitempty"`
	IsFeatured     *bool      `json:"is_featured,omitempty"`
	Status         *string    `json:"status,omitempty"`
}

// PublicationInfo represents publication information for API responses
type PublicationInfo struct {
	ID             string     `json:"id"`
	Title          string     `json:"title"`
	Description    *string    `json:"description,omitempty"`
	Content        string     `json:"content"`
	DOI            *string    `json:"doi,omitempty"`
	Publisher      *string    `json:"publisher,omitempty"`
	PublishedDate  *time.Time `json:"published_date,omitempty"`
	DatasetID      *string    `json:"dataset_id,omitempty"`
	OrganizationID *string    `json:"organization_id,omitempty"`
	Authors        *string    `json:"authors,omitempty"`
	Tags           *string    `json:"tags,omitempty"`
	Status         string     `json:"status"`
	IsFeatured     bool       `json:"is_featured"`
	ViewCount      int64      `json:"view_count"`
	DownloadCount  int64      `json:"download_count"`
	CreatedBy      string     `json:"created_by"`
	CreatedAt      time.Time  `json:"created_at"`
	UpdatedAt      time.Time  `json:"updated_at"`
}

// PublicationListResponse represents paginated publication list
type PublicationListResponse struct {
	Publications []PublicationInfo `json:"publications"`
	Meta         ListMeta          `json:"meta"`
}

// ListMeta represents pagination metadata
type ListMeta struct {
	Page      int `json:"page"`
	Limit     int `json:"limit"`
	Total     int `json:"total"`
	TotalPage int `json:"total_page"`
}
