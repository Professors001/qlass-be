package dtos

import "time"

type CreateSubmissionDto struct {
	ClassMaterialID uint   `json:"class_material_id" binding:"required"`
	StudentComment  string `json:"student_comment"`
	AttachmentIds   []uint `json:"attachment_ids"`
	Action          string `json:"action" binding:"required,oneof=draft submit"`
}

type GetSubmissionResponseDto struct {
	ID               uint                        `json:"id"`
	ClassMaterialID  uint                        `json:"class_material_id"`
	UserID           uint                        `json:"student_id"`
	StudentFirstName string                      `json:"student_first_name"`
	StudentLastName  string                      `json:"student_last_name"`
	StudentImg       string                      `json:"student_profile_img"`
	StudentComment   string                      `json:"student_comment"`
	Status           string                      `json:"status"`
	Score            *int                        `json:"score"`
	TeacherFeedback  string                      `json:"teacher_feedback"`
	CreatedAt        time.Time                   `json:"created_at"`
	UpdatedAt        time.Time                   `json:"updated_at"`
	Attachments      []*GetAttachmentResponseDto `json:"attachments"`
}

type StudentSaveSubmissionDto struct {
	ID             uint   `json:"id" binding:"required"`
	StudentComment string `json:"student_comment"`
	Status         string `json:"status" binding:"required,oneof=draft submit"`
	AttchmentIds   []uint `json:"attachment_ids"`
}

type TeacherSaveSubmissionDto struct {
	SubmissionID uint   `json:"submission_id" binding:"required"`
	Score        int    `json:"score" binding:"required"`
	Feedback     string `json:"feedback"`
}

type TeacherGetSubmissionResponseDto struct {
	GetSubmissionResponseDto

	StudentFirstName    string `json:"student_first_name"`
	StudentLastName     string `json:"student_last_name"`
	StudentUniversityID string `json:"student_university_id"`
	StudentProfileImg   string `json:"student_profile_img"`
}

type GetSubmissionsByClassMaterialResponseDto struct {
	Submissions []*TeacherGetSubmissionResponseDto `json:"submissions"`
}

type TeacherGetStudentScoreListsByUserAndClassIDRequestDto struct {
	ClassID   uint `json:"class_id" binding:"required"`
	StudentID uint `json:"student_id" binding:"required"`
}

type GetStudentScoresResponseDto struct {
	Scores            []*StudentScoreDto `json:"scores"`
	TotalMaxScore     int                `json:"total_max_score"`
	TotalStudentScore int                `json:"total_student_score"`
	StudentID         uint               `json:"student_id"`
	StudentFirstName  string             `json:"student_first_name"`
	StudentLastName   string             `json:"student_last_name"`
	StudentProfileImg string             `json:"student_profile_img"`
}

type StudentScoreDto struct {
	ClassMaterialID    uint       `json:"class_material_id"`
	ClassMaterialTitle string     `json:"class_material_title"`
	ClassMaterialType  string     `json:"class_material_type"`
	Score              *int       `json:"score"`
	MaxScore           *int       `json:"max_score"`
	Status             string     `json:"status"`
	SubmittedAt        *time.Time `json:"submitted_at"`
}

type GetStudentSubmissionsByClassRequestDto struct {
	ClassID uint `form:"class_id" binding:"required"`
}

type GetStudentSubmissionsByClassResponseDto struct {
	Submissions []*StudentSubmissionSummaryDto `json:"submissions"`
	StudentID   uint                           `json:"student_id"`
}

type StudentSubmissionSummaryDto struct {
	SubmissionID       *uint      `json:"submission_id"`
	ClassMaterialID    uint       `json:"class_material_id"`
	ClassMaterialTitle string     `json:"class_material_title"`
	ClassMaterialType  string     `json:"class_material_type"`
	StudentComment     string     `json:"student_comment"`
	TeacherFeedback    string     `json:"teacher_feedback"`
	Score              *int       `json:"score"`
	MaxScore           *int       `json:"max_score"`
	Status             string     `json:"status"`
	SubmittedAt        *time.Time `json:"submitted_at"`
	CreatedAt          *time.Time `json:"created_at"`
	UpdatedAt          *time.Time `json:"updated_at"`
}
