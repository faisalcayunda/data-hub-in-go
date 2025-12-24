package domain

import (
	"context"
)

type Repository interface {
	GetByID(ctx context.Context, id string) (*DataRow, error)
	List(ctx context.Context, filter *DataRowFilter, limit, offset int) ([]*DataRow, int, error)
	Create(ctx context.Context, row *DataRow) error
	BulkCreate(ctx context.Context, rows []*DataRow) error
	Update(ctx context.Context, id string, row *DataRow) error
	Delete(ctx context.Context, id string) error
	DeleteByDatasetID(ctx context.Context, datasetID string) error
	GetByRowIndex(ctx context.Context, datasetID string, rowIndex int) (*DataRow, error)
	GetStats(ctx context.Context, datasetID string) (*DataRowStats, error)
}

type DataRowFilter struct {
	DatasetID string
	Search    string
}
