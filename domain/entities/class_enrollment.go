package entities

import "gorm.io/gorm"

type ClassEnrollment struct {
	gorm.Model
	UserID  uint   `json:"user_id" gorm:"not null"`
	User    User   `json:"user" gorm:"foreignKey:UserID"`
	ClassID uint   `json:"class_id" gorm:"not null"`
	Class   Class  `json:"class" gorm:"foreignKey:ClassID"`
	Role    string `json:"role" gorm:"default:student;type:varchar(50);comment:student, ta (Teaching Assistant)"`
	Status  string `json:"status" gorm:"default:active;type:varchar(50);comment:active, dropped, pending, banned"`
}
