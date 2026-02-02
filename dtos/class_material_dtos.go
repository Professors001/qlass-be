package dtos

import "time"

type CreateClassMaterialDto struct {
	ClassID      uint       `json:"class_id" binding:"required"`
	Type         string     `json:"type" binding:"required,oneof=lecture assignment quiz"`
	Title        string     `json:"title" binding:"required,max=255"`
	Description  string     `json:"description"`
	AttachmentID []uint     `json:"attachment_ids"`
	Points       *int       `json:"points" binding:"omitempty,min=0"`
	DueAt        *time.Time `json:"due_at" binding:"omitempty,gt"`
	Action       string     `json:"action" binding:"required,oneof=draft publish"`
}
