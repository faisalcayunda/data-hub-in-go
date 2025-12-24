package repository

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	pubDomain "portal-data-backend/internal/publication/domain"

	"github.com/jmoiron/sqlx"
)

type publicationPostgresRepository struct {
	db *sqlx.DB
}

func NewPublicationPostgresRepository(db *sqlx.DB) pubDomain.Repository {
	return &publicationPostgresRepository{db: db}
}

func (r *publicationPostgresRepository) GetByID(ctx context.Context, id string) (*pubDomain.Publication, error) {
	query := `
		SELECT id, title, description, content, doi, publisher, published_date, dataset_id, organization_id,
		       authors, tags, status, is_featured, view_count, download_count,
		       created_by, updated_by, created_at, updated_at, deleted_at
		FROM publications
		WHERE id = $1 AND deleted_at IS NULL
	`

	var pub pubDomain.Publication
	err := r.db.GetContext(ctx, &pub, query, id)
	if err != nil {
		return nil, r.handleError(err)
	}
	return &pub, nil
}

func (r *publicationPostgresRepository) List(ctx context.Context, filter *pubDomain.PublicationFilter, limit, offset int) ([]*pubDomain.Publication, int, error) {
	whereClause := "WHERE deleted_at IS NULL"
	args := []interface{}{}
	argCount := 1

	if filter != nil {
		if filter.DatasetID != nil {
			whereClause += fmt.Sprintf(" AND dataset_id = $%d", argCount)
			args = append(args, filter.DatasetID)
			argCount++
		}
		if filter.OrganizationID != nil {
			whereClause += fmt.Sprintf(" AND organization_id = $%d", argCount)
			args = append(args, filter.OrganizationID)
			argCount++
		}
		if filter.Status != nil {
			whereClause += fmt.Sprintf(" AND status = $%d", argCount)
			args = append(args, filter.Status)
			argCount++
		}
		if filter.IsFeatured != nil {
			whereClause += fmt.Sprintf(" AND is_featured = $%d", argCount)
			args = append(args, filter.IsFeatured)
			argCount++
		}
		if filter.Search != "" {
			whereClause += fmt.Sprintf(" AND (title ILIKE $%d OR description ILIKE $%d)", argCount, argCount)
			searchTerm := "%" + filter.Search + "%"
			args = append(args, searchTerm, searchTerm)
			argCount += 2
		}
	}

	countQuery := "SELECT COUNT(*) FROM publications " + whereClause
	var total int
	err := r.db.GetContext(ctx, &total, countQuery, args...)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count publications: %w", err)
	}

	query := `
		SELECT id, title, description, content, doi, publisher, published_date, dataset_id, organization_id,
		       authors, tags, status, is_featured, view_count, download_count,
		       created_by, updated_by, created_at, updated_at, deleted_at
		FROM publications
	` + whereClause + " ORDER BY created_at DESC LIMIT $" + fmt.Sprintf("%d", argCount) + " OFFSET $" + fmt.Sprintf("%d", argCount+1)

	args = append(args, limit, offset)

	var pubs []*pubDomain.Publication
	err = r.db.SelectContext(ctx, &pubs, query, args...)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to list publications: %w", err)
	}

	return pubs, total, nil
}

func (r *publicationPostgresRepository) Create(ctx context.Context, pub *pubDomain.Publication) error {
	query := `
		INSERT INTO publications (
			id, title, description, content, doi, publisher, published_date, dataset_id, organization_id,
			authors, tags, status, is_featured, view_count, download_count,
			created_by, updated_by, created_at, updated_at
		) VALUES (
			:id, :title, :description, :content, :doi, :publisher, :published_date, :dataset_id, :organization_id,
			:authors, :tags, :status, :is_featured, :view_count, :download_count,
			:created_by, :updated_by, :created_at, :updated_at
		)
	`

	_, err := r.db.NamedExecContext(ctx, query, pub)
	if err != nil {
		return fmt.Errorf("failed to create publication: %w", err)
	}
	return nil
}

func (r *publicationPostgresRepository) Update(ctx context.Context, id string, pub *pubDomain.Publication) error {
	query := `
		UPDATE publications
		SET title = :title, description = :description, content = :content, doi = :doi, publisher = :publisher,
		    published_date = :published_date, dataset_id = :dataset_id, organization_id = :organization_id,
		    authors = :authors, tags = :tags, status = :status, is_featured = :is_featured,
		    updated_by = :updated_by, updated_at = :updated_at
		WHERE id = :id
	`

	pub.ID = id
	_, err := r.db.NamedExecContext(ctx, query, pub)
	if err != nil {
		return fmt.Errorf("failed to update publication: %w", err)
	}
	return nil
}

func (r *publicationPostgresRepository) Delete(ctx context.Context, id string) error {
	query := `UPDATE publications SET deleted_at = $1 WHERE id = $2`
	_, err := r.db.ExecContext(ctx, query, time.Now(), id)
	if err != nil {
		return fmt.Errorf("failed to delete publication: %w", err)
	}
	return nil
}

func (r *publicationPostgresRepository) UpdateStatus(ctx context.Context, id string, status string) error {
	query := `UPDATE publications SET status = $1, updated_at = $2 WHERE id = $3`
	_, err := r.db.ExecContext(ctx, query, status, time.Now(), id)
	if err != nil {
		return fmt.Errorf("failed to update publication status: %w", err)
	}
	return nil
}

func (r *publicationPostgresRepository) IncrementViewCount(ctx context.Context, id string) error {
	query := `UPDATE publications SET view_count = view_count + 1 WHERE id = $1`
	_, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to increment view count: %w", err)
	}
	return nil
}

func (r *publicationPostgresRepository) IncrementDownloadCount(ctx context.Context, id string) error {
	query := `UPDATE publications SET download_count = download_count + 1 WHERE id = $1`
	_, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to increment download count: %w", err)
	}
	return nil
}

func (r *publicationPostgresRepository) GetByDatasetID(ctx context.Context, datasetID string, limit, offset int) ([]*pubDomain.Publication, int, error) {
	query := `
		SELECT id, title, description, content, doi, publisher, published_date, dataset_id, organization_id,
		       authors, tags, status, is_featured, view_count, download_count,
		       created_by, updated_by, created_at, updated_at, deleted_at
		FROM publications
		WHERE dataset_id = $1 AND deleted_at IS NULL
		ORDER BY published_date DESC
		LIMIT $2 OFFSET $3
	`

	var pubs []*pubDomain.Publication
	err := r.db.SelectContext(ctx, &pubs, query, datasetID, limit, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to get dataset publications: %w", err)
	}

	countQuery := `SELECT COUNT(*) FROM publications WHERE dataset_id = $1 AND deleted_at IS NULL`
	var total int
	err = r.db.GetContext(ctx, &total, countQuery, datasetID)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count dataset publications: %w", err)
	}

	return pubs, total, nil
}

func (r *publicationPostgresRepository) GetByOrganizationID(ctx context.Context, orgID string, limit, offset int) ([]*pubDomain.Publication, int, error) {
	query := `
		SELECT id, title, description, content, doi, publisher, published_date, dataset_id, organization_id,
		       authors, tags, status, is_featured, view_count, download_count,
		       created_by, updated_by, created_at, updated_at, deleted_at
		FROM publications
		WHERE organization_id = $1 AND deleted_at IS NULL
		ORDER BY published_date DESC
		LIMIT $2 OFFSET $3
	`

	var pubs []*pubDomain.Publication
	err := r.db.SelectContext(ctx, &pubs, query, orgID, limit, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to get organization publications: %w", err)
	}

	countQuery := `SELECT COUNT(*) FROM publications WHERE organization_id = $1 AND deleted_at IS NULL`
	var total int
	err = r.db.GetContext(ctx, &total, countQuery, orgID)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count organization publications: %w", err)
	}

	return pubs, total, nil
}

func (r *publicationPostgresRepository) handleError(err error) error {
	if err == nil {
		return nil
	}
	if err == sql.ErrNoRows {
		return fmt.Errorf("publication not found")
	}
	return fmt.Errorf("database error: %w", err)
}
