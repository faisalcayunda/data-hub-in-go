package usecase

import (
	"context"
	"fmt"
	"math"
	"time"

	"portal-data-backend/internal/settings/domain"

	"github.com/google/uuid"
)

type Usecase interface {
	GetByID(ctx context.Context, id string) (*domain.SettingInfo, error)
	GetByKey(ctx context.Context, key string, userID *string) (*domain.SettingInfo, error)
	List(ctx context.Context, req *domain.ListSettingsRequest) (*domain.SettingListResponse, error)
	Create(ctx context.Context, req *domain.CreateSettingRequest) (*domain.SettingInfo, error)
	Update(ctx context.Context, id string, req *domain.UpdateSettingRequest) (*domain.SettingInfo, error)
	Delete(ctx context.Context, id string) error
	GetByKeys(ctx context.Context, keys []string, userID *string) (map[string]string, error)
	GetByCategory(ctx context.Context, category string, userID *string, page, limit int) (*domain.SettingListResponse, error)
}

type settingsUsecase struct {
	repo domain.Repository
}

func NewSettingsUsecase(repo domain.Repository) Usecase {
	return &settingsUsecase{
		repo: repo,
	}
}

func (u *settingsUsecase) GetByID(ctx context.Context, id string) (*domain.SettingInfo, error) {
	setting, err := u.repo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get setting: %w", err)
	}
	return u.toInfo(setting), nil
}

func (u *settingsUsecase) GetByKey(ctx context.Context, key string, userID *string) (*domain.SettingInfo, error) {
	setting, err := u.repo.GetByKey(ctx, key, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get setting: %w", err)
	}
	return u.toInfo(setting), nil
}

func (u *settingsUsecase) List(ctx context.Context, req *domain.ListSettingsRequest) (*domain.SettingListResponse, error) {
	if req.Page < 1 {
		req.Page = 1
	}
	if req.Limit < 1 {
		req.Limit = 20
	}

	offset := (req.Page - 1) * req.Limit

	filter := &domain.SettingFilter{
		Category: req.Category,
		UserID:   req.UserID,
		Type:     req.Type,
		Search:   req.Search,
	}

	settings, total, err := u.repo.List(ctx, filter, req.Limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to list settings: %w", err)
	}

	infos := make([]domain.SettingInfo, len(settings))
	for i, setting := range settings {
		infos[i] = *u.toInfo(setting)
	}

	totalPage := int(math.Ceil(float64(total) / float64(req.Limit)))

	return &domain.SettingListResponse{
		Settings: infos,
		Meta: domain.ListMeta{
			Page:      req.Page,
			Limit:     req.Limit,
			Total:     total,
			TotalPage: totalPage,
		},
	}, nil
}

func (u *settingsUsecase) Create(ctx context.Context, req *domain.CreateSettingRequest) (*domain.SettingInfo, error) {
	now := time.Now()
	setting := &domain.Setting{
		ID:        uuid.New().String(),
		Key:       req.Key,
		Value:     req.Value,
		Type:      req.Type,
		Category:  req.Category,
		UserID:    req.UserID,
		IsPublic:  req.IsPublic,
		CreatedAt: now,
		UpdatedAt: now,
	}

	if err := u.repo.Create(ctx, setting); err != nil {
		return nil, fmt.Errorf("failed to create setting: %w", err)
	}

	return u.toInfo(setting), nil
}

func (u *settingsUsecase) Update(ctx context.Context, id string, req *domain.UpdateSettingRequest) (*domain.SettingInfo, error) {
	existing, err := u.repo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get setting: %w", err)
	}

	// Update fields
	if req.Value != nil {
		existing.Value = *req.Value
	}
	if req.Type != nil {
		existing.Type = *req.Type
	}
	if req.IsPublic != nil {
		existing.IsPublic = *req.IsPublic
	}
	existing.UpdatedAt = time.Now()

	if err := u.repo.Update(ctx, id, existing); err != nil {
		return nil, fmt.Errorf("failed to update setting: %w", err)
	}

	return u.toInfo(existing), nil
}

func (u *settingsUsecase) Delete(ctx context.Context, id string) error {
	if err := u.repo.Delete(ctx, id); err != nil {
		return fmt.Errorf("failed to delete setting: %w", err)
	}
	return nil
}

func (u *settingsUsecase) GetByKeys(ctx context.Context, keys []string, userID *string) (map[string]string, error) {
	settings, err := u.repo.GetByKeys(ctx, keys, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get settings: %w", err)
	}
	return settings, nil
}

func (u *settingsUsecase) GetByCategory(ctx context.Context, category string, userID *string, page, limit int) (*domain.SettingListResponse, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 20
	}

	offset := (page - 1) * limit

	settings, total, err := u.repo.GetByCategory(ctx, category, userID, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to get settings by category: %w", err)
	}

	infos := make([]domain.SettingInfo, len(settings))
	for i, setting := range settings {
		infos[i] = *u.toInfo(setting)
	}

	totalPage := int(math.Ceil(float64(total) / float64(limit)))

	return &domain.SettingListResponse{
		Settings: infos,
		Meta: domain.ListMeta{
			Page:      page,
			Limit:     limit,
			Total:     total,
			TotalPage: totalPage,
		},
	}, nil
}

func (u *settingsUsecase) toInfo(setting *domain.Setting) *domain.SettingInfo {
	return &domain.SettingInfo{
		ID:        setting.ID,
		Key:       setting.Key,
		Value:     setting.Value,
		Type:      setting.Type,
		Category:  setting.Category,
		UserID:    setting.UserID,
		IsPublic:  setting.IsPublic,
		CreatedAt: setting.CreatedAt,
		UpdatedAt: setting.UpdatedAt,
	}
}
