package usecase

import (
	"context"
	"fmt"
	"math"
	"time"

	"portal-data-backend/internal/unit/domain"

	"github.com/google/uuid"
)

type unitUsecase struct {
	unitRepo domain.Repository
}

func NewUnitUsecase(unitRepo domain.Repository) Usecase {
	return &unitUsecase{unitRepo: unitRepo}
}

func (u *unitUsecase) GetByID(ctx context.Context, id string) (*domain.UnitResponse, error) {
	unit, err := u.unitRepo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get unit: %w", err)
	}
	return u.toResponse(unit), nil
}

func (u *unitUsecase) List(ctx context.Context, req *domain.ListUnitsRequest) (*domain.UnitListResponse, error) {
	if req.Page < 1 {
		req.Page = 1
	}
	if req.Limit < 1 {
		req.Limit = 20
	}

	offset := (req.Page - 1) * req.Limit

	units, total, err := u.unitRepo.List(ctx, req.Search, req.Limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to list units: %w", err)
	}

	responses := make([]domain.UnitResponse, len(units))
	for i, unit := range units {
		responses[i] = *u.toResponse(unit)
	}

	totalPage := int(math.Ceil(float64(total) / float64(req.Limit)))

	return &domain.UnitListResponse{
		Units: responses,
		Meta: domain.ListMeta{
			Page:      req.Page,
			Limit:     req.Limit,
			Total:     total,
			TotalPage: totalPage,
		},
	}, nil
}

func (u *unitUsecase) Create(ctx context.Context, req *domain.CreateUnitRequest) (*domain.UnitResponse, error) {
	unit := &domain.Unit{
		ID:        uuid.New().String(),
		Name:      req.Name,
		Symbol:    req.Symbol,
		CreatedAt: time.Now(),
	}

	if err := u.unitRepo.Create(ctx, unit); err != nil {
		return nil, fmt.Errorf("failed to create unit: %w", err)
	}

	return u.toResponse(unit), nil
}

func (u *unitUsecase) Update(ctx context.Context, id string, req *domain.UpdateUnitRequest) (*domain.UnitResponse, error) {
	unit, err := u.unitRepo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get unit: %w", err)
	}

	unit.Name = req.Name
	unit.Symbol = req.Symbol

	if err := u.unitRepo.Update(ctx, unit); err != nil {
		return nil, fmt.Errorf("failed to update unit: %w", err)
	}

	return u.toResponse(unit), nil
}

func (u *unitUsecase) Delete(ctx context.Context, id string) error {
	if err := u.unitRepo.Delete(ctx, id); err != nil {
		return fmt.Errorf("failed to delete unit: %w", err)
	}
	return nil
}

func (u *unitUsecase) toResponse(unit *domain.Unit) *domain.UnitResponse {
	return &domain.UnitResponse{
		ID:        unit.ID,
		Name:      unit.Name,
		Symbol:    unit.Symbol,
		CreatedAt: unit.CreatedAt,
	}
}
