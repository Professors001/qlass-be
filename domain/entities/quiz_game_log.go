package entities

import (
	"time"

	"gorm.io/gorm"
)

type QuizGameLog struct {
	gorm.Model
	ClassMaterialID uint          `json:"class_material_id" gorm:"not null;index"`
	ClassMaterial   ClassMaterial `json:"class_material" gorm:"foreignKey:ClassMaterialID"`
	UserID          uint          `json:"user_id" gorm:"not null"`
	User            User          `json:"user" gorm:"foreignKey:UserID"`
	TotalScore      int           `json:"total_score" gorm:"default:0"`
	TotalCorrect    int           `json:"total_correct" gorm:"comment:Count of correct answers"`
	AnswersLog      string        `json:"answers_log" gorm:"type:json;comment:Full history of every move"`
	StartedAt       *time.Time    `json:"started_at"`
}
