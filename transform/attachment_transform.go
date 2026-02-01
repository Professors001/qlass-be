package transform

import (
	"qlass-be/domain/entities"
	"qlass-be/dtos"
)

func ToGetAttachmentResponseDto(attachment *entities.Attachment, fileURL string) *dtos.GetAttachmentResponseDto {
	return &dtos.GetAttachmentResponseDto{
		AttachmentID:     attachment.ID,
		FileURL:          fileURL,
		Filename:         attachment.Filename,
		FileSize:         attachment.FileSize,
		FileType:         attachment.FileType,
		UploaderID:       attachment.UploaderID,
		UploaderName:     attachment.Uploader.FirstName + " " + attachment.Uploader.LastName,
		UploaderRole:     attachment.Uploader.Role,
		CourseMaterialID: attachment.CourseMaterialID,
		SubmissionID:     attachment.SubmissionID,
		UploadedAt:       attachment.CreatedAt.Format("2006-01-02 15:04:05"),
	}
}

func ToUploadAttachmentResponseDto(attachment *entities.Attachment, fileURL string) *dtos.UploadAttachmentResponseDto {
	return &dtos.UploadAttachmentResponseDto{
		AttachmentID: attachment.ID,
		FileURL:      fileURL,
	}
}
