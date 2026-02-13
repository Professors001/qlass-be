package entities

import "gorm.io/gorm"

type QuizStudentResponse struct {
	gorm.Model

	QuizGameLogID    uint         `json:"quiz_game_log_id" gorm:"not null;index"`
	QuizGameLog      *QuizGameLog `json:"quiz_game_log" gorm:"foreignKey:QuizGameLogID"`
	UserID           uint         `json:"user_id" gorm:"not null;index"`
	User             *User        `json:"user" gorm:"foreignKey:UserID"`
	QuestionID       uint         `json:"question_id" gorm:"comment:ID from Snapshot"`
	SelectedOptionID uint         `json:"selected_option_id" gorm:"comment:ID from Snapshot"`
	IsCorrect        bool         `json:"is_correct"`
	TimeTakenSeconds int          `json:"time_taken_seconds"`
	PointsEarned     int          `json:"points_earned"`
}
