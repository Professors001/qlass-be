package transforms

import (
	"qlass-be/domain/entities"
	"qlass-be/dtos"
	"time"
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
		StudentFirstName:    student.FirstName,
		StudentLastName:     student.LastName,
		StudentUniversityID: student.UniversityID,
		StudentProfileImg:   profileImgURL,
	}
}

func EntityToStudentScoreDto(material *entities.ClassMaterial, submission *entities.Submission) *dtos.StudentScoreDto {
	status := "draft"
	var score *int
	var submittedAt *time.Time

	if submission != nil {
		status = submission.Status
		score = submission.Score
		submittedAt = submission.SubmittedAt
	}

	return &dtos.StudentScoreDto{
		ClassMaterialID:    material.ID,
		ClassMaterialTitle: material.Title,
		ClassMaterialType:  material.Type,
		Score:              score,
		MaxScore:           material.Points,
		Status:             status,
		SubmittedAt:        submittedAt,
	}
}

func ToGetStudentScoresResponseDto(scores []*dtos.StudentScoreDto, student *entities.User, profileImgURL string, totalMaxScore int, totalStudentScore int) *dtos.GetStudentScoresResponseDto {
	return &dtos.GetStudentScoresResponseDto{
		Scores:            scores,
		TotalMaxScore:     totalMaxScore,
		TotalStudentScore: totalStudentScore,
		StudentID:         student.ID,
		StudentFirstName:  student.FirstName,
		StudentLastName:   student.LastName,
		StudentProfileImg: profileImgURL,
	}
}

func EntityToStudentSubmissionSummaryDto(material *entities.ClassMaterial, submission *entities.Submission) *dtos.StudentSubmissionSummaryDto {
	status := "draft"
	var submissionID *uint
	var score *int
	var submittedAt *time.Time
	var createdAt *time.Time
	var updatedAt *time.Time
	studentComment := ""
	teacherFeedback := ""

	if submission != nil {
		submissionID = &submission.ID
		status = submission.Status
		score = submission.Score
		submittedAt = submission.SubmittedAt
		createdAt = &submission.CreatedAt
		updatedAt = &submission.UpdatedAt
		studentComment = submission.StudentComment
		teacherFeedback = submission.TeacherFeedback
	}

	return &dtos.StudentSubmissionSummaryDto{
		SubmissionID:       submissionID,
		ClassMaterialID:    material.ID,
		ClassMaterialTitle: material.Title,
		ClassMaterialType:  material.Type,
		StudentComment:     studentComment,
		TeacherFeedback:    teacherFeedback,
		Score:              score,
		MaxScore:           material.Points,
		Status:             status,
		SubmittedAt:        submittedAt,
		CreatedAt:          createdAt,
		UpdatedAt:          updatedAt,
	}
}

func ToGetStudentSubmissionsByClassResponseDto(studentID uint, submissions []*dtos.StudentSubmissionSummaryDto) *dtos.GetStudentSubmissionsByClassResponseDto {
	return &dtos.GetStudentSubmissionsByClassResponseDto{
		Submissions: submissions,
		StudentID:   studentID,
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
