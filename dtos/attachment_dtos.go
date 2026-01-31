package dtos

type UpdateAttachmentDto struct {
	AttachmentID uint   `json:"attachment_id" binding:"required"`
	Type         string `json:"type" binding:"required,oneof=course_material submission"`
	LinkID       uint   `json:"link_id" binding:"required"`
}

type UploadAttachmentResponseDto struct {
	AttachmentID uint   `json:"attachment_id"`
	FileURL      string `json:"file_url"`
}
