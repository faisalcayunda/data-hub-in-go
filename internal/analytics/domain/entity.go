package domain

import "time"

// DatasetStats represents dataset statistics
type DatasetStats struct {
	TotalDatasets   int64     `json:"total_datasets"`
	PublishedCount  int64     `json:"published_count"`
	DraftCount      int64     `json:"draft_count"`
	ArchivedCount   int64     `json:"archived_count"`
	TotalDownloads  int64     `json:"total_downloads"`
	TotalViews      int64     `json:"total_views"`
	LastUpdated     time.Time `json:"last_updated"`
}

// OrganizationStats represents organization statistics
type OrganizationStats struct {
	TotalOrganizations int64     `json:"total_organizations"`
	ActiveOrganizations int64    `json:"active_organizations"`
	TotalDatasets      int64     `json:"total_datasets"`
	LastUpdated        time.Time `json:"last_updated"`
}

// UserStats represents user statistics
type UserStats struct {
	TotalUsers      int64     `json:"total_users"`
	ActiveUsers     int64     `json:"active_users"`
	NewUsersThisMonth int64   `json:"new_users_this_month"`
	LastUpdated     time.Time `json:"last_updated"`
}

// PopularDataset represents a popular dataset
type PopularDataset struct {
	ID          string `json:"id"`
	Title       string `json:"title"`
	Organization string `json:"organization"`
	Views       int64  `json:"views"`
	Downloads   int64  `json:"downloads"`
}

// TagStats represents tag statistics
type TagStats struct {
	TagID      string `json:"tag_id"`
	Name       string `json:"name"`
	DatasetCount int64 `json:"dataset_count"`
}

// TimeSeriesData represents time series data point
type TimeSeriesData struct {
	Date  string `json:"date"`
	Count int64  `json:"count"`
}

// DashboardResponse represents the complete dashboard data
type DashboardResponse struct {
	DatasetStats     *DatasetStats       `json:"dataset_stats"`
	OrganizationStats *OrganizationStats `json:"organization_stats"`
	UserStats        *UserStats          `json:"user_stats"`
	PopularDatasets  []PopularDataset    `json:"popular_datasets"`
	PopularTags      []TagStats          `json:"popular_tags"`
	DatasetTrend     []TimeSeriesData    `json:"dataset_trend"`
}

// GetStatsRequest represents query parameters for stats
type GetStatsRequest struct {
	StartDate *string `json:"start_date,omitempty"`
	EndDate   *string `json:"end_date,omitempty"`
	Period    string  `json:"period,omitempty"` // daily, weekly, monthly
}
