package entities

import (
	"time"

	"gorm.io/gorm"
)

type QuizGameLog struct {
	gorm.Model
	CourseMaterialID uint           `json:"course_material_id" gorm:"not null"`
	CourseMaterial   CourseMaterial `json:"course_material" gorm:"foreignKey:CourseMaterialID"`
	UserID           uint           `json:"user_id" gorm:"not null"`
	User             User           `json:"user" gorm:"foreignKey:UserID"`
	TotalScore       int            `json:"total_score" gorm:"default:0"`
	TotalCorrect     int            `json:"total_correct" gorm:"comment:Count of correct answers"`
	AnswersLog       string         `json:"answers_log" gorm:"type:json;comment:Full history of every move"`
	StartedAt        *time.Time     `json:"started_at"`
}
