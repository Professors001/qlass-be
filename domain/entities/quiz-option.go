package entities

type QuizOption struct {
	ID         int    `json:"id" gorm:"primaryKey;autoIncrement"`
	QuestionID int    `json:"question_id" gorm:"not null"`
	OptionText string `json:"option_text" gorm:"not null;type:varchar(500)"`
	IsCorrect  bool   `json:"is_correct" gorm:"default:false"`
}