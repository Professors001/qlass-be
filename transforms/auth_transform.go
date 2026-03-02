package transforms

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

func TempRegisterDataDtoToUserEntity(dto dtos.TempRegisterDataDto) *entities.User {
	return &entities.User{
		UniversityID: dto.UniversityID,
		Email:        dto.Email,
		PasswordHash: dto.PasswordHash,
		FirstName:    dto.FirstName,
		LastName:     dto.LastName,
		Role:         dto.Role,
		IsVerified:   true,
		IsActive:     true,
	}
}

func UserEntityToUserDisplayResponse(u *entities.User, imgUrl string) *dtos.UserDisplayData {
	return &dtos.UserDisplayData{
		UniversityID:  u.UniversityID,
		Email:         u.Email,
		ProfileImgUrl: imgUrl,
		FirstName:     u.FirstName,
		LastName:      u.LastName,
		Role:          u.Role,
	}
}
