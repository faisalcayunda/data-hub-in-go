package repository

import (
	"context"
	"fmt"
	"strings"
	"time"

	"portal-data-backend/internal/user/domain"
	"portal-data-backend/pkg/errors"

	"github.com/jmoiron/sqlx"
)

// userPostgresRepository implements Repository for PostgreSQL
type userPostgresRepository struct {
	db *sqlx.DB
}

// NewUserPostgresRepository creates a new user repository
func NewUserPostgresRepository(db *sqlx.DB) domain.Repository {
	return &userPostgresRepository{db: db}
}

// GetUserByID retrieves a user by ID
func (r *userPostgresRepository) GetUserByID(ctx context.Context, id string) (*domain.User, error) {
	query := `
		SELECT id, organization_id, role_id, name, username, employee_id, position,
		       email, password_hash, address, phone, thumbnail, status, bio, birth_date,
		       created_at, updated_at
		FROM users
		WHERE id = $1 AND status != 'deleted'
	`

	var user domain.User
	err := r.db.GetContext(ctx, &user, query, id)
	if err != nil {
		return nil, r.handleError(err)
	}

	return &user, nil
}

// ListUsers retrieves a list of users with filters and pagination
func (r *userPostgresRepository) ListUsers(ctx context.Context, filter *domain.UserFilter, limit, offset int, sortBy, sortOrder string) ([]*domain.User, int, error) {
	// Build WHERE clause
	whereClause := "WHERE status != 'deleted'"
	args := make([]interface{}, 0)
	argCount := 1

	if filter != nil {
		if filter.OrganizationID != "" {
			whereClause += fmt.Sprintf(" AND organization_id = $%d", argCount)
			args = append(args, filter.OrganizationID)
			argCount++
		}
		if filter.RoleID != "" {
			whereClause += fmt.Sprintf(" AND role_id = $%d", argCount)
			args = append(args, filter.RoleID)
			argCount++
		}
		if filter.Status != "" {
			whereClause += fmt.Sprintf(" AND status = $%d", argCount)
			args = append(args, filter.Status)
			argCount++
		}
		if filter.Search != "" {
			whereClause += fmt.Sprintf(" AND (name ILIKE $%d OR username ILIKE $%d OR email ILIKE $%d)", argCount, argCount, argCount)
			searchTerm := "%" + filter.Search + "%"
			args = append(args, searchTerm, searchTerm, searchTerm)
			argCount += 3
		}
	}

	// Get total count
	countQuery := "SELECT COUNT(*) FROM users " + whereClause
	var total int
	err := r.db.GetContext(ctx, &total, countQuery, args...)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count users: %w", err)
	}

	// Build ORDER BY clause
	orderClause := r.buildOrderClause(sortBy, sortOrder)

	// Get users
	query := `
		SELECT id, organization_id, role_id, name, username, employee_id, position,
		       email, password_hash, address, phone, thumbnail, status, bio, birth_date,
		       created_at, updated_at
		FROM users
	` + whereClause + " " + orderClause + " LIMIT $" + fmt.Sprintf("%d", argCount) + " OFFSET $" + fmt.Sprintf("%d", argCount+1)

	args = append(args, limit, offset)

	var users []*domain.User
	err = r.db.SelectContext(ctx, &users, query, args...)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to list users: %w", err)
	}

	return users, total, nil
}

// UpdateUser updates an existing user
func (r *userPostgresRepository) UpdateUser(ctx context.Context, user *domain.User) error {
	user.UpdatedAt = time.Now()

	query := `
		UPDATE users SET
			organization_id = :organization_id,
			role_id = :role_id,
			name = :name,
			username = :username,
			employee_id = :employee_id,
			position = :position,
			email = :email,
			password_hash = :password_hash,
			address = :address,
			phone = :phone,
			thumbnail = :thumbnail,
			status = :status,
			bio = :bio,
			birth_date = :birth_date,
			updated_at = :updated_at
		WHERE id = :id
	`

	result, err := r.db.NamedExecContext(ctx, query, user)
	if err != nil {
		return fmt.Errorf("failed to update user: %w", err)
	}

	rows, _ := result.RowsAffected()
	if rows == 0 {
		return errors.ErrNotFound
	}

	return nil
}

// DeleteUser soft deletes a user by ID
func (r *userPostgresRepository) DeleteUser(ctx context.Context, id string) error {
	query := `
		UPDATE users
		SET status = 'deleted', updated_at = NOW()
		WHERE id = $1
	`

	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete user: %w", err)
	}

	rows, _ := result.RowsAffected()
	if rows == 0 {
		return errors.ErrNotFound
	}

	return nil
}

// UpdateStatus updates user status
func (r *userPostgresRepository) UpdateStatus(ctx context.Context, id string, status string) error {
	query := `
		UPDATE users
		SET status = $1, updated_at = NOW()
		WHERE id = $2
	`

	result, err := r.db.ExecContext(ctx, query, status, id)
	if err != nil {
		return fmt.Errorf("failed to update user status: %w", err)
	}

	rows, _ := result.RowsAffected()
	if rows == 0 {
		return errors.ErrNotFound
	}

	return nil
}

// buildOrderClause builds a safe ORDER BY clause
func (r *userPostgresRepository) buildOrderClause(sortBy, sortOrder string) string {
	// Whitelist allowed columns
	allowedColumns := map[string]bool{
		"name":        true,
		"username":    true,
		"email":       true,
		"status":      true,
		"created_at":  true,
		"updated_at":  true,
	}

	if !allowedColumns[sortBy] {
		sortBy = "created_at"
	}

	sortOrder = strings.ToUpper(sortOrder)
	if sortOrder != "ASC" && sortOrder != "DESC" {
		sortOrder = "DESC"
	}

	return fmt.Sprintf("ORDER BY %s %s", sortBy, sortOrder)
}

// handleError handles database errors
func (r *userPostgresRepository) handleError(err error) error {
	if err == nil {
		return nil
	}
	return errors.Wrap(err, "database error")
}
