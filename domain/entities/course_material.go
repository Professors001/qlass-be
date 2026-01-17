package entities

import (
	"time"

	"gorm.io/gorm"
)

type CourseMaterial struct {
	gorm.Model
	ClassID             uint       `json:"class_id" gorm:"not null"`
	Class               Class      `json:"class" gorm:"foreignKey:ClassID"`
	Type                string     `json:"type" gorm:"type:varchar(50);comment:lecture, assignment, quiz"`
	Title               string     `json:"title" gorm:"not null;type:varchar(255)"`
	Description         string     `json:"description" gorm:"type:text;comment:HTML or Markdown content"`
	DueAt               *time.Time `json:"due_at" gorm:"comment:Null for lectures. Required for Assignments."`
	Points              *int       `json:"points" gorm:"default:100;comment:Null for lectures."`
	IsPublished         bool       `json:"is_published" gorm:"default:false;comment:Draft mode vs Published"`
	AllowLateSubmission bool       `json:"allow_late_submission" gorm:"default:true"`
	QuizID              *uint      `json:"quiz_id"`
	Quiz                *Quiz      `json:"quiz" gorm:"foreignKey:QuizID"`
	QuizPin             *string    `json:"quiz_pin" gorm:"type:varchar(20);comment:Null if no pin is set"`
	QuizStatus          string     `json:"quiz_status" gorm:"default:idle;type:varchar(50);comment:idle, waiting_lobby, in_progress, finished"`
	TimeLimitSeconds    *int       `json:"time_limit_seconds" gorm:"comment:For the whole quiz or null"`
}
