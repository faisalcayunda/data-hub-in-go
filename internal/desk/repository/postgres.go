package repository

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	deskDomain "portal-data-backend/internal/desk/domain"

	"github.com/jmoiron/sqlx"
)

type deskPostgresRepository struct {
	db *sqlx.DB
}

func NewDeskPostgresRepository(db *sqlx.DB) deskDomain.Repository {
	return &deskPostgresRepository{db: db}
}

func (r *deskPostgresRepository) GetByID(ctx context.Context, id string) (*deskDomain.Ticket, error) {
	query := `
		SELECT id, title, description, status, priority, category, user_id, assigned_to,
		       resolved_at, created_by, created_at, updated_at, deleted_at
		FROM tickets
		WHERE id = $1 AND deleted_at IS NULL
	`

	var ticket deskDomain.Ticket
	err := r.db.GetContext(ctx, &ticket, query, id)
	if err != nil {
		return nil, r.handleError(err)
	}
	return &ticket, nil
}

func (r *deskPostgresRepository) List(ctx context.Context, filter *deskDomain.TicketFilter, limit, offset int) ([]*deskDomain.Ticket, int, error) {
	whereClause := "WHERE deleted_at IS NULL"
	args := []interface{}{}
	argCount := 1

	if filter != nil {
		if filter.UserID != nil {
			whereClause += fmt.Sprintf(" AND user_id = $%d", argCount)
			args = append(args, filter.UserID)
			argCount++
		}
		if filter.AssignedTo != nil {
			whereClause += fmt.Sprintf(" AND assigned_to = $%d", argCount)
			args = append(args, filter.AssignedTo)
			argCount++
		}
		if filter.Status != nil {
			whereClause += fmt.Sprintf(" AND status = $%d", argCount)
			args = append(args, filter.Status)
			argCount++
		}
		if filter.Priority != nil {
			whereClause += fmt.Sprintf(" AND priority = $%d", argCount)
			args = append(args, filter.Priority)
			argCount++
		}
		if filter.Category != nil {
			whereClause += fmt.Sprintf(" AND category = $%d", argCount)
			args = append(args, filter.Category)
			argCount++
		}
		if filter.Search != "" {
			whereClause += fmt.Sprintf(" AND (title ILIKE $%d OR description ILIKE $%d)", argCount, argCount)
			searchTerm := "%" + filter.Search + "%"
			args = append(args, searchTerm, searchTerm)
			argCount += 2
		}
	}

	countQuery := "SELECT COUNT(*) FROM tickets " + whereClause
	var total int
	err := r.db.GetContext(ctx, &total, countQuery, args...)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count tickets: %w", err)
	}

	query := `
		SELECT id, title, description, status, priority, category, user_id, assigned_to,
		       resolved_at, created_by, created_at, updated_at, deleted_at
		FROM tickets
	` + whereClause + " ORDER BY created_at DESC LIMIT $" + fmt.Sprintf("%d", argCount) + " OFFSET $" + fmt.Sprintf("%d", argCount+1)

	args = append(args, limit, offset)

	var tickets []*deskDomain.Ticket
	err = r.db.SelectContext(ctx, &tickets, query, args...)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to list tickets: %w", err)
	}

	return tickets, total, nil
}

func (r *deskPostgresRepository) Create(ctx context.Context, ticket *deskDomain.Ticket) error {
	query := `
		INSERT INTO tickets (id, title, description, status, priority, category, user_id, assigned_to,
		                    created_by, created_at, updated_at)
		VALUES (:id, :title, :description, :status, :priority, :category, :user_id, :assigned_to,
		        :created_by, :created_at, :updated_at)
	`

	_, err := r.db.NamedExecContext(ctx, query, ticket)
	if err != nil {
		return fmt.Errorf("failed to create ticket: %w", err)
	}
	return nil
}

func (r *deskPostgresRepository) Update(ctx context.Context, id string, ticket *deskDomain.Ticket) error {
	query := `
		UPDATE tickets
		SET title = :title, description = :description, status = :status, priority = :priority,
		    category = :category, assigned_to = :assigned_to, resolved_at = :resolved_at, updated_at = :updated_at
		WHERE id = :id
	`

	ticket.ID = id
	_, err := r.db.NamedExecContext(ctx, query, ticket)
	if err != nil {
		return fmt.Errorf("failed to update ticket: %w", err)
	}
	return nil
}

func (r *deskPostgresRepository) Delete(ctx context.Context, id string) error {
	query := `UPDATE tickets SET deleted_at = $1 WHERE id = $2`
	_, err := r.db.ExecContext(ctx, query, time.Now(), id)
	if err != nil {
		return fmt.Errorf("failed to delete ticket: %w", err)
	}
	return nil
}

func (r *deskPostgresRepository) UpdateStatus(ctx context.Context, id string, status string) error {
	query := `UPDATE tickets SET status = $1, updated_at = $2 WHERE id = $3`
	_, err := r.db.ExecContext(ctx, query, status, time.Now(), id)
	if err != nil {
		return fmt.Errorf("failed to update ticket status: %w", err)
	}
	return nil
}

func (r *deskPostgresRepository) AssignTicket(ctx context.Context, id string, assignedTo string) error {
	query := `UPDATE tickets SET assigned_to = $1, updated_at = $2 WHERE id = $3`
	_, err := r.db.ExecContext(ctx, query, assignedTo, time.Now(), id)
	if err != nil {
		return fmt.Errorf("failed to assign ticket: %w", err)
	}
	return nil
}

func (r *deskPostgresRepository) handleError(err error) error {
	if err == nil {
		return nil
	}
	if err == sql.ErrNoRows {
		return fmt.Errorf("ticket not found")
	}
	return fmt.Errorf("database error: %w", err)
}
