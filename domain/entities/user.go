package entities

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	UniversityID           string      `gorm:"uniqueIndex" json:"university_id"`
	Email                  string      `gorm:"uniqueIndex;not null" json:"email"`
	PasswordHash           string      `gorm:"not null" json:"-"`
	FirstName              string      `json:"first_name"`
	LastName               string      `json:"last_name"`
	ProfileImgAttachmentID *uint       `json:"profile_img_attachment_id"`
	ProfileImgAttachment   *Attachment `gorm:"foreignKey:ProfileImgAttachmentID;constraint:-" json:"profile_img_attachment"`
	Role                   string      `gorm:"default:student" json:"role"`
	IsVerified             bool        `gorm:"default:false" json:"is_verified"`
	IsActive               bool        `gorm:"default:true" json:"is_active"`
	LastLoginAt            *time.Time  `json:"last_login_at"`
}
