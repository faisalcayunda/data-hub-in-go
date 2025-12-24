package domain

import (
	"context"
)

type Repository interface {
	GetByID(ctx context.Context, id string) (*Visualization, error)
	List(ctx context.Context, filter *VisualizationFilter, limit, offset int) ([]*Visualization, int, error)
	Create(ctx context.Context, viz *Visualization) error
	Update(ctx context.Context, id string, viz *Visualization) error
	Delete(ctx context.Context, id string) error
	UpdateStatus(ctx context.Context, id string, status string) error
	GetStats(ctx context.Context) (*VisualizationStats, error)
	GetByDatasetID(ctx context.Context, datasetID string, limit, offset int) ([]*Visualization, int, error)
	GetByOrganizationID(ctx context.Context, orgID string, limit, offset int) ([]*Visualization, int, error)
}

type VisualizationFilter struct {
	DatasetID      *string
	OrganizationID *string
	TopicID        *string
	Type           *string
	Status         *string
	IsHighlight    *bool
	Search         string
}
