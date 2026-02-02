package transform

import (
	"qlass-be/domain/entities"
	"qlass-be/dtos"
)

func CreateToEntity(dto *dtos.CreateClassMaterialDto) *entities.ClassMaterial {
	return &entities.ClassMaterial{
		Title:       dto.Title,
		Description: dto.Description,
		ClassID:     dto.ClassID,
		Type:        dto.Type,
		IsPublished: dto.Action == "publish",
		Points:      dto.Points,
		DueAt:       dto.DueAt,
	}
}
