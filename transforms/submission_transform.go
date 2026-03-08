package transforms

import (
	"qlass-be/domain/entities"
	"qlass-be/dtos"
)

func CreateToSubmissionEntity(dto *dtos.CreateSubmissionDto, studentID uint, isLate bool) *entities.Submission {

	baseEntity := &entities.Submission{
		ClassMaterialID: dto.ClassMaterialID,
		UserID:          studentID,
		StudentComment:  dto.StudentComment,
		Status:          dto.Action,
	}

	if isLate && dto.Action == "SUBMIT" {
		baseEntity.Status = "late"
	}

	return baseEntity
}

func EntityToGetSubmissionResponseDto(submission *entities.Submission, attachments []*dtos.GetAttachmentResponseDto, profileImgURL string) *dtos.GetSubmissionResponseDto {
	return &dtos.GetSubmissionResponseDto{
		ID:               submission.ID,
		ClassMaterialID:  submission.ClassMaterialID,
		UserID:           submission.UserID,
		StudentFirstName: submission.User.FirstName,
		StudentLastName:  submission.User.LastName,
		StudentImg:       profileImgURL,
		StudentComment:   submission.StudentComment,
		Score:            submission.Score,
		TeacherFeedback:  submission.TeacherFeedback,
		Status:           submission.Status,
		Attachments:      attachments,
		CreatedAt:        submission.CreatedAt,
		UpdatedAt:        submission.UpdatedAt,
	}
}

func EntityToTeacherGetSubmissionResponseDto(submission *entities.Submission, attachments []*dtos.GetAttachmentResponseDto, student *entities.User, profileImgURL string) *dtos.TeacherGetSubmissionResponseDto {
	return &dtos.TeacherGetSubmissionResponseDto{
		GetSubmissionResponseDto: dtos.GetSubmissionResponseDto{
			ID:              submission.ID,
			ClassMaterialID: submission.ClassMaterialID,
			UserID:          submission.UserID,
			StudentComment:  submission.StudentComment,
			Score:           submission.Score,
			TeacherFeedback: submission.TeacherFeedback,
			Status:          submission.Status,
			Attachments:     attachments,
			CreatedAt:       submission.CreatedAt,
			UpdatedAt:       submission.UpdatedAt,
		},
		StudentFirstName:  student.FirstName,
		StudentLastName:   student.LastName,
		StudentEmail:      student.Email,
		StudentProfileImg: profileImgURL,
	}
}

func CreateBaseSubmissionEntity(studentID uint, classMaterialID uint) *entities.Submission {
	return &entities.Submission{
		ClassMaterialID: classMaterialID,
		UserID:          studentID,
		Status:          "draft",
		IsLate:          false,
	}
}
