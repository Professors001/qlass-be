package transforms

import (
	"qlass-be/domain/entities"
	"qlass-be/dtos"
	"strconv"
)

func EntityToClassDetailsDto(class entities.Class, ownerProfileImg string) dtos.ClassDetailsDto {
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
		OwnerProfileImg: ownerProfileImg,
		CreatedAt:       class.CreatedAt.String(),
		UpdatedAt:       class.UpdatedAt.String(),
	}
}

func EntityToStudentDetailsDto(enrollment entities.ClassEnrollment) dtos.StudentDetailsDto {
	var profileImg string
	if enrollment.User.ProfileImgAttachment != nil {
		profileImg = enrollment.User.ProfileImgAttachment.ObjectKey
	}

	return dtos.StudentDetailsDto{
		UniversityID: enrollment.User.UniversityID,
		FirstName:    enrollment.User.FirstName,
		LastName:     enrollment.User.LastName,
		ProfileImg:   profileImg,
		Email:        enrollment.User.Email,
		EnrolledRole: enrollment.Role,
		EnrolledAt:   enrollment.CreatedAt.String(),
		Status:       enrollment.Status,
	}
}

func UpdateClassReqToClassEntity(req *dtos.UpdateClassRequestDto, class *entities.Class) entities.Class {
	if req.Name != "" && req.Name != class.Name {
		class.Name = req.Name
	}

	if req.Description != "" && req.Description != class.Description {
		class.Description = req.Description
	}

	if req.Room != "" && req.Room != class.Room {
		class.Room = req.Room
	}

	if req.Section != "" && req.Section != class.Section {
		class.Section = req.Section
	}

	if req.Term != "" && req.Term != class.Term {
		class.Term = req.Term
	}

	if req.Hide == true && class.IsArchived == false {
		class.IsArchived = true
	}

	if req.Hide == false && class.IsArchived == true {
		class.IsArchived = false
	}

	return *class
}
