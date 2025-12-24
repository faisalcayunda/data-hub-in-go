package usecase

import (
	"context"
	"fmt"
	"math"
	"time"

	"portal-data-backend/internal/integration/domain"

	"github.com/google/uuid"
)

type Usecase interface {
	GetByID(ctx context.Context, id string) (*domain.IntegrationInfo, error)
	List(ctx context.Context, req *domain.ListIntegrationsRequest) (*domain.IntegrationListResponse, error)
	Create(ctx context.Context, req *domain.CreateIntegrationRequest, userID string) (*domain.IntegrationInfo, error)
	Update(ctx context.Context, id string, req *domain.UpdateIntegrationRequest) (*domain.IntegrationInfo, error)
	Delete(ctx context.Context, id string) error
	UpdateStatus(ctx context.Context, id string, status string) error
	Sync(ctx context.Context, id string) error
}

type integrationUsecase struct {
	repo domain.Repository
}

func NewIntegrationUsecase(repo domain.Repository) Usecase {
	return &integrationUsecase{
		repo: repo,
	}
}

func (u *integrationUsecase) GetByID(ctx context.Context, id string) (*domain.IntegrationInfo, error) {
	integration, err := u.repo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get integration: %w", err)
	}
	return u.toInfo(integration), nil
}

func (u *integrationUsecase) List(ctx context.Context, req *domain.ListIntegrationsRequest) (*domain.IntegrationListResponse, error) {
	if req.Page < 1 {
		req.Page = 1
	}
	if req.Limit < 1 {
		req.Limit = 20
	}

	offset := (req.Page - 1) * req.Limit

	filter := &domain.IntegrationFilter{
		OrganizationID: req.OrganizationID,
		Type:           req.Type,
		Status:         req.Status,
		Search:         req.Search,
	}

	integrations, total, err := u.repo.List(ctx, filter, req.Limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to list integrations: %w", err)
	}

	infos := make([]domain.IntegrationInfo, len(integrations))
	for i, integration := range integrations {
		infos[i] = *u.toInfo(integration)
	}

	totalPage := int(math.Ceil(float64(total) / float64(req.Limit)))

	return &domain.IntegrationListResponse{
		Integrations: infos,
		Meta: domain.ListMeta{
			Page:      req.Page,
			Limit:     req.Limit,
			Total:     total,
			TotalPage: totalPage,
		},
	}, nil
}

func (u *integrationUsecase) Create(ctx context.Context, req *domain.CreateIntegrationRequest, userID string) (*domain.IntegrationInfo, error) {
	now := time.Now()
	integration := &domain.Integration{
		ID:             uuid.New().String(),
		Name:           req.Name,
		Type:           req.Type,
		Description:    req.Description,
		Config:         req.Config,
		Endpoint:       req.Endpoint,
		APIKey:         req.APIKey,
		Status:         string(domain.IntegrationStatusActive),
		OrganizationID: req.OrganizationID,
		CreatedBy:      userID,
		CreatedAt:      now,
		UpdatedAt:      now,
	}

	if err := u.repo.Create(ctx, integration); err != nil {
		return nil, fmt.Errorf("failed to create integration: %w", err)
	}

	return u.toInfo(integration), nil
}

func (u *integrationUsecase) Update(ctx context.Context, id string, req *domain.UpdateIntegrationRequest) (*domain.IntegrationInfo, error) {
	existing, err := u.repo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get integration: %w", err)
	}

	// Update fields
	if req.Name != nil {
		existing.Name = *req.Name
	}
	if req.Description != nil {
		existing.Description = req.Description
	}
	if req.Config != nil {
		existing.Config = *req.Config
	}
	if req.Endpoint != nil {
		existing.Endpoint = req.Endpoint
	}
	if req.APIKey != nil {
		existing.APIKey = req.APIKey
	}
	if req.Status != nil {
		existing.Status = *req.Status
	}
	existing.UpdatedAt = time.Now()

	if err := u.repo.Update(ctx, id, existing); err != nil {
		return nil, fmt.Errorf("failed to update integration: %w", err)
	}

	return u.toInfo(existing), nil
}

func (u *integrationUsecase) Delete(ctx context.Context, id string) error {
	if err := u.repo.Delete(ctx, id); err != nil {
		return fmt.Errorf("failed to delete integration: %w", err)
	}
	return nil
}

func (u *integrationUsecase) UpdateStatus(ctx context.Context, id string, status string) error {
	if err := u.repo.UpdateStatus(ctx, id, status); err != nil {
		return fmt.Errorf("failed to update integration status: %w", err)
	}
	return nil
}

func (u *integrationUsecase) Sync(ctx context.Context, id string) error {
	if err := u.repo.Sync(ctx, id); err != nil {
		return fmt.Errorf("failed to sync integration: %w", err)
	}
	return nil
}

func (u *integrationUsecase) toInfo(integration *domain.Integration) *domain.IntegrationInfo {
	return &domain.IntegrationInfo{
		ID:             integration.ID,
		Name:           integration.Name,
		Type:           integration.Type,
		Description:    integration.Description,
		Endpoint:       integration.Endpoint,
		Status:         integration.Status,
		LastSyncAt:     integration.LastSyncAt,
		OrganizationID: integration.OrganizationID,
		CreatedBy:      integration.CreatedBy,
		CreatedAt:      integration.CreatedAt,
		UpdatedAt:      integration.UpdatedAt,
	}
}
