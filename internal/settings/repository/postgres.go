package repository

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"

	settingsDomain "portal-data-backend/internal/settings/domain"

	"github.com/jmoiron/sqlx"
)

type settingsPostgresRepository struct {
	db *sqlx.DB
}

func NewSettingsPostgresRepository(db *sqlx.DB) settingsDomain.Repository {
	return &settingsPostgresRepository{db: db}
}

func (r *settingsPostgresRepository) GetByID(ctx context.Context, id string) (*settingsDomain.Setting, error) {
	query := `
		SELECT id, key, value, type, category, user_id, is_public, created_at, updated_at, deleted_at
		FROM settings
		WHERE id = $1 AND deleted_at IS NULL
	`

	var setting settingsDomain.Setting
	err := r.db.GetContext(ctx, &setting, query, id)
	if err != nil {
		return nil, r.handleError(err)
	}
	return &setting, nil
}

func (r *settingsPostgresRepository) GetByKey(ctx context.Context, key string, userID *string) (*settingsDomain.Setting, error) {
	query := `
		SELECT id, key, value, type, category, user_id, is_public, created_at, updated_at, deleted_at
		FROM settings
		WHERE key = $1 AND deleted_at IS NULL
	`

	args := []interface{}{key}
	argCount := 2

	if userID != nil {
		query += fmt.Sprintf(" AND user_id = $%d", argCount)
		args = append(args, userID)
		argCount++
	} else {
		query += " AND user_id IS NULL"
	}

	var setting settingsDomain.Setting
	err := r.db.GetContext(ctx, &setting, query, args...)
	if err != nil {
		return nil, r.handleError(err)
	}
	return &setting, nil
}

func (r *settingsPostgresRepository) List(ctx context.Context, filter *settingsDomain.SettingFilter, limit, offset int) ([]*settingsDomain.Setting, int, error) {
	whereClause := "WHERE deleted_at IS NULL"
	args := []interface{}{}
	argCount := 1

	if filter != nil {
		if filter.Category != nil {
			whereClause += fmt.Sprintf(" AND category = $%d", argCount)
			args = append(args, filter.Category)
			argCount++
		}
		if filter.UserID != nil {
			whereClause += fmt.Sprintf(" AND user_id = $%d", argCount)
			args = append(args, filter.UserID)
			argCount++
		}
		if filter.Type != nil {
			whereClause += fmt.Sprintf(" AND type = $%d", argCount)
			args = append(args, filter.Type)
			argCount++
		}
		if filter.Search != "" {
			whereClause += fmt.Sprintf(" AND (key ILIKE $%d OR value ILIKE $%d)", argCount, argCount)
			searchTerm := "%" + filter.Search + "%"
			args = append(args, searchTerm, searchTerm)
			argCount += 2
		}
	}

	countQuery := "SELECT COUNT(*) FROM settings " + whereClause
	var total int
	err := r.db.GetContext(ctx, &total, countQuery, args...)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count settings: %w", err)
	}

	query := `
		SELECT id, key, value, type, category, user_id, is_public, created_at, updated_at, deleted_at
		FROM settings
	` + whereClause + " ORDER BY key ASC LIMIT $" + fmt.Sprintf("%d", argCount) + " OFFSET $" + fmt.Sprintf("%d", argCount+1)

	args = append(args, limit, offset)

	var settings []*settingsDomain.Setting
	err = r.db.SelectContext(ctx, &settings, query, args...)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to list settings: %w", err)
	}

	return settings, total, nil
}

func (r *settingsPostgresRepository) Create(ctx context.Context, setting *settingsDomain.Setting) error {
	query := `
		INSERT INTO settings (id, key, value, type, category, user_id, is_public, created_at, updated_at)
		VALUES (:id, :key, :value, :type, :category, :user_id, :is_public, :created_at, :updated_at)
	`

	_, err := r.db.NamedExecContext(ctx, query, setting)
	if err != nil {
		return fmt.Errorf("failed to create setting: %w", err)
	}
	return nil
}

func (r *settingsPostgresRepository) Update(ctx context.Context, id string, setting *settingsDomain.Setting) error {
	query := `
		UPDATE settings
		SET value = :value, type = :type, is_public = :is_public, updated_at = :updated_at
		WHERE id = :id
	`

	setting.ID = id
	_, err := r.db.NamedExecContext(ctx, query, setting)
	if err != nil {
		return fmt.Errorf("failed to update setting: %w", err)
	}
	return nil
}

func (r *settingsPostgresRepository) Delete(ctx context.Context, id string) error {
	query := `UPDATE settings SET deleted_at = $1 WHERE id = $2`
	_, err := r.db.ExecContext(ctx, query, time.Now(), id)
	if err != nil {
		return fmt.Errorf("failed to delete setting: %w", err)
	}
	return nil
}

func (r *settingsPostgresRepository) GetByKeys(ctx context.Context, keys []string, userID *string) (map[string]string, error) {
	if len(keys) == 0 {
		return make(map[string]string), nil
	}

	query := `
		SELECT key, value
		FROM settings
		WHERE deleted_at IS NULL AND key IN (?)
	`

	args := []interface{}{userID}
	query, args, _ = sqlx.In(query, keys)

	if userID != nil {
		query = strings.Replace(query, "?", "$1", 1)
		query += " AND (user_id = $2 OR user_id IS NULL)"
		args = append([]interface{}{userID}, args...)
	} else {
		query += " AND user_id IS NULL"
	}

	query = r.db.Rebind(query)
	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to get settings by keys: %w", err)
	}
	defer rows.Close()

	result := make(map[string]string)
	for rows.Next() {
		var key, value string
		if err := rows.Scan(&key, &value); err != nil {
			return nil, fmt.Errorf("failed to scan setting: %w", err)
		}
		result[key] = value
	}

	return result, nil
}

func (r *settingsPostgresRepository) GetByCategory(ctx context.Context, category string, userID *string, limit, offset int) ([]*settingsDomain.Setting, int, error) {
	whereClause := "WHERE deleted_at IS NULL AND category = $1"
	args := []interface{}{category}
	argCount := 2

	if userID != nil {
		whereClause += fmt.Sprintf(" AND user_id = $%d", argCount)
		args = append(args, userID)
		argCount++
	} else {
		whereClause += " AND user_id IS NULL"
	}

	countQuery := "SELECT COUNT(*) FROM settings " + whereClause
	var total int
	err := r.db.GetContext(ctx, &total, countQuery, args...)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count settings: %w", err)
	}

	query := `
		SELECT id, key, value, type, category, user_id, is_public, created_at, updated_at, deleted_at
		FROM settings
	` + whereClause + " ORDER BY key ASC LIMIT $" + fmt.Sprintf("%d", argCount) + " OFFSET $" + fmt.Sprintf("%d", argCount+1)

	args = append(args, limit, offset)

	var settings []*settingsDomain.Setting
	err = r.db.SelectContext(ctx, &settings, query, args...)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to get settings by category: %w", err)
	}

	return settings, total, nil
}

func (r *settingsPostgresRepository) handleError(err error) error {
	if err == nil {
		return nil
	}
	if err == sql.ErrNoRows {
		return fmt.Errorf("setting not found")
	}
	return fmt.Errorf("database error: %w", err)
}
