package usecases

import (
	"context"
	"fmt"
	"mime/multipart"
	"qlass-be/infrastructure/storage"
	"time"
)

type FileUseCase interface {
	UploadFile(ctx context.Context, file *multipart.FileHeader, bucketName string) (string, error)
	GetFileUrl(ctx context.Context, objectName string, bucketName string) (string, error)
}

type fileUseCase struct {
	storageService storage.StorageService
}

func NewFileUseCase(storageService storage.StorageService) FileUseCase {
	return &fileUseCase{
		storageService: storageService,
	}
}

func (u *fileUseCase) UploadFile(ctx context.Context, file *multipart.FileHeader, bucketName string) (string, error) {
	// Generate unique object name
	objectName := fmt.Sprintf("%d_%s", time.Now().UnixNano(), file.Filename)

	// Upload using storage service
	err := u.storageService.Upload(ctx, file, bucketName, objectName)
	if err != nil {
		return "", err
	}

	return objectName, nil
}

func (u *fileUseCase) GetFileUrl(ctx context.Context, objectName string, bucketName string) (string, error) {
	// Set expiry for the presigned URL (e.g., 1 hour)
	expiry := time.Hour * 1

	return u.storageService.GetPresignedURL(ctx, bucketName, objectName, expiry)
}
