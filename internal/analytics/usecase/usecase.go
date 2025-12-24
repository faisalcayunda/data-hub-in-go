package usecase

import (
	"context"
	"fmt"

	"portal-data-backend/internal/analytics/domain"
)

type Usecase interface {
	GetDashboard(ctx context.Context) (*domain.DashboardResponse, error)
	GetDatasetStats(ctx context.Context) (*domain.DatasetStats, error)
	GetOrganizationStats(ctx context.Context) (*domain.OrganizationStats, error)
	GetUserStats(ctx context.Context) (*domain.UserStats, error)
	GetPopularDatasets(ctx context.Context, limit int) ([]domain.PopularDataset, error)
	GetPopularTags(ctx context.Context, limit int) ([]domain.TagStats, error)
	GetDatasetTrend(ctx context.Context, period string, limit int) ([]domain.TimeSeriesData, error)
}

type analyticsUsecase struct {
	repo domain.Repository
}

func NewAnalyticsUsecase(repo domain.Repository) Usecase {
	return &analyticsUsecase{
		repo: repo,
	}
}

func (u *analyticsUsecase) GetDashboard(ctx context.Context) (*domain.DashboardResponse, error) {
	// Get all stats in parallel for better performance
	type result struct {
		datasetStats     *domain.DatasetStats
		organizationStats *domain.OrganizationStats
		userStats        *domain.UserStats
		popularDatasets  []domain.PopularDataset
		popularTags      []domain.TagStats
		datasetTrend     []domain.TimeSeriesData
		err              error
	}

	resultChan := make(chan result, 1)

	go func() {
		var r result

		// Get basic stats
		r.datasetStats, r.err = u.repo.GetDatasetStats(ctx)
		if r.err != nil {
			resultChan <- r
			return
		}

		r.organizationStats, r.err = u.repo.GetOrganizationStats(ctx)
		if r.err != nil {
			resultChan <- r
			return
		}

		r.userStats, r.err = u.repo.GetUserStats(ctx)
		if r.err != nil {
			resultChan <- r
			return
		}

		// Get popular items (limit to 10)
		r.popularDatasets, r.err = u.repo.GetPopularDatasets(ctx, 10)
		if r.err != nil {
			resultChan <- r
			return
		}

		r.popularTags, r.err = u.repo.GetPopularTags(ctx, 10)
		if r.err != nil {
			resultChan <- r
			return
		}

		// Get dataset trend for last 30 days
		r.datasetTrend, r.err = u.repo.GetDatasetTrend(ctx, "daily", 30)
		if r.err != nil {
			resultChan <- r
			return
		}

		resultChan <- r
	}()

	r := <-resultChan
	if r.err != nil {
		return nil, fmt.Errorf("failed to get dashboard data: %w", r.err)
	}

	return &domain.DashboardResponse{
		DatasetStats:      r.datasetStats,
		OrganizationStats: r.organizationStats,
		UserStats:         r.userStats,
		PopularDatasets:   r.popularDatasets,
		PopularTags:       r.popularTags,
		DatasetTrend:      r.datasetTrend,
	}, nil
}

func (u *analyticsUsecase) GetDatasetStats(ctx context.Context) (*domain.DatasetStats, error) {
	stats, err := u.repo.GetDatasetStats(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get dataset stats: %w", err)
	}
	return stats, nil
}

func (u *analyticsUsecase) GetOrganizationStats(ctx context.Context) (*domain.OrganizationStats, error) {
	stats, err := u.repo.GetOrganizationStats(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get organization stats: %w", err)
	}
	return stats, nil
}

func (u *analyticsUsecase) GetUserStats(ctx context.Context) (*domain.UserStats, error) {
	stats, err := u.repo.GetUserStats(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get user stats: %w", err)
	}
	return stats, nil
}

func (u *analyticsUsecase) GetPopularDatasets(ctx context.Context, limit int) ([]domain.PopularDataset, error) {
	if limit < 1 {
		limit = 10
	}
	if limit > 100 {
		limit = 100
	}

	datasets, err := u.repo.GetPopularDatasets(ctx, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to get popular datasets: %w", err)
	}
	return datasets, nil
}

func (u *analyticsUsecase) GetPopularTags(ctx context.Context, limit int) ([]domain.TagStats, error) {
	if limit < 1 {
		limit = 10
	}
	if limit > 100 {
		limit = 100
	}

	tags, err := u.repo.GetPopularTags(ctx, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to get popular tags: %w", err)
	}
	return tags, nil
}

func (u *analyticsUsecase) GetDatasetTrend(ctx context.Context, period string, limit int) ([]domain.TimeSeriesData, error) {
	if period == "" {
		period = "daily"
	}
	if limit < 1 {
		limit = 30
	}
	if limit > 365 {
		limit = 365
	}

	trend, err := u.repo.GetDatasetTrend(ctx, period, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to get dataset trend: %w", err)
	}
	return trend, nil
}
