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
}

type gameUseCase struct {
	gameRepo          repositories.GameRepository
	quizGameLogRepo   repositories.QuizGameLogRepository
	classMaterialRepo repositories.ClassMaterialRepository
	classRepo         repositories.ClassRepository
}

func NewGameUseCase(
	gameRepo repositories.GameRepository,
	quizGameLogRepo repositories.QuizGameLogRepository,
	classMaterialRepo repositories.ClassMaterialRepository,
	classRepo repositories.ClassRepository,
) GameUseCase {
	return &gameUseCase{
		gameRepo:          gameRepo,
		quizGameLogRepo:   quizGameLogRepo,
		classMaterialRepo: classMaterialRepo,
		classRepo:         classRepo,
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
