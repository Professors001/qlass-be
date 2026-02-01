package usecases

import (
	"context"
	"fmt"
	"mime/multipart"
	"qlass-be/adapters/storage"
	"qlass-be/config"
	"qlass-be/domain/entities"
	"qlass-be/domain/repositories"
	"qlass-be/dtos"
	"qlass-be/transform"
	"time"
)

type AttachmentUseCase interface {
	UploadAttachment(userID uint, fileHeader *multipart.FileHeader) (*dtos.UploadAttachmentResponseDto, error)
	GetAttachmentByID(id uint) (*dtos.GetAttachmentResponseDto, error)
	GetAttachmentsByCourseMaterialID(courseMaterialID uint) ([]*entities.Attachment, error)
	GetAttachmentsBySubmissionID(submissionID uint) ([]*entities.Attachment, error)
	UpdateAttachment(dto *dtos.UpdateAttachmentDto) error
	DeleteAttachment(id uint) error
}

type attachmentUseCase struct {
	storageService storage.StorageService
	attachmentRepo repositories.AttachmentRepository
	userRepo       repositories.UserRepository
	cfg            *config.Config
}

func NewAttachmentUseCase(storageService storage.StorageService, attachmentRepo repositories.AttachmentRepository, userRepo repositories.UserRepository, cfg *config.Config) AttachmentUseCase {
	return &attachmentUseCase{
		storageService: storageService,
		attachmentRepo: attachmentRepo,
		userRepo:       userRepo,
		cfg:            cfg,
	}
}

func (u *attachmentUseCase) UploadAttachment(userID uint, file *multipart.FileHeader) (*dtos.UploadAttachmentResponseDto, error) {
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
		FileType:   file.Header.Get("Content-Type"),
		FileSize:   int(file.Size),
		UploaderID: userID,
	}

	err = u.attachmentRepo.Create(attachment)
	if err != nil {
		return nil, err
	}

	uploadResponse := transform.ToUploadAttachmentResponseDto(attachment, fileURL)

	return uploadResponse, nil
}

func (u *attachmentUseCase) GetAttachmentByID(id uint) (*dtos.GetAttachmentResponseDto, error) {
	attachment, err := u.attachmentRepo.GetByID(id)
	if err != nil {
		return nil, err
	}

	// Hydrate Uploader if not preloaded by repository
	if attachment.Uploader.ID == 0 {
		user, err := u.userRepo.GetByID(attachment.UploaderID)
		if err == nil {
			attachment.Uploader = *user
		}
	}

	bucketName := u.cfg.MinioBucketName
	fileURL, err := u.storageService.GetPresignedURL(context.Background(), bucketName, attachment.ObjectKey, time.Hour*1)
	if err != nil {
		return nil, err
	}

	getResponse := transform.ToGetAttachmentResponseDto(attachment, fileURL)

	return getResponse, nil
}

func (u *attachmentUseCase) GetAttachmentsByCourseMaterialID(courseMaterialID uint) ([]*entities.Attachment, error) {
	return u.attachmentRepo.GetByCourseMaterialID(courseMaterialID)
}

func (u *attachmentUseCase) GetAttachmentsBySubmissionID(submissionID uint) ([]*entities.Attachment, error) {
	return u.attachmentRepo.GetBySubmissionID(submissionID)
}
func (u *attachmentUseCase) UpdateAttachment(dto *dtos.UpdateAttachmentDto) error {
	attachment, err := u.attachmentRepo.GetByID(dto.AttachmentID)
	if err != nil {
		return err
	}

	switch dto.Type {
	case "course_material":
		attachment.CourseMaterialID = &dto.LinkID
		attachment.SubmissionID = nil
	case "submission":
		attachment.SubmissionID = &dto.LinkID
		attachment.CourseMaterialID = nil
	default:
		return fmt.Errorf("invalid attachment type: %s", dto.Type)
	}

	return u.attachmentRepo.Update(attachment)
}

func (u *attachmentUseCase) DeleteAttachment(id uint) error {
	return u.attachmentRepo.Delete(id)
}
