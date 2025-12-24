package domain

import "time"

// Ticket represents a helpdesk ticket
type Ticket struct {
	ID          string     `db:"id" json:"id"`
	Title       string     `db:"title" json:"title"`
	Description string     `db:"description" json:"description"`
	Status      string     `db:"status" json:"status"`
	Priority    string     `db:"priority" json:"priority"`
	Category    string     `db:"category" json:"category"`
	UserID      string     `db:"user_id" json:"user_id"`
	AssignedTo  *string    `db:"assigned_to" json:"assigned_to,omitempty"`
	ResolvedAt  *time.Time `db:"resolved_at" json:"resolved_at,omitempty"`
	CreatedBy   string     `db:"created_by" json:"created_by"`
	CreatedAt   time.Time  `db:"created_at" json:"created_at"`
	UpdatedAt   time.Time  `db:"updated_at" json:"updated_at"`
	DeletedAt   *time.Time `db:"deleted_at" json:"deleted_at,omitempty"`
}

// TicketStatus represents ticket status
type TicketStatus string

const (
	TicketStatusOpen     TicketStatus = "open"
	TicketStatusInProgress TicketStatus = "in_progress"
	TicketStatusResolved TicketStatus = "resolved"
	TicketStatusClosed   TicketStatus = "closed"
)

// TicketPriority represents ticket priority
type TicketPriority string

const (
	TicketPriorityLow    TicketPriority = "low"
	TicketPriorityMedium TicketPriority = "medium"
	TicketPriorityHigh   TicketPriority = "high"
	TicketPriorityUrgent TicketPriority = "urgent"
)

// TicketCategory represents ticket category
type TicketCategory string

const (
	TicketCategoryTechnical   TicketCategory = "technical"
	TicketCategoryDataRequest TicketCategory = "data_request"
	TicketCategoryReport      TicketCategory = "report"
	TicketCategoryOther       TicketCategory = "other"
)

// ListTicketsRequest represents list tickets input
type ListTicketsRequest struct {
	Page      int     `json:"page" validate:"min=1"`
	Limit     int     `json:"limit" validate:"min=1,max=100"`
	UserID    *string `json:"user_id,omitempty"`
	AssignedTo *string `json:"assigned_to,omitempty"`
	Status    *string `json:"status,omitempty"`
	Priority  *string `json:"priority,omitempty"`
	Category  *string `json:"category,omitempty"`
	Search    string  `json:"search,omitempty"`
}

// CreateTicketRequest represents create ticket input
type CreateTicketRequest struct {
	Title       string  `json:"title" validate:"required,min=2,max=200"`
	Description string  `json:"description" validate:"required"`
	Priority    string  `json:"priority" validate:"required"`
	Category    string  `json:"category" validate:"required"`
	AssignedTo  *string `json:"assigned_to,omitempty"`
}

// UpdateTicketRequest represents update ticket input
type UpdateTicketRequest struct {
	Title       *string `json:"title" validate:"omitempty,min=2,max=200"`
	Description *string `json:"description,omitempty"`
	Status      *string `json:"status,omitempty"`
	Priority    *string `json:"priority,omitempty"`
	Category    *string `json:"category,omitempty"`
	AssignedTo  *string `json:"assigned_to,omitempty"`
}

// TicketInfo represents ticket information for API responses
type TicketInfo struct {
	ID          string     `json:"id"`
	Title       string     `json:"title"`
	Description string     `json:"description"`
	Status      string     `json:"status"`
	Priority    string     `json:"priority"`
	Category    string     `json:"category"`
	UserID      string     `json:"user_id"`
	AssignedTo  *string    `json:"assigned_to,omitempty"`
	ResolvedAt  *time.Time `json:"resolved_at,omitempty"`
	CreatedBy   string     `json:"created_by"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
}

// TicketListResponse represents paginated ticket list
type TicketListResponse struct {
	Tickets []TicketInfo `json:"tickets"`
	Meta    ListMeta     `json:"meta"`
}

// ListMeta represents pagination metadata
type ListMeta struct {
	Page      int `json:"page"`
	Limit     int `json:"limit"`
	Total     int `json:"total"`
	TotalPage int `json:"total_page"`
}
