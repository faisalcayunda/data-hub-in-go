package usecase

import (
	"context"
	"fmt"
	"math"
	"strings"
	"time"

	"portal-data-backend/internal/organization/domain"
	"portal-data-backend/pkg/errors"

	"github.com/google/uuid"
)

// orgUsecase implements Usecase interface
type orgUsecase struct {
	orgRepo domain.Repository
}

// NewOrgUsecase creates a new organization usecase
func NewOrgUsecase(orgRepo domain.Repository) Usecase {
	return &orgUsecase{
		orgRepo: orgRepo,
	}
}

func (u *orgUsecase) GetByID(ctx context.Context, id string) (*domain.OrganizationResponse, error) {
	org, err := u.orgRepo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get organization: %w", err)
	}
	return u.toResponse(org), nil
}

func (u *orgUsecase) GetByCode(ctx context.Context, code string) (*domain.OrganizationResponse, error) {
	org, err := u.orgRepo.GetByCode(ctx, code)
	if err != nil {
		return nil, fmt.Errorf("failed to get organization: %w", err)
	}
	return u.toResponse(org), nil
}

func (u *orgUsecase) GetBySlug(ctx context.Context, slug string) (*domain.OrganizationResponse, error) {
	org, err := u.orgRepo.GetBySlug(ctx, slug)
	if err != nil {
		return nil, fmt.Errorf("failed to get organization: %w", err)
	}
	return u.toResponse(org), nil
}

func (u *orgUsecase) List(ctx context.Context, req *domain.ListOrganizationsRequest) (*domain.OrganizationListResponse, error) {
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

	orgs, total, err := u.orgRepo.List(ctx, req.Status, req.Search, req.Limit, offset, req.SortBy, req.SortOrder)
	if err != nil {
		return nil, fmt.Errorf("failed to list organizations: %w", err)
	}

	responses := make([]domain.OrganizationResponse, len(orgs))
	for i, org := range orgs {
		responses[i] = *u.toResponse(org)
	}

	totalPage := int(math.Ceil(float64(total) / float64(req.Limit)))

	return &domain.OrganizationListResponse{
		Organizations: responses,
		Meta: domain.ListMeta{
			Page:      req.Page,
			Limit:     req.Limit,
			Total:     total,
			TotalPage: totalPage,
		},
	}, nil
}

func (u *orgUsecase) Create(ctx context.Context, req *domain.CreateOrganizationRequest, creatorID string) (*domain.OrganizationResponse, error) {
	// Check if code already exists
	existing, _ := u.orgRepo.GetByCode(ctx, req.Code)
	if existing != nil {
		return nil, errors.ErrAlreadyExists
	}

	now := time.Now()
	org := &domain.Organization{
		ID:        uuid.New().String(),
		Code:      strings.ToUpper(req.Code),
		Name:      req.Name,
		Slug:      u.generateSlug(req.Name),
		Status:    domain.OrgStatusActive,
		CreatedAt: now,
		UpdatedAt: now,
	}

	if creatorID != "" {
		org.CreatedBy = &creatorID
	}
	if req.Description != "" {
		org.Description = &req.Description
	}
	if req.PhoneNumber != "" {
		org.PhoneNumber = &req.PhoneNumber
	}
	if req.Address != "" {
		org.Address = &req.Address
	}
	if req.WebsiteURL != "" {
		org.WebsiteURL = &req.WebsiteURL
	}
	if req.Email != "" {
		org.Email = &req.Email
	}

	if err := u.orgRepo.Create(ctx, org); err != nil {
		return nil, fmt.Errorf("failed to create organization: %w", err)
	}

	return u.toResponse(org), nil
}

func (u *orgUsecase) Update(ctx context.Context, id string, req *domain.UpdateOrganizationRequest, updaterID string) (*domain.OrganizationResponse, error) {
	org, err := u.orgRepo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get organization: %w", err)
	}

	org.Name = req.Name
	org.Slug = u.generateSlug(req.Name)
	org.UpdatedAt = time.Now()

	if updaterID != "" {
		org.UpdatedBy = &updaterID
	}
	if req.Description != "" {
		org.Description = &req.Description
	} else {
		org.Description = nil
	}
	if req.PhoneNumber != "" {
		org.PhoneNumber = &req.PhoneNumber
	} else {
		org.PhoneNumber = nil
	}
	if req.Address != "" {
		org.Address = &req.Address
	} else {
		org.Address = nil
	}
	if req.WebsiteURL != "" {
		org.WebsiteURL = &req.WebsiteURL
	} else {
		org.WebsiteURL = nil
	}
	if req.Email != "" {
		org.Email = &req.Email
	} else {
		org.Email = nil
	}

	if err := u.orgRepo.Update(ctx, org); err != nil {
		return nil, fmt.Errorf("failed to update organization: %w", err)
	}

	return u.toResponse(org), nil
}

func (u *orgUsecase) Delete(ctx context.Context, id string) error {
	if err := u.orgRepo.Delete(ctx, id); err != nil {
		return fmt.Errorf("failed to delete organization: %w", err)
	}
	return nil
}

func (u *orgUsecase) UpdateStatus(ctx context.Context, id string, status domain.OrgStatus) error {
	if err := u.orgRepo.UpdateStatus(ctx, id, status); err != nil {
		return fmt.Errorf("failed to update organization status: %w", err)
	}
	return nil
}

func (u *orgUsecase) toResponse(org *domain.Organization) *domain.OrganizationResponse {
	return &domain.OrganizationResponse{
		ID:             org.ID,
		Code:           org.Code,
		Name:           org.Name,
		Slug:           org.Slug,
		Description:    org.Description,
		LogoURL:        org.LogoURL,
		PhoneNumber:    org.PhoneNumber,
		Address:        org.Address,
		WebsiteURL:     org.WebsiteURL,
		Email:          org.Email,
		TotalDatasets:  org.TotalDatasets,
		PublicDatasets: org.PublicDatasets,
		TotalMapsets:   org.TotalMapsets,
		PublicMapsets:  org.PublicMapsets,
		Status:         string(org.Status),
		CreatedAt:      org.CreatedAt,
		UpdatedAt:      org.UpdatedAt,
	}
}

func (u *orgUsecase) generateSlug(name string) string {
	slug := strings.ToLower(name)
	slug = strings.ReplaceAll(slug, " ", "-")
	slug = strings.ReplaceAll(slug, "_", "-")
	return slug
}
