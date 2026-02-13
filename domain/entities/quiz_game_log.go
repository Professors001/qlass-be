package entities

import (
	"time"

	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type QuizGameLog struct {
	gorm.Model

	ClassMaterialID  uint                  `json:"class_material_id" gorm:"unique;not null"`
	ClassMaterial    *ClassMaterial        `json:"class_material" gorm:"foreignKey:ClassMaterialID"`
	QuizPin          string                `json:"quiz_pin" gorm:"type:varchar(6);comment:only visible while ongoing"`
	QuizSnapshot     datatypes.JSON        `json:"quiz_snapshot" gorm:"type:jsonb;comment:Contains full Q&A structure"`
	Status           string                `json:"status" gorm:"default:not_started;type:varchar(50);comment:not_started, ongoing, finished"`
	AverageScore     float64               `json:"average_score"`
	StartedAt        *time.Time            `json:"started_at"`
	FinishedAt       *time.Time            `json:"finished_at"`
	StudentResponses []QuizStudentResponse `json:"student_responses" gorm:"foreignKey:QuizGameLogID"`
}
