package usecase

import (
	"context"
	"fmt"
	"math"
	"time"

	"portal-data-backend/internal/notification/domain"

	"github.com/google/uuid"
)

type Usecase interface {
	GetByID(ctx context.Context, id string) (*domain.NotificationInfo, error)
	List(ctx context.Context, req *domain.ListNotificationsRequest) (*domain.NotificationListResponse, error)
	Create(ctx context.Context, req *domain.CreateNotificationRequest) (*domain.NotificationInfo, error)
	BulkCreate(ctx context.Context, req *domain.BulkCreateNotificationRequest) error
	MarkAsRead(ctx context.Context, ids []string, userID string) error
	MarkAllAsRead(ctx context.Context, userID string) error
	Delete(ctx context.Context, id string) error
	GetUnreadCount(ctx context.Context, userID string) (int64, error)
}

type notificationUsecase struct {
	repo domain.Repository
}

func NewNotificationUsecase(repo domain.Repository) Usecase {
	return &notificationUsecase{
		repo: repo,
	}
}

func (u *notificationUsecase) GetByID(ctx context.Context, id string) (*domain.NotificationInfo, error) {
	notif, err := u.repo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get notification: %w", err)
	}
	return u.toInfo(notif), nil
}

func (u *notificationUsecase) List(ctx context.Context, req *domain.ListNotificationsRequest) (*domain.NotificationListResponse, error) {
	if req.Page < 1 {
		req.Page = 1
	}
	if req.Limit < 1 {
		req.Limit = 20
	}

	offset := (req.Page - 1) * req.Limit

	filter := &domain.NotificationFilter{
		UserID:    req.UserID,
		Type:      req.Type,
		Category:  req.Category,
		IsRead:    req.IsRead,
		StartDate: req.StartDate,
		EndDate:   req.EndDate,
	}

	notifs, total, err := u.repo.List(ctx, filter, req.Limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to list notifications: %w", err)
	}

	infos := make([]domain.NotificationInfo, len(notifs))
	for i, notif := range notifs {
		infos[i] = *u.toInfo(notif)
	}

	totalPage := int(math.Ceil(float64(total) / float64(req.Limit)))

	return &domain.NotificationListResponse{
		Notifications: infos,
		Meta: domain.ListMeta{
			Page:      req.Page,
			Limit:     req.Limit,
			Total:     total,
			TotalPage: totalPage,
		},
	}, nil
}

func (u *notificationUsecase) Create(ctx context.Context, req *domain.CreateNotificationRequest) (*domain.NotificationInfo, error) {
	now := time.Now()
	notif := &domain.Notification{
		ID:        uuid.New().String(),
		UserID:    req.UserID,
		Title:     req.Title,
		Message:   req.Message,
		Type:      req.Type,
		Category:  req.Category,
		ActionURL: req.ActionURL,
		Read:      false,
		CreatedAt: now,
	}

	if err := u.repo.Create(ctx, notif); err != nil {
		return nil, fmt.Errorf("failed to create notification: %w", err)
	}

	return u.toInfo(notif), nil
}

func (u *notificationUsecase) BulkCreate(ctx context.Context, req *domain.BulkCreateNotificationRequest) error {
	now := time.Now()
	notifs := make([]*domain.Notification, len(req.UserIDs))

	for i, userID := range req.UserIDs {
		notifs[i] = &domain.Notification{
			ID:        uuid.New().String(),
			UserID:    userID,
			Title:     req.Title,
			Message:   req.Message,
			Type:      req.Type,
			Category:  req.Category,
			ActionURL: req.ActionURL,
			Read:      false,
			CreatedAt: now,
		}
	}

	if err := u.repo.BulkCreate(ctx, notifs); err != nil {
		return fmt.Errorf("failed to bulk create notifications: %w", err)
	}

	return nil
}

func (u *notificationUsecase) MarkAsRead(ctx context.Context, ids []string, userID string) error {
	if err := u.repo.MarkAsRead(ctx, ids, userID); err != nil {
		return fmt.Errorf("failed to mark notifications as read: %w", err)
	}
	return nil
}

func (u *notificationUsecase) MarkAllAsRead(ctx context.Context, userID string) error {
	if err := u.repo.MarkAllAsRead(ctx, userID); err != nil {
		return fmt.Errorf("failed to mark all notifications as read: %w", err)
	}
	return nil
}

func (u *notificationUsecase) Delete(ctx context.Context, id string) error {
	if err := u.repo.Delete(ctx, id); err != nil {
		return fmt.Errorf("failed to delete notification: %w", err)
	}
	return nil
}

func (u *notificationUsecase) GetUnreadCount(ctx context.Context, userID string) (int64, error) {
	count, err := u.repo.GetUnreadCount(ctx, userID)
	if err != nil {
		return 0, fmt.Errorf("failed to get unread count: %w", err)
	}
	return count, nil
}

func (u *notificationUsecase) toInfo(notif *domain.Notification) *domain.NotificationInfo {
	return &domain.NotificationInfo{
		ID:        notif.ID,
		UserID:    notif.UserID,
		Title:     notif.Title,
		Message:   notif.Message,
		Type:      notif.Type,
		Category:  notif.Category,
		ActionURL: notif.ActionURL,
		Read:      notif.Read,
		ReadAt:    notif.ReadAt,
		CreatedAt: notif.CreatedAt,
	}
}
