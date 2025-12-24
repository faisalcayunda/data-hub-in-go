package domain

import "time"

// DataRow represents a single row of data in a dataset
type DataRow struct {
	ID        string      `db:"id" json:"id"`
	DatasetID string      `db:"dataset_id" json:"dataset_id"`
	RowIndex  int         `db:"row_index" json:"row_index"`
	Data      string      `db:"data" json:"data"` // JSON data
	CreatedBy string      `db:"created_by" json:"created_by"`
	CreatedAt time.Time   `db:"created_at" json:"created_at"`
	UpdatedAt time.Time   `db:"updated_at" json:"updated_at"`
	DeletedAt *time.Time  `db:"deleted_at" json:"deleted_at,omitempty"`
}

// ListDataRowsRequest represents list data rows input
type ListDataRowsRequest struct {
	Page      int    `json:"page" validate:"min=1"`
	Limit     int    `json:"limit" validate:"min=1,max=1000"`
	DatasetID string `json:"dataset_id" validate:"required"`
	Search    string `json:"search,omitempty"`
}

// CreateDataRowRequest represents create data row input
type CreateDataRowRequest struct {
	DatasetID string `json:"dataset_id" validate:"required"`
	RowIndex  int    `json:"row_index" validate:"required,min=0"`
	Data      string `json:"data" validate:"required"`
}

// BulkCreateDataRowsRequest represents bulk create data rows input
type BulkCreateDataRowsRequest struct {
	DatasetID string             `json:"dataset_id" validate:"required"`
	Rows      []DataRowDataInput `json:"rows" validate:"required,min=1,dive"`
}

// DataRowDataInput represents a single data row input
type DataRowDataInput struct {
	RowIndex int    `json:"row_index" validate:"required,min=0"`
	Data     string `json:"data" validate:"required"`
}

// UpdateDataRowRequest represents update data row input
type UpdateDataRowRequest struct {
	RowIndex *int    `json:"row_index,omitempty"`
	Data     *string `json:"data,omitempty"`
}

// DataRowInfo represents data row information for API responses
type DataRowInfo struct {
	ID        string    `json:"id"`
	DatasetID string    `json:"dataset_id"`
	RowIndex  int       `json:"row_index"`
	Data      string    `json:"data"`
	CreatedBy string    `json:"created_by"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// DataRowListResponse represents paginated data row list
type DataRowListResponse struct {
	Rows []DataRowInfo `json:"rows"`
	Meta ListMeta      `json:"meta"`
}

// ListMeta represents pagination metadata
type ListMeta struct {
	Page      int `json:"page"`
	Limit     int `json:"limit"`
	Total     int `json:"total"`
	TotalPage int `json:"total_page"`
}

// DataRowStats represents data row statistics
type DataRowStats struct {
	TotalRows   int64     `json:"total_rows"`
	LastUpdated time.Time `json:"last_updated"`
}
