package transform

import (
	"qlass-be/domain/entities"
	"qlass-be/dtos"
	"strconv"
)

func EntityToClassDetailsDto(class entities.Class) dtos.ClassDetailsDto {
	return dtos.ClassDetailsDto{
		Id:              strconv.FormatUint(uint64(class.ID), 10),
		Name:            class.Name,
		Description:     class.Description,
		Section:         class.Section,
		Term:            class.Term,
		Room:            class.Room,
		InviteCode:      class.InviteCode,
		IsArchived:      class.IsArchived,
		OwnerID:         strconv.FormatUint(uint64(class.OwnerID), 10),
		OwnerFirstName:  class.Owner.FirstName,
		OwnerLastName:   class.Owner.LastName,
		OwnerProfileImg: class.Owner.ProfileImgURL,
		CreatedAt:       class.CreatedAt.String(),
		UpdatedAt:       class.UpdatedAt.String(),
	}
}

func EntityToStudentDetailsDto(enrollment entities.ClassEnrollment) dtos.StudentDetailsDto {
	return dtos.StudentDetailsDto{
		UniversityID: enrollment.User.UniversityID,
		FirstName:    enrollment.User.FirstName,
		LastName:     enrollment.User.LastName,
		ProfileImg:   enrollment.User.ProfileImgURL,
		Email:        enrollment.User.Email,
		EnrolledRole: enrollment.Role,
		EnrolledAt:   enrollment.CreatedAt.String(),
		Status:       enrollment.Status,
	}
}
