package usecase

import (
	"context"
	"fmt"
	"math"
	"time"

	"portal-data-backend/internal/data_row/domain"

	"github.com/google/uuid"
)

type Usecase interface {
	GetByID(ctx context.Context, id string) (*domain.DataRowInfo, error)
	List(ctx context.Context, req *domain.ListDataRowsRequest) (*domain.DataRowListResponse, error)
	Create(ctx context.Context, req *domain.CreateDataRowRequest, userID string) (*domain.DataRowInfo, error)
	BulkCreate(ctx context.Context, req *domain.BulkCreateDataRowsRequest, userID string) error
	Update(ctx context.Context, id string, req *domain.UpdateDataRowRequest) (*domain.DataRowInfo, error)
	Delete(ctx context.Context, id string) error
	DeleteByDatasetID(ctx context.Context, datasetID string) error
	GetStats(ctx context.Context, datasetID string) (*domain.DataRowStats, error)
}

type dataRowUsecase struct {
	repo domain.Repository
}

func NewDataRowUsecase(repo domain.Repository) Usecase {
	return &dataRowUsecase{
		repo: repo,
	}
}

func (u *dataRowUsecase) GetByID(ctx context.Context, id string) (*domain.DataRowInfo, error) {
	row, err := u.repo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get data row: %w", err)
	}
	return u.toInfo(row), nil
}

func (u *dataRowUsecase) List(ctx context.Context, req *domain.ListDataRowsRequest) (*domain.DataRowListResponse, error) {
	if req.Page < 1 {
		req.Page = 1
	}
	if req.Limit < 1 {
		req.Limit = 20
	}

	offset := (req.Page - 1) * req.Limit

	filter := &domain.DataRowFilter{
		DatasetID: req.DatasetID,
		Search:    req.Search,
	}

	rows, total, err := u.repo.List(ctx, filter, req.Limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to list data rows: %w", err)
	}

	infos := make([]domain.DataRowInfo, len(rows))
	for i, row := range rows {
		infos[i] = *u.toInfo(row)
	}

	totalPage := int(math.Ceil(float64(total) / float64(req.Limit)))

	return &domain.DataRowListResponse{
		Rows: infos,
		Meta: domain.ListMeta{
			Page:      req.Page,
			Limit:     req.Limit,
			Total:     total,
			TotalPage: totalPage,
		},
	}, nil
}

func (u *dataRowUsecase) Create(ctx context.Context, req *domain.CreateDataRowRequest, userID string) (*domain.DataRowInfo, error) {
	now := time.Now()
	row := &domain.DataRow{
		ID:        uuid.New().String(),
		DatasetID: req.DatasetID,
		RowIndex:  req.RowIndex,
		Data:      req.Data,
		CreatedBy: userID,
		CreatedAt: now,
		UpdatedAt: now,
	}

	if err := u.repo.Create(ctx, row); err != nil {
		return nil, fmt.Errorf("failed to create data row: %w", err)
	}

	return u.toInfo(row), nil
}

func (u *dataRowUsecase) BulkCreate(ctx context.Context, req *domain.BulkCreateDataRowsRequest, userID string) error {
	now := time.Now()
	rows := make([]*domain.DataRow, len(req.Rows))

	for i, rowInput := range req.Rows {
		rows[i] = &domain.DataRow{
			ID:        uuid.New().String(),
			DatasetID: req.DatasetID,
			RowIndex:  rowInput.RowIndex,
			Data:      rowInput.Data,
			CreatedBy: userID,
			CreatedAt: now,
			UpdatedAt: now,
		}
	}

	if err := u.repo.BulkCreate(ctx, rows); err != nil {
		return fmt.Errorf("failed to bulk create data rows: %w", err)
	}

	return nil
}

func (u *dataRowUsecase) Update(ctx context.Context, id string, req *domain.UpdateDataRowRequest) (*domain.DataRowInfo, error) {
	existing, err := u.repo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get data row: %w", err)
	}

	// Update fields
	if req.RowIndex != nil {
		existing.RowIndex = *req.RowIndex
	}
	if req.Data != nil {
		existing.Data = *req.Data
	}
	existing.UpdatedAt = time.Now()

	if err := u.repo.Update(ctx, id, existing); err != nil {
		return nil, fmt.Errorf("failed to update data row: %w", err)
	}

	return u.toInfo(existing), nil
}

func (u *dataRowUsecase) Delete(ctx context.Context, id string) error {
	if err := u.repo.Delete(ctx, id); err != nil {
		return fmt.Errorf("failed to delete data row: %w", err)
	}
	return nil
}

func (u *dataRowUsecase) DeleteByDatasetID(ctx context.Context, datasetID string) error {
	if err := u.repo.DeleteByDatasetID(ctx, datasetID); err != nil {
		return fmt.Errorf("failed to delete data rows by dataset: %w", err)
	}
	return nil
}

func (u *dataRowUsecase) GetStats(ctx context.Context, datasetID string) (*domain.DataRowStats, error) {
	stats, err := u.repo.GetStats(ctx, datasetID)
	if err != nil {
		return nil, fmt.Errorf("failed to get data row stats: %w", err)
	}
	return stats, nil
}

func (u *dataRowUsecase) toInfo(row *domain.DataRow) *domain.DataRowInfo {
	return &domain.DataRowInfo{
		ID:        row.ID,
		DatasetID: row.DatasetID,
		RowIndex:  row.RowIndex,
		Data:      row.Data,
		CreatedBy: row.CreatedBy,
		CreatedAt: row.CreatedAt,
		UpdatedAt: row.UpdatedAt,
	}
}
