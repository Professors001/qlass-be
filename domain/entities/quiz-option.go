package entities

import "gorm.io/gorm"

type QuizOption struct {
	gorm.Model
	QuestionID uint         `json:"question_id" gorm:"not null"`
	Question   QuizQuestion `json:"question" gorm:"foreignKey:QuestionID"`
	OptionText string       `json:"option_text" gorm:"not null;type:varchar(500)"`
	IsCorrect  bool         `json:"is_correct" gorm:"default:false"`
	OrderIndex int          `json:"order_index" gorm:"comment:To keep options in order 1, 2, 3..."`
}
