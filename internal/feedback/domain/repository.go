package domain

import (
	"context"
)

type Repository interface {
	GetByID(ctx context.Context, id string) (*Feedback, error)
	List(ctx context.Context, filter *FeedbackFilter, limit, offset int, sortBy, sortOrder string) ([]*Feedback, int, error)
	Create(ctx context.Context, feedback *Feedback) error
	UpdateStatus(ctx context.Context, id string, status FeedbackStatus) error
	Delete(ctx context.Context, id string) error
}

type FeedbackFilter struct {
	DatasetID *string
	Category  *string
	Status    *string
	UserID    *string
	Search    string
}
