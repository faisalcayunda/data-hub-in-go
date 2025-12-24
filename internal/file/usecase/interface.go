package usecase

import (
	"context"
	"io"

	"portal-data-backend/internal/file/domain"
)

type Usecase interface {
	GetByID(ctx context.Context, id string) (*domain.FileInfo, error)
	List(ctx context.Context, req *domain.ListFilesRequest) (*domain.FileListResponse, error)
	Upload(ctx context.Context, fileName string, fileSize int64, mimeType string, reader io.Reader, datasetID *string, userID string) (*domain.UploadResponse, error)
	UpdateStatus(ctx context.Context, id string, status domain.FileStatus) error
	Delete(ctx context.Context, id string) error
	GetByDatasetID(ctx context.Context, datasetID string, page, limit int) (*domain.FileListResponse, error)
}
