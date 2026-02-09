package dtos

type UploadAttachmentResponseDto struct {
	AttachmentID uint   `json:"attachment_id"`
	FileURL      string `json:"file_url"`
}

type GetAttachmentResponseDto struct {
	ID           uint   `json:"id"`
	FileURL      string `json:"file_url"`
	Filename     string `json:"filename"`
	FileSize     int    `json:"file_size"`
	FileType     string `json:"file_type"`
	UploaderID   uint   `json:"uploader_id"`
	UploaderName string `json:"uploader_name"`
	UploadedAt   string `json:"uploaded_at"`

	// Polymorphic Fields (Replaces ClassMaterialID, SubmissionID)
	OwnerID   *uint   `json:"owner_id"`
	OwnerType *string `json:"owner_type"` // e.g. "class_materials", "submissions"
}

type UpdateAttachmentDto struct {
	AttachmentID uint   `json:"attachment_id" binding:"required"`
	LinkID       uint   `json:"link_id" binding:"required"`
	Type         string `json:"type" binding:"required,oneof=class_material submission quiz_question"`
}
