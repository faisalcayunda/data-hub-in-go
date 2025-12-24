package repository

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	notifDomain "portal-data-backend/internal/notification/domain"

	"github.com/jmoiron/sqlx"
)

type notificationPostgresRepository struct {
	db *sqlx.DB
}

func NewNotificationPostgresRepository(db *sqlx.DB) notifDomain.Repository {
	return &notificationPostgresRepository{db: db}
}

func (r *notificationPostgresRepository) GetByID(ctx context.Context, id string) (*notifDomain.Notification, error) {
	query := `
		SELECT id, user_id, title, message, type, category, action_url, read, read_at, created_at, deleted_at
		FROM notifications
		WHERE id = $1 AND deleted_at IS NULL
	`

	var notif notifDomain.Notification
	err := r.db.GetContext(ctx, &notif, query, id)
	if err != nil {
		return nil, r.handleError(err)
	}
	return &notif, nil
}

func (r *notificationPostgresRepository) List(ctx context.Context, filter *notifDomain.NotificationFilter, limit, offset int) ([]*notifDomain.Notification, int, error) {
	whereClause := "WHERE deleted_at IS NULL"
	args := []interface{}{}
	argCount := 1

	if filter != nil {
		if filter.UserID != nil {
			whereClause += fmt.Sprintf(" AND user_id = $%d", argCount)
			args = append(args, filter.UserID)
			argCount++
		}
		if filter.Type != nil {
			whereClause += fmt.Sprintf(" AND type = $%d", argCount)
			args = append(args, filter.Type)
			argCount++
		}
		if filter.Category != nil {
			whereClause += fmt.Sprintf(" AND category = $%d", argCount)
			args = append(args, filter.Category)
			argCount++
		}
		if filter.IsRead != nil {
			whereClause += fmt.Sprintf(" AND read = $%d", argCount)
			args = append(args, filter.IsRead)
			argCount++
		}
		if filter.StartDate != nil {
			whereClause += fmt.Sprintf(" AND created_at >= $%d", argCount)
			args = append(args, filter.StartDate)
			argCount++
		}
		if filter.EndDate != nil {
			whereClause += fmt.Sprintf(" AND created_at <= $%d", argCount)
			args = append(args, filter.EndDate)
			argCount++
		}
	}

	countQuery := "SELECT COUNT(*) FROM notifications " + whereClause
	var total int
	err := r.db.GetContext(ctx, &total, countQuery, args...)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count notifications: %w", err)
	}

	query := `
		SELECT id, user_id, title, message, type, category, action_url, read, read_at, created_at, deleted_at
		FROM notifications
	` + whereClause + " ORDER BY created_at DESC LIMIT $" + fmt.Sprintf("%d", argCount) + " OFFSET $" + fmt.Sprintf("%d", argCount+1)

	args = append(args, limit, offset)

	var notifs []*notifDomain.Notification
	err = r.db.SelectContext(ctx, &notifs, query, args...)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to list notifications: %w", err)
	}

	return notifs, total, nil
}

func (r *notificationPostgresRepository) Create(ctx context.Context, notif *notifDomain.Notification) error {
	query := `
		INSERT INTO notifications (id, user_id, title, message, type, category, action_url, read, created_at)
		VALUES (:id, :user_id, :title, :message, :type, :category, :action_url, :read, :created_at)
	`

	_, err := r.db.NamedExecContext(ctx, query, notif)
	if err != nil {
		return fmt.Errorf("failed to create notification: %w", err)
	}
	return nil
}

func (r *notificationPostgresRepository) BulkCreate(ctx context.Context, notifs []*notifDomain.Notification) error {
	query := `
		INSERT INTO notifications (id, user_id, title, message, type, category, action_url, read, created_at)
		VALUES (:id, :user_id, :title, :message, :type, :category, :action_url, :read, :created_at)
	`

	_, err := r.db.NamedExecContext(ctx, query, notifs)
	if err != nil {
		return fmt.Errorf("failed to bulk create notifications: %w", err)
	}
	return nil
}

func (r *notificationPostgresRepository) MarkAsRead(ctx context.Context, ids []string, userID string) error {
	query := `
		UPDATE notifications
		SET read = true, read_at = $1
		WHERE id = ANY($2) AND user_id = $3
	`
	_, err := r.db.ExecContext(ctx, query, time.Now(), ids, userID)
	if err != nil {
		return fmt.Errorf("failed to mark notifications as read: %w", err)
	}
	return nil
}

func (r *notificationPostgresRepository) MarkAllAsRead(ctx context.Context, userID string) error {
	query := `
		UPDATE notifications
		SET read = true, read_at = $1
		WHERE user_id = $2 AND read = false AND deleted_at IS NULL
	`
	_, err := r.db.ExecContext(ctx, query, time.Now(), userID)
	if err != nil {
		return fmt.Errorf("failed to mark all notifications as read: %w", err)
	}
	return nil
}

func (r *notificationPostgresRepository) Delete(ctx context.Context, id string) error {
	query := `UPDATE notifications SET deleted_at = $1 WHERE id = $2`
	_, err := r.db.ExecContext(ctx, query, time.Now(), id)
	if err != nil {
		return fmt.Errorf("failed to delete notification: %w", err)
	}
	return nil
}

func (r *notificationPostgresRepository) GetUnreadCount(ctx context.Context, userID string) (int64, error) {
	query := `
		SELECT COUNT(*)
		FROM notifications
		WHERE user_id = $1 AND read = false AND deleted_at IS NULL
	`
	var count int64
	err := r.db.GetContext(ctx, &count, query, userID)
	if err != nil {
		return 0, fmt.Errorf("failed to get unread count: %w", err)
	}
	return count, nil
}

func (r *notificationPostgresRepository) DeleteOldReadNotifications(ctx context.Context, olderThan time.Duration) error {
	cutoff := time.Now().Add(-olderThan)
	query := `DELETE FROM notifications WHERE read = true AND read_at < $1`
	_, err := r.db.ExecContext(ctx, query, cutoff)
	if err != nil {
		return fmt.Errorf("failed to delete old notifications: %w", err)
	}
	return nil
}

func (r *notificationPostgresRepository) handleError(err error) error {
	if err == nil {
		return nil
	}
	if err == sql.ErrNoRows {
		return fmt.Errorf("notification not found")
	}
	return fmt.Errorf("database error: %w", err)
}
