package dtos

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
}

type WaitingRoomStatusDto struct {
	Status      string      `json:"status"`
	PlayerCount int         `json:"player_count"`
	Players     []PlayerDto `json:"players"`
}

type PlayerDto struct {
	UserID uint   `json:"user_id"`
	Name   string `json:"name"`
	ImgURL string `json:"img_url"`
	Score  int    `json:"score"`
}

type UpdateGameStatusRequestDto struct {
	Status string `json:"status" binding:"required"`
	PIN    string `json:"pin" binding:"required"`
}

type UpdateGameStatusResponseDto struct {
	Message string `json:"message"`
}

