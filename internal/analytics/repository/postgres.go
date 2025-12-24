package repository

import (
	"context"
	"fmt"

	analyticsDomain "portal-data-backend/internal/analytics/domain"

	"github.com/jmoiron/sqlx"
)

type analyticsPostgresRepository struct {
	db *sqlx.DB
}

func NewAnalyticsPostgresRepository(db *sqlx.DB) analyticsDomain.Repository {
	return &analyticsPostgresRepository{db: db}
}

func (r *analyticsPostgresRepository) GetDatasetStats(ctx context.Context) (*analyticsDomain.DatasetStats, error) {
	query := `
		SELECT
			COUNT(*) as total_datasets,
			COUNT(*) FILTER (WHERE status = 'published') as published_count,
			COUNT(*) FILTER (WHERE status = 'draft') as draft_count,
			COUNT(*) FILTER (WHERE status = 'archived') as archived_count,
			COALESCE(SUM(downloads), 0) as total_downloads,
			COALESCE(SUM(views), 0) as total_views,
			COALESCE(MAX(updated_at), NOW()) as last_updated
		FROM datasets
		WHERE deleted_at IS NULL
	`

	var stats analyticsDomain.DatasetStats
	err := r.db.GetContext(ctx, &stats, query)
	if err != nil {
		return nil, fmt.Errorf("failed to get dataset stats: %w", err)
	}

	return &stats, nil
}

func (r *analyticsPostgresRepository) GetOrganizationStats(ctx context.Context) (*analyticsDomain.OrganizationStats, error) {
	query := `
		SELECT
			COUNT(*) as total_organizations,
			COUNT(*) FILTER (WHERE status = 'active') as active_organizations,
			COALESCE(SUM(dataset_count), 0) as total_datasets,
			NOW() as last_updated
		FROM organizations
		WHERE deleted_at IS NULL
	`

	var stats analyticsDomain.OrganizationStats
	err := r.db.GetContext(ctx, &stats, query)
	if err != nil {
		return nil, fmt.Errorf("failed to get organization stats: %w", err)
	}

	return &stats, nil
}

func (r *analyticsPostgresRepository) GetUserStats(ctx context.Context) (*analyticsDomain.UserStats, error) {
	query := `
		SELECT
			COUNT(*) as total_users,
			COUNT(*) FILTER (WHERE last_login > NOW() - INTERVAL '30 days') as active_users,
			COUNT(*) FILTER (WHERE created_at > DATE_TRUNC('month', NOW())) as new_users_this_month,
			NOW() as last_updated
		FROM users
		WHERE deleted_at IS NULL
	`

	var stats analyticsDomain.UserStats
	err := r.db.GetContext(ctx, &stats, query)
	if err != nil {
		return nil, fmt.Errorf("failed to get user stats: %w", err)
	}

	return &stats, nil
}

func (r *analyticsPostgresRepository) GetPopularDatasets(ctx context.Context, limit int) ([]analyticsDomain.PopularDataset, error) {
	query := `
		SELECT
			d.id,
			d.title,
			o.name as organization,
			COALESCE(d.views, 0) as views,
			COALESCE(d.downloads, 0) as downloads
		FROM datasets d
		JOIN organizations o ON d.organization_id = o.id
		WHERE d.deleted_at IS NULL AND d.status = 'published'
		ORDER BY (d.views + d.downloads) DESC
		LIMIT $1
	`

	var datasets []analyticsDomain.PopularDataset
	err := r.db.SelectContext(ctx, &datasets, query, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to get popular datasets: %w", err)
	}

	return datasets, nil
}

func (r *analyticsPostgresRepository) GetPopularTags(ctx context.Context, limit int) ([]analyticsDomain.TagStats, error) {
	query := `
		SELECT
			t.id as tag_id,
			t.name,
			COUNT(dt.dataset_id) as dataset_count
		FROM tags t
		LEFT JOIN dataset_tags dt ON t.id = dt.tag_id
		LEFT JOIN datasets d ON dt.dataset_id = d.id AND d.deleted_at IS NULL
		WHERE t.deleted_at IS NULL
		GROUP BY t.id, t.name
		ORDER BY dataset_count DESC
		LIMIT $1
	`

	var tags []analyticsDomain.TagStats
	err := r.db.SelectContext(ctx, &tags, query, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to get popular tags: %w", err)
	}

	return tags, nil
}

func (r *analyticsPostgresRepository) GetDatasetTrend(ctx context.Context, period string, limit int) ([]analyticsDomain.TimeSeriesData, error) {
	var interval string
	switch period {
	case "hourly":
		interval = "hour"
	case "daily":
		interval = "day"
	case "weekly":
		interval = "week"
	case "monthly":
		interval = "month"
	default:
		interval = "day"
	}

	query := fmt.Sprintf(`
		SELECT
			DATE_TRUNC('%s', created_at)::date as date,
			COUNT(*) as count
		FROM datasets
		WHERE deleted_at IS NULL
			AND created_at > NOW() - INTERVAL '%d days'
		GROUP BY DATE_TRUNC('%s', created_at)
		ORDER BY date DESC
		LIMIT $1
	`, interval, limit*2, interval)

	var trend []analyticsDomain.TimeSeriesData
	err := r.db.SelectContext(ctx, &trend, query, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to get dataset trend: %w", err)
	}

	return trend, nil
}
