package domain

import (
	"context"
)

// Repository defines the interface for dataset data operations
type Repository interface {
	// GetByID retrieves a dataset by ID
	GetByID(ctx context.Context, id string) (*Dataset, error)

	// GetBySlug retrieves a dataset by slug
	GetBySlug(ctx context.Context, slug string) (*Dataset, error)

	// List retrieves datasets with filters and pagination
	List(ctx context.Context, filter *DatasetFilter, limit, offset int, sortBy, sortOrder string) ([]*Dataset, int, error)

	// Create creates a new dataset
	Create(ctx context.Context, dataset *Dataset, tagIDs []string) error

	// Update updates an existing dataset
	Update(ctx context.Context, dataset *Dataset, tagIDs []string) error

	// Delete soft deletes a dataset
	Delete(ctx context.Context, id string) error

	// UpdateStatus updates dataset status
	UpdateStatus(ctx context.Context, id string, status DatasetStatus) error

	// GetByOrganizationID retrieves datasets by organization ID
	GetByOrganizationID(ctx context.Context, orgID string, limit, offset int) ([]*Dataset, int, error)
}

// DatasetFilter represents filter options for listing datasets
type DatasetFilter struct {
	OrganizationID   string
	TopicID          string
	BusinessFieldID  string
	TagID            string
	Status           string
	ValidationStatus string
	Classification   string
	Search           string
}
