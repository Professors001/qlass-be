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
	"qlass-be/transforms"
	"time"
)

type AttachmentUseCase interface {
	UploadAttachment(userID uint, fileHeader *multipart.FileHeader) (*dtos.UploadAttachmentResponseDto, error)
	GetAttachmentByID(id uint) (*dtos.GetAttachmentResponseDto, error)
	GetAttachmentsByClassMaterialID(classMaterialID uint) ([]*dtos.GetAttachmentResponseDto, error)
	GetAttachmentsBySubmissionID(submissionID uint) ([]*dtos.GetAttachmentResponseDto, error)
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

	uploadResponse := transforms.ToUploadAttachmentResponseDto(attachment, fileURL)

	return uploadResponse, nil
}

func (u *attachmentUseCase) GetAttachmentByID(id uint) (*dtos.GetAttachmentResponseDto, error) {
	attachment, err := u.attachmentRepo.GetByID(id)
	if err != nil {
		return nil, err
	}

	return u.enrichAttachment(attachment)
}

func (u *attachmentUseCase) GetAttachmentsByClassMaterialID(classMaterialID uint) ([]*dtos.GetAttachmentResponseDto, error) {
	// 1. Get raw entities from repo
	attachments, err := u.attachmentRepo.GetByClassMaterialID(classMaterialID)
	if err != nil {
		return nil, err
	}

	// 2. Transforms and enrich each attachment
	response := make([]*dtos.GetAttachmentResponseDto, 0, len(attachments))

	for _, att := range attachments {
		dto, err := u.enrichAttachment(att)
		if err != nil {
			// Option A: Log error and continue (skip broken files)
			// Option B: Return error immediately (fail whole request)
			// Here we choose Option A to ensure the UI gets at least the working files
			continue
		}
		response = append(response, dto)
	}

	return response, nil
}

func (u *attachmentUseCase) GetAttachmentsBySubmissionID(submissionID uint) ([]*dtos.GetAttachmentResponseDto, error) {
	attachments, err := u.attachmentRepo.GetBySubmissionID(submissionID)
	if err != nil {
		return nil, err
	}

	response := make([]*dtos.GetAttachmentResponseDto, 0, len(attachments))

	for _, att := range attachments {
		dto, err := u.enrichAttachment(att)
		if err != nil {
			continue
		}
		response = append(response, dto)
	}

	return response, nil
}
func (u *attachmentUseCase) UpdateAttachment(dto *dtos.UpdateAttachmentDto) error {
	attachment, err := u.attachmentRepo.GetByID(dto.AttachmentID)
	if err != nil {
		return err
	}

	switch dto.Type {
	case "class_material":
		attachment.ClassMaterialID = &dto.LinkID
		attachment.SubmissionID = nil
	case "submission":
		attachment.SubmissionID = &dto.LinkID
		attachment.ClassMaterialID = nil
	default:
		return fmt.Errorf("invalid attachment type: %s", dto.Type)
	}

	return u.attachmentRepo.Update(attachment)
}

func (u *attachmentUseCase) DeleteAttachment(id uint) error {
	return u.attachmentRepo.Delete(id)
}

func (u *attachmentUseCase) enrichAttachment(attachment *entities.Attachment) (*dtos.GetAttachmentResponseDto, error) {
	// 1. Hydrate Uploader if missing
	if attachment.Uploader.ID == 0 {
		user, err := u.userRepo.GetByID(attachment.UploaderID)
		if err == nil {
			attachment.Uploader = *user
		}
	}

	// 2. Generate Presigned URL
	bucketName := u.cfg.MinioBucketName
	fileURL, err := u.storageService.GetPresignedURL(context.Background(), bucketName, attachment.ObjectKey, time.Hour*1)
	if err != nil {
		return nil, err
	}

	// 3. Transforms to DTO
	return transforms.ToGetAttachmentResponseDto(attachment, fileURL), nil
}
