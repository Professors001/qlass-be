package storage

import (
	"context"
	"mime/multipart"
	"net/url"
	"time"

	"github.com/minio/minio-go/v7"
)

type StorageService interface {
	Upload(ctx context.Context, file *multipart.FileHeader, bucketName string, objectName string) error
	GetPresignedURL(ctx context.Context, bucketName string, objectName string, expiry time.Duration) (string, error)
}

type minioStorageService struct {
	client *minio.Client
}

func NewMinioStorageService(client *minio.Client) StorageService {
	return &minioStorageService{
		client: client,
	}
}

func (s *minioStorageService) Upload(ctx context.Context, file *multipart.FileHeader, bucketName string, objectName string) error {
	src, err := file.Open()
	if err != nil {
		return err
	}
	defer src.Close()

	_, err = s.client.PutObject(ctx, bucketName, objectName, src, file.Size, minio.PutObjectOptions{
		ContentType: file.Header.Get("Content-Type"),
	})
	return err
}

func (s *minioStorageService) GetPresignedURL(ctx context.Context, bucketName string, objectName string, expiry time.Duration) (string, error) {
	reqParams := make(url.Values)
	presignedURL, err := s.client.PresignedGetObject(ctx, bucketName, objectName, expiry, reqParams)
	if err != nil {
		return "", err
	}
	return presignedURL.String(), nil
}
