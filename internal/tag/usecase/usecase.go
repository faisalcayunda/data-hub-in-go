package usecase

import (
	"context"
	"fmt"
	"math"
	"strings"
	"time"

	"portal-data-backend/internal/tag/domain"

	"github.com/google/uuid"
)

type tagUsecase struct {
	tagRepo domain.Repository
}

func NewTagUsecase(tagRepo domain.Repository) *tagUsecase {
	return &tagUsecase{tagRepo: tagRepo}
}

func (u *tagUsecase) GetByID(ctx context.Context, id string) (*domain.TagResponse, error) {
	tag, err := u.tagRepo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get tag: %w", err)
	}
	return u.toResponse(tag), nil
}

func (u *tagUsecase) List(ctx context.Context, req *domain.ListTagsRequest) (*domain.TagListResponse, error) {
	if req.Page < 1 {
		req.Page = 1
	}
	if req.Limit < 1 {
		req.Limit = 20
	}

	offset := (req.Page - 1) * req.Limit

	tags, total, err := u.tagRepo.List(ctx, req.Search, req.Limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to list tags: %w", err)
	}

	responses := make([]domain.TagResponse, len(tags))
	for i, tag := range tags {
		responses[i] = *u.toResponse(tag)
	}

	totalPage := int(math.Ceil(float64(total) / float64(req.Limit)))

	return &domain.TagListResponse{
		Tags: responses,
		Meta: domain.ListMeta{
			Page:      req.Page,
			Limit:     req.Limit,
			Total:     total,
			TotalPage: totalPage,
		},
	}, nil
}

func (u *tagUsecase) Create(ctx context.Context, req *domain.CreateTagRequest) (*domain.TagResponse, error) {
	tag := &domain.Tag{
		ID:        uuid.New().String(),
		Name:      req.Name,
		Slug:      u.generateSlug(req.Name),
		CreatedAt: time.Now(),
	}

	if err := u.tagRepo.Create(ctx, tag); err != nil {
		return nil, fmt.Errorf("failed to create tag: %w", err)
	}

	return u.toResponse(tag), nil
}

func (u *tagUsecase) Update(ctx context.Context, id string, req *domain.UpdateTagRequest) (*domain.TagResponse, error) {
	tag, err := u.tagRepo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get tag: %w", err)
	}

	tag.Name = req.Name
	tag.Slug = u.generateSlug(req.Name)

	if err := u.tagRepo.Update(ctx, tag); err != nil {
		return nil, fmt.Errorf("failed to update tag: %w", err)
	}

	return u.toResponse(tag), nil
}

func (u *tagUsecase) Delete(ctx context.Context, id string) error {
	if err := u.tagRepo.Delete(ctx, id); err != nil {
		return fmt.Errorf("failed to delete tag: %w", err)
	}
	return nil
}

func (u *tagUsecase) toResponse(tag *domain.Tag) *domain.TagResponse {
	return &domain.TagResponse{
		ID:        tag.ID,
		Name:      tag.Name,
		Slug:      tag.Slug,
		CreatedAt: tag.CreatedAt,
	}
}

func (u *tagUsecase) generateSlug(name string) string {
	slug := strings.ToLower(name)
	slug = strings.ReplaceAll(slug, " ", "-")
	return slug
}
