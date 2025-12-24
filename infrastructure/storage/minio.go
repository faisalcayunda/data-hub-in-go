package storage

import (
	"context"
	"fmt"
	"io"
	"time"

	"portal-data-backend/internal/file/domain"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

type minioStorage struct {
	client *minio.Client
	bucket string
}

func NewMinIOStorage(endpoint, accessKey, secretKey, bucket string, useSSL bool) (domain.StorageService, error) {
	client, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKey, secretKey, ""),
		Secure: useSSL,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create minio client: %w", err)
	}

	// Create bucket if it doesn't exist
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	exists, err := client.BucketExists(ctx, bucket)
	if err != nil {
		return nil, fmt.Errorf("failed to check bucket existence: %w", err)
	}

	if !exists {
		err = client.MakeBucket(ctx, bucket, minio.MakeBucketOptions{})
		if err != nil {
			return nil, fmt.Errorf("failed to create bucket: %w", err)
		}
	}

	return &minioStorage{
		client: client,
		bucket: bucket,
	}, nil
}

func (s *minioStorage) Upload(ctx context.Context, fileName string, reader io.Reader, contentType string, path string) (string, error) {
	_, err := s.client.PutObject(ctx, s.bucket, path, reader, -1, minio.PutObjectOptions{
		ContentType: contentType,
	})
	if err != nil {
		return "", fmt.Errorf("failed to upload file: %w", err)
	}
	return path, nil
}

func (s *minioStorage) Delete(ctx context.Context, path string) error {
	err := s.client.RemoveObject(ctx, s.bucket, path, minio.RemoveObjectOptions{})
	if err != nil {
		return fmt.Errorf("failed to delete file: %w", err)
	}
	return nil
}

func (s *minioStorage) GetURL(ctx context.Context, path string) (string, error) {
	presignedURL, err := s.client.PresignedGetObject(ctx, s.bucket, path, time.Hour*24, nil)
	if err != nil {
		return "", fmt.Errorf("failed to get presigned URL: %w", err)
	}
	return presignedURL.String(), nil
}
