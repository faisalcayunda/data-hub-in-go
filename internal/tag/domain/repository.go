package domain

import (
	"context"
)

// Repository defines the interface for tag data operations
type Repository interface {
	GetByID(ctx context.Context, id string) (*Tag, error)
	List(ctx context.Context, search string, limit, offset int) ([]*Tag, int, error)
	Create(ctx context.Context, tag *Tag) error
	Update(ctx context.Context, tag *Tag) error
	Delete(ctx context.Context, id string) error
}
