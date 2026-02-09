package usecases

import (
	"errors"
	"qlass-be/domain/repositories"
	"qlass-be/dtos"
	"qlass-be/transforms"
	"qlass-be/utils"
	"time"
)

type SubmissionUseCase interface {
	CreateSubmission(dto dtos.CreateSubmissionDto, studentID uint) error
	GetSubmissionByID(id uint) (*dtos.GetSubmissionResponseDto, error)
	GetSubmissonByMaterialIDAndStudentID(classMaterialID uint, studentID uint) (*dtos.GetSubmissionResponseDto, error)
	GetSubmissionsByMaterialID(classMaterialID uint) ([]*dtos.GetSubmissionResponseDto, error)
	GetSubmissionsByStudentID(studentID uint) ([]*dtos.GetSubmissionResponseDto, error)
}

type submissionUseCase struct {
	submissionRepo    repositories.SubmissionRepository
	classMaterialRepo repositories.ClassMaterialRepository
	attachmentRepo    repositories.AttachmentRepository
	attachmentUseCase AttachmentUseCase
}

func NewSubmissionUseCase(submissionRepo repositories.SubmissionRepository, classMaterialRepo repositories.ClassMaterialRepository, attachmentRepo repositories.AttachmentRepository, attachmentUseCase AttachmentUseCase) SubmissionUseCase {
	return &submissionUseCase{
		submissionRepo:    submissionRepo,
		classMaterialRepo: classMaterialRepo,
		attachmentRepo:    attachmentRepo,
		attachmentUseCase: attachmentUseCase,
	}
}

func (u *submissionUseCase) CreateSubmission(dto dtos.CreateSubmissionDto, studentID uint) error {
	classId := dto.ClassMaterialID

	classMaterial, err := u.classMaterialRepo.GetByID(classId)
	if err != nil {
		return errors.New("class material not found")
	}

	// 1. Get the submission AND the error
	existingSubmission, err := u.submissionRepo.GetByClassMaterialIDAndStudentID(classId, studentID)

	// 2. First, check for actual database crashes (Connection failed, SQL error)
	if err != nil {
		return err
	}

	// 3. Then, check if we actually FOUND data
	// If existingSubmission is NOT nil, THAT is a duplicate.
	if existingSubmission != nil {
		return errors.New("submission already exists")
	}

	// 3. Fix the Potential Crash (Panic) here
	// If DueAt is nil, the code below will crash your server.
	isLate := false
	if classMaterial.DueAt != nil {
		isLate = time.Now().After(*classMaterial.DueAt)
	}

	submission := transforms.CreateToSubmissionEntity(&dto, studentID, isLate)

	err = u.submissionRepo.Create(submission)
	if err != nil {
		return err
	}

	for _, attachmentID := range dto.AttachmentIds {
		attachment, err := u.attachmentRepo.GetByID(attachmentID)
		if err != nil {
			return err
		}

		attachment.OwnerType = utils.Ptr("submission")
		attachment.OwnerID = &submission.ID

		err = u.attachmentRepo.Update(attachment)
		if err != nil {
			return err
		}
	}

	return nil
}

func (u *submissionUseCase) GetSubmissionByID(id uint) (*dtos.GetSubmissionResponseDto, error) {

	val, err := u.submissionRepo.GetByID(id)
	if err != nil {
		return nil, err
	}

	attachments, err := u.attachmentUseCase.GetAttachmentsByOwner("submission", val.ID)
	if err != nil {
		return nil, err
	}

	return transforms.EntityToGetSubmissionResponseDto(val, attachments), nil
}

func (u *submissionUseCase) GetSubmissonByMaterialIDAndStudentID(classMaterialID uint, studentID uint) (*dtos.GetSubmissionResponseDto, error) {
	val, err := u.submissionRepo.GetByClassMaterialIDAndStudentID(classMaterialID, studentID)

	if err != nil {
		return nil, err
	}

	if val == nil {
		return nil, nil
	}

	attachments, err := u.attachmentUseCase.GetAttachmentsByOwner("submission", val.ID)
	if err != nil {
		return nil, err
	}

	return transforms.EntityToGetSubmissionResponseDto(val, attachments), nil
}

func (u *submissionUseCase) GetSubmissionsByMaterialID(classMaterialID uint) ([]*dtos.GetSubmissionResponseDto, error) {
	return nil, nil
}

func (u *submissionUseCase) GetSubmissionsByStudentID(studentID uint) ([]*dtos.GetSubmissionResponseDto, error) {
	return nil, nil
}
