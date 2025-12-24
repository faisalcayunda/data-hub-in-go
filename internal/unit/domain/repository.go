package domain

import (
	"context"
)

type Repository interface {
	GetByID(ctx context.Context, id string) (*Unit, error)
	List(ctx context.Context, search string, limit, offset int) ([]*Unit, int, error)
	Create(ctx context.Context, unit *Unit) error
	Update(ctx context.Context, unit *Unit) error
	Delete(ctx context.Context, id string) error
}
