package repository

import (
	"context"
	"fmt"

	"portal-data-backend/internal/business_field/domain"
	"portal-data-backend/pkg/errors"

	"github.com/jmoiron/sqlx"
)

type businessFieldPostgresRepository struct {
	db *sqlx.DB
}

func NewBusinessFieldPostgresRepository(db *sqlx.DB) domain.Repository {
	return &businessFieldPostgresRepository{db: db}
}

func (r *businessFieldPostgresRepository) GetByID(ctx context.Context, id string) (*domain.BusinessField, error) {
	query := `SELECT id, name, slug, created_at FROM business_fields WHERE id = $1`
	var bf domain.BusinessField
	err := r.db.GetContext(ctx, &bf, query, id)
	if err != nil {
		return nil, r.handleError(err)
	}
	return &bf, nil
}

func (r *businessFieldPostgresRepository) List(ctx context.Context, search string, limit, offset int) ([]*domain.BusinessField, int, error) {
	whereClause := "WHERE 1=1"
	args := []interface{}{}
	argCount := 1

	if search != "" {
		whereClause += fmt.Sprintf(" AND (name ILIKE $%d OR slug ILIKE $%d)", argCount, argCount)
		searchTerm := "%" + search + "%"
		args = append(args, searchTerm, searchTerm)
		argCount += 2
	}

	countQuery := "SELECT COUNT(*) FROM business_fields " + whereClause
	var total int
	err := r.db.GetContext(ctx, &total, countQuery, args...)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count business fields: %w", err)
	}

	query := "SELECT id, name, slug, created_at FROM business_fields " + whereClause + " ORDER BY name ASC LIMIT $" + fmt.Sprintf("%d", argCount) + " OFFSET $" + fmt.Sprintf("%d", argCount+1)
	args = append(args, limit, offset)

	var bfs []*domain.BusinessField
	err = r.db.SelectContext(ctx, &bfs, query, args...)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to list business fields: %w", err)
	}

	return bfs, total, nil
}

func (r *businessFieldPostgresRepository) Create(ctx context.Context, bf *domain.BusinessField) error {
	query := `INSERT INTO business_fields (id, name, slug, created_at) VALUES (:id, :name, :slug, :created_at)`
	_, err := r.db.NamedExecContext(ctx, query, bf)
	if err != nil {
		return fmt.Errorf("failed to create business field: %w", err)
	}
	return nil
}

func (r *businessFieldPostgresRepository) Update(ctx context.Context, bf *domain.BusinessField) error {
	query := `UPDATE business_fields SET name = :name, slug = :slug WHERE id = :id`
	result, err := r.db.NamedExecContext(ctx, query, bf)
	if err != nil {
		return fmt.Errorf("failed to update business field: %w", err)
	}

	rows, _ := result.RowsAffected()
	if rows == 0 {
		return errors.ErrNotFound
	}
	return nil
}

func (r *businessFieldPostgresRepository) Delete(ctx context.Context, id string) error {
	query := `DELETE FROM business_fields WHERE id = $1`
	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete business field: %w", err)
	}

	rows, _ := result.RowsAffected()
	if rows == 0 {
		return errors.ErrNotFound
	}
	return nil
}

func (r *businessFieldPostgresRepository) handleError(err error) error {
	if err == nil {
		return nil
	}
	return errors.Wrap(err, "database error")
}
