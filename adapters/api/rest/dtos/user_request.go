package dtos

// RegisterRequest defines the payload for user registration
type RegisterRequestDto struct {
	Email     string `json:"email" binding:"required,email"`
	Password  string `json:"password" binding:"required,min=6"`
	FirstName string `json:"first_name" binding:"required"`
	LastName  string `json:"last_name" binding:"required"`
}

// UpdateUserRequest (Example for later)
type UpdateUserRequestDto struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	ProfileImg string `json:"profile_img_url"`
}