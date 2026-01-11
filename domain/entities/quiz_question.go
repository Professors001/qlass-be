package entities

type QuizQuestion struct {
	ID               int    `json:"id" gorm:"primaryKey;autoIncrement"`
	QuizID           int    `json:"quiz_id" gorm:"not null"`
	QuestionText     string `json:"question_text" gorm:"not null;type:text"`
	MediaURL         string `json:"media_url" gorm:"type:varchar(500);comment:Image or Video for the question"`
	PointsMultiplier int    `json:"points_multiplier" gorm:"default:1;comment:1x, 2x points"`
	TimeLimitSeconds int    `json:"time_limit_seconds" gorm:"default:30"`
	OrderIndex       int    `json:"order_index" gorm:"comment:To keep questions in order 1, 2, 3..."`
}