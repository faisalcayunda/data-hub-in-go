package domain

import (
	"context"
)

type Repository interface {
	GetByID(ctx context.Context, id string) (*Ticket, error)
	List(ctx context.Context, filter *TicketFilter, limit, offset int) ([]*Ticket, int, error)
	Create(ctx context.Context, ticket *Ticket) error
	Update(ctx context.Context, id string, ticket *Ticket) error
	Delete(ctx context.Context, id string) error
	UpdateStatus(ctx context.Context, id string, status string) error
	AssignTicket(ctx context.Context, id string, assignedTo string) error
}

type TicketFilter struct {
	UserID     *string
	AssignedTo *string
	Status     *string
	Priority   *string
	Category   *string
	Search     string
}
