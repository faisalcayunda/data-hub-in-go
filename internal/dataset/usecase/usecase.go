package usecase

import (
	"context"

	"portal-data-backend/internal/dataset/domain"
)

// Usecase defines the interface for dataset business logic
type Usecase interface {
	// GetByID retrieves a dataset by ID
	GetByID(ctx context.Context, id string) (*domain.DatasetResponse, error)

	// GetBySlug retrieves a dataset by slug
	GetBySlug(ctx context.Context, slug string) (*domain.DatasetResponse, error)

	// List retrieves a paginated list of datasets
	List(ctx context.Context, req *domain.ListDatasetsRequest) (*domain.DatasetListResponse, error)

	// Create creates a new dataset
	Create(ctx context.Context, req *domain.CreateDatasetRequest, creatorID, orgID string) (*domain.DatasetResponse, error)

	// Update updates an existing dataset
	Update(ctx context.Context, id string, req *domain.UpdateDatasetRequest, updaterID string) (*domain.DatasetResponse, error)

	// Delete soft deletes a dataset
	Delete(ctx context.Context, id string) error

	// UpdateStatus updates dataset status
	UpdateStatus(ctx context.Context, id string, status domain.DatasetStatus) error

	// GetByOrganizationID retrieves datasets by organization ID
	GetByOrganizationID(ctx context.Context, orgID string, page, limit int) (*domain.DatasetListResponse, error)
}
