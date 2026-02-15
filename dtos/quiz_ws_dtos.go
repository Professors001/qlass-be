package dtos

// ==========================================
// 1. HTTP REST (Requests & Responses)
// ==========================================

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
	Role    string `json:"role"`  // "host" or "player"
}

// Used by Host to control the game flow
type UpdateGameStatusRequestDto struct {
	PIN    string `json:"pin" binding:"required"`
	Status string `json:"status" binding:"required,oneof=start next_question reveal_answer end"`
}

type UpdateGameStatusResponseDto struct {
	Message string `json:"message"`
}

// Used by Student to submit an answer
type StudentAnswerRequestDto struct {
	PIN      string `json:"pin" binding:"required"`
	OptionID uint   `json:"option_id" binding:"required"`
}

type StudentAnswerResponseDto struct {
	Message string `json:"message"`
}

// For Reconnecting / Refreshing Page
type GameInfoResponseDto struct {
	PIN           string `json:"pin"`
	Status        string `json:"status"`         // waiting|running|finished
	QuestionState string `json:"question_state"` // hold|answering|revealed
	QuizTitle     string `json:"quiz_title"`

	CurrentQuestion int `json:"current_question"`
	TotalQuestions  int `json:"total_questions"`
	TotalPlayers    int `json:"total_players"`

	// Unix Milliseconds for easy JS countdown sync
	QuestionStartedAt int64 `json:"question_started_at"`
	QuestionEndsAt    int64 `json:"question_ends_at"`
}

// ==========================================
// 2. WEBSOCKET ENVELOPE (The Container)
// ==========================================

// WSEventDto is the ONLY struct sent over the socket.
// The 'Payload' changes based on the 'Type'.
type WSEventDto struct {
	Type    string      `json:"type"`    // e.g. "LOBBY_UPDATE", "NEXT_QUESTION", "ROUND_RESULT"
	Payload interface{} `json:"payload"` // One of the Payload structs below
}

// ==========================================
// 3. WEBSOCKET PAYLOADS (The Data)
// ==========================================

// --- LOBBY PHASE ---
type LobbyUpdatePayload struct {
	PlayerCount int         `json:"player_count"`
	Players     []PlayerDto `json:"players"`              // List of avatars in lobby
	NewPlayer   *PlayerDto  `json:"new_player,omitempty"` // Optimization: Only send who just joined
}

// --- QUESTION PHASE (State: "answering") ---
type QuestionPayload struct {
	QuestionIndex    int    `json:"question_index"` // 1-based index
	TotalQuestions   int    `json:"total_questions"`
	TimeLimitSeconds int    `json:"time_limit_seconds"`
	QuestionText     string `json:"question_text"`
	ImageURL         string `json:"image_url,omitempty"`
	PointsMultiplier int    `json:"points_multiplier"`

	// The 4 Choices (A, B, C, D)
	Options []WSQuizOptionDto `json:"options"`
}

// --- LIVE STATS PHASE (State: "answering" - Periodic updates) ---
type LiveStatsPayload struct {
	TotalPlayers  int `json:"total_players"`
	AnsweredCount int `json:"answered_count"`
	// Simple bar chart data
	OptionACount int `json:"option_a_count"`
	OptionBCount int `json:"option_b_count"`
	OptionCCount int `json:"option_c_count"`
	OptionDCount int `json:"option_d_count"`
}

// --- RESULT PHASE (State: "revealed") ---

// 1. Broadcast to Everyone (Host & Players see this for the big screen)
type RoundResultPayload struct {
	CorrectOptionID uint             `json:"correct_option_id"`
	Stats           LiveStatsPayload `json:"stats"`       // Final bar chart
	Leaderboard     []PlayerDto      `json:"leaderboard"` // Top 5
}

// 2. Sent Private to Specific Student (Pop-up on their phone)
type StudentPersonalResultPayload struct {
	IsCorrect    bool `json:"is_correct"`
	PointsEarned int  `json:"points_earned"`
	TotalScore   int  `json:"total_score"`
	Streak       int  `json:"streak"`
	Rank         int  `json:"rank"`
}

// --- GAME OVER PHASE ---
type GameOverPayload struct {
	Winner PlayerDto   `json:"winner"`
	Top3   []PlayerDto `json:"top_3"`
	MyRank int         `json:"my_rank,omitempty"`
}

// ==========================================
// 4. SHARED STRUCTS
// ==========================================

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
