package repository

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	integrationDomain "portal-data-backend/internal/integration/domain"

	"github.com/jmoiron/sqlx"
)

type integrationPostgresRepository struct {
	db *sqlx.DB
}

func NewIntegrationPostgresRepository(db *sqlx.DB) integrationDomain.Repository {
	return &integrationPostgresRepository{db: db}
}

func (r *integrationPostgresRepository) GetByID(ctx context.Context, id string) (*integrationDomain.Integration, error) {
	query := `
		SELECT id, name, type, description, config, endpoint, api_key, status, last_sync_at,
		       organization_id, created_by, created_at, updated_at, deleted_at
		FROM integrations
		WHERE id = $1 AND deleted_at IS NULL
	`

	var integration integrationDomain.Integration
	err := r.db.GetContext(ctx, &integration, query, id)
	if err != nil {
		return nil, r.handleError(err)
	}
	return &integration, nil
}

func (r *integrationPostgresRepository) List(ctx context.Context, filter *integrationDomain.IntegrationFilter, limit, offset int) ([]*integrationDomain.Integration, int, error) {
	whereClause := "WHERE deleted_at IS NULL"
	args := []interface{}{}
	argCount := 1

	if filter != nil {
		if filter.OrganizationID != nil {
			whereClause += fmt.Sprintf(" AND organization_id = $%d", argCount)
			args = append(args, filter.OrganizationID)
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
		if filter.Search != "" {
			whereClause += fmt.Sprintf(" AND (name ILIKE $%d OR description ILIKE $%d)", argCount, argCount)
			searchTerm := "%" + filter.Search + "%"
			args = append(args, searchTerm, searchTerm)
			argCount += 2
		}
	}

	countQuery := "SELECT COUNT(*) FROM integrations " + whereClause
	var total int
	err := r.db.GetContext(ctx, &total, countQuery, args...)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count integrations: %w", err)
	}

	query := `
		SELECT id, name, type, description, config, endpoint, api_key, status, last_sync_at,
		       organization_id, created_by, created_at, updated_at, deleted_at
		FROM integrations
	` + whereClause + " ORDER BY created_at DESC LIMIT $" + fmt.Sprintf("%d", argCount) + " OFFSET $" + fmt.Sprintf("%d", argCount+1)

	args = append(args, limit, offset)

	var integrations []*integrationDomain.Integration
	err = r.db.SelectContext(ctx, &integrations, query, args...)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to list integrations: %w", err)
	}

	return integrations, total, nil
}

func (r *integrationPostgresRepository) Create(ctx context.Context, integration *integrationDomain.Integration) error {
	query := `
		INSERT INTO integrations (id, name, type, description, config, endpoint, api_key, status,
		                        organization_id, created_by, created_at, updated_at)
		VALUES (:id, :name, :type, :description, :config, :endpoint, :api_key, :status,
		        :organization_id, :created_by, :created_at, :updated_at)
	`

	_, err := r.db.NamedExecContext(ctx, query, integration)
	if err != nil {
		return fmt.Errorf("failed to create integration: %w", err)
	}
	return nil
}

func (r *integrationPostgresRepository) Update(ctx context.Context, id string, integration *integrationDomain.Integration) error {
	query := `
		UPDATE integrations
		SET name = :name, description = :description, config = :config, endpoint = :endpoint,
		    api_key = :api_key, status = :status, updated_at = :updated_at
		WHERE id = :id
	`

	integration.ID = id
	_, err := r.db.NamedExecContext(ctx, query, integration)
	if err != nil {
		return fmt.Errorf("failed to update integration: %w", err)
	}
	return nil
}

func (r *integrationPostgresRepository) Delete(ctx context.Context, id string) error {
	query := `UPDATE integrations SET deleted_at = $1 WHERE id = $2`
	_, err := r.db.ExecContext(ctx, query, time.Now(), id)
	if err != nil {
		return fmt.Errorf("failed to delete integration: %w", err)
	}
	return nil
}

func (r *integrationPostgresRepository) UpdateStatus(ctx context.Context, id string, status string) error {
	query := `UPDATE integrations SET status = $1, updated_at = $2 WHERE id = $3`
	_, err := r.db.ExecContext(ctx, query, status, time.Now(), id)
	if err != nil {
		return fmt.Errorf("failed to update integration status: %w", err)
	}
	return nil
}

func (r *integrationPostgresRepository) Sync(ctx context.Context, id string) error {
	// Update last_sync_at
	now := time.Now()
	query := `UPDATE integrations SET last_sync_at = $1, updated_at = $2 WHERE id = $3`
	_, err := r.db.ExecContext(ctx, query, now, now, id)
	if err != nil {
		return fmt.Errorf("failed to sync integration: %w", err)
	}
	return nil
}

func (r *integrationPostgresRepository) handleError(err error) error {
	if err == nil {
		return nil
	}
	if err == sql.ErrNoRows {
		return fmt.Errorf("integration not found")
	}
	return fmt.Errorf("database error: %w", err)
}
