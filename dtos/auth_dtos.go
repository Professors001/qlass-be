package dtos

// RegisterRequest defines the payload for user registration
type RegisterRequestDto struct {
	UniversityID string `json:"university_id"`
	Email        string `json:"email" binding:"required,email"`
	Password     string `json:"password" binding:"required,min=6"`
	FirstName    string `json:"first_name" binding:"required"`
	LastName     string `json:"last_name" binding:"required"`
	Role         string `json:"role"`
}

