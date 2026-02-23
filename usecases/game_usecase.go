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
	StartGame(ctx context.Context, pin string, hostID uint) error
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
	// 1. Check Game State
	state, err := u.gameRepo.GetGameState(ctx, pin)
	if err != nil {
		return nil, errors.New("game session not found")
	}

	if state.Status != "waiting" {
		inLobby, err := u.gameRepo.IsPlayerInLobby(ctx, pin, userID)
		if err != nil || !inLobby {
			return nil, errors.New("game has already started")
		}
	}

	// 2. Get User Details
	user, err := u.userRepo.GetByID(userID)
	if err != nil {
		return nil, err
	}

	// 3. Save to Redis (Preserve stats if reconnecting)
	existingData, _ := u.gameRepo.GetPlayerData(ctx, pin, userID)

	playerData := &entities.PlayerDataRedis{
		Name:      user.FirstName + " " + user.LastName,
		AvatarURL: user.ProfileImgURL,
		Score:     0,
		Correct:   0,
		Streak:    0,
	}

	// If player already exists (reconnect), keep their score
	if existingData != nil && existingData.Name != "" {
		playerData.Score = existingData.Score
		playerData.Correct = existingData.Correct
		playerData.Streak = existingData.Streak
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
	// 1. Check Game State
	state, err := u.gameRepo.GetGameState(ctx, pin)
	if err != nil {
		return nil, errors.New("game session not found")
	}

	// 2. Only remove player if game is waiting (Lock list if running)
	if state.Status == "waiting" {
		if err := u.gameRepo.RemovePlayerFromLobby(ctx, pin, userID); err != nil {
			return nil, err
		}
		_ = u.gameRepo.DeletePlayerData(ctx, pin, userID)
	}

	// 3. Get updated list for broadcast
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

func (u *gameUseCase) StartGame(ctx context.Context, pin string, hostID uint) error {
	state, err := u.gameRepo.GetGameState(ctx, pin)
	if err != nil {
		return err
	}

	if state.HostID != hostID {
		return errors.New("unauthorized: only host can start the game")
	}

	// Update status to running
	return u.gameRepo.UpdateGameState(ctx, pin, map[string]interface{}{
		"status": "running",
	})
}

func (u *gameUseCase) NextQuestion(ctx context.Context, pin string, hostID uint) (*dtos.QuestionPayload, error) {
	// 1. Get current Game State
	state, err := u.gameRepo.GetGameState(ctx, pin)
	if err != nil {
		return nil, errors.New("game session not found")
	}

	// 2. Security Check: Only the Host can change the question
	if state.HostID != hostID {
		return nil, errors.New("unauthorized: only host can control the game")
	}

	// 3. Get the Quiz Snapshot from Redis
	quizDataJSON, err := u.gameRepo.GetQuizData(ctx, pin)
	if err != nil {
		return nil, errors.New("failed to load quiz data")
	}

	var quizSnapshot dtos.GetQuizResponseDto
	if err := json.Unmarshal([]byte(quizDataJSON), &quizSnapshot); err != nil {
		return nil, errors.New("failed to parse quiz data")
	}

	// 4. Increment Question Index
	state.CurrentQuestion++
	if state.CurrentQuestion > state.TotalQuestions {
		return nil, errors.New("no more questions available")
	}

	// 5. Extract the specific question (0-indexed array, so minus 1)
	currentQ := quizSnapshot.Questions[state.CurrentQuestion-1]

	// 6. Set Timers and Update State
	now := time.Now()
	endsAt := now.Add(time.Duration(currentQ.TimeLimitSeconds) * time.Second)

	state.Status = "running"
	state.QuestionState = "answering"
	state.QuestionStartedAt = now
	state.QuestionEndsAt = endsAt

	// Save the new state back to Redis
	// (Assuming you have a method like SaveGameState or UpdateGameState)
	if err := u.gameRepo.UpdateGameState(ctx, pin, map[string]interface{}{
		"status":              state.Status,
		"question_state":      state.QuestionState,
		"current_question":    state.CurrentQuestion,
		"question_started_at": state.QuestionStartedAt,
		"question_ends_at":    state.QuestionEndsAt,
	}); err != nil {
		return nil, err
	}

	// 7. Map options for the payload (hide which one is correct!)
	var options []dtos.WSQuizOptionDto
	for _, opt := range currentQ.Options {
		var label string
		switch opt.OrderIndex {
		case 1:
			label = "A"
		case 2:
			label = "B"
		case 3:
			label = "C"
		case 4:
			label = "D"
		}
		options = append(options, dtos.WSQuizOptionDto{
			ID:         opt.ID,
			OptionText: opt.OptionText,
			Label:      label,
		})
	}

	// 8. Return the Payload
	return &dtos.QuestionPayload{
		QuestionIndex:    state.CurrentQuestion,
		TotalQuestions:   state.TotalQuestions,
		TimeLimitSeconds: currentQ.TimeLimitSeconds,
		QuestionText:     currentQ.QuestionText,
		PointsMultiplier: currentQ.PointsMultiplier,
		Options:          options,
	}, nil
}
