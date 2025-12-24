package usecase

import (
	"context"
	"fmt"
	"math"
	"time"

	"portal-data-backend/internal/feedback/domain"

	"github.com/google/uuid"
)

type feedbackUsecase struct {
	feedbackRepo domain.Repository
}

func NewFeedbackUsecase(feedbackRepo domain.Repository) Usecase {
	return &feedbackUsecase{feedbackRepo: feedbackRepo}
}

func (u *feedbackUsecase) GetByID(ctx context.Context, id string) (*domain.FeedbackResponse, error) {
	feedback, err := u.feedbackRepo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get feedback: %w", err)
	}
	return u.toResponse(feedback), nil
}

func (u *feedbackUsecase) List(ctx context.Context, req *domain.ListFeedbacksRequest) (*domain.FeedbackListResponse, error) {
	if req.Page < 1 {
		req.Page = 1
	}
	if req.Limit < 1 {
		req.Limit = 20
	}

	offset := (req.Page - 1) * req.Limit

	filter := &domain.FeedbackFilter{
		DatasetID: req.DatasetID,
		Category:  req.Category,
		Status:    req.Status,
		UserID:    req.UserID,
		Search:    req.Search,
	}

	sortBy := req.SortBy
	if sortBy == "" {
		sortBy = "created_at"
	}
	sortOrder := req.SortOrder
	if sortOrder == "" {
		sortOrder = "DESC"
	}

	feedbacks, total, err := u.feedbackRepo.List(ctx, filter, req.Limit, offset, sortBy, sortOrder)
	if err != nil {
		return nil, fmt.Errorf("failed to list feedbacks: %w", err)
	}

	responses := make([]domain.FeedbackResponse, len(feedbacks))
	for i, fb := range feedbacks {
		responses[i] = *u.toResponse(fb)
	}

	totalPage := int(math.Ceil(float64(total) / float64(req.Limit)))

	return &domain.FeedbackListResponse{
		Feedbacks: responses,
		Meta: domain.ListMeta{
			Page:      req.Page,
			Limit:     req.Limit,
			Total:     total,
			TotalPage: totalPage,
		},
	}, nil
}

func (u *feedbackUsecase) Create(ctx context.Context, req *domain.CreateFeedbackRequest, userID string) (*domain.FeedbackResponse, error) {
	feedback := &domain.Feedback{
		ID:        uuid.New().String(),
		UserID:    userID,
		DatasetID: req.DatasetID,
		Rating:    req.Rating,
		Comment:   req.Comment,
		Category:  req.Category,
		Status:    domain.FeedbackStatusPending,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if err := u.feedbackRepo.Create(ctx, feedback); err != nil {
		return nil, fmt.Errorf("failed to create feedback: %w", err)
	}

	return u.toResponse(feedback), nil
}

func (u *feedbackUsecase) UpdateStatus(ctx context.Context, id string, status domain.FeedbackStatus) error {
	if err := u.feedbackRepo.UpdateStatus(ctx, id, status); err != nil {
		return fmt.Errorf("failed to update feedback status: %w", err)
	}
	return nil
}

func (u *feedbackUsecase) Delete(ctx context.Context, id string) error {
	if err := u.feedbackRepo.Delete(ctx, id); err != nil {
		return fmt.Errorf("failed to delete feedback: %w", err)
	}
	return nil
}

func (u *feedbackUsecase) toResponse(feedback *domain.Feedback) *domain.FeedbackResponse {
	return &domain.FeedbackResponse{
		ID:        feedback.ID,
		UserID:    feedback.UserID,
		DatasetID: feedback.DatasetID,
		Rating:    feedback.Rating,
		Comment:   feedback.Comment,
		Category:  string(feedback.Category),
		Status:    string(feedback.Status),
		CreatedAt: feedback.CreatedAt,
		UpdatedAt: feedback.UpdatedAt,
	}
}
