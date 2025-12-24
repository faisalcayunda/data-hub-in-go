package repository

import (
	"context"
	"fmt"
	"time"

	"portal-data-backend/internal/auth/domain"
	"portal-data-backend/pkg/errors"

	"github.com/jmoiron/sqlx"
)

// userPostgresRepository implements UserRepository for PostgreSQL
type userPostgresRepository struct {
	db *sqlx.DB
}

// NewUserPostgresRepository creates a new user repository
func NewUserPostgresRepository(db *sqlx.DB) domain.UserRepository {
	return &userPostgresRepository{db: db}
}

// GetUserByID retrieves a user by ID
func (r *userPostgresRepository) GetUserByID(ctx context.Context, id string) (*domain.User, error) {
	query := `
		SELECT id, organization_id, role_id, name, username, employee_id, position,
		       email, password_hash, address, phone, thumbnail, status, created_at, updated_at
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

// GetUserByEmail retrieves a user by email
func (r *userPostgresRepository) GetUserByEmail(ctx context.Context, email string) (*domain.User, error) {
	query := `
		SELECT id, organization_id, role_id, name, username, employee_id, position,
		       email, password_hash, address, phone, thumbnail, status, created_at, updated_at
		FROM users
		WHERE email = $1 AND status != 'deleted'
	`

	var user domain.User
	err := r.db.GetContext(ctx, &user, query, email)
	if err != nil {
		return nil, r.handleError(err)
	}

	return &user, nil
}

// GetUserByUsername retrieves a user by username
func (r *userPostgresRepository) GetUserByUsername(ctx context.Context, username string) (*domain.User, error) {
	query := `
		SELECT id, organization_id, role_id, name, username, employee_id, position,
		       email, password_hash, address, phone, thumbnail, status, created_at, updated_at
		FROM users
		WHERE username = $1 AND status != 'deleted'
	`

	var user domain.User
	err := r.db.GetContext(ctx, &user, query, username)
	if err != nil {
		return nil, r.handleError(err)
	}

	return &user, nil
}

// CreateUser creates a new user
func (r *userPostgresRepository) CreateUser(ctx context.Context, user *domain.User) error {
	query := `
		INSERT INTO users (
			id, organization_id, role_id, name, username, employee_id, position,
			email, password_hash, address, phone, thumbnail, status, created_at, updated_at
		) VALUES (
			:id, :organization_id, :role_id, :name, :username, :employee_id, :position,
			:email, :password_hash, :address, :phone, :thumbnail, :status, :created_at, :updated_at
		)
	`

	_, err := r.db.NamedExecContext(ctx, query, user)
	if err != nil {
		return fmt.Errorf("failed to create user: %w", err)
	}

	return nil
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

// ListUsers retrieves a list of users with pagination
func (r *userPostgresRepository) ListUsers(ctx context.Context, limit, offset int) ([]*domain.User, int, error) {
	// Get total count
	countQuery := `SELECT COUNT(*) FROM users WHERE status != 'deleted'`
	var total int
	err := r.db.GetContext(ctx, &total, countQuery)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count users: %w", err)
	}

	// Get users
	query := `
		SELECT id, organization_id, role_id, name, username, employee_id, position,
		       email, password_hash, address, phone, thumbnail, status, created_at, updated_at
		FROM users
		WHERE status != 'deleted'
		ORDER BY created_at DESC
		LIMIT $1 OFFSET $2
	`

	var users []*domain.User
	err = r.db.SelectContext(ctx, &users, query, limit, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to list users: %w", err)
	}

	return users, total, nil
}

// IsEmailExists checks if email already exists
func (r *userPostgresRepository) IsEmailExists(ctx context.Context, email string) (bool, error) {
	query := `SELECT EXISTS(SELECT 1 FROM users WHERE email = $1 AND status != 'deleted')`

	var exists bool
	err := r.db.GetContext(ctx, &exists, query, email)
	if err != nil {
		return false, fmt.Errorf("failed to check email existence: %w", err)
	}

	return exists, nil
}

// IsUsernameExists checks if username already exists
func (r *userPostgresRepository) IsUsernameExists(ctx context.Context, username string) (bool, error) {
	query := `SELECT EXISTS(SELECT 1 FROM users WHERE username = $1 AND status != 'deleted')`

	var exists bool
	err := r.db.GetContext(ctx, &exists, query, username)
	if err != nil {
		return false, fmt.Errorf("failed to check username existence: %w", err)
	}

	return exists, nil
}

// handleError handles database errors
func (r *userPostgresRepository) handleError(err error) error {
	if err == nil {
		return nil
	}

	return errors.Wrap(err, "database error")
}

// tokenPostgresRepository implements TokenRepository for PostgreSQL
type tokenPostgresRepository struct {
	db *sqlx.DB
}

// NewTokenPostgresRepository creates a new token repository
func NewTokenPostgresRepository(db *sqlx.DB) domain.TokenRepository {
	return &tokenPostgresRepository{db: db}
}

// CreateToken creates a new token
func (r *tokenPostgresRepository) CreateToken(ctx context.Context, token *domain.Token) error {
	query := `
		INSERT INTO tokens (id, user_id, access_token, refresh_token, expires_at, revoked, created_at)
		VALUES (:id, :user_id, :access_token, :refresh_token, :expires_at, :revoked, :created_at)
	`

	_, err := r.db.NamedExecContext(ctx, query, token)
	if err != nil {
		return fmt.Errorf("failed to create token: %w", err)
	}

	return nil
}

// GetTokenByRefreshToken retrieves a token by refresh token
func (r *tokenPostgresRepository) GetTokenByRefreshToken(ctx context.Context, refreshToken string) (*domain.Token, error) {
	query := `
		SELECT id, user_id, access_token, refresh_token, expires_at, revoked, created_at
		FROM tokens
		WHERE refresh_token = $1
		ORDER BY created_at DESC
		LIMIT 1
	`

	var token domain.Token
	err := r.db.GetContext(ctx, &token, query, refreshToken)
	if err != nil {
		return nil, r.handleError(err)
	}

	return &token, nil
}

// GetTokenByAccessToken retrieves a token by access token
func (r *tokenPostgresRepository) GetTokenByAccessToken(ctx context.Context, accessToken string) (*domain.Token, error) {
	query := `
		SELECT id, user_id, access_token, refresh_token, expires_at, revoked, created_at
		FROM tokens
		WHERE access_token = $1
		ORDER BY created_at DESC
		LIMIT 1
	`

	var token domain.Token
	err := r.db.GetContext(ctx, &token, query, accessToken)
	if err != nil {
		return nil, r.handleError(err)
	}

	return &token, nil
}

// RevokeToken revokes a token by ID
func (r *tokenPostgresRepository) RevokeToken(ctx context.Context, id string) error {
	query := `UPDATE tokens SET revoked = true WHERE id = $1`

	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to revoke token: %w", err)
	}

	rows, _ := result.RowsAffected()
	if rows == 0 {
		return errors.ErrNotFound
	}

	return nil
}

// RevokeUserTokens revokes all tokens for a user
func (r *tokenPostgresRepository) RevokeUserTokens(ctx context.Context, userID string) error {
	query := `UPDATE tokens SET revoked = true WHERE user_id = $1`

	_, err := r.db.ExecContext(ctx, query, userID)
	if err != nil {
		return fmt.Errorf("failed to revoke user tokens: %w", err)
	}

	return nil
}

// DeleteToken deletes a token by ID
func (r *tokenPostgresRepository) DeleteToken(ctx context.Context, id string) error {
	query := `DELETE FROM tokens WHERE id = $1`

	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete token: %w", err)
	}

	rows, _ := result.RowsAffected()
	if rows == 0 {
		return errors.ErrNotFound
	}

	return nil
}

// CleanupExpiredTokens deletes expired tokens
func (r *tokenPostgresRepository) CleanupExpiredTokens(ctx context.Context) error {
	query := `DELETE FROM tokens WHERE expires_at < NOW() OR revoked = true`

	_, err := r.db.ExecContext(ctx, query)
	if err != nil {
		return fmt.Errorf("failed to cleanup expired tokens: %w", err)
	}

	return nil
}

// handleError handles database errors for token repository
func (r *tokenPostgresRepository) handleError(err error) error {
	if err == nil {
		return nil
	}

	return errors.Wrap(err, "database error")
}
