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

type GetAttachmentResponseDto struct {
	ID              uint   `json:"id"`
	FileURL         string `json:"file_url"`
	Filename        string `json:"filename"`
	FileSize        int    `json:"file_size"`
	FileType        string `json:"file_type"`
	UploaderID      uint   `json:"uploader_id"`
	UploaderName    string `json:"uploader_name"`
	ClassMaterialID *uint  `json:"class_material_id,omitempty"`
	SubmissionID    *uint  `json:"submission_id,omitempty"`
	UploadedAt      string `json:"uploaded_at"`
}
