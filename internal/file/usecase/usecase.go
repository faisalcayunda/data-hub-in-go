package usecase

import (
	"context"
	"fmt"
	"io"
	"math"
	"path/filepath"
	"strings"
	"time"

	"portal-data-backend/internal/file/domain"

	"github.com/google/uuid"
)

type fileUsecase struct {
	fileRepo    domain.Repository
	storage     domain.StorageService
	baseStoragePath string
}

func NewFileUsecase(fileRepo domain.Repository, storage domain.StorageService, basePath string) Usecase {
	return &fileUsecase{
		fileRepo:        fileRepo,
		storage:         storage,
		baseStoragePath: basePath,
	}
}

func (u *fileUsecase) GetByID(ctx context.Context, id string) (*domain.FileInfo, error) {
	file, err := u.fileRepo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get file: %w", err)
	}
	return u.toInfo(file), nil
}

func (u *fileUsecase) List(ctx context.Context, req *domain.ListFilesRequest) (*domain.FileListResponse, error) {
	if req.Page < 1 {
		req.Page = 1
	}
	if req.Limit < 1 {
		req.Limit = 20
	}

	offset := (req.Page - 1) * req.Limit

	filter := &domain.FileFilter{
		DatasetID: req.DatasetID,
		Status:    req.Status,
		Search:    req.Search,
	}

	files, total, err := u.fileRepo.List(ctx, filter, req.Limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to list files: %w", err)
	}

	infos := make([]domain.FileInfo, len(files))
	for i, file := range files {
		infos[i] = *u.toInfo(file)
	}

	totalPage := int(math.Ceil(float64(total) / float64(req.Limit)))

	return &domain.FileListResponse{
		Files: infos,
		Meta: domain.ListMeta{
			Page:      req.Page,
			Limit:     req.Limit,
			Total:     total,
			TotalPage: totalPage,
		},
	}, nil
}

func (u *fileUsecase) Upload(ctx context.Context, fileName string, fileSize int64, mimeType string, reader io.Reader, datasetID *string, userID string) (*domain.UploadResponse, error) {
	// Generate file ID and path
	ext := filepath.Ext(fileName)
	fileID := uuid.New().String()
	storagePath := fmt.Sprintf("%s/%s%s", u.baseStoragePath, fileID, ext)

	// Upload to storage
	uploadedPath, err := u.storage.Upload(ctx, fileName, reader, mimeType, storagePath)
	if err != nil {
		return nil, fmt.Errorf("failed to upload file: %w", err)
	}

	// Create file record
	file := &domain.File{
		ID:           fileID,
		Name:         strings.TrimSuffix(fileName, ext),
		OriginalName: fileName,
		Extension:    ext,
		Size:         fileSize,
		MimeType:     mimeType,
		Path:         uploadedPath,
		StoragePath:  storagePath,
		StorageType:   domain.StorageTypeMinIO,
		DatasetID:    datasetID,
		UploadedBy:   userID,
		Status:       domain.FileStatusReady,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	if err := u.fileRepo.Create(ctx, file); err != nil {
		// Rollback storage upload
		_ = u.storage.Delete(ctx, uploadedPath)
		return nil, fmt.Errorf("failed to create file record: %w", err)
	}

	return &domain.UploadResponse{
		ID:        file.ID,
		Name:      file.Name,
		Size:      file.Size,
		MimeType:  file.MimeType,
		Path:      file.Path,
		Status:    string(file.Status),
		CreatedAt: file.CreatedAt.Format(time.RFC3339),
	}, nil
}

func (u *fileUsecase) UpdateStatus(ctx context.Context, id string, status domain.FileStatus) error {
	if err := u.fileRepo.UpdateStatus(ctx, id, status); err != nil {
		return fmt.Errorf("failed to update file status: %w", err)
	}
	return nil
}

func (u *fileUsecase) Delete(ctx context.Context, id string) error {
	file, err := u.fileRepo.GetByID(ctx, id)
	if err != nil {
		return fmt.Errorf("failed to get file: %w", err)
	}

	// Delete from storage
	if err := u.storage.Delete(ctx, file.Path); err != nil {
		return fmt.Errorf("failed to delete from storage: %w", err)
	}

	// Delete record
	if err := u.fileRepo.Delete(ctx, id); err != nil {
		return fmt.Errorf("failed to delete file record: %w", err)
	}

	return nil
}

func (u *fileUsecase) GetByDatasetID(ctx context.Context, datasetID string, page, limit int) (*domain.FileListResponse, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 20
	}

	offset := (page - 1) * limit

	files, total, err := u.fileRepo.GetByDatasetID(ctx, datasetID, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to get dataset files: %w", err)
	}

	infos := make([]domain.FileInfo, len(files))
	for i, file := range files {
		infos[i] = *u.toInfo(file)
	}

	totalPage := int(math.Ceil(float64(total) / float64(limit)))

	return &domain.FileListResponse{
		Files: infos,
		Meta: domain.ListMeta{
			Page:      page,
			Limit:     limit,
			Total:     total,
			TotalPage: totalPage,
		},
	}, nil
}

func (u *fileUsecase) toInfo(file *domain.File) *domain.FileInfo {
	return &domain.FileInfo{
		ID:           file.ID,
		Name:         file.Name,
		OriginalName: file.OriginalName,
		Extension:    file.Extension,
		Size:         file.Size,
		MimeType:     file.MimeType,
		Path:         file.Path,
		DatasetID:    file.DatasetID,
		UploadedBy:   file.UploadedBy,
		Status:       string(file.Status),
		CreatedAt:    file.CreatedAt,
	}
}
