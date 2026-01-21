package usecases

import (
	"context"
	"fmt"
	"mime/multipart"
	"net/url"
	"time"

	"github.com/minio/minio-go/v7"
)

type FileUseCase interface {
	UploadFile(ctx context.Context, file *multipart.FileHeader, bucketName string) (string, error)
	GetFileUrl(ctx context.Context, objectName string, bucketName string) (string, error)
}

type fileUseCase struct {
	minioClient *minio.Client
}

func NewFileUseCase(minioClient *minio.Client) FileUseCase {
	return &fileUseCase{
		minioClient: minioClient,
	}
}

func (u *fileUseCase) UploadFile(ctx context.Context, file *multipart.FileHeader, bucketName string) (string, error) {
	src, err := file.Open()
	if err != nil {
		return "", err
	}
	defer src.Close()

	// Generate unique object name
	objectName := fmt.Sprintf("%d_%s", time.Now().UnixNano(), file.Filename)

	// Upload to MinIO
	_, err = u.minioClient.PutObject(ctx, bucketName, objectName, src, file.Size, minio.PutObjectOptions{
		ContentType: file.Header.Get("Content-Type"),
	})
	if err != nil {
		return "", err
	}

	return objectName, nil
}

func (u *fileUseCase) GetFileUrl(ctx context.Context, objectName string, bucketName string) (string, error) {
	// Set expiry for the presigned URL (e.g., 1 hour)
	expiry := time.Hour * 1

	reqParams := make(url.Values)
	// Optional: Force download
	// reqParams.Set("response-content-disposition", "attachment; filename=\""+objectName+"\"")

	presignedURL, err := u.minioClient.PresignedGetObject(ctx, bucketName, objectName, expiry, reqParams)
	if err != nil {
		return "", err
	}

	return presignedURL.String(), nil
}
