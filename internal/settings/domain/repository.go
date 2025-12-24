package domain

import (
	"context"
)

type Repository interface {
	GetByID(ctx context.Context, id string) (*Setting, error)
	GetByKey(ctx context.Context, key string, userID *string) (*Setting, error)
	List(ctx context.Context, filter *SettingFilter, limit, offset int) ([]*Setting, int, error)
	Create(ctx context.Context, setting *Setting) error
	Update(ctx context.Context, id string, setting *Setting) error
	Delete(ctx context.Context, id string) error
	GetByKeys(ctx context.Context, keys []string, userID *string) (map[string]string, error)
	GetByCategory(ctx context.Context, category string, userID *string, limit, offset int) ([]*Setting, int, error)
}

type SettingFilter struct {
	Category *string
	UserID   *string
	Type     *string
	Search   string
}
