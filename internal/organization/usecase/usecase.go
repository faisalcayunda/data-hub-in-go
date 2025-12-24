package usecase

import (
	"context"

	"portal-data-backend/internal/organization/domain"
)

// Usecase defines the interface for organization business logic
type Usecase interface {
	// GetByID retrieves an organization by ID
	GetByID(ctx context.Context, id string) (*domain.OrganizationResponse, error)

	// GetByCode retrieves an organization by code
	GetByCode(ctx context.Context, code string) (*domain.OrganizationResponse, error)

	// GetBySlug retrieves an organization by slug
	GetBySlug(ctx context.Context, slug string) (*domain.OrganizationResponse, error)

	// List retrieves a paginated list of organizations
	List(ctx context.Context, req *domain.ListOrganizationsRequest) (*domain.OrganizationListResponse, error)

	// Create creates a new organization
	Create(ctx context.Context, req *domain.CreateOrganizationRequest, creatorID string) (*domain.OrganizationResponse, error)

	// Update updates an existing organization
	Update(ctx context.Context, id string, req *domain.UpdateOrganizationRequest, updaterID string) (*domain.OrganizationResponse, error)

	// Delete soft deletes an organization
	Delete(ctx context.Context, id string) error

	// UpdateStatus updates organization status
	UpdateStatus(ctx context.Context, id string, status domain.OrgStatus) error
}
