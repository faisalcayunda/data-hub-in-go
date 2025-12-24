package domain

import (
	"time"
)

// Organization represents an organization entity
type Organization struct {
	ID              string     `db:"id" json:"id"`
	Code            string     `db:"code" json:"code"`
	Name            string     `db:"name" json:"name"`
	Slug            string     `db:"slug" json:"slug"`
	Description     *string    `db:"description" json:"description,omitempty"`
	LogoURL         *string    `db:"logo_url" json:"logo_url,omitempty"`
	PhoneNumber     *string    `db:"phone_number" json:"phone_number,omitempty"`
	Address         *string    `db:"address" json:"address,omitempty"`
	WebsiteURL      *string    `db:"website_url" json:"website_url,omitempty"`
	Email           *string    `db:"email" json:"email,omitempty"`
	TotalDatasets   int        `db:"total_datasets" json:"total_datasets"`
	PublicDatasets  int        `db:"public_datasets" json:"public_datasets"`
	TotalMapsets    int        `db:"total_mapsets" json:"total_mapsets"`
	PublicMapsets   int        `db:"public_mapsets" json:"public_mapsets"`
	Status          OrgStatus  `db:"status" json:"status"`
	CreatedBy       *string    `db:"created_by" json:"created_by,omitempty"`
	CreatedAt       time.Time  `db:"created_at" json:"created_at"`
	UpdatedBy       *string    `db:"updated_by" json:"updated_by,omitempty"`
	UpdatedAt       time.Time  `db:"updated_at" json:"updated_at"`
}

// OrgStatus represents organization status
type OrgStatus string

const (
	OrgStatusActive   OrgStatus = "active"
	OrgStatusInactive OrgStatus = "inactive"
	OrgStatusSuspended OrgStatus = "suspended"
)

// CreateOrganizationRequest represents organization creation input
type CreateOrganizationRequest struct {
	Code        string `json:"code" validate:"required,min=2,max=20"`
	Name        string `json:"name" validate:"required,min=2"`
	Description string `json:"description,omitempty"`
	PhoneNumber string `json:"phone_number,omitempty"`
	Address     string `json:"address,omitempty"`
	WebsiteURL  string `json:"website_url,omitempty"`
	Email       string `json:"email,omitempty"`
}

// UpdateOrganizationRequest represents organization update input
type UpdateOrganizationRequest struct {
	Name        string `json:"name" validate:"required,min=2"`
	Description string `json:"description,omitempty"`
	PhoneNumber string `json:"phone_number,omitempty"`
	Address     string `json:"address,omitempty"`
	WebsiteURL  string `json:"website_url,omitempty"`
	Email       string `json:"email,omitempty"`
}

// ListOrganizationsRequest represents list organizations input
type ListOrganizationsRequest struct {
	Page      int    `json:"page" validate:"min=1"`
	Limit     int    `json:"limit" validate:"min=1,max=100"`
	Status    string `json:"status,omitempty"`
	Search    string `json:"search,omitempty"`
	SortBy    string `json:"sort_by,omitempty"`
	SortOrder string `json:"sort_order,omitempty"`
}

// OrganizationResponse represents organization response
type OrganizationResponse struct {
	ID             string     `json:"id"`
	Code           string     `json:"code"`
	Name           string     `json:"name"`
	Slug           string     `json:"slug"`
	Description    *string    `json:"description,omitempty"`
	LogoURL        *string    `json:"logo_url,omitempty"`
	PhoneNumber    *string    `json:"phone_number,omitempty"`
	Address        *string    `json:"address,omitempty"`
	WebsiteURL     *string    `json:"website_url,omitempty"`
	Email          *string    `json:"email,omitempty"`
	TotalDatasets  int        `json:"total_datasets"`
	PublicDatasets int        `json:"public_datasets"`
	TotalMapsets   int        `json:"total_mapsets"`
	PublicMapsets  int        `json:"public_mapsets"`
	Status         string     `json:"status"`
	CreatedAt      time.Time  `json:"created_at"`
	UpdatedAt      time.Time  `json:"updated_at"`
}

// OrganizationListResponse represents paginated organization list
type OrganizationListResponse struct {
	Organizations []OrganizationResponse `json:"organizations"`
	Meta          ListMeta               `json:"meta"`
}

// ListMeta represents pagination metadata
type ListMeta struct {
	Page      int `json:"page"`
	Limit     int `json:"limit"`
	Total     int `json:"total"`
	TotalPage int `json:"total_page"`
}
