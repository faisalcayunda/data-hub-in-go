package domain

import (
	"context"
)

type Repository interface {
	GetByID(ctx context.Context, id string) (*Topic, error)
	List(ctx context.Context, search string, limit, offset int) ([]*Topic, int, error)
	Create(ctx context.Context, topic *Topic) error
	Update(ctx context.Context, topic *Topic) error
	Delete(ctx context.Context, id string) error
}
