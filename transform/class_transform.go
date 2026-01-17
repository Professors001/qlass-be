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
		OwnerName:       "STILL NOT IMPLEMENT",
		OwnerProfileImg: "STILL NOT IMPLEMENT",
		CreatedAt:       class.CreatedAt.String(),
		UpdatedAt:       class.UpdatedAt.String(),
	}
}
