package usecase

import (
	"context"
	"fmt"
	"math"
	"strings"
	"time"

	"portal-data-backend/internal/topic/domain"

	"github.com/google/uuid"
)

type topicUsecase struct {
	topicRepo domain.Repository
}

func NewTopicUsecase(topicRepo domain.Repository) Usecase {
	return &topicUsecase{topicRepo: topicRepo}
}

func (u *topicUsecase) GetByID(ctx context.Context, id string) (*domain.TopicResponse, error) {
	topic, err := u.topicRepo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get topic: %w", err)
	}
	return u.toResponse(topic), nil
}

func (u *topicUsecase) List(ctx context.Context, req *domain.ListTopicsRequest) (*domain.TopicListResponse, error) {
	if req.Page < 1 {
		req.Page = 1
	}
	if req.Limit < 1 {
		req.Limit = 20
	}

	offset := (req.Page - 1) * req.Limit

	topics, total, err := u.topicRepo.List(ctx, req.Search, req.Limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to list topics: %w", err)
	}

	responses := make([]domain.TopicResponse, len(topics))
	for i, topic := range topics {
		responses[i] = *u.toResponse(topic)
	}

	totalPage := int(math.Ceil(float64(total) / float64(req.Limit)))

	return &domain.TopicListResponse{
		Topics: responses,
		Meta: domain.ListMeta{
			Page:      req.Page,
			Limit:     req.Limit,
			Total:     total,
			TotalPage: totalPage,
		},
	}, nil
}

func (u *topicUsecase) Create(ctx context.Context, req *domain.CreateTopicRequest) (*domain.TopicResponse, error) {
	topic := &domain.Topic{
		ID:        uuid.New().String(),
		Name:      req.Name,
		Slug:      u.generateSlug(req.Name),
		CreatedAt: time.Now(),
	}

	if err := u.topicRepo.Create(ctx, topic); err != nil {
		return nil, fmt.Errorf("failed to create topic: %w", err)
	}

	return u.toResponse(topic), nil
}

func (u *topicUsecase) Update(ctx context.Context, id string, req *domain.UpdateTopicRequest) (*domain.TopicResponse, error) {
	topic, err := u.topicRepo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get topic: %w", err)
	}

	topic.Name = req.Name
	topic.Slug = u.generateSlug(req.Name)

	if err := u.topicRepo.Update(ctx, topic); err != nil {
		return nil, fmt.Errorf("failed to update topic: %w", err)
	}

	return u.toResponse(topic), nil
}

func (u *topicUsecase) Delete(ctx context.Context, id string) error {
	if err := u.topicRepo.Delete(ctx, id); err != nil {
		return fmt.Errorf("failed to delete topic: %w", err)
	}
	return nil
}

func (u *topicUsecase) toResponse(topic *domain.Topic) *domain.TopicResponse {
	return &domain.TopicResponse{
		ID:        topic.ID,
		Name:      topic.Name,
		Slug:      topic.Slug,
		CreatedAt: topic.CreatedAt,
	}
}

func (u *topicUsecase) generateSlug(name string) string {
	slug := strings.ToLower(name)
	slug = strings.ReplaceAll(slug, " ", "-")
	return slug
}
