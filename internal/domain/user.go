package domain

import (
	"time"
	"gorm.io/gorm"
)

// 1. The Entity (Database Schema)
type User struct {
	ID            uint           `gorm:"primaryKey" json:"id"`
	UUID          string         `gorm:"uniqueIndex;not null;type:uuid" json:"uuid"`
	UniversityID  *string        `gorm:"uniqueIndex" json:"university_id"`
	Email         string         `gorm:"uniqueIndex;not null" json:"email"`
	PasswordHash  string         `gorm:"not null" json:"-"`
	FirstName     string         `json:"first_name"`
	LastName      string         `json:"last_name"`
	ProfileImgURL string         `json:"profile_img_url"`
	Role          string         `gorm:"default:student" json:"role"`
	IsVerified    bool           `gorm:"default:false" json:"is_verified"`
	IsActive      bool           `gorm:"default:true" json:"is_active"`
	
	CreatedAt     time.Time      `json:"created_at"`
	UpdatedAt     time.Time      `json:"updated_at"`
	LastLoginAt   *time.Time     `json:"last_login_at"`
	DeletedAt     gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
}

// 2. The Repository Interface (The Contract)
// Other layers (Usecase) will talk to this, not directly to GORM.
type UserRepository interface {
	Create(user *User) error
	GetByEmail(email string) (*User, error)
	GetByID(id uint) (*User, error)
}