package entities

import "time"

type ClassEnrollment struct {
	ID       int       `json:"id" gorm:"primaryKey;autoIncrement"`
	UserID   int       `json:"user_id" gorm:"not null"`
	ClassID  int       `json:"class_id" gorm:"not null"`
	Role     string    `json:"role" gorm:"default:student;type:varchar(50);comment:student, ta (Teaching Assistant)"`
	Status   string    `json:"status" gorm:"default:active;type:varchar(50);comment:active, dropped, pending, banned"`
	JoinedAt time.Time `json:"joined_at" gorm:"autoCreateTime"`
}