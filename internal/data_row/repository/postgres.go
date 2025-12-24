package repository

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	dataRowDomain "portal-data-backend/internal/data_row/domain"

	"github.com/jmoiron/sqlx"
)

type dataRowPostgresRepository struct {
	db *sqlx.DB
}

func NewDataRowPostgresRepository(db *sqlx.DB) dataRowDomain.Repository {
	return &dataRowPostgresRepository{db: db}
}

func (r *dataRowPostgresRepository) GetByID(ctx context.Context, id string) (*dataRowDomain.DataRow, error) {
	query := `
		SELECT id, dataset_id, row_index, data, created_by, created_at, updated_at, deleted_at
		FROM data_rows
		WHERE id = $1 AND deleted_at IS NULL
	`

	var row dataRowDomain.DataRow
	err := r.db.GetContext(ctx, &row, query, id)
	if err != nil {
		return nil, r.handleError(err)
	}
	return &row, nil
}

func (r *dataRowPostgresRepository) List(ctx context.Context, filter *dataRowDomain.DataRowFilter, limit, offset int) ([]*dataRowDomain.DataRow, int, error) {
	whereClause := "WHERE deleted_at IS NULL AND dataset_id = $1"
	args := []interface{}{filter.DatasetID}
	argCount := 2

	if filter.Search != "" {
		whereClause += fmt.Sprintf(" AND data::text ILIKE $%d", argCount)
		args = append(args, "%"+filter.Search+"%")
		argCount++
	}

	countQuery := "SELECT COUNT(*) FROM data_rows " + whereClause
	var total int
	err := r.db.GetContext(ctx, &total, countQuery, args...)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count data rows: %w", err)
	}

	query := `
		SELECT id, dataset_id, row_index, data, created_by, created_at, updated_at, deleted_at
		FROM data_rows
	` + whereClause + " ORDER BY row_index ASC LIMIT $" + fmt.Sprintf("%d", argCount) + " OFFSET $" + fmt.Sprintf("%d", argCount+1)

	args = append(args, limit, offset)

	var rows []*dataRowDomain.DataRow
	err = r.db.SelectContext(ctx, &rows, query, args...)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to list data rows: %w", err)
	}

	return rows, total, nil
}

func (r *dataRowPostgresRepository) Create(ctx context.Context, row *dataRowDomain.DataRow) error {
	query := `
		INSERT INTO data_rows (id, dataset_id, row_index, data, created_by, created_at, updated_at)
		VALUES (:id, :dataset_id, :row_index, :data, :created_by, :created_at, :updated_at)
	`

	_, err := r.db.NamedExecContext(ctx, query, row)
	if err != nil {
		return fmt.Errorf("failed to create data row: %w", err)
	}
	return nil
}

func (r *dataRowPostgresRepository) BulkCreate(ctx context.Context, rows []*dataRowDomain.DataRow) error {
	query := `
		INSERT INTO data_rows (id, dataset_id, row_index, data, created_by, created_at, updated_at)
		VALUES (:id, :dataset_id, :row_index, :data, :created_by, :created_at, :updated_at)
	`

	_, err := r.db.NamedExecContext(ctx, query, rows)
	if err != nil {
		return fmt.Errorf("failed to bulk create data rows: %w", err)
	}
	return nil
}

func (r *dataRowPostgresRepository) Update(ctx context.Context, id string, row *dataRowDomain.DataRow) error {
	query := `
		UPDATE data_rows
		SET row_index = :row_index, data = :data, updated_at = :updated_at
		WHERE id = :id
	`

	row.ID = id
	_, err := r.db.NamedExecContext(ctx, query, row)
	if err != nil {
		return fmt.Errorf("failed to update data row: %w", err)
	}
	return nil
}

func (r *dataRowPostgresRepository) Delete(ctx context.Context, id string) error {
	query := `UPDATE data_rows SET deleted_at = $1 WHERE id = $2`
	_, err := r.db.ExecContext(ctx, query, time.Now(), id)
	if err != nil {
		return fmt.Errorf("failed to delete data row: %w", err)
	}
	return nil
}

func (r *dataRowPostgresRepository) DeleteByDatasetID(ctx context.Context, datasetID string) error {
	query := `UPDATE data_rows SET deleted_at = $1 WHERE dataset_id = $2`
	_, err := r.db.ExecContext(ctx, query, time.Now(), datasetID)
	if err != nil {
		return fmt.Errorf("failed to delete data rows by dataset: %w", err)
	}
	return nil
}

func (r *dataRowPostgresRepository) GetByRowIndex(ctx context.Context, datasetID string, rowIndex int) (*dataRowDomain.DataRow, error) {
	query := `
		SELECT id, dataset_id, row_index, data, created_by, created_at, updated_at, deleted_at
		FROM data_rows
		WHERE dataset_id = $1 AND row_index = $2 AND deleted_at IS NULL
	`

	var row dataRowDomain.DataRow
	err := r.db.GetContext(ctx, &row, query, datasetID, rowIndex)
	if err != nil {
		return nil, r.handleError(err)
	}
	return &row, nil
}

func (r *dataRowPostgresRepository) GetStats(ctx context.Context, datasetID string) (*dataRowDomain.DataRowStats, error) {
	query := `
		SELECT
			COUNT(*) as total_rows,
			COALESCE(MAX(updated_at), NOW()) as last_updated
		FROM data_rows
		WHERE dataset_id = $1 AND deleted_at IS NULL
	`

	var stats dataRowDomain.DataRowStats
	err := r.db.GetContext(ctx, &stats, query, datasetID)
	if err != nil {
		return nil, fmt.Errorf("failed to get data row stats: %w", err)
	}

	return &stats, nil
}

func (r *dataRowPostgresRepository) handleError(err error) error {
	if err == nil {
		return nil
	}
	if err == sql.ErrNoRows {
		return fmt.Errorf("data row not found")
	}
	return fmt.Errorf("database error: %w", err)
}
