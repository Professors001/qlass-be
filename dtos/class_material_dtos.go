package dtos

import (
	"time"
)

type CreateClassMaterialDto struct {
	ClassID       uint       `json:"class_id" binding:"required"`
	Type          string     `json:"type" binding:"required,oneof=lecture assignment quiz"`
	Title         string     `json:"title" binding:"required,max=255"`
	Description   string     `json:"description"`
	AttachmentIds []uint     `json:"attachment_ids"`
	Points        *int       `json:"points" binding:"omitempty,min=0"`
	DueAt         *time.Time `json:"due_at" binding:"omitempty,gt"`
	Action        string     `json:"action" binding:"required,oneof=draft publish"`
}

type GetThumnailClassMaterialDto struct {
	ID        uint       `json:"id"`
	Title     string     `json:"title"`
	Type      string     `json:"type"`
	CreatedAt time.Time  `json:"created_at"`
	DueAt     *time.Time `json:"due_at"`
}

type GetClassMaterialDto struct {
	ID          uint                        `json:"id"`
	ClassID     uint                        `json:"class_id"`
	Type        string                      `json:"type"`
	Title       string                      `json:"title"`
	Description string                      `json:"description"`
	PublishedAt *time.Time                  `json:"published_at"`
	Attachments []*GetAttachmentResponseDto `json:"attachments"`
	CreatedAt   time.Time                   `json:"created_at"`
	Points      *int                        `json:"points"`
	DueAt       *time.Time                  `json:"due_at"`
}
