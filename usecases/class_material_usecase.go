package usecases

import (
	"errors"
	"qlass-be/domain/repositories"
	"qlass-be/dtos"
	"qlass-be/transform"
	"time"
)

type ClassMaterialUseCase interface {
	CreateClassMaterial(dto *dtos.CreateClassMaterialDto, ownerID uint) error
	GetMaterialByID(id uint) (*dtos.GetClassMaterialDto, error)
	GetMaterialsByClassID(classID uint) ([]*dtos.GetThumnailClassMaterialDto, error)
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

	classMaterial := transform.CreateToEntity(dto)

	if dto.Action == "publish" {
		classMaterial.IsPublished = true
		now := time.Now()
		classMaterial.PublishedAt = &now
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

func (u *classMaterialUseCase) GetMaterialsByClassID(classID uint) ([]*dtos.GetThumnailClassMaterialDto, error) {
	materials, err := u.classMaterialRepo.GetByClassID(classID)
	if err != nil {
		return nil, err
	}

	response := make([]*dtos.GetThumnailClassMaterialDto, 0, len(materials))
	for _, material := range materials {
		response = append(response, transform.EntityToGetThumnailClassMaterialDto(material))
	}

	return response, nil
}
