package usecase

import (
	"context"

	"portal-data-backend/internal/business_field/domain"
)

// Usecase defines the interface for business field business logic
type Usecase interface {
	GetByID(ctx context.Context, id string) (*domain.BusinessFieldResponse, error)
	List(ctx context.Context, req *domain.ListBusinessFieldsRequest) (*domain.BusinessFieldListResponse, error)
	Create(ctx context.Context, req *domain.CreateBusinessFieldRequest) (*domain.BusinessFieldResponse, error)
	Update(ctx context.Context, id string, req *domain.UpdateBusinessFieldRequest) (*domain.BusinessFieldResponse, error)
	Delete(ctx context.Context, id string) error
}
