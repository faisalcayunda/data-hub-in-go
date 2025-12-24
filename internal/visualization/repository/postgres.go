package repository

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	visualizationDomain "portal-data-backend/internal/visualization/domain"

	"github.com/jmoiron/sqlx"
)

type visualizationPostgresRepository struct {
	db *sqlx.DB
}

func NewVisualizationPostgresRepository(db *sqlx.DB) visualizationDomain.Repository {
	return &visualizationPostgresRepository{db: db}
}

func (r *visualizationPostgresRepository) GetByID(ctx context.Context, id string) (*visualizationDomain.Visualization, error) {
	query := `
		SELECT id, title, description, type, config, dataset_id, organization_id, topic_id,
		       is_highlight, status, created_by, updated_by, created_at, updated_at, deleted_at
		FROM visualizations
		WHERE id = $1 AND deleted_at IS NULL
	`

	var viz visualizationDomain.Visualization
	err := r.db.GetContext(ctx, &viz, query, id)
	if err != nil {
		return nil, r.handleError(err)
	}
	return &viz, nil
}

func (r *visualizationPostgresRepository) List(ctx context.Context, filter *visualizationDomain.VisualizationFilter, limit, offset int) ([]*visualizationDomain.Visualization, int, error) {
	whereClause := "WHERE deleted_at IS NULL"
	args := []interface{}{}
	argCount := 1

	if filter != nil {
		if filter.DatasetID != nil {
			whereClause += fmt.Sprintf(" AND dataset_id = $%d", argCount)
			args = append(args, filter.DatasetID)
			argCount++
		}
		if filter.OrganizationID != nil {
			whereClause += fmt.Sprintf(" AND organization_id = $%d", argCount)
			args = append(args, filter.OrganizationID)
			argCount++
		}
		if filter.TopicID != nil {
			whereClause += fmt.Sprintf(" AND topic_id = $%d", argCount)
			args = append(args, filter.TopicID)
			argCount++
		}
		if filter.Type != nil {
			whereClause += fmt.Sprintf(" AND type = $%d", argCount)
			args = append(args, filter.Type)
			argCount++
		}
		if filter.Status != nil {
			whereClause += fmt.Sprintf(" AND status = $%d", argCount)
			args = append(args, filter.Status)
			argCount++
		}
		if filter.IsHighlight != nil {
			whereClause += fmt.Sprintf(" AND is_highlight = $%d", argCount)
			args = append(args, filter.IsHighlight)
			argCount++
		}
		if filter.Search != "" {
			whereClause += fmt.Sprintf(" AND (title ILIKE $%d OR description ILIKE $%d)", argCount, argCount)
			searchTerm := "%" + filter.Search + "%"
			args = append(args, searchTerm, searchTerm)
			argCount += 2
		}
	}

	countQuery := "SELECT COUNT(*) FROM visualizations " + whereClause
	var total int
	err := r.db.GetContext(ctx, &total, countQuery, args...)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count visualizations: %w", err)
	}

	query := `
		SELECT id, title, description, type, config, dataset_id, organization_id, topic_id,
		       is_highlight, status, created_by, updated_by, created_at, updated_at, deleted_at
		FROM visualizations
	` + whereClause + " ORDER BY created_at DESC LIMIT $" + fmt.Sprintf("%d", argCount) + " OFFSET $" + fmt.Sprintf("%d", argCount+1)

	args = append(args, limit, offset)

	var vizs []*visualizationDomain.Visualization
	err = r.db.SelectContext(ctx, &vizs, query, args...)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to list visualizations: %w", err)
	}

	return vizs, total, nil
}

func (r *visualizationPostgresRepository) Create(ctx context.Context, viz *visualizationDomain.Visualization) error {
	query := `
		INSERT INTO visualizations (
			id, title, description, type, config, dataset_id, organization_id, topic_id,
			is_highlight, status, created_by, updated_by, created_at, updated_at
		) VALUES (
			:id, :title, :description, :type, :config, :dataset_id, :organization_id, :topic_id,
			:is_highlight, :status, :created_by, :updated_by, :created_at, :updated_at
		)
	`

	_, err := r.db.NamedExecContext(ctx, query, viz)
	if err != nil {
		return fmt.Errorf("failed to create visualization: %w", err)
	}
	return nil
}

func (r *visualizationPostgresRepository) Update(ctx context.Context, id string, viz *visualizationDomain.Visualization) error {
	query := `
		UPDATE visualizations
		SET title = :title, description = :description, type = :type, config = :config,
		    dataset_id = :dataset_id, organization_id = :organization_id, topic_id = :topic_id,
		    is_highlight = :is_highlight, status = :status, updated_by = :updated_by,
		    updated_at = :updated_at
		WHERE id = :id
	`

	viz.ID = id
	_, err := r.db.NamedExecContext(ctx, query, viz)
	if err != nil {
		return fmt.Errorf("failed to update visualization: %w", err)
	}
	return nil
}

func (r *visualizationPostgresRepository) Delete(ctx context.Context, id string) error {
	query := `UPDATE visualizations SET deleted_at = $1 WHERE id = $2`
	_, err := r.db.ExecContext(ctx, query, time.Now(), id)
	if err != nil {
		return fmt.Errorf("failed to delete visualization: %w", err)
	}
	return nil
}

func (r *visualizationPostgresRepository) UpdateStatus(ctx context.Context, id string, status string) error {
	query := `UPDATE visualizations SET status = $1, updated_at = $2 WHERE id = $3`
	_, err := r.db.ExecContext(ctx, query, status, time.Now(), id)
	if err != nil {
		return fmt.Errorf("failed to update visualization status: %w", err)
	}
	return nil
}

func (r *visualizationPostgresRepository) GetStats(ctx context.Context) (*visualizationDomain.VisualizationStats, error) {
	query := `
		SELECT
			COUNT(*) as total_count,
			COUNT(*) FILTER (WHERE status = 'published') as published_count,
			COUNT(*) FILTER (WHERE status = 'draft') as draft_count,
			COUNT(*) FILTER (WHERE is_highlight = true) as highlight_count,
			COALESCE(MAX(updated_at), NOW()) as last_updated
		FROM visualizations
		WHERE deleted_at IS NULL
	`

	var stats visualizationDomain.VisualizationStats
	err := r.db.GetContext(ctx, &stats, query)
	if err != nil {
		return nil, fmt.Errorf("failed to get visualization stats: %w", err)
	}

	return &stats, nil
}

func (r *visualizationPostgresRepository) GetByDatasetID(ctx context.Context, datasetID string, limit, offset int) ([]*visualizationDomain.Visualization, int, error) {
	query := `
		SELECT id, title, description, type, config, dataset_id, organization_id, topic_id,
		       is_highlight, status, created_by, updated_by, created_at, updated_at, deleted_at
		FROM visualizations
		WHERE dataset_id = $1 AND deleted_at IS NULL
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3
	`

	var vizs []*visualizationDomain.Visualization
	err := r.db.SelectContext(ctx, &vizs, query, datasetID, limit, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to get dataset visualizations: %w", err)
	}

	countQuery := `SELECT COUNT(*) FROM visualizations WHERE dataset_id = $1 AND deleted_at IS NULL`
	var total int
	err = r.db.GetContext(ctx, &total, countQuery, datasetID)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count dataset visualizations: %w", err)
	}

	return vizs, total, nil
}

func (r *visualizationPostgresRepository) GetByOrganizationID(ctx context.Context, orgID string, limit, offset int) ([]*visualizationDomain.Visualization, int, error) {
	query := `
		SELECT id, title, description, type, config, dataset_id, organization_id, topic_id,
		       is_highlight, status, created_by, updated_by, created_at, updated_at, deleted_at
		FROM visualizations
		WHERE organization_id = $1 AND deleted_at IS NULL
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3
	`

	var vizs []*visualizationDomain.Visualization
	err := r.db.SelectContext(ctx, &vizs, query, orgID, limit, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to get organization visualizations: %w", err)
	}

	countQuery := `SELECT COUNT(*) FROM visualizations WHERE organization_id = $1 AND deleted_at IS NULL`
	var total int
	err = r.db.GetContext(ctx, &total, countQuery, orgID)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count organization visualizations: %w", err)
	}

	return vizs, total, nil
}

func (r *visualizationPostgresRepository) handleError(err error) error {
	if err == nil {
		return nil
	}
	if err == sql.ErrNoRows {
		return fmt.Errorf("visualization not found")
	}
	return fmt.Errorf("database error: %w", err)
}
