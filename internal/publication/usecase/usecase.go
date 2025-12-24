package usecase

import (
	"context"
	"fmt"
	"math"
	"time"

	"portal-data-backend/internal/publication/domain"

	"github.com/google/uuid"
)

type Usecase interface {
	GetByID(ctx context.Context, id string) (*domain.PublicationInfo, error)
	List(ctx context.Context, req *domain.ListPublicationsRequest) (*domain.PublicationListResponse, error)
	Create(ctx context.Context, req *domain.CreatePublicationRequest, userID string) (*domain.PublicationInfo, error)
	Update(ctx context.Context, id string, req *domain.UpdatePublicationRequest, userID string) (*domain.PublicationInfo, error)
	Delete(ctx context.Context, id string) error
	UpdateStatus(ctx context.Context, id string, status string) error
	IncrementViewCount(ctx context.Context, id string) error
	IncrementDownloadCount(ctx context.Context, id string) error
	GetByDatasetID(ctx context.Context, datasetID string, page, limit int) (*domain.PublicationListResponse, error)
	GetByOrganizationID(ctx context.Context, orgID string, page, limit int) (*domain.PublicationListResponse, error)
}

type publicationUsecase struct {
	repo domain.Repository
}

func NewPublicationUsecase(repo domain.Repository) Usecase {
	return &publicationUsecase{
		repo: repo,
	}
}

func (u *publicationUsecase) GetByID(ctx context.Context, id string) (*domain.PublicationInfo, error) {
	pub, err := u.repo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get publication: %w", err)
	}
	// Increment view count
	go u.repo.IncrementViewCount(ctx, id)
	return u.toInfo(pub), nil
}

func (u *publicationUsecase) List(ctx context.Context, req *domain.ListPublicationsRequest) (*domain.PublicationListResponse, error) {
	if req.Page < 1 {
		req.Page = 1
	}
	if req.Limit < 1 {
		req.Limit = 20
	}

	offset := (req.Page - 1) * req.Limit

	filter := &domain.PublicationFilter{
		DatasetID:      req.DatasetID,
		OrganizationID: req.OrganizationID,
		Status:         req.Status,
		IsFeatured:     req.IsFeatured,
		Search:         req.Search,
	}

	pubs, total, err := u.repo.List(ctx, filter, req.Limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to list publications: %w", err)
	}

	infos := make([]domain.PublicationInfo, len(pubs))
	for i, pub := range pubs {
		infos[i] = *u.toInfo(pub)
	}

	totalPage := int(math.Ceil(float64(total) / float64(req.Limit)))

	return &domain.PublicationListResponse{
		Publications: infos,
		Meta: domain.ListMeta{
			Page:      req.Page,
			Limit:     req.Limit,
			Total:     total,
			TotalPage: totalPage,
		},
	}, nil
}

func (u *publicationUsecase) Create(ctx context.Context, req *domain.CreatePublicationRequest, userID string) (*domain.PublicationInfo, error) {
	now := time.Now()
	pub := &domain.Publication{
		ID:             uuid.New().String(),
		Title:          req.Title,
		Description:    req.Description,
		Content:        req.Content,
		DOI:            req.DOI,
		Publisher:      req.Publisher,
		PublishedDate:  req.PublishedDate,
		DatasetID:      req.DatasetID,
		OrganizationID: req.OrganizationID,
		Authors:        req.Authors,
		Tags:           req.Tags,
		Status:         string(domain.PublicationStatusDraft),
		IsFeatured:     req.IsFeatured,
		ViewCount:      0,
		DownloadCount:  0,
		CreatedBy:      userID,
		UpdatedBy:      userID,
		CreatedAt:      now,
		UpdatedAt:      now,
	}

	if err := u.repo.Create(ctx, pub); err != nil {
		return nil, fmt.Errorf("failed to create publication: %w", err)
	}

	return u.toInfo(pub), nil
}

func (u *publicationUsecase) Update(ctx context.Context, id string, req *domain.UpdatePublicationRequest, userID string) (*domain.PublicationInfo, error) {
	existing, err := u.repo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get publication: %w", err)
	}

	// Update fields
	if req.Title != nil {
		existing.Title = *req.Title
	}
	if req.Description != nil {
		existing.Description = req.Description
	}
	if req.Content != nil {
		existing.Content = *req.Content
	}
	if req.DOI != nil {
		existing.DOI = req.DOI
	}
	if req.Publisher != nil {
		existing.Publisher = req.Publisher
	}
	if req.PublishedDate != nil {
		existing.PublishedDate = req.PublishedDate
	}
	if req.DatasetID != nil {
		existing.DatasetID = req.DatasetID
	}
	if req.OrganizationID != nil {
		existing.OrganizationID = req.OrganizationID
	}
	if req.Authors != nil {
		existing.Authors = req.Authors
	}
	if req.Tags != nil {
		existing.Tags = req.Tags
	}
	if req.IsFeatured != nil {
		existing.IsFeatured = *req.IsFeatured
	}
	if req.Status != nil {
		existing.Status = *req.Status
	}
	existing.UpdatedBy = userID
	existing.UpdatedAt = time.Now()

	if err := u.repo.Update(ctx, id, existing); err != nil {
		return nil, fmt.Errorf("failed to update publication: %w", err)
	}

	return u.toInfo(existing), nil
}

func (u *publicationUsecase) Delete(ctx context.Context, id string) error {
	if err := u.repo.Delete(ctx, id); err != nil {
		return fmt.Errorf("failed to delete publication: %w", err)
	}
	return nil
}

func (u *publicationUsecase) UpdateStatus(ctx context.Context, id string, status string) error {
	if err := u.repo.UpdateStatus(ctx, id, status); err != nil {
		return fmt.Errorf("failed to update publication status: %w", err)
	}
	return nil
}

func (u *publicationUsecase) IncrementViewCount(ctx context.Context, id string) error {
	if err := u.repo.IncrementViewCount(ctx, id); err != nil {
		return fmt.Errorf("failed to increment view count: %w", err)
	}
	return nil
}

func (u *publicationUsecase) IncrementDownloadCount(ctx context.Context, id string) error {
	if err := u.repo.IncrementDownloadCount(ctx, id); err != nil {
		return fmt.Errorf("failed to increment download count: %w", err)
	}
	return nil
}

func (u *publicationUsecase) GetByDatasetID(ctx context.Context, datasetID string, page, limit int) (*domain.PublicationListResponse, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 20
	}

	offset := (page - 1) * limit

	pubs, total, err := u.repo.GetByDatasetID(ctx, datasetID, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to get dataset publications: %w", err)
	}

	infos := make([]domain.PublicationInfo, len(pubs))
	for i, pub := range pubs {
		infos[i] = *u.toInfo(pub)
	}

	totalPage := int(math.Ceil(float64(total) / float64(limit)))

	return &domain.PublicationListResponse{
		Publications: infos,
		Meta: domain.ListMeta{
			Page:      page,
			Limit:     limit,
			Total:     total,
			TotalPage: totalPage,
		},
	}, nil
}

func (u *publicationUsecase) GetByOrganizationID(ctx context.Context, orgID string, page, limit int) (*domain.PublicationListResponse, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 20
	}

	offset := (page - 1) * limit

	pubs, total, err := u.repo.GetByOrganizationID(ctx, orgID, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to get organization publications: %w", err)
	}

	infos := make([]domain.PublicationInfo, len(pubs))
	for i, pub := range pubs {
		infos[i] = *u.toInfo(pub)
	}

	totalPage := int(math.Ceil(float64(total) / float64(limit)))

	return &domain.PublicationListResponse{
		Publications: infos,
		Meta: domain.ListMeta{
			Page:      page,
			Limit:     limit,
			Total:     total,
			TotalPage: totalPage,
		},
	}, nil
}

func (u *publicationUsecase) toInfo(pub *domain.Publication) *domain.PublicationInfo {
	return &domain.PublicationInfo{
		ID:            pub.ID,
		Title:         pub.Title,
		Description:   pub.Description,
		Content:       pub.Content,
		DOI:           pub.DOI,
		Publisher:     pub.Publisher,
		PublishedDate: pub.PublishedDate,
		DatasetID:     pub.DatasetID,
		OrganizationID: pub.OrganizationID,
		Authors:       pub.Authors,
		Tags:          pub.Tags,
		Status:        pub.Status,
		IsFeatured:    pub.IsFeatured,
		ViewCount:     pub.ViewCount,
		DownloadCount: pub.DownloadCount,
		CreatedBy:     pub.CreatedBy,
		CreatedAt:     pub.CreatedAt,
		UpdatedAt:     pub.UpdatedAt,
	}
}
