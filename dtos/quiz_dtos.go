package dtos

type CreateQuizDto struct {
	ClassMaterialID        uint                 `json:"class_material_id" binding:"required"`
	Title                  string               `json:"title" binding:"required"`
	Description            string               `json:"description"`
	DefaultTimePerQuestion int                  `json:"default_time_per_question" binding:"required,min=10"`
	Questions              []CreateQuizQuestion `json:"questions" binding:"required,dive"`
}

type CreateQuizQuestion struct {
	QuestionText      string             `json:"question_text" binding:"required"`
	Options           []CreateQuizOption `json:"options" binding:"required,dive"`
	MediaAttachmentID *uint              `json:"media_attachment_id"`
	PointsMultiplier  int                `json:"points_multiplier" binding:"required,min=1"`
	TimeLimitSeconds  int                `json:"time_limit_seconds" binding:"required,min=10"`
	OrderIndex        int                `json:"order_index" binding:"required,min=1"`
}

type CreateQuizOption struct {
	OptionText string `json:"option_text" binding:"required"`
	IsCorrect  bool   `json:"is_correct"`
	OrderIndex int    `json:"order_index" binding:"required,min=1"`
}

type GetQuizResponseDto struct {
	ID                     uint                      `json:"id"`
	ClassMaterialID        uint                      `json:"class_material_id"`
	Title                  string                    `json:"title"`
	Description            string                    `json:"description"`
	DefaultTimePerQuestion int                       `json:"default_time_per_question"`
	Questions              []GetQuizQuestionResponse `json:"questions"`
}

type GetQuizQuestionResponse struct {
	ID               uint                      `json:"id"`
	QuestionText     string                    `json:"question_text"`
	MediaAttachment  *GetAttachmentResponseDto `json:"media_attachment,omitempty"`
	PointsMultiplier int                       `json:"points_multiplier"`
	TimeLimitSeconds int                       `json:"time_limit_seconds"`
	OrderIndex       int                       `json:"order_index"`
	Options          []GetQuizOptionResponse   `json:"options"`
}

type GetQuizOptionResponse struct {
	ID         uint   `json:"id"`
	OptionText string `json:"option_text"`
	IsCorrect  bool   `json:"is_correct"`
	OrderIndex int    `json:"order_index"`
}
