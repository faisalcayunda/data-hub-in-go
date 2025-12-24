package domain

import (
	"context"
)

type Repository interface {
	GetByID(ctx context.Context, id string) (*Integration, error)
	List(ctx context.Context, filter *IntegrationFilter, limit, offset int) ([]*Integration, int, error)
	Create(ctx context.Context, integration *Integration) error
	Update(ctx context.Context, id string, integration *Integration) error
	Delete(ctx context.Context, id string) error
	UpdateStatus(ctx context.Context, id string, status string) error
	Sync(ctx context.Context, id string) error
}

type IntegrationFilter struct {
	OrganizationID *string
	Type           *string
	Status         *string
	Search         string
}
