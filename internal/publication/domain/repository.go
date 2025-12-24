package domain

import (
	"context"
)

type Repository interface {
	GetByID(ctx context.Context, id string) (*Publication, error)
	List(ctx context.Context, filter *PublicationFilter, limit, offset int) ([]*Publication, int, error)
	Create(ctx context.Context, pub *Publication) error
	Update(ctx context.Context, id string, pub *Publication) error
	Delete(ctx context.Context, id string) error
	UpdateStatus(ctx context.Context, id string, status string) error
	IncrementViewCount(ctx context.Context, id string) error
	IncrementDownloadCount(ctx context.Context, id string) error
	GetByDatasetID(ctx context.Context, datasetID string, limit, offset int) ([]*Publication, int, error)
	GetByOrganizationID(ctx context.Context, orgID string, limit, offset int) ([]*Publication, int, error)
}

type PublicationFilter struct {
	DatasetID      *string
	OrganizationID *string
	Status         *string
	IsFeatured     *bool
	Search         string
}
