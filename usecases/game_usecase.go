package usecases

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"math"
	"math/rand"
	"qlass-be/domain/entities"
	"qlass-be/domain/repositories"
	"qlass-be/dtos"
	"time"
)

type GameUseCase interface {
	StartGameSession(ctx context.Context, teacherID uint, dto dtos.CreateGameRequestDto) (*dtos.CreateGameResponseDto, error)
	JoinGame(ctx context.Context, pin string, userID uint) (*dtos.GameInfoResponseDto, *dtos.LobbyUpdatePayload, error)
	LeaveGame(ctx context.Context, pin string, userID uint) (*dtos.LobbyUpdatePayload, error)
	StartGame(ctx context.Context, pin string, hostID uint) error
	NextStep(ctx context.Context, pin string, hostID uint) (*dtos.WSEventDto, error)
	TimeoutQuestion(ctx context.Context, pin string, questionIndex int) (*dtos.WSEventDto, error)
	SubmitAnswer(ctx context.Context, pin string, userID uint, optionID uint) (*dtos.StudentAnswerResponseDto, *dtos.LiveStatsPayload, error)
	GetCurrentQuestion(ctx context.Context, pin string) (*dtos.QuestionPayload, error)
	GetLiveStats(ctx context.Context, pin string) (*dtos.LiveStatsPayload, error)
	HasUserAnswered(ctx context.Context, pin string, userID uint) (bool, error)
	GetRoundResult(ctx context.Context, pin string) (*dtos.RoundResultPayload, error)
	StoreGameData(ctx context.Context, pin string) error
}

type gameUseCase struct {
	gameRepo                repositories.GameRepository
	quizGameLogRepo         repositories.QuizGameLogRepository
	classMaterialRepo       repositories.ClassMaterialRepository
	classRepo               repositories.ClassRepository
	userRepo                repositories.UserRepository
	submissionRepo          repositories.SubmissionRepository
	quizStudentResponseRepo repositories.QuizStudentResponseRepository
}

func NewGameUseCase(
	gameRepo repositories.GameRepository,
	quizGameLogRepo repositories.QuizGameLogRepository,
	classMaterialRepo repositories.ClassMaterialRepository,
	classRepo repositories.ClassRepository,
	userRepo repositories.UserRepository,
	submissionRepo repositories.SubmissionRepository,
	quizStudentResponseRepo repositories.QuizStudentResponseRepository,
) GameUseCase {
	return &gameUseCase{
		gameRepo:                gameRepo,
		quizGameLogRepo:         quizGameLogRepo,
		classMaterialRepo:       classMaterialRepo,
		classRepo:               classRepo,
		userRepo:                userRepo,
		submissionRepo:          submissionRepo,
		quizStudentResponseRepo: quizStudentResponseRepo,
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

func (u *gameUseCase) JoinGame(ctx context.Context, pin string, userID uint) (*dtos.GameInfoResponseDto, *dtos.LobbyUpdatePayload, error) {
	// 1. Check Game State
	state, err := u.gameRepo.GetGameState(ctx, pin)
	if err != nil {
		return nil, nil, errors.New("game session not found")
	}

	inLobby, err := u.gameRepo.IsPlayerInLobby(ctx, pin, userID)
	if err != nil {
		return nil, nil, err
	}

	if state.Status != "waiting" {
		if !inLobby {
			return nil, nil, errors.New("game has already started")
		}
	}

	// 2. Get User Details
	user, err := u.userRepo.GetByID(userID)
	if err != nil {
		return nil, nil, err
	}

	// 3. Save to Redis (Preserve stats if reconnecting)
	existingData, _ := u.gameRepo.GetPlayerData(ctx, pin, userID)

	playerData := &entities.PlayerDataRedis{
		Name:      user.FirstName + " " + user.LastName,
		AvatarURL: "----",
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
		return nil, nil, err
	}
	if err := u.gameRepo.AddPlayerToLobby(ctx, pin, userID); err != nil {
		return nil, nil, err
	}

	// 4. Get Current Lobby State for Broadcast
	playerIDs, err := u.gameRepo.GetLobbyPlayers(ctx, pin)
	if err != nil {
		return nil, nil, err
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

	gameInfo := &dtos.GameInfoResponseDto{
		PIN:               state.Pin,
		Status:            state.Status,
		QuestionState:     state.QuestionState,
		QuizTitle:         state.QuizTitle,
		CurrentQuestion:   state.CurrentQuestion,
		TotalQuestions:    state.TotalQuestions,
		TotalPlayers:      len(players),
		QuestionStartedAt: state.QuestionStartedAt.UnixMilli(),
		QuestionEndsAt:    state.QuestionEndsAt.UnixMilli(),
	}

	var lobbyUpdate *dtos.LobbyUpdatePayload
	if !inLobby {
		lobbyUpdate = &dtos.LobbyUpdatePayload{
			PlayerCount: len(players),
			Players:     players,
			NewPlayer: &dtos.PlayerDto{
				UserID:    user.ID,
				Name:      playerData.Name,
				AvatarURL: playerData.AvatarURL,
			},
		}
	}

	return gameInfo, lobbyUpdate, nil
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
	} else {
		return nil, nil
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

func (u *gameUseCase) NextStep(ctx context.Context, pin string, hostID uint) (*dtos.WSEventDto, error) {
	// 1. Get current Game State
	state, err := u.gameRepo.GetGameState(ctx, pin)
	if err != nil {
		return nil, errors.New("game session not found")
	}

	// 2. Security Check
	if state.HostID != hostID {
		return nil, errors.New("unauthorized: only host can control the game")
	}

	if state.Status == "finished" {
		return nil, errors.New("game is already finished")
	}

	// 3. Handle State Transitions
	// Case 1: Waiting -> Start Game (Go to Q1 Hold)
	if state.Status == "waiting" {
		if err := u.gameRepo.UpdateGameState(ctx, pin, map[string]interface{}{
			"status":           "running",
			"current_question": 1,
			"question_state":   "hold",
			// Reset counters
			"option_a_count": 0, "option_b_count": 0, "option_c_count": 0, "option_d_count": 0,
		}); err != nil {
			return nil, err
		}
		return &dtos.WSEventDto{
			Type: "GET_READY",
			Payload: map[string]interface{}{
				"question_index":  1,
				"total_questions": state.TotalQuestions,
			},
		}, nil
	}

	// Case 2: Running - Cycle through question states
	switch state.QuestionState {
	case "hold":
		// Hold -> Answering (Start Question)
		quizDataJSON, err := u.gameRepo.GetQuizData(ctx, pin)
		if err != nil {
			return nil, err
		}
		var quizSnapshot dtos.GetQuizResponseDto
		json.Unmarshal([]byte(quizDataJSON), &quizSnapshot)

		currentQ := quizSnapshot.Questions[state.CurrentQuestion-1]
		now := time.Now()
		endsAt := now.Add(time.Duration(currentQ.TimeLimitSeconds) * time.Second)

		updates := map[string]interface{}{
			"question_state":      "answering",
			"question_started_at": now,
			"question_ends_at":    endsAt,
		}
		// Map Option IDs for counters
		for _, opt := range currentQ.Options {
			switch opt.OrderIndex {
			case 1:
				updates["option_a_id"] = opt.ID
			case 2:
				updates["option_b_id"] = opt.ID
			case 3:
				updates["option_c_id"] = opt.ID
			case 4:
				updates["option_d_id"] = opt.ID
			}
		}
		if err := u.gameRepo.UpdateGameState(ctx, pin, updates); err != nil {
			return nil, err
		}

		// Build Options Payload
		var options []dtos.WSQuizOptionDto
		for _, opt := range currentQ.Options {
			label := ""
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

		return &dtos.WSEventDto{
			Type: "NEXT_QUESTION",
			Payload: dtos.QuestionPayload{
				QuestionIndex:    state.CurrentQuestion,
				TotalQuestions:   state.TotalQuestions,
				TimeLimitSeconds: currentQ.TimeLimitSeconds,
				QuestionText:     currentQ.QuestionText,
				ImageURL:         currentQ.MediaAttachment.FileURL,
				PointsMultiplier: currentQ.PointsMultiplier,
				Options:          options,
			},
		}, nil

	case "answering":
		// Answering -> Time Up (Show Stats/Correct Answer)
		return u.finishQuestion(ctx, pin, state)

	case "time_up":
		// Time Up -> Revealed (Show Leaderboard)
		if err := u.gameRepo.UpdateGameState(ctx, pin, map[string]interface{}{
			"question_state": "revealed",
		}); err != nil {
			return nil, err
		}

		// Get Correct Option & Stats
		quizDataJSON, err := u.gameRepo.GetQuizData(ctx, pin)
		if err != nil {
			return nil, err
		}
		var quizSnapshot dtos.GetQuizResponseDto
		if err := json.Unmarshal([]byte(quizDataJSON), &quizSnapshot); err != nil {
			return nil, err
		}
		currentQ := quizSnapshot.Questions[state.CurrentQuestion-1]

		var correctOptionID uint
		for _, o := range currentQ.Options {
			if o.IsCorrect {
				correctOptionID = o.ID
				break
			}
		}

		stats := dtos.LiveStatsPayload{
			TotalPlayers:  state.TotalPlayers,
			OptionACount:  state.OptionACount,
			OptionBCount:  state.OptionBCount,
			OptionCCount:  state.OptionCCount,
			OptionDCount:  state.OptionDCount,
			AnsweredCount: state.OptionACount + state.OptionBCount + state.OptionCCount + state.OptionDCount,
			OptionAID:     state.OptionAID,
			OptionBID:     state.OptionBID,
			OptionCID:     state.OptionCID,
			OptionDID:     state.OptionDID,
		}

		top5 := u.getLeaderboardDtos(ctx, pin, 5)
		return &dtos.WSEventDto{
			Type: "ROUND_RESULT",
			Payload: dtos.RoundResultPayload{
				CorrectOptionID: correctOptionID,
				Stats:           stats,
				Leaderboard:     top5,
			},
		}, nil

	case "revealed":
		// Revealed -> Next Question (Hold) OR Game Over
		if state.CurrentQuestion >= state.TotalQuestions {
			// Game Over
			if err := u.gameRepo.UpdateGameState(ctx, pin, map[string]interface{}{
				"status": "finished",
			}); err != nil {
				return nil, err
			}

			// Log Data for future saving process
			u.StoreGameData(ctx, pin)

			top3 := u.getLeaderboardDtos(ctx, pin, 3)
			var winner dtos.PlayerDto
			if len(top3) > 0 {
				winner = top3[0]
			}

			return &dtos.WSEventDto{
				Type: "GAME_OVER",
				Payload: dtos.GameOverPayload{
					Winner: winner,
					Top3:   top3,
				},
			}, nil
		} else {
			// Next Question
			nextQ := state.CurrentQuestion + 1
			if err := u.gameRepo.UpdateGameState(ctx, pin, map[string]interface{}{
				"current_question": nextQ,
				"question_state":   "hold",
				// Reset counters
				"option_a_count": 0, "option_b_count": 0, "option_c_count": 0, "option_d_count": 0,
			}); err != nil {
				return nil, err
			}

			return &dtos.WSEventDto{
				Type: "GET_READY",
				Payload: map[string]interface{}{
					"question_index":  nextQ,
					"total_questions": state.TotalQuestions,
				},
			}, nil
		}
	}

	return nil, errors.New("unknown game state")
}

func (u *gameUseCase) TimeoutQuestion(ctx context.Context, pin string, questionIndex int) (*dtos.WSEventDto, error) {
	state, err := u.gameRepo.GetGameState(ctx, pin)
	if err != nil {
		return nil, err
	}

	// Ensure the timeout is for the current question (prevents stale timers from closing new questions)
	if state.CurrentQuestion != questionIndex {
		return nil, nil
	}

	// Only act if we are still in the answering state
	if state.QuestionState != "answering" {
		return nil, nil
	}

	return u.finishQuestion(ctx, pin, state)
}

func (u *gameUseCase) finishQuestion(ctx context.Context, pin string, state *entities.GameStateRedis) (*dtos.WSEventDto, error) {
	if err := u.gameRepo.UpdateGameState(ctx, pin, map[string]interface{}{
		"question_state": "time_up",
	}); err != nil {
		return nil, err
	}

	// Get Correct Option
	quizDataJSON, _ := u.gameRepo.GetQuizData(ctx, pin)
	var quizSnapshot dtos.GetQuizResponseDto
	json.Unmarshal([]byte(quizDataJSON), &quizSnapshot)
	currentQ := quizSnapshot.Questions[state.CurrentQuestion-1]

	var correctOptionID uint
	for _, o := range currentQ.Options {
		if o.IsCorrect {
			correctOptionID = o.ID
			break
		}
	}

	stats := dtos.LiveStatsPayload{
		TotalPlayers:  state.TotalPlayers,
		OptionACount:  state.OptionACount,
		OptionBCount:  state.OptionBCount,
		OptionCCount:  state.OptionCCount,
		OptionDCount:  state.OptionDCount,
		AnsweredCount: state.OptionACount + state.OptionBCount + state.OptionCCount + state.OptionDCount,
		OptionAID:     state.OptionAID,
		OptionBID:     state.OptionBID,
		OptionCID:     state.OptionCID,
		OptionDID:     state.OptionDID,
	}

	return &dtos.WSEventDto{
		Type: "QUESTION_TIME_UP",
		Payload: map[string]interface{}{
			"correct_option_id": correctOptionID,
			"stats":             stats,
		},
	}, nil
}

// Helper to get enriched leaderboard
func (u *gameUseCase) getLeaderboardDtos(ctx context.Context, pin string, limit int) []dtos.PlayerDto {
	leaderboard, _ := u.gameRepo.GetLeaderboard(ctx, pin, limit)
	var result []dtos.PlayerDto
	for _, p := range leaderboard {
		pData, _ := u.gameRepo.GetPlayerData(ctx, pin, p.UserID)
		name := "Unknown"
		avatar := ""
		streak := 0
		if pData != nil {
			name = pData.Name
			avatar = pData.AvatarURL
			streak = pData.Streak
		}
		result = append(result, dtos.PlayerDto{
			UserID:    p.UserID,
			Name:      name,
			AvatarURL: avatar,
			Score:     p.Score,
			Rank:      p.Rank,
			Streak:    streak,
		})
	}
	return result
}

func (u *gameUseCase) SubmitAnswer(ctx context.Context, pin string, userID uint, optionID uint) (*dtos.StudentAnswerResponseDto, *dtos.LiveStatsPayload, error) {
	// 1. Get Game State
	state, err := u.gameRepo.GetGameState(ctx, pin)
	if err != nil {
		return nil, nil, errors.New("game session not found")
	}

	if state.QuestionState != "answering" {
		return nil, nil, errors.New("question is not open for answers")
	}

	// 3. Calculate Time Taken
	now := time.Now()
	timeTaken := now.Sub(state.QuestionStartedAt).Seconds()

	// 4. Get Quiz Data
	quizDataJSON, err := u.gameRepo.GetQuizData(ctx, pin)
	if err != nil {
		return nil, nil, err
	}
	var quizSnapshot dtos.GetQuizResponseDto
	if err := json.Unmarshal([]byte(quizDataJSON), &quizSnapshot); err != nil {
		return nil, nil, err
	}

	if state.CurrentQuestion < 1 || state.CurrentQuestion > len(quizSnapshot.Questions) {
		return nil, nil, errors.New("invalid question index")
	}
	question := quizSnapshot.Questions[state.CurrentQuestion-1]

	// 5. Validate Option & Calculate Score
	var selectedOption *dtos.GetQuizOptionResponse
	for _, opt := range question.Options {
		if opt.ID == optionID {
			selectedOption = &opt
			break
		}
	}
	if selectedOption == nil {
		return nil, nil, errors.New("invalid option id")
	}

	// 2. Prevent Duplicate Answer
	isNew, err := u.gameRepo.MarkUserAnswered(ctx, pin, state.CurrentQuestion, userID)
	if err != nil {
		return nil, nil, err
	}
	if !isNew {
		return nil, nil, errors.New("already answered this question")
	}

	points := 0
	if selectedOption.IsCorrect {
		ratio := timeTaken / float64(question.TimeLimitSeconds)

		if ratio > 1 {
			ratio = 1
		}

		baseScore := (1 - (ratio / 2)) * 1000.0

		roundedScore := math.Round(baseScore)

		points = int(roundedScore) * question.PointsMultiplier
	}

	// 6. Update Player Stats
	pData, _ := u.gameRepo.GetPlayerData(ctx, pin, userID)
	if pData == nil {
		pData = &entities.PlayerDataRedis{}
	}

	if selectedOption.IsCorrect {
		pData.Score += points
		pData.Correct++
		pData.Streak++
	} else {
		pData.Streak = 0
	}

	if err := u.gameRepo.SavePlayerData(ctx, pin, userID, pData); err != nil {
		return nil, nil, err
	}
	if err := u.gameRepo.UpdateScore(ctx, pin, userID, float64(pData.Score)); err != nil {
		return nil, nil, err
	}

	// 7. Update Game Counters
	var fieldToInc string
	switch selectedOption.OrderIndex {
	case 1:
		fieldToInc = "option_a_count"
	case 2:
		fieldToInc = "option_b_count"
	case 3:
		fieldToInc = "option_c_count"
	case 4:
		fieldToInc = "option_d_count"
	}
	if fieldToInc != "" {
		u.gameRepo.IncrementField(ctx, pin, fieldToInc, 1)
	}

	// 8. Log Answer
	logEntry := &entities.AnswerLog{
		OptionID:  optionID,
		TimeMs:    int(timeTaken * 1000),
		Points:    points,
		IsCorrect: selectedOption.IsCorrect,
	}
	u.gameRepo.SaveAnswerDetail(ctx, pin, state.CurrentQuestion, userID, logEntry)

	// 9. Get Updated Stats for Broadcast
	updatedState, err := u.gameRepo.GetGameState(ctx, pin)
	var stats *dtos.LiveStatsPayload
	if err == nil {
		stats = &dtos.LiveStatsPayload{
			TotalPlayers:  updatedState.TotalPlayers,
			OptionACount:  updatedState.OptionACount,
			OptionBCount:  updatedState.OptionBCount,
			OptionCCount:  updatedState.OptionCCount,
			OptionDCount:  updatedState.OptionDCount,
			AnsweredCount: updatedState.OptionACount + updatedState.OptionBCount + updatedState.OptionCCount + updatedState.OptionDCount,
			OptionAID:     updatedState.OptionAID,
			OptionBID:     updatedState.OptionBID,
			OptionCID:     updatedState.OptionCID,
			OptionDID:     updatedState.OptionDID,
		}
	}

	return &dtos.StudentAnswerResponseDto{
		Message: "Answer submitted successfully",
	}, stats, nil
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

func (u *gameUseCase) GetCurrentQuestion(ctx context.Context, pin string) (*dtos.QuestionPayload, error) {
	state, err := u.gameRepo.GetGameState(ctx, pin)
	if err != nil {
		return nil, err
	}

	quizDataJSON, err := u.gameRepo.GetQuizData(ctx, pin)
	if err != nil {
		return nil, err
	}
	var quizSnapshot dtos.GetQuizResponseDto
	if err := json.Unmarshal([]byte(quizDataJSON), &quizSnapshot); err != nil {
		return nil, err
	}

	if state.CurrentQuestion < 1 || state.CurrentQuestion > len(quizSnapshot.Questions) {
		return nil, errors.New("invalid question index")
	}
	currentQ := quizSnapshot.Questions[state.CurrentQuestion-1]

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

	return &dtos.QuestionPayload{
		QuestionIndex:    state.CurrentQuestion,
		TotalQuestions:   state.TotalQuestions,
		TimeLimitSeconds: currentQ.TimeLimitSeconds,
		QuestionText:     currentQ.QuestionText,
		PointsMultiplier: currentQ.PointsMultiplier,
		Options:          options,
	}, nil
}

func (u *gameUseCase) GetLiveStats(ctx context.Context, pin string) (*dtos.LiveStatsPayload, error) {
	state, err := u.gameRepo.GetGameState(ctx, pin)
	if err != nil {
		return nil, err
	}
	return &dtos.LiveStatsPayload{
		TotalPlayers:  state.TotalPlayers,
		OptionACount:  state.OptionACount,
		OptionBCount:  state.OptionBCount,
		OptionCCount:  state.OptionCCount,
		OptionDCount:  state.OptionDCount,
		AnsweredCount: state.OptionACount + state.OptionBCount + state.OptionCCount + state.OptionDCount,
		OptionAID:     state.OptionAID,
		OptionBID:     state.OptionBID,
		OptionCID:     state.OptionCID,
		OptionDID:     state.OptionDID,
	}, nil
}

func (u *gameUseCase) HasUserAnswered(ctx context.Context, pin string, userID uint) (bool, error) {
	state, err := u.gameRepo.GetGameState(ctx, pin)
	if err != nil {
		return false, err
	}
	// Assuming gameRepo has a method to check set membership without adding
	return u.gameRepo.HasUserAnswered(ctx, pin, state.CurrentQuestion, userID)
}

func (u *gameUseCase) GetRoundResult(ctx context.Context, pin string) (*dtos.RoundResultPayload, error) {
	state, err := u.gameRepo.GetGameState(ctx, pin)
	if err != nil {
		return nil, err
	}

	// Get Correct Option & Stats
	quizDataJSON, err := u.gameRepo.GetQuizData(ctx, pin)
	if err != nil {
		return nil, err
	}
	var quizSnapshot dtos.GetQuizResponseDto
	if err := json.Unmarshal([]byte(quizDataJSON), &quizSnapshot); err != nil {
		return nil, err
	}

	if state.CurrentQuestion < 1 || state.CurrentQuestion > len(quizSnapshot.Questions) {
		return nil, errors.New("invalid question index")
	}
	currentQ := quizSnapshot.Questions[state.CurrentQuestion-1]

	var correctOptionID uint
	for _, o := range currentQ.Options {
		if o.IsCorrect {
			correctOptionID = o.ID
			break
		}
	}

	stats := dtos.LiveStatsPayload{
		TotalPlayers:  state.TotalPlayers,
		OptionACount:  state.OptionACount,
		OptionBCount:  state.OptionBCount,
		OptionCCount:  state.OptionCCount,
		OptionDCount:  state.OptionDCount,
		AnsweredCount: state.OptionACount + state.OptionBCount + state.OptionCCount + state.OptionDCount,
		OptionAID:     state.OptionAID,
		OptionBID:     state.OptionBID,
		OptionCID:     state.OptionCID,
		OptionDID:     state.OptionDID,
	}

	top5 := u.getLeaderboardDtos(ctx, pin, 5)
	return &dtos.RoundResultPayload{
		CorrectOptionID: correctOptionID,
		Stats:           stats,
		Leaderboard:     top5,
	}, nil
}

func (u *gameUseCase) StoreGameData(ctx context.Context, pin string) error {
	state, err := u.gameRepo.GetGameState(ctx, pin)
	if err != nil {
		log.Println("Failed to get game state for logging:", err)
		return err
	}

	// 1. Get Quiz Data (Snapshot)
	quizDataJSON, err := u.gameRepo.GetQuizData(ctx, pin)
	if err != nil {
		return err
	}
	var quizSnapshot dtos.GetQuizResponseDto
	json.Unmarshal([]byte(quizDataJSON), &quizSnapshot)

	// 2. Get Class Material & Game Log
	material, err := u.classMaterialRepo.GetByID(state.ClassMaterialID)
	if err != nil {
		return err
	}

	logs, err := u.quizGameLogRepo.GetByClassMaterialID(material.ID)
	if err != nil || len(logs) == 0 {
		return errors.New("game log not found")
	}
	gameLog := logs[0]

	// 3. Update Game Log
	now := time.Now()
	gameLog.Status = "finished"
	gameLog.FinishedAt = &now
	if err := u.quizGameLogRepo.Update(gameLog); err != nil {
		return err
	}

	// 4. Calculate Max Possible Game Score (Sum of 1000 * Multiplier for all questions)
	var maxGameScore float64 = 0
	for _, q := range quizSnapshot.Questions {
		maxGameScore += float64(1000 * q.PointsMultiplier)
	}
	if maxGameScore == 0 {
		maxGameScore = 1 // Avoid division by zero
	}

	// 5. Process Players
	playerIDs, err := u.gameRepo.GetLobbyPlayers(ctx, pin)
	if err != nil {
		return err
	}

	for _, userID := range playerIDs {
		// A. Get Player Stats
		pData, err := u.gameRepo.GetPlayerData(ctx, pin, userID)
		if err != nil {
			continue
		}

		// B. Create/Update Submission
		// Formula: (PlayerScore / MaxGameScore) * MaterialPoints
		rawScore := float64(pData.Score)
		var materialPoints float64
		if material.Points != nil {
			materialPoints = float64(*material.Points)
		}

		finalScore := (rawScore / maxGameScore) * materialPoints
		finalScoreInt := int(finalScore)

		// Check if late
		isLate := false
		if material.DueAt != nil && now.After(*material.DueAt) {
			isLate = true
		}

		submission := &entities.Submission{
			ClassMaterialID: material.ID,
			UserID:          userID,
			Score:           &finalScoreInt,
			SubmittedAt:     &now,
			IsLate:          isLate,
		}

		existingSub, _ := u.submissionRepo.GetByClassMaterialIDAndStudentID(material.ID, userID)
		if existingSub != nil {
			existingSub.Score = &finalScoreInt
			existingSub.SubmittedAt = &now
			existingSub.IsLate = isLate
			u.submissionRepo.Update(existingSub)
		} else {
			u.submissionRepo.Create(submission)
		}

		// C. Create Student Response Logs
		for i, q := range quizSnapshot.Questions {
			qIndex := i + 1
			// Assuming GetAnswerLog exists in GameRepository to fetch specific answer
			answerLog, err := u.gameRepo.GetAnswerLog(ctx, pin, qIndex, userID)

			response := &entities.QuizStudentResponse{
				QuizGameLogID: gameLog.ID,
				UserID:        userID,
				QuestionID:    q.ID,
				IsCorrect:     false,
				Points:        0,
				TimeTaken:     0,
			}

			if err == nil && answerLog != nil {
				response.OptionID = &answerLog.OptionID
				response.IsCorrect = answerLog.IsCorrect
				response.Points = answerLog.Points
				response.TimeTaken = answerLog.TimeMs
			}

			u.quizStudentResponseRepo.Create(response)
		}
	}

	return nil
}
