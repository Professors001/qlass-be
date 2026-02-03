package usecases

import (
	"errors"
	"qlass-be/domain/entities"
	"qlass-be/domain/repositories"
	"qlass-be/dtos"
	"qlass-be/transform"
)

type ClassMaterialUseCase interface {
	CreateClassMaterial(dto *dtos.CreateClassMaterialDto, ownerID uint) error
	GetMaterialByID(id uint) (*dtos.GetClassMaterialDto, error)
}

type classMaterialUseCase struct {
	classMaterialRepo repositories.ClassMaterialRepository
	classRepo         repositories.ClassRepository
	attachmentRepo    repositories.AttachmentRepository
	attachmentUseCase AttachmentUseCase
}

func NewClassMaterialUseCase(classMaterialRepo repositories.ClassMaterialRepository, classRepo repositories.ClassRepository, attachmentRepo repositories.AttachmentRepository, attachmentUseCase AttachmentUseCase) ClassMaterialUseCase {
	return &classMaterialUseCase{
		classMaterialRepo: classMaterialRepo,
		classRepo:         classRepo,
		attachmentRepo:    attachmentRepo,
		attachmentUseCase: attachmentUseCase,
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

func (u *classMaterialUseCase) GetMaterialByID(id uint) (*dtos.GetClassMaterialDto, error) {
	material, err := u.classMaterialRepo.GetByID(id)
	if err != nil {
		return nil, err
	}

	attachmentDtos, err := u.attachmentUseCase.GetAttachmentsByClassMaterialID(material.ID)
	if err != nil {
		return nil, err
	}

	return transform.EntityToGetClassMaterialDtoWithAttachments(material, attachmentDtos), nil
}
