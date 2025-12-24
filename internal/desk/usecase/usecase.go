package usecase

import (
	"context"
	"fmt"
	"math"
	"time"

	"portal-data-backend/internal/desk/domain"

	"github.com/google/uuid"
)

type Usecase interface {
	GetByID(ctx context.Context, id string) (*domain.TicketInfo, error)
	List(ctx context.Context, req *domain.ListTicketsRequest) (*domain.TicketListResponse, error)
	Create(ctx context.Context, req *domain.CreateTicketRequest, userID string) (*domain.TicketInfo, error)
	Update(ctx context.Context, id string, req *domain.UpdateTicketRequest) (*domain.TicketInfo, error)
	Delete(ctx context.Context, id string) error
	UpdateStatus(ctx context.Context, id string, status string) error
	AssignTicket(ctx context.Context, id string, assignedTo string) error
}

type deskUsecase struct {
	repo domain.Repository
}

func NewDeskUsecase(repo domain.Repository) Usecase {
	return &deskUsecase{
		repo: repo,
	}
}

func (u *deskUsecase) GetByID(ctx context.Context, id string) (*domain.TicketInfo, error) {
	ticket, err := u.repo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get ticket: %w", err)
	}
	return u.toInfo(ticket), nil
}

func (u *deskUsecase) List(ctx context.Context, req *domain.ListTicketsRequest) (*domain.TicketListResponse, error) {
	if req.Page < 1 {
		req.Page = 1
	}
	if req.Limit < 1 {
		req.Limit = 20
	}

	offset := (req.Page - 1) * req.Limit

	filter := &domain.TicketFilter{
		UserID:     req.UserID,
		AssignedTo: req.AssignedTo,
		Status:     req.Status,
		Priority:   req.Priority,
		Category:   req.Category,
		Search:     req.Search,
	}

	tickets, total, err := u.repo.List(ctx, filter, req.Limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to list tickets: %w", err)
	}

	infos := make([]domain.TicketInfo, len(tickets))
	for i, ticket := range tickets {
		infos[i] = *u.toInfo(ticket)
	}

	totalPage := int(math.Ceil(float64(total) / float64(req.Limit)))

	return &domain.TicketListResponse{
		Tickets: infos,
		Meta: domain.ListMeta{
			Page:      req.Page,
			Limit:     req.Limit,
			Total:     total,
			TotalPage: totalPage,
		},
	}, nil
}

func (u *deskUsecase) Create(ctx context.Context, req *domain.CreateTicketRequest, userID string) (*domain.TicketInfo, error) {
	now := time.Now()
	ticket := &domain.Ticket{
		ID:          uuid.New().String(),
		Title:       req.Title,
		Description: req.Description,
		Status:      string(domain.TicketStatusOpen),
		Priority:    req.Priority,
		Category:    req.Category,
		UserID:      userID,
		AssignedTo:  req.AssignedTo,
		CreatedBy:   userID,
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	if err := u.repo.Create(ctx, ticket); err != nil {
		return nil, fmt.Errorf("failed to create ticket: %w", err)
	}

	return u.toInfo(ticket), nil
}

func (u *deskUsecase) Update(ctx context.Context, id string, req *domain.UpdateTicketRequest) (*domain.TicketInfo, error) {
	existing, err := u.repo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get ticket: %w", err)
	}

	// Update fields
	if req.Title != nil {
		existing.Title = *req.Title
	}
	if req.Description != nil {
		existing.Description = *req.Description
	}
	if req.Status != nil {
		existing.Status = *req.Status
		// Set resolved_at when status is resolved
		if *req.Status == string(domain.TicketStatusResolved) && existing.ResolvedAt == nil {
			now := time.Now()
			existing.ResolvedAt = &now
		}
	}
	if req.Priority != nil {
		existing.Priority = *req.Priority
	}
	if req.Category != nil {
		existing.Category = *req.Category
	}
	if req.AssignedTo != nil {
		existing.AssignedTo = req.AssignedTo
	}
	existing.UpdatedAt = time.Now()

	if err := u.repo.Update(ctx, id, existing); err != nil {
		return nil, fmt.Errorf("failed to update ticket: %w", err)
	}

	return u.toInfo(existing), nil
}

func (u *deskUsecase) Delete(ctx context.Context, id string) error {
	if err := u.repo.Delete(ctx, id); err != nil {
		return fmt.Errorf("failed to delete ticket: %w", err)
	}
	return nil
}

func (u *deskUsecase) UpdateStatus(ctx context.Context, id string, status string) error {
	if err := u.repo.UpdateStatus(ctx, id, status); err != nil {
		return fmt.Errorf("failed to update ticket status: %w", err)
	}
	return nil
}

func (u *deskUsecase) AssignTicket(ctx context.Context, id string, assignedTo string) error {
	if err := u.repo.AssignTicket(ctx, id, assignedTo); err != nil {
		return fmt.Errorf("failed to assign ticket: %w", err)
	}
	return nil
}

func (u *deskUsecase) toInfo(ticket *domain.Ticket) *domain.TicketInfo {
	return &domain.TicketInfo{
		ID:          ticket.ID,
		Title:       ticket.Title,
		Description: ticket.Description,
		Status:      ticket.Status,
		Priority:    ticket.Priority,
		Category:    ticket.Category,
		UserID:      ticket.UserID,
		AssignedTo:  ticket.AssignedTo,
		ResolvedAt:  ticket.ResolvedAt,
		CreatedBy:   ticket.CreatedBy,
		CreatedAt:   ticket.CreatedAt,
		UpdatedAt:   ticket.UpdatedAt,
	}
}
