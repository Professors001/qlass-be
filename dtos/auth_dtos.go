package dtos

import "qlass-be/domain/entities"

type RegisterRequestStepOneDto struct {
	UniversityID string `json:"university_id" binding:"required"`
	Email        string `json:"email" binding:"required,email"`
	Password     string `json:"password" binding:"required,min=6"`
	FirstName    string `json:"first_name" binding:"required"`
	LastName     string `json:"last_name" binding:"required"`
	Role         string `json:"role" binding:"required,oneof=student teacher"`
}

type ResponseRegisterStepOneDto struct {
	Message          string `json:"message"`
	RefEmail         string `json:"ref_email"`
	ExpiresInSeconds int64  `json:"expires_in_seconds"`
}

type TempRegisterDataDto struct {
	UniversityID string `json:"university_id"`
	Email        string `json:"email"`
	PasswordHash string `json:"password_hash"`
	FirstName    string `json:"first_name"`
	LastName     string `json:"last_name"`
	Role         string `json:"role"`
	OTP          string `json:"otp"`
}

type RegisterRequestStepTwoDto struct {
	Email string `json:"email" validate:"required,email"`
	OTP   string `json:"otp" validate:"required"`
}

type ResponseRegisterStepTwoDto struct {
	Message string `json:"message"`
}

type UserDisplayData struct {
	UniversityID  string `json:"university_id"`
	Email         string `json:"email"`
	ProfileImgUrl string `json:"profile_img_url"`
	FirstName     string `json:"first_name"`
	LastName      string `json:"last_name"`
	Role          string `json:"role"`
}

type LoginRequestDto struct {
	UniversityID string `json:"university_id" validate:"required"`
	Password     string `json:"password" validate:"required"`
}

type LoginResponseDto struct {
	Message string          `json:"message"`
	Token   string          `json:"token"`
	User    UserDisplayData `json:"user"`
}

type CreateTeacherRequestDto struct {
	UniversityID string `json:"university_id" binding:"required"`
	Email        string `json:"email" binding:"required,email"`
	Password     string `json:"password" binding:"required,min=6"`
	FirstName    string `json:"first_name" binding:"required"`
	LastName     string `json:"last_name" binding:"required"`
}

type CreateTeacherResponseDto struct {
	Message string `json:"message"`
	UserID  uint   `json:"user_id"`
}

type UpdateUserRequestDto struct {
	FirstName              string `json:"first_name"`
	LastName               string `json:"last_name"`
	ProfileImgAttachmentID uint   `json:"profile_img_attachment_id"`
}

type ChangePasswordRequestDto struct {
	CurrentPassword string `json:"current_password" binding:"required"`
	NewPassword     string `json:"new_password" binding:"required,min=6"`
}

type ChangePasswordResponseDto struct {
	Message string `json:"message"`
}

type ForgetPasswordStep1RequestDto struct {
	UniversityID string `json:"university_id" binding:"required"`
}

type ForgetPasswordStep1ResponseDto struct {
	Message string `json:"message"`
	Email   string `json:"email"`
}

type TempForgetPasswordData struct {
	UniversityID string        `json:"university_id" binding:"required"`
	OTP          string        `json:"otp" binding:"required"`
	User         entities.User `json:"user"`
}

type ForgetPasswordStep2RequestDto struct {
	UniversityID string `json:"university_id" binding:"required"`
	OTP          string `json:"otp" binding:"required"`
	NewPassword  string `json:"new_password" binding:"required,min=6"`
}

type ForgetPasswordStep2ResponseDto struct {
	Message string `json:"message"`
}
