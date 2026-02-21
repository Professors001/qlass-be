package usecases

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"math/rand"
	"qlass-be/domain/entities"
	"qlass-be/domain/repositories"
	"qlass-be/dtos"
	"time"
)

type GameUseCase interface {
	StartGameSession(ctx context.Context, teacherID uint, dto dtos.CreateGameRequestDto) (*dtos.CreateGameResponseDto, error)
	JoinGame(ctx context.Context, pin string, userID uint) (*dtos.LobbyUpdatePayload, error)
	LeaveGame(ctx context.Context, pin string, userID uint) (*dtos.LobbyUpdatePayload, error)
}

type gameUseCase struct {
	gameRepo          repositories.GameRepository
	quizGameLogRepo   repositories.QuizGameLogRepository
	classMaterialRepo repositories.ClassMaterialRepository
	classRepo         repositories.ClassRepository
	userRepo          repositories.UserRepository
}

func NewGameUseCase(
	gameRepo repositories.GameRepository,
	quizGameLogRepo repositories.QuizGameLogRepository,
	classMaterialRepo repositories.ClassMaterialRepository,
	classRepo repositories.ClassRepository,
	userRepo repositories.UserRepository,
) GameUseCase {
	return &gameUseCase{
		gameRepo:          gameRepo,
		quizGameLogRepo:   quizGameLogRepo,
		classMaterialRepo: classMaterialRepo,
		classRepo:         classRepo,
		userRepo:          userRepo,
	}
}

func (u *gameUseCase) StartGameSession(ctx context.Context, teacherID uint, dto dtos.CreateGameRequestDto) (*dtos.CreateGameResponseDto, error) {
	// 1. Get Class Material to check ownership
	material, err := u.classMaterialRepo.GetByID(dto.ClassMaterialID)
	if err != nil {
		return nil, err
	}

	class, err := u.classRepo.GetByID(material.ClassID)
	if err != nil {
		return nil, err
	}

	if class.OwnerID != teacherID {
		return nil, errors.New("unauthorized: only class owner can start the game")
	}

	// 2. Get Quiz Game Log
	logs, err := u.quizGameLogRepo.GetByClassMaterialID(material.ID)
	if err != nil || len(logs) == 0 {
		return nil, errors.New("quiz game log not found")
	}
	gameLog := logs[0]

	// 3. Parse Snapshot
	var quizSnapshot dtos.GetQuizResponseDto
	if err := json.Unmarshal(gameLog.QuizSnapshot, &quizSnapshot); err != nil {
		return nil, fmt.Errorf("failed to parse quiz snapshot: %v", err)
	}

	// 4. Generate Unique PIN
	rand.Seed(time.Now().UnixNano())
	pin := ""
	for {
		pin = fmt.Sprintf("%06d", rand.Intn(1000000))
		exists, err := u.gameRepo.KeyExists(ctx, fmt.Sprintf("game:%s:state", pin))
		if err != nil {
			return nil, err
		}
		if !exists {
			break
		}
	}

	// 5. Initialize Redis Game State
	gameState := &entities.GameStateRedis{
		Pin:             pin,
		Status:          "waiting",
		QuestionState:   "hold",
		ClassMaterialID: material.ID,
		QuizTitle:       quizSnapshot.Title,
		HostID:          teacherID,
		CurrentQuestion: 0,
		TotalQuestions:  len(quizSnapshot.Questions),
		TotalPlayers:    0,
	}

	if err := u.gameRepo.CreateGameState(ctx, pin, gameState); err != nil {
		return nil, err
	}

	// 6. Store Quiz Data in Redis
	snapshotJSON, _ := json.Marshal(quizSnapshot)
	if err := u.gameRepo.SetQuizData(ctx, pin, string(snapshotJSON)); err != nil {
		return nil, err
	}

	// 7. Set Teacher Session
	if err := u.gameRepo.SetTeacherSession(ctx, teacherID, pin); err != nil {
		return nil, err
	}

	// 8. Update DB
	gameLog.QuizPin = pin
	gameLog.Status = "waiting"
	if err := u.quizGameLogRepo.Update(gameLog); err != nil {
		return nil, err
	}

	return &dtos.CreateGameResponseDto{
		Message: "Game session started",
		GamePIN: pin,
	}, nil
}

func (u *gameUseCase) JoinGame(ctx context.Context, pin string, userID uint) (*dtos.LobbyUpdatePayload, error) {
	// 1. Check Game Exists
	exists, err := u.gameRepo.KeyExists(ctx, fmt.Sprintf("game:%s:state", pin))
	if err != nil || !exists {
		return nil, errors.New("game session not found")
	}

	// 2. Get User Details
	user, err := u.userRepo.GetByID(userID)
	if err != nil {
		return nil, err
	}

	// 3. Save to Redis
	playerData := &entities.PlayerDataRedis{
		Name:      user.FirstName + " " + user.LastName,
		AvatarURL: user.ProfileImgURL,
		Score:     0,
		Correct:   0,
		Streak:    0,
	}
	if err := u.gameRepo.SavePlayerData(ctx, pin, userID, playerData); err != nil {
		return nil, err
	}
	if err := u.gameRepo.AddPlayerToLobby(ctx, pin, userID); err != nil {
		return nil, err
	}

	// 4. Get Current Lobby State for Broadcast
	playerIDs, err := u.gameRepo.GetLobbyPlayers(ctx, pin)
	if err != nil {
		return nil, err
	}

	var players []dtos.PlayerDto
	for _, pid := range playerIDs {
		pData, _ := u.gameRepo.GetPlayerData(ctx, pin, pid)
		if pData != nil {
			players = append(players, dtos.PlayerDto{
				UserID:    pid,
				Name:      pData.Name,
				AvatarURL: pData.AvatarURL,
				Score:     pData.Score,
			})
		}
	}

	return &dtos.LobbyUpdatePayload{
		PlayerCount: len(players),
		Players:     players,
		NewPlayer: &dtos.PlayerDto{
			UserID:    user.ID,
			Name:      playerData.Name,
			AvatarURL: playerData.AvatarURL,
		},
	}, nil
}

func (u *gameUseCase) LeaveGame(ctx context.Context, pin string, userID uint) (*dtos.LobbyUpdatePayload, error) {
	// 1. Remove from Redis Set
	if err := u.gameRepo.RemovePlayerFromLobby(ctx, pin, userID); err != nil {
		return nil, err
	}

	// 2. Get updated list for broadcast
	playerIDs, err := u.gameRepo.GetLobbyPlayers(ctx, pin)
	if err != nil {
		return nil, err
	}

	var players []dtos.PlayerDto
	for _, pid := range playerIDs {
		pData, _ := u.gameRepo.GetPlayerData(ctx, pin, pid)
		if pData != nil {
			players = append(players, dtos.PlayerDto{
				UserID:    pid,
				Name:      pData.Name,
				AvatarURL: pData.AvatarURL,
				Score:     pData.Score,
			})
		}
	}

	return &dtos.LobbyUpdatePayload{
		PlayerCount: len(players),
		Players:     players,
	}, nil
}
