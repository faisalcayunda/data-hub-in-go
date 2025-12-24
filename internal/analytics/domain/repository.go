package domain

import (
	"context"
)

type Repository interface {
	GetDatasetStats(ctx context.Context) (*DatasetStats, error)
	GetOrganizationStats(ctx context.Context) (*OrganizationStats, error)
	GetUserStats(ctx context.Context) (*UserStats, error)
	GetPopularDatasets(ctx context.Context, limit int) ([]PopularDataset, error)
	GetPopularTags(ctx context.Context, limit int) ([]TagStats, error)
	GetDatasetTrend(ctx context.Context, period string, limit int) ([]TimeSeriesData, error)
}
