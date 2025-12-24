package repository

import (
	"context"
	"fmt"

	"portal-data-backend/internal/unit/domain"
	"portal-data-backend/pkg/errors"

	"github.com/jmoiron/sqlx"
)

type unitPostgresRepository struct {
	db *sqlx.DB
}

func NewUnitPostgresRepository(db *sqlx.DB) domain.Repository {
	return &unitPostgresRepository{db: db}
}

func (r *unitPostgresRepository) GetByID(ctx context.Context, id string) (*domain.Unit, error) {
	query := `SELECT id, name, symbol, created_at FROM units WHERE id = $1`
	var unit domain.Unit
	err := r.db.GetContext(ctx, &unit, query, id)
	if err != nil {
		return nil, r.handleError(err)
	}
	return &unit, nil
}

func (r *unitPostgresRepository) List(ctx context.Context, search string, limit, offset int) ([]*domain.Unit, int, error) {
	whereClause := "WHERE 1=1"
	args := []interface{}{}
	argCount := 1

	if search != "" {
		whereClause += fmt.Sprintf(" AND (name ILIKE $%d OR symbol ILIKE $%d)", argCount, argCount)
		searchTerm := "%" + search + "%"
		args = append(args, searchTerm, searchTerm)
		argCount += 2
	}

	countQuery := "SELECT COUNT(*) FROM units " + whereClause
	var total int
	err := r.db.GetContext(ctx, &total, countQuery, args...)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count units: %w", err)
	}

	query := "SELECT id, name, symbol, created_at FROM units " + whereClause + " ORDER BY name ASC LIMIT $" + fmt.Sprintf("%d", argCount) + " OFFSET $" + fmt.Sprintf("%d", argCount+1)
	args = append(args, limit, offset)

	var units []*domain.Unit
	err = r.db.SelectContext(ctx, &units, query, args...)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to list units: %w", err)
	}

	return units, total, nil
}

func (r *unitPostgresRepository) Create(ctx context.Context, unit *domain.Unit) error {
	query := `INSERT INTO units (id, name, symbol, created_at) VALUES (:id, :name, :symbol, :created_at)`
	_, err := r.db.NamedExecContext(ctx, query, unit)
	if err != nil {
		return fmt.Errorf("failed to create unit: %w", err)
	}
	return nil
}

func (r *unitPostgresRepository) Update(ctx context.Context, unit *domain.Unit) error {
	query := `UPDATE units SET name = :name, symbol = :symbol WHERE id = :id`
	result, err := r.db.NamedExecContext(ctx, query, unit)
	if err != nil {
		return fmt.Errorf("failed to update unit: %w", err)
	}

	rows, _ := result.RowsAffected()
	if rows == 0 {
		return errors.ErrNotFound
	}
	return nil
}

func (r *unitPostgresRepository) Delete(ctx context.Context, id string) error {
	query := `DELETE FROM units WHERE id = $1`
	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete unit: %w", err)
	}

	rows, _ := result.RowsAffected()
	if rows == 0 {
		return errors.ErrNotFound
	}
	return nil
}

func (r *unitPostgresRepository) handleError(err error) error {
	if err == nil {
		return nil
	}
	return errors.Wrap(err, "database error")
}
