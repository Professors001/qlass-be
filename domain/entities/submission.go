package entities

import "gorm.io/gorm"

type Submission struct {
	gorm.Model
	CourseMaterialID int    `json:"course_material_id"`
	UserID           int    `json:"user_id"`
	QuizLogID        *int   `json:"quiz_log_id" gorm:"comment:Null if assignment"`
	StudentComment   string `json:"student_comment" gorm:"type:text"`
	Score            *int   `json:"score" gorm:"comment:e.g. 85/100"`
	TeacherFeedback  string `json:"teacher_feedback" gorm:"type:text"`
	Status           string `json:"status" gorm:"default:submitted;type:varchar(50);comment:submitted, graded, returned, late, draft"`
}
