package entities

import "time"

type QuizStudentResponse struct {
	ID            uint      `gorm:"primaryKey" json:"id"`
	QuizGameLogID uint      `gorm:"index" json:"quiz_game_log_id"`
	UserID        uint      `gorm:"index" json:"user_id"`
	QuestionID    uint      `json:"question_id"`
	OptionID      *uint     `json:"option_id"` // Nullable if they didn't answer
	IsCorrect     bool      `json:"is_correct"`
	TimeTaken     int       `json:"time_taken"` // In milliseconds
	Points        int       `json:"points"`
	CreatedAt     time.Time `json:"created_at"`
}
