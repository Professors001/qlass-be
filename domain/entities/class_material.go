package entities

import (
	"time"

	"gorm.io/gorm"
)

type ClassMaterial struct {
	gorm.Model
	ClassID             uint         `json:"class_id" gorm:"not null;index"`
	Class               Class        `json:"class" gorm:"foreignKey:ClassID"`
	Type                string       `json:"type" gorm:"type:varchar(50);comment:lecture, assignment, quiz"`
	Title               string       `json:"title" gorm:"not null;type:varchar(255)"`
	Description         string       `json:"description" gorm:"type:text;comment:HTML or Markdown content"`
	DueAt               *time.Time   `json:"due_at" gorm:"comment:Null for lectures."`
	Points              *int         `json:"points" gorm:"default:100;comment:Null for lectures."`
	IsPublished         bool         `json:"is_published" gorm:"default:false"`
	PublishedAt         *time.Time   `json:"published_at"`
	AllowLateSubmission bool         `json:"allow_late_submission" gorm:"default:true"`
	QuizPin             *string      `json:"quiz_pin" gorm:"type:varchar(20)"`
	QuizStatus          string       `json:"quiz_status" gorm:"default:idle;type:varchar(50)"`
	Attachments         []Attachment `json:"attachments" gorm:"polymorphic:Owner;"`
}
