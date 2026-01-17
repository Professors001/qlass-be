package entities

import (
	"gorm.io/gorm"
)

type Quiz struct {
	gorm.Model
	CourseMaterialID       uint           `json:"course_material_id" gorm:"unique;not null"`
	CourseMaterial         CourseMaterial `json:"course_material" gorm:"foreignKey:CourseMaterialID"`
	Title                  string         `json:"title" gorm:"type:varchar(255)"`
	Description            string         `json:"description" gorm:"type:text"`
	DefaultTimePerQuestion int            `json:"default_time_per_question" gorm:"default:30"`
}
