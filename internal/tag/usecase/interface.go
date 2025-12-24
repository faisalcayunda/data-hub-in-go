package usecase

import (
	"context"

	"portal-data-backend/internal/tag/domain"
)

// Usecase defines the interface for tag business logic
type Usecase interface {
	GetByID(ctx context.Context, id string) (*domain.TagResponse, error)
	List(ctx context.Context, req *domain.ListTagsRequest) (*domain.TagListResponse, error)
	Create(ctx context.Context, req *domain.CreateTagRequest) (*domain.TagResponse, error)
	Update(ctx context.Context, id string, req *domain.UpdateTagRequest) (*domain.TagResponse, error)
	Delete(ctx context.Context, id string) error
}

var _ Usecase = (*tagUsecase)(nil)
