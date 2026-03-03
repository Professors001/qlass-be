package entities

import (
	"time"

	"gorm.io/gorm"
)

type Submission struct {
	gorm.Model
	ClassMaterialID uint          `json:"class_material_id" gorm:"not null"`
	ClassMaterial   ClassMaterial `json:"class_material" gorm:"foreignKey:ClassMaterialID"`
	UserID          uint          `json:"user_id"`
	User            User          `json:"user" gorm:"foreignKey:UserID"`
	StudentComment  string        `json:"student_comment" gorm:"type:text"`
	Score           *int          `json:"score" gorm:"comment:e.g. 85/100"`
	SubmittedAt     *time.Time    `json:"submitted_at"`
	IsLate          bool          `json:"is_late" gorm:"default:false"`
	TeacherFeedback string        `json:"teacher_feedback" gorm:"type:text"`
	Status          string        `json:"status" gorm:"default:submitted;type:varchar(50);comment:submit, graded, return, late, draft"`
	Attachments     []Attachment  `json:"attachments" gorm:"polymorphic:Owner;"`
}
