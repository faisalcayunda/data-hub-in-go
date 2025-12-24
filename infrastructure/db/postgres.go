package db

import (
	"context"
	"fmt"
	"time"

	"portal-data-backend/infrastructure/config"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

// Postgres wraps sqlx.DB with additional functionality
type Postgres struct {
	DB *sqlx.DB
}

// NewPostgres creates a new PostgreSQL connection
func NewPostgres(cfg *config.DatabaseConfig) (*Postgres, error) {
	dsn := cfg.DSN()

	db, err := sqlx.Connect("postgres", dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	// Configure connection pool
	db.SetMaxOpenConns(cfg.MaxOpenConns)
	db.SetMaxIdleConns(cfg.MaxIdleConns)
	db.SetConnMaxLifetime(cfg.MaxLifetime)

	// Test connection
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := db.PingContext(ctx); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	return &Postgres{DB: db}, nil
}

// Close closes the database connection
func (p *Postgres) Close() error {
	return p.DB.Close()
}

// Health checks the database connection
func (p *Postgres) Health(ctx context.Context) error {
	return p.DB.PingContext(ctx)
}

// Transaction executes a function within a transaction
func (p *Postgres) Transaction(ctx context.Context, fn func(*sqlx.Tx) error) error {
	tx, err := p.DB.BeginTxx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}

	defer func() {
		if p := recover(); p != nil {
			_ = tx.Rollback()
			panic(p)
		}
	}()

	if err := fn(tx); err != nil {
		if rbErr := tx.Rollback(); rbErr != nil {
			return fmt.Errorf("tx error: %v, rollback error: %w", err, rbErr)
		}
		return err
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

// QueryHelper provides common query building helpers
type QueryHelper struct{}

// NewQueryHelper creates a new query helper
func NewQueryHelper() *QueryHelper {
	return &QueryHelper{}
}

// BuildWhereClause builds a WHERE clause with dynamic conditions
func (q *QueryHelper) BuildWhereClause(conditions map[string]interface{}) (string, []interface{}, error) {
	if len(conditions) == 0 {
		return "", nil, nil
	}

	whereClause := "WHERE "
	args := make([]interface{}, 0, len(conditions))
	i := 1

	for key, value := range conditions {
		if value != nil && value != "" {
			if i > 1 {
				whereClause += " AND "
			}
			whereClause += fmt.Sprintf("%s = $%d", key, i)
			args = append(args, value)
			i++
		}
	}

	if i == 1 {
		return "", nil, nil
	}

	return whereClause, args, nil
}

// BuildPaginationClause builds ORDER BY and LIMIT/OFFSET clauses
func (q *QueryHelper) BuildPaginationClause(sortBy, sortOrder string, limit, offset int) string {
	clause := ""

	if sortBy != "" {
		clause = fmt.Sprintf("ORDER BY %s", sortBy)
		if sortOrder == "DESC" {
			clause += " DESC"
		} else {
			clause += " ASC"
		}
	}

	if limit > 0 {
		clause += fmt.Sprintf(" LIMIT %d", limit)
	}

	if offset > 0 {
		clause += fmt.Sprintf(" OFFSET %d", offset)
	}

	return clause
}

// CountHelper executes a COUNT query
func (p *Postgres) CountHelper(ctx context.Context, query string, args ...interface{}) (int, error) {
	var count int
	err := p.DB.GetContext(ctx, &count, query, args...)
	if err != nil {
		return 0, fmt.Errorf("failed to execute count query: %w", err)
	}
	return count, nil
}
