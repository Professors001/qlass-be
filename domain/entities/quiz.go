package entities

import (
	"gorm.io/gorm"
)

type Quiz struct {
	gorm.Model
	ClassID                uint           `json:"class_id" gorm:"not null"`
	Class                  Class          `json:"class" gorm:"foreignKey:ClassID"`
	Title                  string         `json:"title" gorm:"type:varchar(255)"`
	Description            string         `json:"description" gorm:"type:text"`
	DefaultTimePerQuestion int            `json:"default_time_per_question" gorm:"default:30"`
	Questions              []QuizQuestion `json:"questions" gorm:"foreignKey:QuizID"`
}
