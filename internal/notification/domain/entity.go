package domain

import "time"

// Notification represents a user notification
type Notification struct {
	ID        string     `db:"id" json:"id"`
	UserID    string     `db:"user_id" json:"user_id"`
	Title     string     `db:"title" json:"title"`
	Message   string     `db:"message" json:"message"`
	Type      string     `db:"type" json:"type"` // info, warning, error, success
	Category  string     `db:"category" json:"category"` // system, dataset, publication, etc.
	ActionURL *string    `db:"action_url" json:"action_url,omitempty"`
	Read      bool       `db:"read" json:"read"`
	ReadAt    *time.Time `db:"read_at" json:"read_at,omitempty"`
	CreatedAt time.Time  `db:"created_at" json:"created_at"`
	DeletedAt *time.Time `db:"deleted_at" json:"deleted_at,omitempty"`
}

// NotificationType represents notification type
type NotificationType string

const (
	NotificationTypeInfo    NotificationType = "info"
	NotificationTypeWarning NotificationType = "warning"
	NotificationTypeError   NotificationType = "error"
	NotificationTypeSuccess NotificationType = "success"
)

// NotificationCategory represents notification category
type NotificationCategory string

const (
	NotificationCategorySystem     NotificationCategory = "system"
	NotificationCategoryDataset    NotificationCategory = "dataset"
	NotificationCategoryPublication NotificationCategory = "publication"
	NotificationCategoryUser       NotificationCategory = "user"
	NotificationCategoryFeedback   NotificationCategory = "feedback"
)

// ListNotificationsRequest represents list notifications input
type ListNotificationsRequest struct {
	Page      int     `json:"page" validate:"min=1"`
	Limit     int     `json:"limit" validate:"min=1,max=100"`
	UserID    *string `json:"user_id,omitempty"`
	Type      *string `json:"type,omitempty"`
	Category  *string `json:"category,omitempty"`
	IsRead    *bool   `json:"is_read,omitempty"`
	StartDate *string `json:"start_date,omitempty"`
	EndDate   *string `json:"end_date,omitempty"`
}

// CreateNotificationRequest represents create notification input
type CreateNotificationRequest struct {
	UserID    string  `json:"user_id" validate:"required"`
	Title     string  `json:"title" validate:"required,min=2,max=200"`
	Message   string  `json:"message" validate:"required"`
	Type      string  `json:"type" validate:"required"`
	Category  string  `json:"category" validate:"required"`
	ActionURL *string `json:"action_url,omitempty"`
}

// BulkCreateNotificationRequest represents bulk create notification input
type BulkCreateNotificationRequest struct {
	UserIDs   []string `json:"user_ids" validate:"required,min=1"`
	Title     string   `json:"title" validate:"required,min=2,max=200"`
	Message   string   `json:"message" validate:"required"`
	Type      string   `json:"type" validate:"required"`
	Category  string   `json:"category" validate:"required"`
	ActionURL *string  `json:"action_url,omitempty"`
}

// MarkAsReadRequest represents mark as read input
type MarkAsReadRequest struct {
	NotificationIDs []string `json:"notification_ids" validate:"required,min=1"`
}

// NotificationInfo represents notification information for API responses
type NotificationInfo struct {
	ID        string     `json:"id"`
	UserID    string     `json:"user_id"`
	Title     string     `json:"title"`
	Message   string     `json:"message"`
	Type      string     `json:"type"`
	Category  string     `json:"category"`
	ActionURL *string    `json:"action_url,omitempty"`
	Read      bool       `json:"read"`
	ReadAt    *time.Time `json:"read_at,omitempty"`
	CreatedAt time.Time  `json:"created_at"`
}

// NotificationListResponse represents paginated notification list
type NotificationListResponse struct {
	Notifications []NotificationInfo `json:"notifications"`
	Meta          ListMeta           `json:"meta"`
}

// ListMeta represents pagination metadata
type ListMeta struct {
	Page      int `json:"page"`
	Limit     int `json:"limit"`
	Total     int `json:"total"`
	TotalPage int `json:"total_page"`
}

// UnreadCountResponse represents unread count response
type UnreadCountResponse struct {
	Count int64 `json:"count"`
}
