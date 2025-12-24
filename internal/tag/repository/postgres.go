package repository

import (
	"context"
	"fmt"

	"portal-data-backend/internal/tag/domain"
	"portal-data-backend/pkg/errors"

	"github.com/jmoiron/sqlx"
)

type tagPostgresRepository struct {
	db *sqlx.DB
}

func NewTagPostgresRepository(db *sqlx.DB) domain.Repository {
	return &tagPostgresRepository{db: db}
}

func (r *tagPostgresRepository) GetByID(ctx context.Context, id string) (*domain.Tag, error) {
	query := `SELECT id, name, slug, created_at FROM tags WHERE id = $1`
	var tag domain.Tag
	err := r.db.GetContext(ctx, &tag, query, id)
	if err != nil {
		return nil, r.handleError(err)
	}
	return &tag, nil
}

func (r *tagPostgresRepository) List(ctx context.Context, search string, limit, offset int) ([]*domain.Tag, int, error) {
	whereClause := "WHERE 1=1"
	args := []interface{}{}
	argCount := 1

	if search != "" {
		whereClause += fmt.Sprintf(" AND (name ILIKE $%d OR slug ILIKE $%d)", argCount, argCount)
		searchTerm := "%" + search + "%"
		args = append(args, searchTerm, searchTerm)
		argCount += 2
	}

	countQuery := "SELECT COUNT(*) FROM tags " + whereClause
	var total int
	err := r.db.GetContext(ctx, &total, countQuery, args...)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count tags: %w", err)
	}

	query := "SELECT id, name, slug, created_at FROM tags " + whereClause + " ORDER BY name ASC LIMIT $" + fmt.Sprintf("%d", argCount) + " OFFSET $" + fmt.Sprintf("%d", argCount+1)
	args = append(args, limit, offset)

	var tags []*domain.Tag
	err = r.db.SelectContext(ctx, &tags, query, args...)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to list tags: %w", err)
	}

	return tags, total, nil
}

func (r *tagPostgresRepository) Create(ctx context.Context, tag *domain.Tag) error {
	query := `INSERT INTO tags (id, name, slug, created_at) VALUES (:id, :name, :slug, :created_at)`
	_, err := r.db.NamedExecContext(ctx, query, tag)
	if err != nil {
		return fmt.Errorf("failed to create tag: %w", err)
	}
	return nil
}

func (r *tagPostgresRepository) Update(ctx context.Context, tag *domain.Tag) error {
	query := `UPDATE tags SET name = :name, slug = :slug WHERE id = :id`
	result, err := r.db.NamedExecContext(ctx, query, tag)
	if err != nil {
		return fmt.Errorf("failed to update tag: %w", err)
	}

	rows, _ := result.RowsAffected()
	if rows == 0 {
		return errors.ErrNotFound
	}
	return nil
}

func (r *tagPostgresRepository) Delete(ctx context.Context, id string) error {
	query := `DELETE FROM tags WHERE id = $1`
	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete tag: %w", err)
	}

	rows, _ := result.RowsAffected()
	if rows == 0 {
		return errors.ErrNotFound
	}
	return nil
}

func (r *tagPostgresRepository) handleError(err error) error {
	if err == nil {
		return nil
	}
	return errors.Wrap(err, "database error")
}
