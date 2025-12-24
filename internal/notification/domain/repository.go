package domain

import (
	"context"
	"time"
)

type Repository interface {
	GetByID(ctx context.Context, id string) (*Notification, error)
	List(ctx context.Context, filter *NotificationFilter, limit, offset int) ([]*Notification, int, error)
	Create(ctx context.Context, notif *Notification) error
	BulkCreate(ctx context.Context, notifs []*Notification) error
	MarkAsRead(ctx context.Context, ids []string, userID string) error
	MarkAllAsRead(ctx context.Context, userID string) error
	Delete(ctx context.Context, id string) error
	GetUnreadCount(ctx context.Context, userID string) (int64, error)
	DeleteOldReadNotifications(ctx context.Context, olderThan time.Duration) error
}

type NotificationFilter struct {
	UserID    *string
	Type      *string
	Category  *string
	IsRead    *bool
	StartDate *string
	EndDate   *string
}
