package usecase

import (
	"context"
	"fmt"
	"math"
	"strings"
	"time"

	"portal-data-backend/internal/dataset/domain"

	"github.com/google/uuid"
)

// datasetUsecase implements Usecase interface
type datasetUsecase struct {
	datasetRepo domain.Repository
}

// NewDatasetUsecase creates a new dataset usecase
func NewDatasetUsecase(datasetRepo domain.Repository) Usecase {
	return &datasetUsecase{
		datasetRepo: datasetRepo,
	}
}

func (u *datasetUsecase) GetByID(ctx context.Context, id string) (*domain.DatasetResponse, error) {
	dataset, err := u.datasetRepo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get dataset: %w", err)
	}
	return u.toResponse(dataset), nil
}

func (u *datasetUsecase) GetBySlug(ctx context.Context, slug string) (*domain.DatasetResponse, error) {
	dataset, err := u.datasetRepo.GetBySlug(ctx, slug)
	if err != nil {
		return nil, fmt.Errorf("failed to get dataset: %w", err)
	}
	return u.toResponse(dataset), nil
}

func (u *datasetUsecase) List(ctx context.Context, req *domain.ListDatasetsRequest) (*domain.DatasetListResponse, error) {
	if req.Page < 1 {
		req.Page = 1
	}
	if req.Limit < 1 {
		req.Limit = 20
	}
	if req.Limit > 100 {
		req.Limit = 100
	}

	offset := (req.Page - 1) * req.Limit

	filter := &domain.DatasetFilter{
		OrganizationID:   req.OrganizationID,
		TopicID:          req.TopicID,
		BusinessFieldID:  req.BusinessFieldID,
		TagID:            req.TagID,
		Status:           req.Status,
		ValidationStatus: req.ValidationStatus,
		Classification:   req.Classification,
		Search:           req.Search,
	}

	sortBy := req.SortBy
	if sortBy == "" {
		sortBy = "created_at"
	}
	sortOrder := req.SortOrder
	if sortOrder == "" {
		sortOrder = "DESC"
	}

	datasets, total, err := u.datasetRepo.List(ctx, filter, req.Limit, offset, sortBy, sortOrder)
	if err != nil {
		return nil, fmt.Errorf("failed to list datasets: %w", err)
	}

	responses := make([]domain.DatasetResponse, len(datasets))
	for i, ds := range datasets {
		responses[i] = *u.toResponse(ds)
	}

	totalPage := int(math.Ceil(float64(total) / float64(req.Limit)))

	return &domain.DatasetListResponse{
		Datasets: responses,
		Meta: domain.ListMeta{
			Page:      req.Page,
			Limit:     req.Limit,
			Total:     total,
			TotalPage: totalPage,
		},
	}, nil
}

func (u *datasetUsecase) Create(ctx context.Context, req *domain.CreateDatasetRequest, creatorID, orgID string) (*domain.DatasetResponse, error) {
	now := time.Now()

	validationStatus := domain.ValidationStatusPending
	if req.ValidationStatus != "" {
		validationStatus = domain.ValidationStatus(req.ValidationStatus)
	}

	dataset := &domain.Dataset{
		ID:               uuid.New().String(),
		Name:             req.Name,
		Slug:             u.generateSlug(req.Name),
		OrganizationID:   orgID,
		Classification:   req.Classification,
		Category:         req.Category,
		DataFixed:        req.DataFixed,
		ValidationStatus: validationStatus,
		CreatedBy:        creatorID,
		IsHighlight:      req.IsHighlight,
		Status:           domain.DatasetStatusDraft,
		CreatedAt:        now,
		UpdatedAt:        now,
	}

	if req.Description != "" {
		dataset.Description = &req.Description
	}
	if req.Period != "" {
		dataset.Period = &req.Period
	}
	if req.UnitID != "" {
		dataset.UnitID = &req.UnitID
	}
	if req.BusinessFieldID != "" {
		dataset.BusinessFieldID = &req.BusinessFieldID
	}
	if req.Image != "" {
		dataset.Image = &req.Image
	}
	if req.TopicID != "" {
		dataset.TopicID = &req.TopicID
	}
	if req.ReferenceID != "" {
		dataset.ReferenceID = &req.ReferenceID
	}
	if req.Metadata != "" {
		dataset.Metadata = &req.Metadata
	}

	if err := u.datasetRepo.Create(ctx, dataset, req.TagIDs); err != nil {
		return nil, fmt.Errorf("failed to create dataset: %w", err)
	}

	// Fetch full dataset with relations
	fullDataset, err := u.datasetRepo.GetByID(ctx, dataset.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch created dataset: %w", err)
	}

	return u.toResponse(fullDataset), nil
}

func (u *datasetUsecase) Update(ctx context.Context, id string, req *domain.UpdateDatasetRequest, updaterID string) (*domain.DatasetResponse, error) {
	dataset, err := u.datasetRepo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get dataset: %w", err)
	}

	dataset.Name = req.Name
	dataset.Slug = u.generateSlug(req.Name)
	dataset.Classification = req.Classification
	dataset.Category = req.Category
	dataset.DataFixed = req.DataFixed
	dataset.IsHighlight = req.IsHighlight
	dataset.UpdatedAt = time.Now()

	if updaterID != "" {
		dataset.UpdatedBy = &updaterID
	}
	if req.Description != "" {
		dataset.Description = &req.Description
	} else {
		dataset.Description = nil
	}
	if req.Period != "" {
		dataset.Period = &req.Period
	} else {
		dataset.Period = nil
	}
	if req.UnitID != "" {
		dataset.UnitID = &req.UnitID
	} else {
		dataset.UnitID = nil
	}
	if req.BusinessFieldID != "" {
		dataset.BusinessFieldID = &req.BusinessFieldID
	} else {
		dataset.BusinessFieldID = nil
	}
	if req.Image != "" {
		dataset.Image = &req.Image
	} else {
		dataset.Image = nil
	}
	if req.TopicID != "" {
		dataset.TopicID = &req.TopicID
	} else {
		dataset.TopicID = nil
	}
	if req.ReferenceID != "" {
		dataset.ReferenceID = &req.ReferenceID
	} else {
		dataset.ReferenceID = nil
	}
	if req.Metadata != "" {
		dataset.Metadata = &req.Metadata
	} else {
		dataset.Metadata = nil
	}

	if req.ValidationStatus != "" {
		dataset.ValidationStatus = domain.ValidationStatus(req.ValidationStatus)
	}

	if err := u.datasetRepo.Update(ctx, dataset, req.TagIDs); err != nil {
		return nil, fmt.Errorf("failed to update dataset: %w", err)
	}

	// Fetch full dataset with relations
	fullDataset, err := u.datasetRepo.GetByID(ctx, dataset.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch updated dataset: %w", err)
	}

	return u.toResponse(fullDataset), nil
}

func (u *datasetUsecase) Delete(ctx context.Context, id string) error {
	if err := u.datasetRepo.Delete(ctx, id); err != nil {
		return fmt.Errorf("failed to delete dataset: %w", err)
	}
	return nil
}

func (u *datasetUsecase) UpdateStatus(ctx context.Context, id string, status domain.DatasetStatus) error {
	if err := u.datasetRepo.UpdateStatus(ctx, id, status); err != nil {
		return fmt.Errorf("failed to update dataset status: %w", err)
	}
	return nil
}

func (u *datasetUsecase) GetByOrganizationID(ctx context.Context, orgID string, page, limit int) (*domain.DatasetListResponse, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 20
	}

	offset := (page - 1) * limit

	datasets, total, err := u.datasetRepo.GetByOrganizationID(ctx, orgID, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to get organization datasets: %w", err)
	}

	responses := make([]domain.DatasetResponse, len(datasets))
	for i, ds := range datasets {
		responses[i] = *u.toResponse(ds)
	}

	totalPage := int(math.Ceil(float64(total) / float64(limit)))

	return &domain.DatasetListResponse{
		Datasets: responses,
		Meta: domain.ListMeta{
			Page:      page,
			Limit:     limit,
			Total:     total,
			TotalPage: totalPage,
		},
	}, nil
}

func (u *datasetUsecase) toResponse(dataset *domain.Dataset) *domain.DatasetResponse {
	resp := &domain.DatasetResponse{
		ID:               dataset.ID,
		Name:             dataset.Name,
		Slug:             dataset.Slug,
		Description:      dataset.Description,
		Period:           dataset.Period,
		OrganizationID:   dataset.OrganizationID,
		ReferenceID:      dataset.ReferenceID,
		Classification:   dataset.Classification,
		Category:         dataset.Category,
		DataFixed:        dataset.DataFixed,
		ValidationStatus: string(dataset.ValidationStatus),
		Metadata:         dataset.Metadata,
		CreatedBy:        dataset.CreatedBy,
		UpdatedBy:        dataset.UpdatedBy,
		CreatedAt:        dataset.CreatedAt,
		UpdatedAt:        dataset.UpdatedAt,
		IsHighlight:      dataset.IsHighlight,
		Status:           string(dataset.Status),
		Tags:             dataset.Tags,
		Unit:             dataset.Unit,
		BusinessField:    dataset.BusinessField,
		Topic:            dataset.Topic,
		Image:            dataset.Image,
	}

	return resp
}

func (u *datasetUsecase) generateSlug(name string) string {
	slug := strings.ToLower(name)
	slug = strings.ReplaceAll(slug, " ", "-")
	slug = strings.ReplaceAll(slug, "_", "-")
	slug = strings.ReplaceAll(slug, "/", "-")
	return slug
}
