package usecase

import (
	"context"

	"portal-data-backend/internal/feedback/domain"
)

type Usecase interface {
	GetByID(ctx context.Context, id string) (*domain.FeedbackResponse, error)
	List(ctx context.Context, req *domain.ListFeedbacksRequest) (*domain.FeedbackListResponse, error)
	Create(ctx context.Context, req *domain.CreateFeedbackRequest, userID string) (*domain.FeedbackResponse, error)
	UpdateStatus(ctx context.Context, id string, status domain.FeedbackStatus) error
	Delete(ctx context.Context, id string) error
}
