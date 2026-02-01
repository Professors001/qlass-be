package entities

import "gorm.io/gorm"

type Submission struct {
	gorm.Model
	ClassMaterialID uint          `json:"class_material_id" gorm:"not null"`
	ClassMaterial   ClassMaterial `json:"class_material" gorm:"foreignKey:ClassMaterialID"`
	UserID          uint          `json:"user_id"`
	User            User          `json:"user" gorm:"foreignKey:UserID"`
	QuizLogID       *uint         `json:"quiz_log_id" gorm:"comment:Null if assignment"`
	QuizGameLog     *QuizGameLog  `json:"quiz_game_log" gorm:"foreignKey:QuizLogID"`
	StudentComment  string        `json:"student_comment" gorm:"type:text"`
	Score           *int          `json:"score" gorm:"comment:e.g. 85/100"`
	TeacherFeedback string        `json:"teacher_feedback" gorm:"type:text"`
	Status          string        `json:"status" gorm:"default:submitted;type:varchar(50);comment:submitted, graded, returned, late, draft"`
}
