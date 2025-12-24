package domain

import (
	"context"
	"io"
)

type Repository interface {
	GetByID(ctx context.Context, id string) (*File, error)
	List(ctx context.Context, filter *FileFilter, limit, offset int) ([]*File, int, error)
	Create(ctx context.Context, file *File) error
	UpdateStatus(ctx context.Context, id string, status FileStatus) error
	Delete(ctx context.Context, id string) error
	GetByDatasetID(ctx context.Context, datasetID string, limit, offset int) ([]*File, int, error)
}

type FileFilter struct {
	DatasetID *string
	Status    *string
	Search    string
}

// StorageService defines interface for file storage operations
type StorageService interface {
	Upload(ctx context.Context, fileName string, reader io.Reader, contentType string, path string) (string, error)
	Delete(ctx context.Context, path string) error
	GetURL(ctx context.Context, path string) (string, error)
}
