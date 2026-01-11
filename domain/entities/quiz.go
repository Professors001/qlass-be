package entities

import "time"

type Quiz struct {
	ID                     int        `json:"id" gorm:"primaryKey;autoIncrement"`
	CourseMaterialID       int        `json:"course_material_id" gorm:"unique;not null"`
	Title                  string     `json:"title" gorm:"type:varchar(255)"`
	Description            string     `json:"description" gorm:"type:text"`
	DefaultTimePerQuestion int        `json:"default_time_per_question" gorm:"default:30"`
	CreatedAt              time.Time  `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt              *time.Time `json:"updated_at" gorm:"autoUpdateTime"`
}
