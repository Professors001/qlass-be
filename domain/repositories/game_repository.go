package repositories

import (
	"context"
	"qlass-be/domain/entities"
)

type GameRepository interface {
	// 1. Teacher Session
	SetTeacherSession(ctx context.Context, teacherID uint, pin string) error
	GetTeacherSession(ctx context.Context, teacherID uint) (string, error)

	// 2. Game State & Quiz Data
	CreateGameState(ctx context.Context, pin string, state *entities.GameStateRedis) error
	GetGameState(ctx context.Context, pin string) (*entities.GameStateRedis, error)
	UpdateGameState(ctx context.Context, pin string, fields map[string]interface{}) error
	IncrementField(ctx context.Context, pin string, field string, amount int) error

	SetQuizData(ctx context.Context, pin string, quizJSON string) error
	GetQuizData(ctx context.Context, pin string) (string, error)

	// 3. Players (Lobby & Allowed)
	AddPlayerToLobby(ctx context.Context, pin string, userID uint) error
	RemovePlayerFromLobby(ctx context.Context, pin string, userID uint) error
	AddAllowedPlayer(ctx context.Context, pin string, userID uint) error
	IsPlayerAllowed(ctx context.Context, pin string, userID uint) (bool, error)
	IsPlayerInLobby(ctx context.Context, pin string, userID uint) (bool, error)
	GetLobbyPlayers(ctx context.Context, pin string) ([]uint, error)

	SavePlayerData(ctx context.Context, pin string, userID uint, data *entities.PlayerDataRedis) error
	GetPlayerData(ctx context.Context, pin string, userID uint) (*entities.PlayerDataRedis, error)
	DeletePlayerData(ctx context.Context, pin string, userID uint) error

	// 4. Answers & Logic
	// Returns true if added new, false if already existed (duplicate submit)
	MarkUserAnswered(ctx context.Context, pin string, questionIndex int, userID uint) (bool, error)
	SaveAnswerDetail(ctx context.Context, pin string, questionIndex int, userID uint, answerLog *entities.AnswerLog) error

	// 5. Leaderboard (ZSET)
	UpdateScore(ctx context.Context, pin string, userID uint, totalScore float64) error
	GetLeaderboard(ctx context.Context, pin string, limit int) ([]entities.PlayerScore, error) // You might need a struct for this

	// Utility
	KeyExists(ctx context.Context, key string) (bool, error)
}
