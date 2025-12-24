package repository

import (
	"context"
	"fmt"
	"time"

	"portal-data-backend/internal/file/domain"
	"portal-data-backend/pkg/errors"

	"github.com/jmoiron/sqlx"
)

type filePostgresRepository struct {
	db *sqlx.DB
}

func NewFilePostgresRepository(db *sqlx.DB) domain.Repository {
	return &filePostgresRepository{db: db}
}

func (r *filePostgresRepository) GetByID(ctx context.Context, id string) (*domain.File, error) {
	query := `
		SELECT id, name, original_name, extension, size, mime_type, path, storage_path,
		       storage_type, dataset_id, uploaded_by, status, created_at, updated_at
		FROM files
		WHERE id = $1
	`

	var file domain.File
	err := r.db.GetContext(ctx, &file, query, id)
	if err != nil {
		return nil, r.handleError(err)
	}
	return &file, nil
}

func (r *filePostgresRepository) List(ctx context.Context, filter *domain.FileFilter, limit, offset int) ([]*domain.File, int, error) {
	whereClause := "WHERE status != 'deleted'"
	args := []interface{}{}
	argCount := 1

	if filter != nil {
		if filter.DatasetID != nil {
			whereClause += fmt.Sprintf(" AND dataset_id = $%d", argCount)
			args = append(args, filter.DatasetID)
			argCount++
		}
		if filter.Status != nil {
			whereClause += fmt.Sprintf(" AND status = $%d", argCount)
			args = append(args, filter.Status)
			argCount++
		}
		if filter.Search != "" {
			whereClause += fmt.Sprintf(" AND (name ILIKE $%d OR original_name ILIKE $%d)", argCount, argCount)
			searchTerm := "%" + filter.Search + "%"
			args = append(args, searchTerm, searchTerm)
			argCount += 2
		}
	}

	countQuery := "SELECT COUNT(*) FROM files " + whereClause
	var total int
	err := r.db.GetContext(ctx, &total, countQuery, args...)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count files: %w", err)
	}

	query := `
		SELECT id, name, original_name, extension, size, mime_type, path, storage_path,
		       storage_type, dataset_id, uploaded_by, status, created_at, updated_at
		FROM files
	` + whereClause + " ORDER BY created_at DESC LIMIT $" + fmt.Sprintf("%d", argCount) + " OFFSET $" + fmt.Sprintf("%d", argCount+1)

	args = append(args, limit, offset)

	var files []*domain.File
	err = r.db.SelectContext(ctx, &files, query, args...)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to list files: %w", err)
	}

	return files, total, nil
}

func (r *filePostgresRepository) Create(ctx context.Context, file *domain.File) error {
	query := `
		INSERT INTO files (
			id, name, original_name, extension, size, mime_type, path, storage_path,
			storage_type, dataset_id, uploaded_by, status, created_at, updated_at
		) VALUES (
			:id, :name, :original_name, :extension, :size, :mime_type, :path, :storage_path,
			:storage_type, :dataset_id, :uploaded_by, :status, :created_at, :updated_at
		)
	`

	_, err := r.db.NamedExecContext(ctx, query, file)
	if err != nil {
		return fmt.Errorf("failed to create file: %w", err)
	}
	return nil
}

func (r *filePostgresRepository) UpdateStatus(ctx context.Context, id string, status domain.FileStatus) error {
	query := `UPDATE files SET status = $1, updated_at = $2 WHERE id = $3`
	result, err := r.db.ExecContext(ctx, query, status, time.Now(), id)
	if err != nil {
		return fmt.Errorf("failed to update file status: %w", err)
	}

	rows, _ := result.RowsAffected()
	if rows == 0 {
		return errors.ErrNotFound
	}
	return nil
}

func (r *filePostgresRepository) Delete(ctx context.Context, id string) error {
	query := `DELETE FROM files WHERE id = $1`
	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete file: %w", err)
	}

	rows, _ := result.RowsAffected()
	if rows == 0 {
		return errors.ErrNotFound
	}
	return nil
}

func (r *filePostgresRepository) GetByDatasetID(ctx context.Context, datasetID string, limit, offset int) ([]*domain.File, int, error) {
	query := `
		SELECT id, name, original_name, extension, size, mime_type, path, storage_path,
		       storage_type, dataset_id, uploaded_by, status, created_at, updated_at
		FROM files
		WHERE dataset_id = $1 AND status != 'deleted'
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3
	`

	var files []*domain.File
	err := r.db.SelectContext(ctx, &files, query, datasetID, limit, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to get dataset files: %w", err)
	}

	countQuery := `SELECT COUNT(*) FROM files WHERE dataset_id = $1 AND status != 'deleted'`
	var total int
	err = r.db.GetContext(ctx, &total, countQuery, datasetID)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count dataset files: %w", err)
	}

	return files, total, nil
}

func (r *filePostgresRepository) handleError(err error) error {
	if err == nil {
		return nil
	}
	return errors.Wrap(err, "database error")
}
