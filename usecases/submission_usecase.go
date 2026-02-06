package usecases

import (
	"errors"
	"qlass-be/domain/repositories"
	"qlass-be/dtos"
	"qlass-be/transform"
	"strings"
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

	// 1. Check for duplicate
	_, err = u.submissionRepo.GetByClassMaterialIDAndStudentID(classId, studentID)

	// 2. Handle the error logic properly
	if err == nil {
		// If err is NIL, it means a record WAS found -> Duplicate!
		return errors.New("submission already exists")
	} else {
		// If there IS an error, we check if it's strictly "Record Not Found"
		// Note: You might need to use `gorm.ErrRecordNotFound` or your repo's specific error constant
		if err.Error() != "record not found" && !strings.Contains(err.Error(), "record not found") {
			// It's a real database error (e.g., connection died), so return it
			return err
		}
		// If the error IS "record not found", we do nothing and proceed!
	}

	// 3. Fix the Potential Crash (Panic) here
	// If DueAt is nil, the code below will crash your server.
	isLate := false
	if classMaterial.DueAt != nil {
		isLate = time.Now().After(*classMaterial.DueAt)
	}

	submission := transform.CreateToSubmissionEntity(&dto, studentID, isLate)

	err = u.submissionRepo.Create(submission)
	if err != nil {
		return err
	}

	for _, attachmentID := range dto.AttachmentIds {
		attachment, err := u.attachmentRepo.GetByID(attachmentID)
		if err != nil {
			return err
		}

		attachment.SubmissionID = &submission.ID

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

	attachments, err := u.attachmentUseCase.GetAttachmentsBySubmissionID(val.ID)
	if err != nil {
		return nil, err
	}

	return transform.EntityToGetSubmissionResponseDto(val, attachments), nil
}

func (u *submissionUseCase) GetSubmissonByMaterialIDAndStudentID(classMaterialID uint, studentID uint) (*dtos.GetSubmissionResponseDto, error) {
	val, err := u.submissionRepo.GetByClassMaterialIDAndStudentID(classMaterialID, studentID)
	if err != nil {
		return nil, err
	}

	attachments, err := u.attachmentUseCase.GetAttachmentsBySubmissionID(val.ID)
	if err != nil {
		return nil, err
	}

	return transform.EntityToGetSubmissionResponseDto(val, attachments), nil
}

func (u *submissionUseCase) GetSubmissionsByMaterialID(classMaterialID uint) ([]*dtos.GetSubmissionResponseDto, error) {
	return nil, nil
}

func (u *submissionUseCase) GetSubmissionsByStudentID(studentID uint) ([]*dtos.GetSubmissionResponseDto, error) {
	return nil, nil
}
