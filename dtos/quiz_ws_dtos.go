package dtos

type CreateGameRequestDto struct {
	ClassMaterialID uint `json:"class_material_id" binding:"required"`
}
type CreateGameResponseDto struct {
	Message string `json:"message"`
	GamePIN string `json:"game_pin"`
}

type JoinGameResponseDto struct {
	Message string `json:"message"`
	Role    string `json:"role"` // "host" or "player"
}

type JoinGameRequestDto struct {
	GamePIN string `json:"game_pin" binding:"required"`
}

type UpdateGameStatusRequestDto struct {
	PIN    string `json:"pin" binding:"required"`
	Status string `json:"status" binding:"required,oneof=start next_question reveal_answer end"`
}

type StudentAnswerRequestDto struct {
	PIN      string `json:"pin" binding:"required"`
	OptionID uint   `json:"option_id" binding:"required"`
}

type StudentAnswerResponseDto struct {
	Message string `json:"message"`
}

type GameSyncResponseDto struct {
	PIN           string `json:"pin"`
	QuizTitle     string `json:"quiz_title"`
	Role          string `json:"role"`
	Status        string `json:"status"`
	QuestionState string `json:"question_state"`

	ServerTimeMs      int64 `json:"server_time_ms"`
	QuestionStartedAt int64 `json:"question_started_at,omitempty"`
	QuestionEndsAt    int64 `json:"question_ends_at,omitempty"`

	LobbyStateObject    *LobbyStateDto    `json:"lobby_state_object,omitempty"`
	QuestionStateObject *QuestionStateDto `json:"question_state_object,omitempty"`
	ResultStateObject   *ResultStateDto   `json:"result_state_object,omitempty"`
	Leaderboard         []PlayerDto       `json:"leaderboard,omitempty"`

	MyState *MyPlayerStateDto `json:"my_state,omitempty"`
}

type LobbyStateDto struct {
	TotalPlayers int         `json:"total_players"`
	Players      []PlayerDto `json:"players"`
}

type QuestionStateDto struct {
	CurrentQuestion  int               `json:"current_question"`
	TotalQuestions   int               `json:"total_questions"`
	TimeLimitSeconds int               `json:"time_limit_seconds"`
	QuestionText     string            `json:"question_text"`
	ImageURL         string            `json:"image_url,omitempty"`
	PointsMultiplier int               `json:"points_multiplier"`
	Options          []WSQuizOptionDto `json:"options"`

	AnsweredCount int `json:"answered_count,omitempty"`
	TotalPlayers  int `json:"total_players,omitempty"`
}

type ResultStateDto struct {
	CorrectOptionID uint             `json:"correct_option_id"`
	Stats           LiveStatsPayload `json:"stats"`
	Leaderboard     []PlayerDto      `json:"leaderboard,omitempty"`
}

type MyPlayerStateDto struct {
	UserID    uint   `json:"user_id"`
	Name      string `json:"name"`
	AvatarURL string `json:"avatar_url"`
	Score     int    `json:"score"`
	Rank      int    `json:"rank"`
	Streak    int    `json:"streak"`

	HasAnswered      bool `json:"has_answered"`
	SelectedOptionID uint `json:"selected_option_id,omitempty"`

	LastResult *StudentPersonalResultPayload `json:"last_result,omitempty"`
}

type WSEventDto struct {
	Type    string      `json:"type"`
	Payload interface{} `json:"payload"`
}

type AnswerCountUpdatePayload struct {
	AnsweredCount int `json:"answered_count"`
	TotalPlayers  int `json:"total_players"`
}

type PlayerDto struct {
	UserID    uint   `json:"user_id"`
	Name      string `json:"name"`
	AvatarURL string `json:"avatar_url"`
	Score     int    `json:"score"`
	Rank      int    `json:"rank,omitempty"`
	Streak    int    `json:"streak,omitempty"`
}

type WSQuizOptionDto struct {
	ID         uint   `json:"id"`
	OptionText string `json:"option_text"`
	Label      string `json:"label"` // "A", "B", "C", "D"
}

type LiveStatsPayload struct {
	OptionACount int  `json:"option_a_count"`
	OptionBCount int  `json:"option_b_count"`
	OptionCCount int  `json:"option_c_count"`
	OptionDCount int  `json:"option_d_count"`
	OptionAID    uint `json:"option_a_id"`
	OptionBID    uint `json:"option_b_id"`
	OptionCID    uint `json:"option_c_id"`
	OptionDID    uint `json:"option_d_id"`
}

type StudentPersonalResultPayload struct {
	IsCorrect    bool `json:"is_correct"`
	PointsEarned int  `json:"points_earned"`
	TotalScore   int  `json:"total_score"`
	Streak       int  `json:"streak"`
	Rank         int  `json:"rank"`
}
