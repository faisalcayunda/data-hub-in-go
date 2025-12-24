package repository

import (
	"context"
	"fmt"
	"strings"
	"time"

	"portal-data-backend/internal/organization/domain"
	"portal-data-backend/pkg/errors"

	"github.com/jmoiron/sqlx"
)

// orgPostgresRepository implements Repository for PostgreSQL
type orgPostgresRepository struct {
	db *sqlx.DB
}

// NewOrgPostgresRepository creates a new organization repository
func NewOrgPostgresRepository(db *sqlx.DB) domain.Repository {
	return &orgPostgresRepository{db: db}
}

func (r *orgPostgresRepository) GetByID(ctx context.Context, id string) (*domain.Organization, error) {
	query := `
		SELECT id, code, name, slug, description, logo_url, phone_number, address,
		       website_url, email, total_datasets, public_datasets, total_mapsets,
		       public_mapsets, status, created_by, created_at, updated_by, updated_at
		FROM organizations
		WHERE id = $1
	`

	var org domain.Organization
	err := r.db.GetContext(ctx, &org, query, id)
	if err != nil {
		return nil, r.handleError(err)
	}
	return &org, nil
}

func (r *orgPostgresRepository) GetByCode(ctx context.Context, code string) (*domain.Organization, error) {
	query := `
		SELECT id, code, name, slug, description, logo_url, phone_number, address,
		       website_url, email, total_datasets, public_datasets, total_mapsets,
		       public_mapsets, status, created_by, created_at, updated_by, updated_at
		FROM organizations
		WHERE code = $1
	`

	var org domain.Organization
	err := r.db.GetContext(ctx, &org, query, code)
	if err != nil {
		return nil, r.handleError(err)
	}
	return &org, nil
}

func (r *orgPostgresRepository) GetBySlug(ctx context.Context, slug string) (*domain.Organization, error) {
	query := `
		SELECT id, code, name, slug, description, logo_url, phone_number, address,
		       website_url, email, total_datasets, public_datasets, total_mapsets,
		       public_mapsets, status, created_by, created_at, updated_by, updated_at
		FROM organizations
		WHERE slug = $1
	`

	var org domain.Organization
	err := r.db.GetContext(ctx, &org, query, slug)
	if err != nil {
		return nil, r.handleError(err)
	}
	return &org, nil
}

func (r *orgPostgresRepository) List(ctx context.Context, status, search string, limit, offset int, sortBy, sortOrder string) ([]*domain.Organization, int, error) {
	whereClause := "WHERE 1=1"
	args := []interface{}{}
	argCount := 1

	if status != "" {
		whereClause += fmt.Sprintf(" AND status = $%d", argCount)
		args = append(args, status)
		argCount++
	}

	if search != "" {
		whereClause += fmt.Sprintf(" AND (name ILIKE $%d OR code ILIKE $%d)", argCount, argCount)
		searchTerm := "%" + search + "%"
		args = append(args, searchTerm, searchTerm)
		argCount++
	}

	countQuery := "SELECT COUNT(*) FROM organizations " + whereClause
	var total int
	err := r.db.GetContext(ctx, &total, countQuery, args...)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count organizations: %w", err)
	}

	orderClause := r.buildOrderClause(sortBy, sortOrder)
	query := `
		SELECT id, code, name, slug, description, logo_url, phone_number, address,
		       website_url, email, total_datasets, public_datasets, total_mapsets,
		       public_mapsets, status, created_by, created_at, updated_by, updated_at
		FROM organizations
	` + whereClause + " " + orderClause + " LIMIT $" + fmt.Sprintf("%d", argCount) + " OFFSET $" + fmt.Sprintf("%d", argCount+1)

	args = append(args, limit, offset)

	var orgs []*domain.Organization
	err = r.db.SelectContext(ctx, &orgs, query, args...)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to list organizations: %w", err)
	}

	return orgs, total, nil
}

func (r *orgPostgresRepository) Create(ctx context.Context, org *domain.Organization) error {
	query := `
		INSERT INTO organizations (
			id, code, name, slug, description, logo_url, phone_number, address,
			website_url, email, total_datasets, public_datasets, total_mapsets,
			public_mapsets, status, created_by, created_at, updated_by, updated_at
		) VALUES (
			:id, :code, :name, :slug, :description, :logo_url, :phone_number, :address,
			:website_url, :email, :total_datasets, :public_datasets, :total_mapsets,
			:public_mapsets, :status, :created_by, :created_at, :updated_by, :updated_at
		)
	`

	_, err := r.db.NamedExecContext(ctx, query, org)
	if err != nil {
		return fmt.Errorf("failed to create organization: %w", err)
	}
	return nil
}

func (r *orgPostgresRepository) Update(ctx context.Context, org *domain.Organization) error {
	org.UpdatedAt = time.Now()

	query := `
		UPDATE organizations SET
			code = :code, name = :name, slug = :slug, description = :description,
			logo_url = :logo_url, phone_number = :phone_number, address = :address,
			website_url = :website_url, email = :email, status = :status,
			updated_by = :updated_by, updated_at = :updated_at
		WHERE id = :id
	`

	result, err := r.db.NamedExecContext(ctx, query, org)
	if err != nil {
		return fmt.Errorf("failed to update organization: %w", err)
	}

	rows, _ := result.RowsAffected()
	if rows == 0 {
		return errors.ErrNotFound
	}
	return nil
}

func (r *orgPostgresRepository) Delete(ctx context.Context, id string) error {
	query := `DELETE FROM organizations WHERE id = $1`
	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete organization: %w", err)
	}

	rows, _ := result.RowsAffected()
	if rows == 0 {
		return errors.ErrNotFound
	}
	return nil
}

func (r *orgPostgresRepository) UpdateStatus(ctx context.Context, id string, status domain.OrgStatus) error {
	query := `UPDATE organizations SET status = $1, updated_at = NOW() WHERE id = $2`
	result, err := r.db.ExecContext(ctx, query, status, id)
	if err != nil {
		return fmt.Errorf("failed to update organization status: %w", err)
	}

	rows, _ := result.RowsAffected()
	if rows == 0 {
		return errors.ErrNotFound
	}
	return nil
}

func (r *orgPostgresRepository) IncrementDatasetCount(ctx context.Context, id string, isPublic bool) error {
	if isPublic {
		query := `
			UPDATE organizations
			SET total_datasets = total_datasets + 1,
			    public_datasets = public_datasets + 1,
			    updated_at = NOW()
			WHERE id = $1
		`
		_, err := r.db.ExecContext(ctx, query, id)
		return err
	}

	query := `
		UPDATE organizations
		SET total_datasets = total_datasets + 1, updated_at = NOW()
		WHERE id = $1
	`
	_, err := r.db.ExecContext(ctx, query, id)
	return err
}

func (r *orgPostgresRepository) DecrementDatasetCount(ctx context.Context, id string, isPublic bool) error {
	if isPublic {
		query := `
			UPDATE organizations
			SET total_datasets = GREATEST(total_datasets - 1, 0),
			    public_datasets = GREATEST(public_datasets - 1, 0),
			    updated_at = NOW()
			WHERE id = $1
		`
		_, err := r.db.ExecContext(ctx, query, id)
		return err
	}

	query := `
		UPDATE organizations
		SET total_datasets = GREATEST(total_datasets - 1, 0),
		    updated_at = NOW()
		WHERE id = $1
	`
	_, err := r.db.ExecContext(ctx, query, id)
	return err
}

func (r *orgPostgresRepository) buildOrderClause(sortBy, sortOrder string) string {
	allowedColumns := map[string]bool{
		"name":        true,
		"code":        true,
		"status":      true,
		"created_at":  true,
		"updated_at":  true,
	}

	if !allowedColumns[sortBy] {
		sortBy = "created_at"
	}

	sortOrder = strings.ToUpper(sortOrder)
	if sortOrder != "ASC" && sortOrder != "DESC" {
		sortOrder = "DESC"
	}

	return fmt.Sprintf("ORDER BY %s %s", sortBy, sortOrder)
}

func (r *orgPostgresRepository) handleError(err error) error {
	if err == nil {
		return nil
	}
	return errors.Wrap(err, "database error")
}
