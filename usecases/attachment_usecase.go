package usecases

import (
	"context"
	"fmt"
	"mime/multipart"
	"qlass-be/config"
	"qlass-be/domain/entities"
	"qlass-be/domain/repositories"
	"qlass-be/adapters/storage"
	"time"
)

type AttachmentUseCase interface {
	UploadAttachment(userID uint, fileHeader *multipart.FileHeader) (*entities.Attachment, error)
	// GetAttachmentByID(id uint) (*entities.Attachment, error)
	// GetAttachmentsByCourseMaterialID(courseMaterialID uint) ([]*entities.Attachment, error)
	// GetAttachmentsBySubmissionID(submissionID uint) ([]*entities.Attachment, error)
	// UpdateAttachment(dto *dtos.UpdateAttachmentDto, claims *middleware.JWTCustomClaims) (*entities.Attachment, error)
	// DeleteAttachment(id uint) error
}

type attachmentUseCase struct {
	storageService storage.StorageService
	attachmentRepo repositories.AttachmentRepository
	cfg            *config.Config
}

func NewAttachmentUseCase(storageService storage.StorageService, attachmentRepo repositories.AttachmentRepository, cfg *config.Config) AttachmentUseCase {
	return &attachmentUseCase{
		storageService: storageService,
		attachmentRepo: attachmentRepo,
		cfg:            cfg,
	}
}

func (u *attachmentUseCase) UploadAttachment(userID uint, file *multipart.FileHeader) (*entities.Attachment, error) {
	bucketName := u.cfg.MinioBucketName

	objectName := fmt.Sprintf("%d_%s", time.Now().UnixNano(), file.Filename)

	err := u.storageService.Upload(context.Background(), file, bucketName, objectName)
	if err != nil {
		return nil, err
	}

	// Generate URL for immediate response (optional), but store ObjectKey in DB
	fileURL, err := u.storageService.GetPresignedURL(context.Background(), bucketName, objectName, time.Hour*1)
	if err != nil {
		return nil, err
	}

	attachment := &entities.Attachment{
		Filename:   file.Filename,
		ObjectKey:  objectName, // Store the key, not the expiring URL
		FileURL:    fileURL,    // Populate transient field for response
		FileType:   file.Header.Get("Content-Type"),
		FileSize:   int(file.Size),
		UploaderID: userID,
	}

	err = u.attachmentRepo.Create(attachment)
	if err != nil {
		return nil, err
	}

	return attachment, nil
}
