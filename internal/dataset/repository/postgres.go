package repository

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"portal-data-backend/internal/dataset/domain"
	"portal-data-backend/pkg/errors"

	"github.com/jmoiron/sqlx"
)

// datasetPostgresRepository implements Repository for PostgreSQL
type datasetPostgresRepository struct {
	db *sqlx.DB
}

// NewDatasetPostgresRepository creates a new dataset repository
func NewDatasetPostgresRepository(db *sqlx.DB) domain.Repository {
	return &datasetPostgresRepository{db: db}
}

func (r *datasetPostgresRepository) GetByID(ctx context.Context, id string) (*domain.Dataset, error) {
	query := `
		SELECT
			d.id, d.name, d.slug, d.description, d.period, d.unit_id, d.business_field_id,
			d.image, d.topic_id, d.organization_id, d.reference_id, d.classification,
			d.category, d.data_fixed, d.validation_status, d.metadatas, d.created_by,
			d.updated_by, d.created_at, d.updated_at, d.is_highlight, d.status,
			o.id as org_id, o.name as org_name, o.slug as org_slug,
			u.id as unit_id, u.name as unit_name, u.symbol as unit_symbol,
			bf.id as bf_id, bf.name as bf_name, bf.slug as bf_slug,
			t.id as topic_id, t.name as topic_name, t.slug as topic_slug
		FROM datasets d
		LEFT JOIN organizations o ON d.organization_id = o.id
		LEFT JOIN units u ON d.unit_id = u.id
		LEFT JOIN business_fields bf ON d.business_field_id = bf.id
		LEFT JOIN topics t ON d.topic_id = t.id
		WHERE d.id = $1
	`

	dataset, err := r.scanDataset(ctx, query, id)
	if err != nil {
		return nil, err
	}

	// Get tags
	tags, err := r.getTagsByDatasetID(ctx, id)
	if err == nil {
		dataset.Tags = tags
	}

	return dataset, nil
}

func (r *datasetPostgresRepository) GetBySlug(ctx context.Context, slug string) (*domain.Dataset, error) {
	query := `
		SELECT
			d.id, d.name, d.slug, d.description, d.period, d.unit_id, d.business_field_id,
			d.image, d.topic_id, d.organization_id, d.reference_id, d.classification,
			d.category, d.data_fixed, d.validation_status, d.metadatas, d.created_by,
			d.updated_by, d.created_at, d.updated_at, d.is_highlight, d.status,
			o.id as org_id, o.name as org_name, o.slug as org_slug,
			u.id as unit_id, u.name as unit_name, u.symbol as unit_symbol,
			bf.id as bf_id, bf.name as bf_name, bf.slug as bf_slug,
			t.id as topic_id, t.name as topic_name, t.slug as topic_slug
		FROM datasets d
		LEFT JOIN organizations o ON d.organization_id = o.id
		LEFT JOIN units u ON d.unit_id = u.id
		LEFT JOIN business_fields bf ON d.business_field_id = bf.id
		LEFT JOIN topics t ON d.topic_id = t.id
		WHERE d.slug = $1
	`

	dataset, err := r.scanDataset(ctx, query, slug)
	if err != nil {
		return nil, err
	}

	tags, err := r.getTagsByDatasetID(ctx, dataset.ID)
	if err == nil {
		dataset.Tags = tags
	}

	return dataset, nil
}

func (r *datasetPostgresRepository) List(ctx context.Context, filter *domain.DatasetFilter, limit, offset int, sortBy, sortOrder string) ([]*domain.Dataset, int, error) {
	whereClause, args := r.buildWhereClause(filter)

	countQuery := "SELECT COUNT(*) FROM datasets d " + whereClause
	var total int
	err := r.db.GetContext(ctx, &total, countQuery, args...)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count datasets: %w", err)
	}

	orderClause := r.buildOrderClause(sortBy, sortOrder)
	query := `
		SELECT
			d.id, d.name, d.slug, d.description, d.period, d.unit_id, d.business_field_id,
			d.image, d.topic_id, d.organization_id, d.reference_id, d.classification,
			d.category, d.data_fixed, d.validation_status, d.metadatas, d.created_by,
			d.updated_by, d.created_at, d.updated_at, d.is_highlight, d.status,
			o.id as org_id, o.name as org_name, o.slug as org_slug,
			u.id as unit_id, u.name as unit_name, u.symbol as unit_symbol,
			bf.id as bf_id, bf.name as bf_name, bf.slug as bf_slug,
			t.id as topic_id, t.name as topic_name, t.slug as topic_slug
		FROM datasets d
		LEFT JOIN organizations o ON d.organization_id = o.id
		LEFT JOIN units u ON d.unit_id = u.id
		LEFT JOIN business_fields bf ON d.business_field_id = bf.id
		LEFT JOIN topics t ON d.topic_id = t.id
	` + whereClause + " " + orderClause + " LIMIT $" + fmt.Sprintf("%d", len(args)+1) + " OFFSET $" + fmt.Sprintf("%d", len(args)+2)

	args = append(args, limit, offset)

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to list datasets: %w", err)
	}
	defer rows.Close()

	var datasets []*domain.Dataset
	for rows.Next() {
		dataset, err := r.scanRow(rows)
		if err != nil {
			return nil, 0, err
		}
		datasets = append(datasets, dataset)
	}

	return datasets, total, nil
}

func (r *datasetPostgresRepository) Create(ctx context.Context, dataset *domain.Dataset, tagIDs []string) error {
	tx, err := r.db.BeginTxx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	insertQuery := `
		INSERT INTO datasets (
			id, name, slug, description, period, unit_id, business_field_id, image,
			topic_id, organization_id, reference_id, classification, category,
			data_fixed, validation_status, metadatas, created_by, updated_by,
			created_at, updated_at, is_highlight, status
		) VALUES (
			:id, :name, :slug, :description, :period, :unit_id, :business_field_id, :image,
			:topic_id, :organization_id, :reference_id, :classification, :category,
			:data_fixed, :validation_status, :metadatas, :created_by, :updated_by,
			:created_at, :updated_at, :is_highlight, :status
		)
	`

	_, err = tx.NamedExecContext(ctx, insertQuery, dataset)
	if err != nil {
		return fmt.Errorf("failed to create dataset: %w", err)
	}

	// Insert tags
	if len(tagIDs) > 0 {
		for _, tagID := range tagIDs {
			_, err = tx.ExecContext(ctx, `INSERT INTO dataset_tag_link (dataset_id, tag_id) VALUES ($1, $2)`, dataset.ID, tagID)
			if err != nil {
				return fmt.Errorf("failed to link tags: %w", err)
			}
		}
	}

	return tx.Commit()
}

func (r *datasetPostgresRepository) Update(ctx context.Context, dataset *domain.Dataset, tagIDs []string) error {
	tx, err := r.db.BeginTxx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	updateQuery := `
		UPDATE datasets SET
			name = :name, slug = :slug, description = :description, period = :period,
			unit_id = :unit_id, business_field_id = :business_field_id, image = :image,
			topic_id = :topic_id, reference_id = :reference_id, classification = :classification,
			category = :category, data_fixed = :data_fixed, validation_status = :validation_status,
			metadatas = :metadatas, updated_by = :updated_by, updated_at = :updated_at,
			is_highlight = :is_highlight, status = :status
		WHERE id = :id
	`

	result, err := tx.NamedExecContext(ctx, updateQuery, dataset)
	if err != nil {
		return fmt.Errorf("failed to update dataset: %w", err)
	}

	rows, _ := result.RowsAffected()
	if rows == 0 {
		return errors.ErrNotFound
	}

	// Update tags - delete existing and insert new
	_, err = tx.ExecContext(ctx, `DELETE FROM dataset_tag_link WHERE dataset_id = $1`, dataset.ID)
	if err != nil {
		return fmt.Errorf("failed to delete existing tags: %w", err)
	}

	if len(tagIDs) > 0 {
		for _, tagID := range tagIDs {
			_, err = tx.ExecContext(ctx, `INSERT INTO dataset_tag_link (dataset_id, tag_id) VALUES ($1, $2)`, dataset.ID, tagID)
			if err != nil {
				return fmt.Errorf("failed to link tags: %w", err)
			}
		}
	}

	return tx.Commit()
}

func (r *datasetPostgresRepository) Delete(ctx context.Context, id string) error {
	query := `UPDATE datasets SET status = 'archived', updated_at = NOW() WHERE id = $1`
	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete dataset: %w", err)
	}

	rows, _ := result.RowsAffected()
	if rows == 0 {
		return errors.ErrNotFound
	}
	return nil
}

func (r *datasetPostgresRepository) UpdateStatus(ctx context.Context, id string, status domain.DatasetStatus) error {
	query := `UPDATE datasets SET status = $1, updated_at = NOW() WHERE id = $2`
	result, err := r.db.ExecContext(ctx, query, status, id)
	if err != nil {
		return fmt.Errorf("failed to update dataset status: %w", err)
	}

	rows, _ := result.RowsAffected()
	if rows == 0 {
		return errors.ErrNotFound
	}
	return nil
}

func (r *datasetPostgresRepository) GetByOrganizationID(ctx context.Context, orgID string, limit, offset int) ([]*domain.Dataset, int, error) {
	filter := &domain.DatasetFilter{OrganizationID: orgID}
	return r.List(ctx, filter, limit, offset, "created_at", "DESC")
}

// Helper functions

func (r *datasetPostgresRepository) scanDataset(ctx context.Context, query string, arg interface{}) (*domain.Dataset, error) {
	row := r.db.QueryRowxContext(ctx, query, arg)
	dataset, err := r.scanRowFromQueryx(row)
	if err != nil {
		return nil, r.handleError(err)
	}
	return dataset, nil
}

func (r *datasetPostgresRepository) scanRow(rows *sql.Rows) (*domain.Dataset, error) {
	var dataset domain.Dataset
	var orgName, orgSlug *string
	var unitName, unitSymbol *string
	var bfName, bfSlug *string
	var topicName, topicSlug *string

	err := rows.Scan(
		&dataset.ID, &dataset.Name, &dataset.Slug, &dataset.Description, &dataset.Period,
		&dataset.UnitID, &dataset.BusinessFieldID, &dataset.Image, &dataset.TopicID,
		&dataset.OrganizationID, &dataset.ReferenceID, &dataset.Classification,
		&dataset.Category, &dataset.DataFixed, &dataset.ValidationStatus, &dataset.Metadata,
		&dataset.CreatedBy, &dataset.UpdatedBy, &dataset.CreatedAt, &dataset.UpdatedAt,
		&dataset.IsHighlight, &dataset.Status,
		&orgName, &orgSlug, &unitName, &unitSymbol, &bfName, &bfSlug, &topicName, &topicSlug,
	)
	if err != nil {
		return nil, err
	}

	// Populate relations
	if orgName != nil {
		dataset.Organization = &domain.OrganizationSummary{
			ID:   dataset.OrganizationID,
			Name: *orgName,
			Slug: *orgSlug,
		}
	}
	if unitName != nil {
		dataset.Unit = &domain.Unit{
			ID:    *dataset.UnitID,
			Name:  *unitName,
			Symbol: *unitSymbol,
		}
	}
	if bfName != nil {
		dataset.BusinessField = &domain.BusinessField{
			ID:   *dataset.BusinessFieldID,
			Name: *bfName,
			Slug: *bfSlug,
		}
	}
	if topicName != nil {
		dataset.Topic = &domain.Topic{
			ID:   *dataset.TopicID,
			Name: *topicName,
			Slug: *topicSlug,
		}
	}

	return &dataset, nil
}

func (r *datasetPostgresRepository) scanRowFromQueryx(row *sqlx.Row) (*domain.Dataset, error) {
	var dataset domain.Dataset
	var orgName, orgSlug *string
	var unitName, unitSymbol *string
	var bfName, bfSlug *string
	var topicName, topicSlug *string

	err := row.Scan(
		&dataset.ID, &dataset.Name, &dataset.Slug, &dataset.Description, &dataset.Period,
		&dataset.UnitID, &dataset.BusinessFieldID, &dataset.Image, &dataset.TopicID,
		&dataset.OrganizationID, &dataset.ReferenceID, &dataset.Classification,
		&dataset.Category, &dataset.DataFixed, &dataset.ValidationStatus, &dataset.Metadata,
		&dataset.CreatedBy, &dataset.UpdatedBy, &dataset.CreatedAt, &dataset.UpdatedAt,
		&dataset.IsHighlight, &dataset.Status,
		&orgName, &orgSlug, &unitName, &unitSymbol, &bfName, &bfSlug, &topicName, &topicSlug,
	)
	if err != nil {
		return nil, err
	}

	// Populate relations
	if orgName != nil {
		dataset.Organization = &domain.OrganizationSummary{
			ID:   dataset.OrganizationID,
			Name: *orgName,
			Slug: *orgSlug,
		}
	}
	if unitName != nil {
		dataset.Unit = &domain.Unit{
			ID:    *dataset.UnitID,
			Name:  *unitName,
			Symbol: *unitSymbol,
		}
	}
	if bfName != nil {
		dataset.BusinessField = &domain.BusinessField{
			ID:   *dataset.BusinessFieldID,
			Name: *bfName,
			Slug: *bfSlug,
		}
	}
	if topicName != nil {
		dataset.Topic = &domain.Topic{
			ID:   *dataset.TopicID,
			Name: *topicName,
			Slug: *topicSlug,
		}
	}

	return &dataset, nil
}

func (r *datasetPostgresRepository) getTagsByDatasetID(ctx context.Context, datasetID string) ([]domain.Tag, error) {
	query := `
		SELECT t.id, t.name, t.slug
		FROM tags t
		INNER JOIN dataset_tag_link dtl ON t.id = dtl.tag_id
		WHERE dtl.dataset_id = $1
	`

	var tags []domain.Tag
	err := r.db.SelectContext(ctx, &tags, query, datasetID)
	if err != nil {
		return nil, err
	}
	return tags, nil
}

func (r *datasetPostgresRepository) buildWhereClause(filter *domain.DatasetFilter) (string, []interface{}) {
	whereClause := "WHERE 1=1"
	args := []interface{}{}
	argCount := 1

	if filter == nil {
		return whereClause, args
	}

	if filter.OrganizationID != "" {
		whereClause += fmt.Sprintf(" AND d.organization_id = $%d", argCount)
		args = append(args, filter.OrganizationID)
		argCount++
	}
	if filter.TopicID != "" {
		whereClause += fmt.Sprintf(" AND d.topic_id = $%d", argCount)
		args = append(args, filter.TopicID)
		argCount++
	}
	if filter.BusinessFieldID != "" {
		whereClause += fmt.Sprintf(" AND d.business_field_id = $%d", argCount)
		args = append(args, filter.BusinessFieldID)
		argCount++
	}
	if filter.Status != "" {
		whereClause += fmt.Sprintf(" AND d.status = $%d", argCount)
		args = append(args, filter.Status)
		argCount++
	}
	if filter.ValidationStatus != "" {
		whereClause += fmt.Sprintf(" AND d.validation_status = $%d", argCount)
		args = append(args, filter.ValidationStatus)
		argCount++
	}
	if filter.Classification != "" {
		whereClause += fmt.Sprintf(" AND d.classification = $%d", argCount)
		args = append(args, filter.Classification)
		argCount++
	}
	if filter.Search != "" {
		whereClause += fmt.Sprintf(" AND (d.name ILIKE $%d OR d.description ILIKE $%d)", argCount, argCount)
		searchTerm := "%" + filter.Search + "%"
		args = append(args, searchTerm, searchTerm)
		argCount += 2
	}
	if filter.TagID != "" {
		whereClause += fmt.Sprintf(" AND EXISTS(SELECT 1 FROM dataset_tag_link dtl WHERE dtl.dataset_id = d.id AND dtl.tag_id = $%d)", argCount)
		args = append(args, filter.TagID)
		argCount++
	}

	return whereClause, args
}

func (r *datasetPostgresRepository) buildOrderClause(sortBy, sortOrder string) string {
	allowedColumns := map[string]bool{
		"name":        true,
		"created_at":  true,
		"updated_at":  true,
		"category":    true,
		"classification": true,
	}

	if !allowedColumns[sortBy] {
		sortBy = "created_at"
	}

	sortOrder = strings.ToUpper(sortOrder)
	if sortOrder != "ASC" && sortOrder != "DESC" {
		sortOrder = "DESC"
	}

	return fmt.Sprintf("ORDER BY d.%s %s", sortBy, sortOrder)
}

func (r *datasetPostgresRepository) handleError(err error) error {
	if err == nil {
		return nil
	}
	return errors.Wrap(err, "database error")
}
