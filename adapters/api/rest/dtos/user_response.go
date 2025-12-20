package dtos

import (
	"qlass-be/domain/entities"
	"time"
)

// UserResponse sanitizes the domain entity for public view
type UserResponseDto struct {
	UniversityID  *string   `json:"university_id,omitempty"` // omitempty hides it if null
	Email         string    `json:"email"`
	FirstName     string    `json:"first_name"`
	LastName      string    `json:"last_name"`
	Role          string    `json:"role"`
	ProfileImgURL string    `json:"profile_img_url"`
	IsVerified    bool      `json:"is_verified"`
	CreatedAt     time.Time `json:"created_at"`
}

// Helper function to convert Domain -> DTO
// This keeps the conversion logic out of the handler!
func ToUserResponse(u *entities.User) UserResponseDto {
	return UserResponseDto{
		UniversityID:  u.UniversityID,
		Email:         u.Email,
		FirstName:     u.FirstName,
		LastName:      u.LastName,
		Role:          u.Role,
		ProfileImgURL: u.ProfileImgURL,
		IsVerified:    u.IsVerified,
		CreatedAt:     u.CreatedAt,
	}
}
