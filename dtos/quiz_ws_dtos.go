package dtos

import "time"

type CreateGameRequestDto struct {
	ClassMaterialID uint `json:"class_material_id" binding:"required"`
}

type CreateGameResponseDto struct {
	Message string `json:"message"`
	GamePIN string `json:"game_pin"`
}

type JoinGameRequestDto struct {
	GamePIN string `json:"game_pin" binding:"required"`
}

type JoinGameResponseDto struct {
	Message string `json:"message"`
	Role    string `json:"role"` // Moderator or Player
}

type GameInfoRequestDto struct {
	GamePIN string `json:"game_pin"`
}

type GameInfoResponseDto struct {
	PIN             string `json:"pin"`
	Status          string `json:"status"`         // waiting|running|finished
	QuestionState   string `json:"question_state"` // hold|answering|time_up|revealed
	QuizTitle       string `json:"quiz_title"`
	CurrentQuestion int    `json:"current_question"`
	TotalQuestions  int    `json:"total_questions"`
	TotalPlayers    int    `json:"total_players"`

	QuestionStartedAt time.Time `json:"question_started_at"`
	QuestionEndsAt    time.Time `json:"question_ends_at"`
}

type WaitingRoomPlayersDto struct {
	PlayerCount int         `json:"player_count"`
	Players     []PlayerDto `json:"players"`
}

type PlayerDto struct {
	UserID  uint   `json:"user_id"`
	Name    string `json:"name"`
	ImgURL  string `json:"img_url"`
	Score   int    `json:"score"`
	Correct int    `json:"correct"`
	Streak  int    `json:"streak"`
}

type UpdateGameStatusRequestDto struct {
	Status string `json:"status" binding:"required"`
	PIN    string `json:"pin" binding:"required"`
}

type UpdateGameStatusResponseDto struct {
	Message string `json:"message"`
}

type WSQuizQuestionDto struct {
	QuestionText     string                    `json:"question_text"`
	MediaAttachment  *GetAttachmentResponseDto `json:"media_attachment,omitempty"`
	PointsMultiplier int                       `json:"points_multiplier"`
	TimeLimitSeconds int                       `json:"time_limit_seconds"`
	OrderIndex       int                       `json:"order_index"`
	Options          []WSQuizOptionDto         `json:"options"`
}

type WSQuizOptionDto struct {
	ID         uint   `json:"id"`
	OptionText string `json:"option_text"`
	OrderIndex int    `json:"order_index"`
}

type StudentAnswerRequestDto struct {
	PIN      string `json:"pin" binding:"required"`
	OptionID uint   `json:"option_id" binding:"required"`
}

type StudentAnswerResponseDto struct {
	Message string `json:"message"`
}

type AnsweringStatasDto struct {
	OptionACount int `json:"option_a_count"`
	OptionBCount int `json:"option_b_count"`
	OptionCCount int `json:"option_c_count"`
	OptionDCount int `json:"option_d_count"`
}

type RevealedAnswerDto struct {
	CorrectOptionID uint        `json:"correct_option_id"`
	Leaderboard     []PlayerDto `json:"leaderboard"` // Sorted by score
}

type StudentRevealedAnswerResponseDto struct {
	Points     int       `json:"points"` // points that get added to the player
	PlayerData PlayerDto `json:"player_data"`
}
