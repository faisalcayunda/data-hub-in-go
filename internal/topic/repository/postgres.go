package repository

import (
	"context"
	"fmt"

	"portal-data-backend/internal/topic/domain"
	"portal-data-backend/pkg/errors"

	"github.com/jmoiron/sqlx"
)

type topicPostgresRepository struct {
	db *sqlx.DB
}

func NewTopicPostgresRepository(db *sqlx.DB) domain.Repository {
	return &topicPostgresRepository{db: db}
}

func (r *topicPostgresRepository) GetByID(ctx context.Context, id string) (*domain.Topic, error) {
	query := `SELECT id, name, slug, created_at FROM topics WHERE id = $1`
	var topic domain.Topic
	err := r.db.GetContext(ctx, &topic, query, id)
	if err != nil {
		return nil, r.handleError(err)
	}
	return &topic, nil
}

func (r *topicPostgresRepository) List(ctx context.Context, search string, limit, offset int) ([]*domain.Topic, int, error) {
	whereClause := "WHERE 1=1"
	args := []interface{}{}
	argCount := 1

	if search != "" {
		whereClause += fmt.Sprintf(" AND (name ILIKE $%d OR slug ILIKE $%d)", argCount, argCount)
		searchTerm := "%" + search + "%"
		args = append(args, searchTerm, searchTerm)
		argCount += 2
	}

	countQuery := "SELECT COUNT(*) FROM topics " + whereClause
	var total int
	err := r.db.GetContext(ctx, &total, countQuery, args...)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count topics: %w", err)
	}

	query := "SELECT id, name, slug, created_at FROM topics " + whereClause + " ORDER BY name ASC LIMIT $" + fmt.Sprintf("%d", argCount) + " OFFSET $" + fmt.Sprintf("%d", argCount+1)
	args = append(args, limit, offset)

	var topics []*domain.Topic
	err = r.db.SelectContext(ctx, &topics, query, args...)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to list topics: %w", err)
	}

	return topics, total, nil
}

func (r *topicPostgresRepository) Create(ctx context.Context, topic *domain.Topic) error {
	query := `INSERT INTO topics (id, name, slug, created_at) VALUES (:id, :name, :slug, :created_at)`
	_, err := r.db.NamedExecContext(ctx, query, topic)
	if err != nil {
		return fmt.Errorf("failed to create topic: %w", err)
	}
	return nil
}

func (r *topicPostgresRepository) Update(ctx context.Context, topic *domain.Topic) error {
	query := `UPDATE topics SET name = :name, slug = :slug WHERE id = :id`
	result, err := r.db.NamedExecContext(ctx, query, topic)
	if err != nil {
		return fmt.Errorf("failed to update topic: %w", err)
	}

	rows, _ := result.RowsAffected()
	if rows == 0 {
		return errors.ErrNotFound
	}
	return nil
}

func (r *topicPostgresRepository) Delete(ctx context.Context, id string) error {
	query := `DELETE FROM topics WHERE id = $1`
	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete topic: %w", err)
	}

	rows, _ := result.RowsAffected()
	if rows == 0 {
		return errors.ErrNotFound
	}
	return nil
}

func (r *topicPostgresRepository) handleError(err error) error {
	if err == nil {
		return nil
	}
	return errors.Wrap(err, "database error")
}
