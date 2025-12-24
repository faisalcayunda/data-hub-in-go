package usecase

import (
	"context"

	"portal-data-backend/internal/unit/domain"
)

type Usecase interface {
	GetByID(ctx context.Context, id string) (*domain.UnitResponse, error)
	List(ctx context.Context, req *domain.ListUnitsRequest) (*domain.UnitListResponse, error)
	Create(ctx context.Context, req *domain.CreateUnitRequest) (*domain.UnitResponse, error)
	Update(ctx context.Context, id string, req *domain.UpdateUnitRequest) (*domain.UnitResponse, error)
	Delete(ctx context.Context, id string) error
}
