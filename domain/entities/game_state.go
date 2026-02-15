package entities

// Maps to HASH: game:{pin}:state
type GameStateRedis struct {
	Pin             string `redis:"pin"`
	Status          string `redis:"status"`         // waiting|running|finished
	QuestionState   string `redis:"question_state"` // hold|answering|time_up|revealed
	ClassMaterialID uint   `redis:"class_material_id"`
	QuizTitle       string `redis:"quiz_title"`
	HostID          uint   `redis:"host_id"`
	CurrentQuestion int    `redis:"current_question"`
	TotalQuestions  int    `redis:"total_questions"`
	TotalPlayers    int    `redis:"total_players"`

	// Timestamps (Unix Milliseconds or RFC3339 string depending on preference)
	QuestionStartedAt int64 `redis:"question_started_at"`
	QuestionEndsAt    int64 `redis:"question_ends_at"`

	// Current Question Stats (Reset every question)
	CorrectOptionID int `redis:"correct_option_id"`
	OptionACount    int `redis:"option_a_count"`
	OptionBCount    int `redis:"option_b_count"`
	OptionCCount    int `redis:"option_c_count"`
	OptionDCount    int `redis:"option_d_count"`
}

// Maps to HASH: game:{pin}:player:{user_id}
type PlayerDataRedis struct {
	Name     string `redis:"name"`
	Avatar   string `redis:"avatar"`
	Score    int    `redis:"score"`
	Correct  int    `redis:"correct"`
	Streak   int    `redis:"streak"`
	IsOnline bool   `redis:"is_online"` // Optional: for UI status
}

// PlayerScore represents a single row in the leaderboard
type PlayerScore struct {
	UserID   uint   `json:"user_id"`
	Username string `json:"username"`
	Avatar   string `json:"avatar"`
	Score    int    `json:"score"`
	Rank     int    `json:"rank"`
}

// Maps to Value inside HASH: game:{pin}:answers:{q_index}
// Field: user_id, Value: JSON(AnswerLog)
type AnswerLog struct {
	OptionID  int  `json:"opt_id"`
	TimeMs    int  `json:"time_ms"`
	Points    int  `json:"points"`
	IsCorrect bool `json:"correct"`
}
