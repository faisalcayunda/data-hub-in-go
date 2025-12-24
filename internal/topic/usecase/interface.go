package usecase

import (
	"context"

	"portal-data-backend/internal/topic/domain"
)

type Usecase interface {
	GetByID(ctx context.Context, id string) (*domain.TopicResponse, error)
	List(ctx context.Context, req *domain.ListTopicsRequest) (*domain.TopicListResponse, error)
	Create(ctx context.Context, req *domain.CreateTopicRequest) (*domain.TopicResponse, error)
	Update(ctx context.Context, id string, req *domain.UpdateTopicRequest) (*domain.TopicResponse, error)
	Delete(ctx context.Context, id string) error
}
