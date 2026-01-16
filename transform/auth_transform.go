package transform

import (
	"qlass-be/domain/entities"
	"qlass-be/dtos"
)

func ToUserResponse(u *entities.User) string {
	return "TODO: implement this function"
}

func RequestToTempRegisterDataDto(req *dtos.RegisterRequestStepOneDto,
	passwordHash string,
	otp string) dtos.TempRegisterDataDto {

	var dto dtos.TempRegisterDataDto
	dto.UniversityID = req.UniversityID
	dto.Email = req.Email
	dto.PasswordHash = passwordHash
	dto.FirstName = req.FirstName
	dto.LastName = req.LastName
	dto.Role = req.Role
	dto.OTP = otp
	return dto

}
