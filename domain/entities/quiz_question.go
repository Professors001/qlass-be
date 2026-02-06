package entities

import "gorm.io/gorm"

type QuizQuestion struct {
	gorm.Model
	QuizID            uint        `json:"quiz_id" gorm:"not null"`
	Quiz              Quiz        `json:"quiz" gorm:"foreignKey:QuizID"`
	QuestionText      string      `json:"question_text" gorm:"not null;type:text"`
	MediaAttachmentID *uint       `json:"media_attachment_id" gorm:"comment:Optional media attachment for the question"`
	MediaAttachment   *Attachment `json:"media_attachment" gorm:"foreignKey:MediaAttachmentID"`
	PointsMultiplier  int         `json:"points_multiplier" gorm:"default:1;comment:1x, 2x points"`
	TimeLimitSeconds  int         `json:"time_limit_seconds" gorm:"default:30"`
	OrderIndex        int         `json:"order_index" gorm:"comment:To keep questions in order 1, 2, 3..."`
}
