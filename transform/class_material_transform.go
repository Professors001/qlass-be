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

func EntityToGetClassMaterialDtoWithAttachments(material *entities.ClassMaterial, attachments []*dtos.GetAttachmentResponseDto) *dtos.GetClassMaterialDto {
	return &dtos.GetClassMaterialDto{
		ID:          material.ID,
		ClassID:     material.ClassID,
		Type:        material.Type,
		Title:       material.Title,
		Description: material.Description,
		Attachments: attachments,
		CreatedAt:   material.CreatedAt,
		Points:      material.Points,
		DueAt:       material.DueAt,
	}
}

func EntityToGetThumnailClassMaterialDto(material *entities.ClassMaterial) *dtos.GetThumnailClassMaterialDto {
	return &dtos.GetThumnailClassMaterialDto{
		ID:        material.ID,
		Title:     material.Title,
		Type:      material.Type,
		CreatedAt: material.CreatedAt,
		DueAt:     material.DueAt,
	}
}
