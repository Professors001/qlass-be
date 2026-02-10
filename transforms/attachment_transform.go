package transforms

import (
    "qlass-be/domain/entities"
    "qlass-be/dtos"
)

func ToGetAttachmentResponseDto(attachment *entities.Attachment, fileURL string) *dtos.GetAttachmentResponseDto {
    return &dtos.GetAttachmentResponseDto{
        ID:           attachment.ID,
        FileURL:      fileURL,
        Filename:     attachment.Filename,
        FileSize:     attachment.FileSize,
        FileType:     attachment.FileType,
        UploaderID:   attachment.UploaderID,
        // Safe check for Uploader to prevent panic if it wasn't preloaded
        UploaderName: attachment.Uploader.FirstName + " " + attachment.Uploader.LastName,
        UploadedAt:   attachment.CreatedAt.Format("2006-01-02 15:04:05"),
        
        // Direct mapping
        OwnerID:   attachment.OwnerID,
        OwnerType: attachment.OwnerType,
    }
}

func ToUploadAttachmentResponseDto(attachment *entities.Attachment, fileURL string) *dtos.UploadAttachmentResponseDto {
    return &dtos.UploadAttachmentResponseDto{
        AttachmentID: attachment.ID,
        FileURL:      fileURL,
    }
}