package entities

import "gorm.io/gorm"

type QuizOption struct {
	gorm.Model
	QuestionID int    `json:"question_id" gorm:"not null"`
	OptionText string `json:"option_text" gorm:"not null;type:varchar(500)"`
	IsCorrect  bool   `json:"is_correct" gorm:"default:false"`
}
