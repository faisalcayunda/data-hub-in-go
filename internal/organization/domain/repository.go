package domain

import (
	"context"
)

// Repository defines the interface for organization data operations
type Repository interface {
	// GetByID retrieves an organization by ID
	GetByID(ctx context.Context, id string) (*Organization, error)

	// GetByCode retrieves an organization by code
	GetByCode(ctx context.Context, code string) (*Organization, error)

	// GetBySlug retrieves an organization by slug
	GetBySlug(ctx context.Context, slug string) (*Organization, error)

	// List retrieves organizations with filters and pagination
	List(ctx context.Context, status, search string, limit, offset int, sortBy, sortOrder string) ([]*Organization, int, error)

	// Create creates a new organization
	Create(ctx context.Context, org *Organization) error

	// Update updates an existing organization
	Update(ctx context.Context, org *Organization) error

	// Delete soft deletes an organization
	Delete(ctx context.Context, id string) error

	// UpdateStatus updates organization status
	UpdateStatus(ctx context.Context, id string, status OrgStatus) error

	// IncrementDatasetCount increments dataset counters
	IncrementDatasetCount(ctx context.Context, id string, isPublic bool) error

	// DecrementDatasetCount decrements dataset counters
	DecrementDatasetCount(ctx context.Context, id string, isPublic bool) error
}
