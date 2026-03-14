package dtos

import (
	"time"

	"gorm.io/datatypes"
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

type CreateQuizClassMaterialDto struct {
	ClassID     uint       `json:"class_id" binding:"required"`
	Type        string     `json:"type" binding:"required,oneof=lecture assignment quiz"`
	Title       string     `json:"title" binding:"required,max=255"`
	Description string     `json:"description"`
	QuizID      uint       `json:"quiz_id"`
	Points      *int       `json:"points" binding:"omitempty,min=0"`
	DueAt       *time.Time `json:"due_at" binding:"omitempty,gt"`
	Action      string     `json:"action" binding:"required,oneof=draft publish"`
}

type GetThumnailClassMaterialDto struct {
	ID          uint       `json:"id"`
	ClassID     uint       `json:"class_id"`
	Type        string     `json:"type"`
	Title       string     `json:"title"`
	PublishedAt *time.Time `json:"published_at"`
	CreatedAt   time.Time  `json:"created_at"`
	Points      *int       `json:"points"`
	DueAt       *time.Time `json:"due_at"`
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
	QuizGameLog *QuizGameLogDto             `json:"quiz_game_log,omitempty"`
	CreatedBy   CreatedByDto                `json:"created_by"`
}

type CreatedByDto struct {
	ID       uint   `json:"id"`
	FullName string `json:"full_name"`
	ImgURL   string `json:"img_url,omitempty"`
}

type QuizGameLogDto struct {
	ID              uint           `json:"id"`
	ClassMaterialID uint           `json:"class_material_id"`
	QuizPin         string         `json:"quiz_pin"`
	Status          string         `json:"status"`
	StartedAt       *time.Time     `json:"started_at"`
	FinishedAt      *time.Time     `json:"finished_at"`
	QuizSnapshot    datatypes.JSON `json:"quiz_snapshot"`
}

type UpdatePostClassMaterialDto struct {
	ClassMaterialID uint   `json:"class_material_id" binding:"required"`
	Title           string `json:"title"`
	Description     string `json:"description"`
	AttachmentIds   []uint `json:"attachment_ids"`
	Published       bool   `json:"published"`
}

type UpdateAssignmentClassMaterialDto struct {
	ClassMaterialID uint       `json:"class_material_id" binding:"required"`
	Title           string     `json:"title"`
	Description     string     `json:"description"`
	AttachmentIds   []uint     `json:"attachment_ids"`
	Published       bool       `json:"published"`
	Points          *int       `json:"points" binding:"omitempty,min=0"`
	DueAt           *time.Time `json:"due_at" binding:"omitempty,gt"`
}

type UpdateQuizClassMaterialDto struct {
	ClassMaterialID uint   `json:"class_material_id" binding:"required"`
	Title           string `json:"title"`
	Description     string `json:"description"`
	Published       bool   `json:"published"`
	Points          *int   `json:"points" binding:"omitempty,min=0"`
}
