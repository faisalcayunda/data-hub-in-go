package usecase

import (
	"context"
	"fmt"
	"math"
	"strings"
	"time"

	"portal-data-backend/internal/business_field/domain"

	"github.com/google/uuid"
)

type businessFieldUsecase struct {
	bfRepo domain.Repository
}

func NewBusinessFieldUsecase(bfRepo domain.Repository) Usecase {
	return &businessFieldUsecase{bfRepo: bfRepo}
}

func (u *businessFieldUsecase) GetByID(ctx context.Context, id string) (*domain.BusinessFieldResponse, error) {
	bf, err := u.bfRepo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get business field: %w", err)
	}
	return u.toResponse(bf), nil
}

func (u *businessFieldUsecase) List(ctx context.Context, req *domain.ListBusinessFieldsRequest) (*domain.BusinessFieldListResponse, error) {
	if req.Page < 1 {
		req.Page = 1
	}
	if req.Limit < 1 {
		req.Limit = 20
	}

	offset := (req.Page - 1) * req.Limit

	bfs, total, err := u.bfRepo.List(ctx, req.Search, req.Limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to list business fields: %w", err)
	}

	responses := make([]domain.BusinessFieldResponse, len(bfs))
	for i, bf := range bfs {
		responses[i] = *u.toResponse(bf)
	}

	totalPage := int(math.Ceil(float64(total) / float64(req.Limit)))

	return &domain.BusinessFieldListResponse{
		BusinessFields: responses,
		Meta: domain.ListMeta{
			Page:      req.Page,
			Limit:     req.Limit,
			Total:     total,
			TotalPage: totalPage,
		},
	}, nil
}

func (u *businessFieldUsecase) Create(ctx context.Context, req *domain.CreateBusinessFieldRequest) (*domain.BusinessFieldResponse, error) {
	bf := &domain.BusinessField{
		ID:        uuid.New().String(),
		Name:      req.Name,
		Slug:      u.generateSlug(req.Name),
		CreatedAt: time.Now(),
	}

	if err := u.bfRepo.Create(ctx, bf); err != nil {
		return nil, fmt.Errorf("failed to create business field: %w", err)
	}

	return u.toResponse(bf), nil
}

func (u *businessFieldUsecase) Update(ctx context.Context, id string, req *domain.UpdateBusinessFieldRequest) (*domain.BusinessFieldResponse, error) {
	bf, err := u.bfRepo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get business field: %w", err)
	}

	bf.Name = req.Name
	bf.Slug = u.generateSlug(req.Name)

	if err := u.bfRepo.Update(ctx, bf); err != nil {
		return nil, fmt.Errorf("failed to update business field: %w", err)
	}

	return u.toResponse(bf), nil
}

func (u *businessFieldUsecase) Delete(ctx context.Context, id string) error {
	if err := u.bfRepo.Delete(ctx, id); err != nil {
		return fmt.Errorf("failed to delete business field: %w", err)
	}
	return nil
}

func (u *businessFieldUsecase) toResponse(bf *domain.BusinessField) *domain.BusinessFieldResponse {
	return &domain.BusinessFieldResponse{
		ID:        bf.ID,
		Name:      bf.Name,
		Slug:      bf.Slug,
		CreatedAt: bf.CreatedAt,
	}
}

func (u *businessFieldUsecase) generateSlug(name string) string {
	slug := strings.ToLower(name)
	slug = strings.ReplaceAll(slug, " ", "-")
	return slug
}
