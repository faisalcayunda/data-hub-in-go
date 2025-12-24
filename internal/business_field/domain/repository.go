package domain

import (
	"context"
)

// Repository defines the interface for business field data operations
type Repository interface {
	GetByID(ctx context.Context, id string) (*BusinessField, error)
	List(ctx context.Context, search string, limit, offset int) ([]*BusinessField, int, error)
	Create(ctx context.Context, bf *BusinessField) error
	Update(ctx context.Context, bf *BusinessField) error
	Delete(ctx context.Context, id string) error
}
