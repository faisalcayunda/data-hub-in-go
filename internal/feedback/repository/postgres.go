package repository

import (
	"context"
	"fmt"
	"strings"
	"time"

	"portal-data-backend/internal/feedback/domain"
	"portal-data-backend/pkg/errors"

	"github.com/jmoiron/sqlx"
)

type feedbackPostgresRepository struct {
	db *sqlx.DB
}

func NewFeedbackPostgresRepository(db *sqlx.DB) domain.Repository {
	return &feedbackPostgresRepository{db: db}
}

func (r *feedbackPostgresRepository) GetByID(ctx context.Context, id string) (*domain.Feedback, error) {
	query := `
		SELECT id, user_id, dataset_id, rating, comment, category, status, created_at, updated_at
		FROM feedbacks
		WHERE id = $1
	`

	var feedback domain.Feedback
	err := r.db.GetContext(ctx, &feedback, query, id)
	if err != nil {
		return nil, r.handleError(err)
	}
	return &feedback, nil
}

func (r *feedbackPostgresRepository) List(ctx context.Context, filter *domain.FeedbackFilter, limit, offset int, sortBy, sortOrder string) ([]*domain.Feedback, int, error) {
	whereClause := "WHERE 1=1"
	args := []interface{}{}
	argCount := 1

	if filter != nil {
		if filter.DatasetID != nil {
			whereClause += fmt.Sprintf(" AND dataset_id = $%d", argCount)
			args = append(args, filter.DatasetID)
			argCount++
		}
		if filter.Category != nil {
			whereClause += fmt.Sprintf(" AND category = $%d", argCount)
			args = append(args, filter.Category)
			argCount++
		}
		if filter.Status != nil {
			whereClause += fmt.Sprintf(" AND status = $%d", argCount)
			args = append(args, filter.Status)
			argCount++
		}
		if filter.UserID != nil {
			whereClause += fmt.Sprintf(" AND user_id = $%d", argCount)
			args = append(args, filter.UserID)
			argCount++
		}
		if filter.Search != "" {
			whereClause += fmt.Sprintf(" AND (comment ILIKE $%d)", argCount)
			searchTerm := "%" + filter.Search + "%"
			args = append(args, searchTerm)
			argCount++
		}
	}

	countQuery := "SELECT COUNT(*) FROM feedbacks " + whereClause
	var total int
	err := r.db.GetContext(ctx, &total, countQuery, args...)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count feedbacks: %w", err)
	}

	orderClause := r.buildOrderClause(sortBy, sortOrder)
	query := `
		SELECT id, user_id, dataset_id, rating, comment, category, status, created_at, updated_at
		FROM feedbacks
	` + whereClause + " " + orderClause + " LIMIT $" + fmt.Sprintf("%d", argCount) + " OFFSET $" + fmt.Sprintf("%d", argCount+1)

	args = append(args, limit, offset)

	var feedbacks []*domain.Feedback
	err = r.db.SelectContext(ctx, &feedbacks, query, args...)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to list feedbacks: %w", err)
	}

	return feedbacks, total, nil
}

func (r *feedbackPostgresRepository) Create(ctx context.Context, feedback *domain.Feedback) error {
	query := `
		INSERT INTO feedbacks (
			id, user_id, dataset_id, rating, comment, category, status, created_at, updated_at
		) VALUES (
			:id, :user_id, :dataset_id, :rating, :comment, :category, :status, :created_at, :updated_at
		)
	`

	_, err := r.db.NamedExecContext(ctx, query, feedback)
	if err != nil {
		return fmt.Errorf("failed to create feedback: %w", err)
	}
	return nil
}

func (r *feedbackPostgresRepository) UpdateStatus(ctx context.Context, id string, status domain.FeedbackStatus) error {
	query := `UPDATE feedbacks SET status = $1, updated_at = $2 WHERE id = $3`
	result, err := r.db.ExecContext(ctx, query, status, time.Now(), id)
	if err != nil {
		return fmt.Errorf("failed to update feedback status: %w", err)
	}

	rows, _ := result.RowsAffected()
	if rows == 0 {
		return errors.ErrNotFound
	}
	return nil
}

func (r *feedbackPostgresRepository) Delete(ctx context.Context, id string) error {
	query := `DELETE FROM feedbacks WHERE id = $1`
	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete feedback: %w", err)
	}

	rows, _ := result.RowsAffected()
	if rows == 0 {
		return errors.ErrNotFound
	}
	return nil
}

func (r *feedbackPostgresRepository) buildOrderClause(sortBy, sortOrder string) string {
	allowedColumns := map[string]bool{
		"rating":     true,
		"created_at": true,
		"updated_at": true,
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

func (r *feedbackPostgresRepository) handleError(err error) error {
	if err == nil {
		return nil
	}
	return errors.Wrap(err, "database error")
}
