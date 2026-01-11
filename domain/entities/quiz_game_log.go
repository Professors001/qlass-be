package entities

import "time"

type QuizGameLog struct {
	ID               int        `json:"id" gorm:"primaryKey;autoIncrement"`
	CourseMaterialID int        `json:"course_material_id" gorm:"not null"`
	UserID           int        `json:"user_id" gorm:"not null"`
	TotalScore       int        `json:"total_score" gorm:"default:0"`
	TotalCorrect     int        `json:"total_correct" gorm:"comment:Count of correct answers"`
	AnswersLog       string     `json:"answers_log" gorm:"type:json;comment:Full history of every move"`
	StartedAt        *time.Time `json:"started_at"`
	FinishedAt       time.Time  `json:"finished_at" gorm:"autoCreateTime"`
}