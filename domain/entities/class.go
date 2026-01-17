package entities

import "gorm.io/gorm"

type Class struct {
	gorm.Model
	Name        string `json:"name" gorm:"type:varchar(255);comment:e.g. Intro to Computer Science"`
	Description string `json:"description" gorm:"type:text;comment:Course syllabus or details"`
	Section     string `json:"section" gorm:"type:varchar(100);comment:e.g. Section 1 or Group A"`
	Term        string `json:"term" gorm:"type:varchar(100);comment:e.g. 1/2025 or Fall 2025"`
	Room        string `json:"room" gorm:"type:varchar(100);comment:e.g. Lab 402"`
	InviteCode  string `json:"invite_code" gorm:"unique;type:varchar(6);comment:Random 6-char code for students to join"`
	IsArchived  bool   `json:"is_archived" gorm:"default:false;comment:Hide old classes without deleting"`
	OwnerID     uint   `json:"owner_id" gorm:"not null;comment:The Teacher"`
	Owner       User   `json:"owner" gorm:"foreignKey:OwnerID"`
}
