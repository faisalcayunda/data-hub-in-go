package usecase

import (
	"context"
	"fmt"
	"math"
	"time"

	"portal-data-backend/internal/visualization/domain"

	"github.com/google/uuid"
)

type Usecase interface {
	GetByID(ctx context.Context, id string) (*domain.VisualizationInfo, error)
	List(ctx context.Context, req *domain.ListVisualizationsRequest) (*domain.VisualizationListResponse, error)
	Create(ctx context.Context, req *domain.CreateVisualizationRequest, userID string) (*domain.VisualizationInfo, error)
	Update(ctx context.Context, id string, req *domain.UpdateVisualizationRequest, userID string) (*domain.VisualizationInfo, error)
	Delete(ctx context.Context, id string) error
	UpdateStatus(ctx context.Context, id string, status string) error
	GetStats(ctx context.Context) (*domain.VisualizationStats, error)
	GetByDatasetID(ctx context.Context, datasetID string, page, limit int) (*domain.VisualizationListResponse, error)
	GetByOrganizationID(ctx context.Context, orgID string, page, limit int) (*domain.VisualizationListResponse, error)
}

type visualizationUsecase struct {
	repo domain.Repository
}

func NewVisualizationUsecase(repo domain.Repository) Usecase {
	return &visualizationUsecase{
		repo: repo,
	}
}

func (u *visualizationUsecase) GetByID(ctx context.Context, id string) (*domain.VisualizationInfo, error) {
	viz, err := u.repo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get visualization: %w", err)
	}
	return u.toInfo(viz), nil
}

func (u *visualizationUsecase) List(ctx context.Context, req *domain.ListVisualizationsRequest) (*domain.VisualizationListResponse, error) {
	if req.Page < 1 {
		req.Page = 1
	}
	if req.Limit < 1 {
		req.Limit = 20
	}

	offset := (req.Page - 1) * req.Limit

	filter := &domain.VisualizationFilter{
		DatasetID:      req.DatasetID,
		OrganizationID: req.OrganizationID,
		TopicID:        req.TopicID,
		Type:           req.Type,
		Status:         req.Status,
		IsHighlight:    req.IsHighlight,
		Search:         req.Search,
	}

	vizs, total, err := u.repo.List(ctx, filter, req.Limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to list visualizations: %w", err)
	}

	infos := make([]domain.VisualizationInfo, len(vizs))
	for i, viz := range vizs {
		infos[i] = *u.toInfo(viz)
	}

	totalPage := int(math.Ceil(float64(total) / float64(req.Limit)))

	return &domain.VisualizationListResponse{
		Visualizations: infos,
		Meta: domain.ListMeta{
			Page:      req.Page,
			Limit:     req.Limit,
			Total:     total,
			TotalPage: totalPage,
		},
	}, nil
}

func (u *visualizationUsecase) Create(ctx context.Context, req *domain.CreateVisualizationRequest, userID string) (*domain.VisualizationInfo, error) {
	now := time.Now()
	viz := &domain.Visualization{
		ID:             uuid.New().String(),
		Title:          req.Title,
		Description:    req.Description,
		Type:           req.Type,
		Config:         req.Config,
		DatasetID:      req.DatasetID,
		OrganizationID: req.OrganizationID,
		TopicID:        req.TopicID,
		IsHighlight:    req.IsHighlight,
		Status:         string(domain.VisualizationStatusDraft),
		CreatedBy:      userID,
		UpdatedBy:      userID,
		CreatedAt:      now,
		UpdatedAt:      now,
	}

	if err := u.repo.Create(ctx, viz); err != nil {
		return nil, fmt.Errorf("failed to create visualization: %w", err)
	}

	return u.toInfo(viz), nil
}

func (u *visualizationUsecase) Update(ctx context.Context, id string, req *domain.UpdateVisualizationRequest, userID string) (*domain.VisualizationInfo, error) {
	existing, err := u.repo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get visualization: %w", err)
	}

	// Update fields
	if req.Title != nil {
		existing.Title = *req.Title
	}
	if req.Description != nil {
		existing.Description = req.Description
	}
	if req.Type != nil {
		existing.Type = *req.Type
	}
	if req.Config != nil {
		existing.Config = *req.Config
	}
	if req.DatasetID != nil {
		existing.DatasetID = req.DatasetID
	}
	if req.OrganizationID != nil {
		existing.OrganizationID = req.OrganizationID
	}
	if req.TopicID != nil {
		existing.TopicID = req.TopicID
	}
	if req.IsHighlight != nil {
		existing.IsHighlight = *req.IsHighlight
	}
	if req.Status != nil {
		existing.Status = *req.Status
	}
	existing.UpdatedBy = userID
	existing.UpdatedAt = time.Now()

	if err := u.repo.Update(ctx, id, existing); err != nil {
		return nil, fmt.Errorf("failed to update visualization: %w", err)
	}

	return u.toInfo(existing), nil
}

func (u *visualizationUsecase) Delete(ctx context.Context, id string) error {
	if err := u.repo.Delete(ctx, id); err != nil {
		return fmt.Errorf("failed to delete visualization: %w", err)
	}
	return nil
}

func (u *visualizationUsecase) UpdateStatus(ctx context.Context, id string, status string) error {
	if err := u.repo.UpdateStatus(ctx, id, status); err != nil {
		return fmt.Errorf("failed to update visualization status: %w", err)
	}
	return nil
}

func (u *visualizationUsecase) GetStats(ctx context.Context) (*domain.VisualizationStats, error) {
	stats, err := u.repo.GetStats(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get visualization stats: %w", err)
	}
	return stats, nil
}

func (u *visualizationUsecase) GetByDatasetID(ctx context.Context, datasetID string, page, limit int) (*domain.VisualizationListResponse, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 20
	}

	offset := (page - 1) * limit

	vizs, total, err := u.repo.GetByDatasetID(ctx, datasetID, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to get dataset visualizations: %w", err)
	}

	infos := make([]domain.VisualizationInfo, len(vizs))
	for i, viz := range vizs {
		infos[i] = *u.toInfo(viz)
	}

	totalPage := int(math.Ceil(float64(total) / float64(limit)))

	return &domain.VisualizationListResponse{
		Visualizations: infos,
		Meta: domain.ListMeta{
			Page:      page,
			Limit:     limit,
			Total:     total,
			TotalPage: totalPage,
		},
	}, nil
}

func (u *visualizationUsecase) GetByOrganizationID(ctx context.Context, orgID string, page, limit int) (*domain.VisualizationListResponse, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 20
	}

	offset := (page - 1) * limit

	vizs, total, err := u.repo.GetByOrganizationID(ctx, orgID, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to get organization visualizations: %w", err)
	}

	infos := make([]domain.VisualizationInfo, len(vizs))
	for i, viz := range vizs {
		infos[i] = *u.toInfo(viz)
	}

	totalPage := int(math.Ceil(float64(total) / float64(limit)))

	return &domain.VisualizationListResponse{
		Visualizations: infos,
		Meta: domain.ListMeta{
			Page:      page,
			Limit:     limit,
			Total:     total,
			TotalPage: totalPage,
		},
	}, nil
}

func (u *visualizationUsecase) toInfo(viz *domain.Visualization) *domain.VisualizationInfo {
	return &domain.VisualizationInfo{
		ID:             viz.ID,
		Title:          viz.Title,
		Description:    viz.Description,
		Type:           viz.Type,
		Config:         viz.Config,
		DatasetID:      viz.DatasetID,
		OrganizationID: viz.OrganizationID,
		TopicID:        viz.TopicID,
		IsHighlight:    viz.IsHighlight,
		Status:         viz.Status,
		CreatedBy:      viz.CreatedBy,
		CreatedAt:      viz.CreatedAt,
		UpdatedAt:      viz.UpdatedAt,
	}
}
