package dtos

import "time"

type CreateSubmissionDto struct {
	ClassMaterialID uint   `json:"class_material_id" binding:"required"`
	StudentComment  string `json:"student_comment"`
	AttachmentIds   []uint `json:"attachment_ids"`
	Action          string `json:"action" binding:"required,oneof=DRAFT SUBMIT UNSUBMIT"`
}

type GetSubmissionResponseDto struct {
	ID              uint                        `json:"id"`
	ClassMaterialID uint                        `json:"class_material_id"`
	UserID          uint                        `json:"student_id"`
	StudentComment  string                      `json:"student_comment"`
	Status          string                      `json:"status"`
	Score           *int                        `json:"score"`
	TeacherFeedback string                      `json:"teacher_feedback"`
	CreatedAt       time.Time                   `json:"created_at"`
	UpdatedAt       time.Time                   `json:"updated_at"`
	Attachments     []*GetAttachmentResponseDto `json:"attachments"`
}
