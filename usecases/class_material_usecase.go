package usecases

import (
	"errors"
	"qlass-be/domain/entities"
	"qlass-be/domain/repositories"
	"qlass-be/dtos"
)

type ClassMaterialUseCase interface {
	CreateClassMaterial(dto *dtos.CreateClassMaterialDto, ownerID uint) error
}

type classMaterialUseCase struct {
	classMaterialRepo repositories.ClassMaterialRepository
	classRepo         repositories.ClassRepository
	attachmentRepo    repositories.AttachmentRepository
}

func NewClassMaterialUseCase(classMaterialRepo repositories.ClassMaterialRepository, classRepo repositories.ClassRepository, attachmentRepo repositories.AttachmentRepository) ClassMaterialUseCase {
	return &classMaterialUseCase{
		classMaterialRepo: classMaterialRepo,
		classRepo:         classRepo,
		attachmentRepo:    attachmentRepo,
	}
}

func (u *classMaterialUseCase) CreateClassMaterial(dto *dtos.CreateClassMaterialDto, ownerID uint) error {

	classId := dto.ClassID

	class, err := u.classRepo.GetByID(classId)
	if err != nil {
		return err
	}

	if class.OwnerID != ownerID {
		return errors.New("only class owner can create class material")
	}

	classMaterial := &entities.ClassMaterial{
		Title:       dto.Title,
		Description: dto.Description,
		ClassID:     dto.ClassID,
		Type:        dto.Type,
		IsPublished: dto.Action == "publish",
		Points:      dto.Points,
		DueAt:       dto.DueAt,
	}

	err = u.classMaterialRepo.Create(classMaterial)
	if err != nil {
		return err
	}

	for _, attachmentID := range dto.AttachmentIds {
		attachment, err := u.attachmentRepo.GetByID(attachmentID)
		if err != nil {
			return err
		}

		attachment.ClassMaterialID = &classMaterial.ID

		err = u.attachmentRepo.Update(attachment)
		if err != nil {
			return err
		}
	}

	return nil
}
