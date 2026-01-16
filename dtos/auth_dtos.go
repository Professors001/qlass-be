package dtos

// RegisterRequest defines the payload for user registration
type RegisterRequestStepOneDto struct {
	UniversityID string `json:"university_id"`
	Email        string `json:"email" binding:"required,email"`
	Password     string `json:"password" binding:"required,min=6"`
	FirstName    string `json:"first_name" binding:"required"`
	LastName     string `json:"last_name" binding:"required"`
	Role         string `json:"role"`
}

type ResponseRegisterStepOneDto struct {
	Message string `json:"message"`
	RefEmail string `json:"ref_email"`
	ExpiresInSeconds int64 `json:"expires_in_seconds"`
}

type RegisterRequestStepTwoDto struct {
	Email string `json:"email" binding:"required,email"`
	OTP   string `json:"otp" binding:"required"`
}

type ResponseRegisterStepTwoDto struct {
	Message string `json:"message"`
}