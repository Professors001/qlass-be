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
	"qlass-be/utils"
	"time"
)

// Define constants to avoid typos in magic strings
const (
	OwnerTypeClassMaterial = "class_materials"
	OwnerTypeSubmission    = "submissions"
	OwnerTypeQuizQuestion  = "quiz_questions"
)

type AttachmentUseCase interface {
	UploadAttachment(userID uint, fileHeader *multipart.FileHeader) (*dtos.UploadAttachmentResponseDto, error)
	GetAttachmentByID(id uint) (*dtos.GetAttachmentResponseDto, error)
	GetAttachmentsByOwner(ownerType string, ownerID uint) ([]*dtos.GetAttachmentResponseDto, error)
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

	fileURL, err := u.storageService.GetPresignedURL(context.Background(), bucketName, objectName, time.Hour*1)
	if err != nil {
		return nil, err
	}

	attachment := &entities.Attachment{
		Filename:   file.Filename,
		ObjectKey:  objectName,
		FileType:   file.Header.Get("Content-Type"),
		FileSize:   int(file.Size),
		UploaderID: userID,
		// OwnerID/Type are 0/"" initially (unlinked) until UpdateAttachment is called
	}

	err = u.attachmentRepo.Create(attachment)
	if err != nil {
		return nil, err
	}

	return transforms.ToUploadAttachmentResponseDto(attachment, fileURL), nil
}

func (u *attachmentUseCase) GetAttachmentByID(id uint) (*dtos.GetAttachmentResponseDto, error) {
	attachment, err := u.attachmentRepo.GetByID(id)
	if err != nil {
		return nil, err
	}
	return u.enrichAttachment(attachment)
}

// New Helper Method to handle the loop
func (u *attachmentUseCase) GetAttachmentsByOwner(ownerType string, ownerID uint) ([]*dtos.GetAttachmentResponseDto, error) {
	// Assuming your Repo has been updated to support GetByOwner(id, type)
	attachments, err := u.attachmentRepo.GetByOwnerTypeAndOwnerID(ownerType, ownerID)
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

	// Map the incoming "type" string to our database OwnerType
	switch dto.Type {
	case "class_material":
		attachment.OwnerID = &dto.LinkID
		attachment.OwnerType = utils.Ptr("class_material")
	case "submission":
		attachment.OwnerID = &dto.LinkID
		attachment.OwnerType = utils.Ptr("submission")
	case "quiz_question":
		attachment.OwnerID = &dto.LinkID
		attachment.OwnerType = utils.Ptr("quiz_question")
	default:
		return fmt.Errorf("invalid attachment type: %s", dto.Type)
	}

	return u.attachmentRepo.Update(attachment)
}

func (u *attachmentUseCase) DeleteAttachment(id uint) error {
	return u.attachmentRepo.Delete(id)
}

func (u *attachmentUseCase) enrichAttachment(attachment *entities.Attachment) (*dtos.GetAttachmentResponseDto, error) {
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

	return transforms.ToGetAttachmentResponseDto(attachment, fileURL), nil
}
