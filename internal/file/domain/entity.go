package domain

import (
	"time"
)

// File represents uploaded file metadata
type File struct {
	ID            string       `db:"id" json:"id"`
	Name          string       `db:"name" json:"name"`
	OriginalName  string       `db:"original_name" json:"original_name"`
	Extension     string       `db:"extension" json:"extension"`
	Size          int64        `db:"size" json:"size"`
	MimeType      string       `db:"mime_type" json:"mime_type"`
	Path          string       `db:"path" json:"path"`
	StoragePath   string       `db:"storage_path" json:"storage_path"`
	StorageType   StorageType `db:"storage_type" json:"storage_type"`
	DatasetID     *string      `db:"dataset_id" json:"dataset_id,omitempty"`
	UploadedBy    string       `db:"uploaded_by" json:"uploaded_by"`
	Status        FileStatus   `db:"status" json:"status"`
	CreatedAt     time.Time    `db:"created_at" json:"created_at"`
	UpdatedAt     time.Time    `db:"updated_at" json:"updated_at"`
}

// StorageType represents where file is stored
type StorageType string

const (
	StorageTypeLocal    StorageType = "local"
	StorageTypeS3       StorageType = "s3"
	StorageTypeMinIO    StorageType = "minio"
)

// FileStatus represents file status
type FileStatus string

const (
	FileStatusUploading   FileStatus = "uploading"
	FileStatusProcessing FileStatus = "processing"
	FileStatusReady       FileStatus = "ready"
	FileStatusFailed      FileStatus = "failed"
	FileStatusDeleted     FileStatus = "deleted"
)

// UploadResponse represents file upload response
type UploadResponse struct {
	ID         string `json:"id"`
	Name       string `json:"name"`
	Size       int64  `json:"size"`
	MimeType   string `json:"mime_type"`
	Path       string `json:"path"`
	Status     string `json:"status"`
	CreatedAt  string `json:"created_at"`
}

// FileInfo represents file information
type FileInfo struct {
	ID           string    `json:"id"`
	Name         string    `json:"name"`
	OriginalName string    `json:"original_name"`
	Extension    string    `json:"extension"`
	Size         int64     `json:"size"`
	MimeType     string    `json:"mime_type"`
	Path         string    `json:"path"`
	DatasetID    *string   `json:"dataset_id,omitempty"`
	UploadedBy   string    `json:"uploaded_by"`
	Status       string    `json:"status"`
	CreatedAt    time.Time `json:"created_at"`
}

// ListFilesRequest represents list files input
type ListFilesRequest struct {
	Page      int    `json:"page" validate:"min=1"`
	Limit     int    `json:"limit" validate:"min=1,max=100"`
	DatasetID *string `json:"dataset_id,omitempty"`
	Status    *string `json:"status,omitempty"`
	Search    string `json:"search,omitempty"`
}

// FileListResponse represents paginated file list
type FileListResponse struct {
	Files []FileInfo `json:"files"`
	Meta  ListMeta    `json:"meta"`
}

// ListMeta represents pagination metadata
type ListMeta struct {
	Page      int `json:"page"`
	Limit     int `json:"limit"`
	Total     int `json:"total"`
	TotalPage int `json:"total_page"`
}
