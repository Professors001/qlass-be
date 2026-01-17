package entities

import "gorm.io/gorm"

type ClassEnrollment struct {
	gorm.Model
	UserID  int    `json:"user_id" gorm:"not null"`
	ClassID int    `json:"class_id" gorm:"not null"`
	Role    string `json:"role" gorm:"default:student;type:varchar(50);comment:student, ta (Teaching Assistant)"`
	Status  string `json:"status" gorm:"default:active;type:varchar(50);comment:active, dropped, pending, banned"`
}
