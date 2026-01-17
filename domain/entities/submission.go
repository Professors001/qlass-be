package entities

import "gorm.io/gorm"

type Submission struct {
	gorm.Model
	CourseMaterialID uint           `json:"course_material_id"`
	CourseMaterial   CourseMaterial `json:"course_material" gorm:"foreignKey:CourseMaterialID"`
	UserID           uint           `json:"user_id"`
	User             User           `json:"user" gorm:"foreignKey:UserID"`
	QuizLogID        *uint          `json:"quiz_log_id" gorm:"comment:Null if assignment"`
	QuizGameLog      *QuizGameLog   `json:"quiz_game_log" gorm:"foreignKey:QuizLogID"`
	StudentComment   string         `json:"student_comment" gorm:"type:text"`
	Score            *int           `json:"score" gorm:"comment:e.g. 85/100"`
	TeacherFeedback  string         `json:"teacher_feedback" gorm:"type:text"`
	Status           string         `json:"status" gorm:"default:submitted;type:varchar(50);comment:submitted, graded, returned, late, draft"`
}
