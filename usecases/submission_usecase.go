package usecases

import (
	"errors"
	"log"
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
	GetSubmissionsByMaterialID(classMaterialID uint, teacherId uint) (*dtos.GetSubmissionsByClassMaterialResponseDto, error)
	GetSubmissionsByStudentID(studentID uint) (*dtos.GetSubmissionsByClassMaterialResponseDto, error)
	StudentSaveSubmission(dto dtos.StudentSaveSubmissionDto, studentID uint) error
	TeacherSaveSubmission(dto dtos.TeacherSaveSubmissionDto, teacherID uint) error
}
type submissionUseCase struct {
	submissionRepo    repositories.SubmissionRepository
	classMaterialRepo repositories.ClassMaterialRepository
	attachmentRepo    repositories.AttachmentRepository
	classRepo         repositories.ClassRepository
	userRepo          repositories.UserRepository
	userUsecase       UserUseCase
	attachmentUseCase AttachmentUseCase
}

func NewSubmissionUseCase(
	submissionRepo repositories.SubmissionRepository,
	classMaterialRepo repositories.ClassMaterialRepository,
	attachmentRepo repositories.AttachmentRepository,
	attachmentUseCase AttachmentUseCase,
	classRepo repositories.ClassRepository,
	userUsecase UserUseCase,
	userRepo repositories.UserRepository) SubmissionUseCase {
	return &submissionUseCase{
		submissionRepo:    submissionRepo,
		classMaterialRepo: classMaterialRepo,
		attachmentRepo:    attachmentRepo,
		attachmentUseCase: attachmentUseCase,
		classRepo:         classRepo,
		userUsecase:       userUsecase,
		userRepo:          userRepo,
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

	imgUrl := u.userUsecase.GetProfileImgUrlByUserID(val.UserID)

	return transforms.EntityToGetSubmissionResponseDto(val, attachments, imgUrl), nil
}

func (u *submissionUseCase) GetSubmissonByMaterialIDAndStudentID(classMaterialID uint, studentID uint) (*dtos.GetSubmissionResponseDto, error) {
	val, err := u.submissionRepo.GetByClassMaterialIDAndStudentID(classMaterialID, studentID)
	if err != nil {
		return nil, err
	}

	if val == nil {
		// Create a new "Draft/Assigned" submission
		submission := transforms.CreateBaseSubmissionEntity(studentID, classMaterialID)

		err = u.submissionRepo.Create(submission)
		if err != nil {
			val, retryErr := u.submissionRepo.GetByClassMaterialIDAndStudentID(classMaterialID, studentID)
			if retryErr != nil || val == nil {
				return nil, err
			}
		} else {
			val = submission
		}
	}

	imgUrl := u.userUsecase.GetProfileImgUrlByUserID(val.UserID)

	attachments, err := u.attachmentUseCase.GetAttachmentsByOwner("submission", val.ID)
	if err != nil {
		return nil, err
	}

	return transforms.EntityToGetSubmissionResponseDto(val, attachments, imgUrl), nil
}

func (u *submissionUseCase) GetSubmissionsByMaterialID(classMaterialID uint, teacherId uint) (
	*dtos.GetSubmissionsByClassMaterialResponseDto, error) {
	submissions, err := u.submissionRepo.GetByClassMaterialID(classMaterialID)
	if err != nil {
		return nil, err
	}

	var submissionDtos []*dtos.TeacherGetSubmissionResponseDto
	for _, submission := range submissions {
		if submission.Status != "draft" && submission.SubmittedAt != nil {
			attachments, err := u.attachmentUseCase.GetAttachmentsByOwner("submission", submission.ID)
			if err != nil {
				return nil, err
			}

			log.Println("Submission ID:", submission.ID, "has", len(attachments), "attachments")

			student, err := u.userRepo.GetByID(submission.UserID)
			if err != nil {
				return nil, err
			}

			var profileImgURL string
			if student.ProfileImgAttachment != nil && student.ProfileImgAttachment.ObjectKey != "" {
				url, err := u.attachmentUseCase.GenerateFileURL(student.ProfileImgAttachment.ObjectKey)
				if err == nil {
					profileImgURL = url
				} else {
					log.Println("Error generating profile image URL:", err)
				}
			}

			if profileImgURL == "" {
				profileImgURL = "https://ui-avatars.com/api/?name=" + student.FirstName + "+" + student.LastName
			}

			submissionDtos = append(submissionDtos, transforms.EntityToTeacherGetSubmissionResponseDto(submission, attachments, student, profileImgURL))
		}
	}

	res := &dtos.GetSubmissionsByClassMaterialResponseDto{
		Submissions: submissionDtos,
	}

	return res, nil
}

func (u *submissionUseCase) GetSubmissionsByStudentID(studentID uint) (*dtos.GetSubmissionsByClassMaterialResponseDto, error) {
	return nil, nil
}

func (u *submissionUseCase) StudentSaveSubmission(req dtos.StudentSaveSubmissionDto, studentID uint) error {
	submission, err := u.submissionRepo.GetByID(req.ID)
	if err != nil {
		return err
	}

	if submission == nil {
		return errors.New("submission not found")
	}

	if submission.UserID != studentID {
		return errors.New("unauthorized")
	}

	if req.StudentComment != "" && req.StudentComment != submission.StudentComment {
		submission.StudentComment = req.StudentComment
	}

	if submission.Status == "graded" {
		return errors.New("submission already graded")
	}

	if req.Status == "draft" && submission.Status != "draft" {
		submission.Status = "draft"
	}

	if req.Status == "submit" && (submission.Status != "submit" && submission.SubmittedAt == nil) {
		submission.Status = "submit"
		now := time.Now()
		submission.SubmittedAt = &now

		if submission.ClassMaterial.DueAt != nil && submission.SubmittedAt.After(*submission.ClassMaterial.DueAt) {
			submission.Status = "late"
			submission.IsLate = true
		}
	}

	// UnAttach all old Attachment
	attachments, err := u.attachmentRepo.GetByOwnerTypeAndOwnerID("submission", submission.ID)
	if err == nil {
		for _, att := range attachments {
			att.OwnerID = nil
			att.OwnerType = nil
			if err := u.attachmentRepo.Update(att); err != nil {
				return err
			}
		}
	}

	for _, attachmentID := range req.AttchmentIds {
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

	err = u.submissionRepo.Update(submission)
	if err != nil {
		return err
	}

	return nil
}

func (u *submissionUseCase) TeacherSaveSubmission(req dtos.TeacherSaveSubmissionDto, teacherID uint) error {
	submission, err := u.submissionRepo.GetByID(req.SubmissionID)
	if err != nil {
		return err
	}

	if submission == nil {
		return errors.New("submission not found")
	}

	classMaterial, err := u.classMaterialRepo.GetByID(submission.ClassMaterialID)
	if err != nil {
		return errors.New("class material not found")
	}

	class, err := u.classRepo.GetByID(classMaterial.ClassID)
	if err != nil {
		return errors.New("class not found")
	}

	if class.OwnerID != teacherID {
		log.Println(classMaterial.Class.OwnerID, teacherID)
		return errors.New("unauthorized : not owner of class")
	}

	if submission.Status != "draft" && submission.SubmittedAt == nil {
		return errors.New("submission is not submitted")
	}

	submission.Score = &req.Score
	submission.TeacherFeedback = req.Feedback
	submission.Status = "graded"

	err = u.submissionRepo.Update(submission)
	if err != nil {
		return err
	}

	return nil
}
