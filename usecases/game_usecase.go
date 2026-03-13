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
	JoinGame(ctx context.Context, pin string, userID uint) (*dtos.GameSyncResponseDto, *dtos.WSEventDto, error)
	LeaveGame(ctx context.Context, pin string, userID uint) (*dtos.WSEventDto, error)
	StartGame(ctx context.Context, pin string, hostID uint) error
	NextStep(ctx context.Context, pin string, hostID uint) (*dtos.WSEventDto, error)
	TimeoutQuestion(ctx context.Context, pin string, questionIndex int) (*dtos.WSEventDto, error)
	SubmitAnswer(ctx context.Context, pin string, userID uint, optionID uint) (*dtos.StudentAnswerResponseDto, *dtos.WSEventDto, error)
	StoreGameData(ctx context.Context, pin string) error

	// Master Sync Function
	GetSyncState(ctx context.Context, pin string, userID uint) (*dtos.GameSyncResponseDto, error)
}

type gameUseCase struct {
	gameRepo                repositories.GameRepository
	quizGameLogRepo         repositories.QuizGameLogRepository
	classMaterialRepo       repositories.ClassMaterialRepository
	classRepo               repositories.ClassRepository
	enrollRepo              repositories.EnrollRepository
	userRepo                repositories.UserRepository
	userUseCase             UserUseCase
	submissionRepo          repositories.SubmissionRepository
	quizStudentResponseRepo repositories.QuizStudentResponseRepository
}

func NewGameUseCase(
	gameRepo repositories.GameRepository,
	quizGameLogRepo repositories.QuizGameLogRepository,
	classMaterialRepo repositories.ClassMaterialRepository,
	classRepo repositories.ClassRepository,
	enrollRepo repositories.EnrollRepository,
	userRepo repositories.UserRepository,
	userUseCase UserUseCase,
	submissionRepo repositories.SubmissionRepository,
	quizStudentResponseRepo repositories.QuizStudentResponseRepository,
) GameUseCase {
	return &gameUseCase{
		gameRepo:                gameRepo,
		quizGameLogRepo:         quizGameLogRepo,
		classMaterialRepo:       classMaterialRepo,
		classRepo:               classRepo,
		enrollRepo:              enrollRepo,
		userRepo:                userRepo,
		userUseCase:             userUseCase,
		submissionRepo:          submissionRepo,
		quizStudentResponseRepo: quizStudentResponseRepo,
	}
}

// =========================================================================
// THE MASTER SYNC FUNCTION (Single Source of Truth)
// =========================================================================
func (u *gameUseCase) GetSyncState(ctx context.Context, pin string, userID uint) (*dtos.GameSyncResponseDto, error) {
	state, err := u.gameRepo.GetGameState(ctx, pin)
	if err != nil {
		return nil, errors.New("game session not found")
	}

	// 1. ตรวจสอบ Role
	role := "player"
	if state.HostID == userID {
		role = "host"
	}

	// 2. สร้างโครงร่างพื้นฐาน
	syncData := &dtos.GameSyncResponseDto{
		PIN:               state.Pin,
		QuizTitle:         state.QuizTitle,
		Role:              role,
		Status:            state.Status,
		QuestionState:     state.QuestionState,
		ServerTimeMs:      time.Now().UnixMilli(),
		QuestionStartedAt: state.QuestionStartedAt.UnixMilli(),
		QuestionEndsAt:    state.QuestionEndsAt.UnixMilli(),
	}

	// 3. แนบข้อมูลตาม Phase
	if state.Status == "waiting" {
		players := u.getLeaderboardDtos(ctx, pin, 0)
		syncData.LobbyStateObject = &dtos.LobbyStateDto{
			TotalPlayers: len(players),
			Players:      players,
		}
	} else if state.Status == "running" || state.Status == "finished" {
		quizDataJSON, _ := u.gameRepo.GetQuizData(ctx, pin)
		var quizSnapshot dtos.GetQuizResponseDto
		json.Unmarshal([]byte(quizDataJSON), &quizSnapshot)

		if state.CurrentQuestion > 0 && state.CurrentQuestion <= len(quizSnapshot.Questions) {
			currentQ := quizSnapshot.Questions[state.CurrentQuestion-1]

			var options []dtos.WSQuizOptionDto
			for _, opt := range currentQ.Options {
				options = append(options, dtos.WSQuizOptionDto{
					ID:         opt.ID,
					OptionText: opt.OptionText,
					Label:      getOptionLabel(opt.OrderIndex),
				})
			}

			syncData.QuestionStateObject = &dtos.QuestionStateDto{
				CurrentQuestion:  state.CurrentQuestion,
				TotalQuestions:   state.TotalQuestions,
				TimeLimitSeconds: currentQ.TimeLimitSeconds,
				QuestionText:     currentQ.QuestionText,
				PointsMultiplier: currentQ.PointsMultiplier,
				Options:          options,
			}

			if currentQ.MediaAttachment != nil {
				syncData.QuestionStateObject.ImageURL = currentQ.MediaAttachment.FileURL
			}

			syncData.QuestionStateObject.TotalPlayers = state.TotalPlayers
			syncData.QuestionStateObject.AnsweredCount = state.OptionACount + state.OptionBCount + state.OptionCCount + state.OptionDCount
		}

		if state.QuestionState == "revealed" || state.Status == "finished" {
			resultLeaderboard := u.getLeaderboardDtos(ctx, pin, 5)
			syncData.ResultStateObject = &dtos.ResultStateDto{
				CorrectOptionID: state.CorrectOptionID,
				Stats: dtos.LiveStatsPayload{
					OptionACount: state.OptionACount,
					OptionBCount: state.OptionBCount,
					OptionCCount: state.OptionCCount,
					OptionDCount: state.OptionDCount,
					OptionAID:    state.OptionAID,
					OptionBID:    state.OptionBID,
					OptionCID:    state.OptionCID,
					OptionDID:    state.OptionDID,
				},
				Leaderboard: resultLeaderboard,
			}
		}

		if state.Status == "finished" {
			syncData.Leaderboard = u.getLeaderboardDtos(ctx, pin, 5)
		}
	}

	// 4. แนบข้อมูลส่วนตัวของนักเรียน
	if role == "player" {
		pData, _ := u.gameRepo.GetPlayerData(ctx, pin, userID)
		if pData != nil {
			myState := &dtos.MyPlayerStateDto{
				UserID:       userID,
				Name:         pData.Name,
				UniversityID: pData.UniversityID,
				AvatarURL:    pData.AvatarURL,
				Score:        pData.Score,
				Streak:       pData.Streak,
			}

			if state.CurrentQuestion > 0 {
				hasAnswered, _ := u.gameRepo.HasUserAnswered(ctx, pin, state.CurrentQuestion, userID)
				myState.HasAnswered = hasAnswered

				if hasAnswered {
					ansLog, err := u.gameRepo.GetAnswerLog(ctx, pin, state.CurrentQuestion, userID)
					if err == nil && ansLog != nil {
						myState.SelectedOptionID = ansLog.OptionID

						if state.QuestionState == "revealed" || state.Status == "finished" {
							myState.LastResult = &dtos.StudentPersonalResultPayload{
								IsCorrect:    ansLog.IsCorrect,
								PointsEarned: ansLog.Points,
								TotalScore:   pData.Score,
								Streak:       pData.Streak,
							}
						}
					}
				}
			}

			syncData.MyState = myState
		}
	}

	return syncData, nil
}

// =========================================================================
// GAME FLOW LOGIC
// =========================================================================

func (u *gameUseCase) StartGameSession(ctx context.Context, teacherID uint, dto dtos.CreateGameRequestDto) (*dtos.CreateGameResponseDto, error) {
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

	logs, err := u.quizGameLogRepo.GetByClassMaterialID(material.ID)
	if err != nil || len(logs) == 0 {
		return nil, errors.New("quiz game log not found")
	}
	gameLog := logs[0]

	var quizSnapshot dtos.GetQuizResponseDto
	if err := json.Unmarshal(gameLog.QuizSnapshot, &quizSnapshot); err != nil {
		return nil, fmt.Errorf("failed to parse quiz snapshot: %v", err)
	}

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

	gameState := &entities.GameStateRedis{
		Pin:             pin,
		Status:          "waiting",
		QuestionState:   "none",
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

	snapshotJSON, _ := json.Marshal(quizSnapshot)
	if err := u.gameRepo.SetQuizData(ctx, pin, string(snapshotJSON)); err != nil {
		return nil, err
	}

	if err := u.gameRepo.SetTeacherSession(ctx, teacherID, pin); err != nil {
		return nil, err
	}

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

func (u *gameUseCase) JoinGame(ctx context.Context, pin string, userID uint) (*dtos.GameSyncResponseDto, *dtos.WSEventDto, error) {
	state, err := u.gameRepo.GetGameState(ctx, pin)
	if err != nil {
		return nil, nil, errors.New("game session not found")
	}

	// Skip enrollment check for the host
	if state.HostID != userID {
		material, err := u.classMaterialRepo.GetByID(state.ClassMaterialID)
		if err != nil {
			return nil, nil, errors.New("class material not found")
		}
		enrolled, err := u.enrollRepo.IsEnrolled(material.ClassID, userID)
		if err != nil {
			return nil, nil, err
		}
		if !enrolled {
			return nil, nil, errors.New("you are not enrolled in this class")
		}
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

	user, err := u.userRepo.GetByID(userID)
	if err != nil {
		return nil, nil, err
	}
	user_img_url := u.userUseCase.GetProfileImgUrlByUserID(user.ID)

	existingData, _ := u.gameRepo.GetPlayerData(ctx, pin, userID)
	playerData := &entities.PlayerDataRedis{
		Name:         user.FirstName + " " + user.LastName,
		UniversityID: user.UniversityID,
		AvatarURL:    user_img_url,
		Score:        0,
		Correct:      0,
		Streak:       0,
	}

	if existingData != nil && existingData.Name != "" {
		playerData.Score = existingData.Score
		playerData.Correct = existingData.Correct
		playerData.Streak = existingData.Streak
	}

	if err := u.gameRepo.SavePlayerData(ctx, pin, userID, playerData); err != nil {
		return nil, nil, err
	}

	if !inLobby {
		if err := u.gameRepo.AddPlayerToLobby(ctx, pin, userID); err != nil {
			return nil, nil, err
		}
		u.gameRepo.IncrementField(ctx, pin, "total_players", 1)
	}

	syncState, err := u.GetSyncState(ctx, pin, userID)
	if err != nil {
		return nil, nil, err
	}

	var wsEvent *dtos.WSEventDto
	if !inLobby {
		wsEvent = &dtos.WSEventDto{
			Type: "PLAYER_JOINED",
			Payload: map[string]interface{}{
				"new_player": dtos.PlayerDto{
					UserID:       user.ID,
					Name:         playerData.Name,
					UniversityID: playerData.UniversityID,
					AvatarURL:    playerData.AvatarURL,
				},
				"total_players": syncState.LobbyStateObject.TotalPlayers,
			},
		}
	}

	return syncState, wsEvent, nil
}

func (u *gameUseCase) LeaveGame(ctx context.Context, pin string, userID uint) (*dtos.WSEventDto, error) {
	state, err := u.gameRepo.GetGameState(ctx, pin)
	if err != nil {
		return nil, errors.New("game session not found")
	}

	if state.Status == "waiting" {
		if err := u.gameRepo.RemovePlayerFromLobby(ctx, pin, userID); err != nil {
			return nil, err
		}
		_ = u.gameRepo.DeletePlayerData(ctx, pin, userID)
		u.gameRepo.IncrementField(ctx, pin, "total_players", -1)

		players := u.getLeaderboardDtos(ctx, pin, 0)
		return &dtos.WSEventDto{
			Type: "LOBBY_UPDATE",
			Payload: dtos.LobbyStateDto{
				TotalPlayers: len(players),
				Players:      players,
			},
		}, nil
	}
	return nil, nil
}

func (u *gameUseCase) StartGame(ctx context.Context, pin string, hostID uint) error {
	state, err := u.gameRepo.GetGameState(ctx, pin)
	if err != nil {
		return err
	}

	if state.HostID != hostID {
		return errors.New("unauthorized: only host can start the game")
	}
	return u.gameRepo.UpdateGameState(ctx, pin, map[string]interface{}{"status": "running"})
}

func (u *gameUseCase) NextStep(ctx context.Context, pin string, hostID uint) (*dtos.WSEventDto, error) {
	state, err := u.gameRepo.GetGameState(ctx, pin)
	if err != nil {
		return nil, errors.New("game session not found")
	}

	if state.HostID != hostID {
		return nil, errors.New("unauthorized")
	}
	if state.Status == "finished" {
		return nil, errors.New("game is already finished")
	}

	if state.Status == "waiting" {
		if err := u.gameRepo.UpdateGameState(ctx, pin, map[string]interface{}{
			"status":           "running",
			"current_question": 1,
			"question_state":   "hold",
			"option_a_count":   0, "option_b_count": 0, "option_c_count": 0, "option_d_count": 0,
		}); err != nil {
			return nil, err
		}
	} else if state.Status == "running" {
		switch state.QuestionState {
		case "hold":
			quizDataJSON, _ := u.gameRepo.GetQuizData(ctx, pin)
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

			var correctID uint
			for _, opt := range currentQ.Options {
				if opt.IsCorrect {
					correctID = opt.ID
				}
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
			updates["correct_option_id"] = correctID
			if err := u.gameRepo.UpdateGameState(ctx, pin, updates); err != nil {
				return nil, err
			}

		case "answering":
			return u.finishQuestion(ctx, pin, state)
		case "time_up":
			if err := u.gameRepo.UpdateGameState(ctx, pin, map[string]interface{}{"question_state": "revealed"}); err != nil {
				return nil, err
			}
		case "revealed":
			if state.CurrentQuestion >= state.TotalQuestions {
				if err := u.gameRepo.UpdateGameState(ctx, pin, map[string]interface{}{"status": "finished"}); err != nil {
					return nil, err
				}
				u.StoreGameData(ctx, pin)
			} else {
				if err := u.gameRepo.UpdateGameState(ctx, pin, map[string]interface{}{
					"current_question": state.CurrentQuestion + 1,
					"question_state":   "hold",
					"option_a_count":   0, "option_b_count": 0, "option_c_count": 0, "option_d_count": 0,
				}); err != nil {
					return nil, err
				}
			}
		}
	}

	return &dtos.WSEventDto{
		Type:    "SYNC_TRIGGER",
		Payload: nil,
	}, nil
}

func (u *gameUseCase) TimeoutQuestion(ctx context.Context, pin string, questionIndex int) (*dtos.WSEventDto, error) {
	state, err := u.gameRepo.GetGameState(ctx, pin)
	if err != nil || state.CurrentQuestion != questionIndex || state.QuestionState != "answering" {
		return nil, nil
	}
	return u.finishQuestion(ctx, pin, state)
}

func (u *gameUseCase) finishQuestion(ctx context.Context, pin string, state *entities.GameStateRedis) (*dtos.WSEventDto, error) {
	if err := u.gameRepo.UpdateGameState(ctx, pin, map[string]interface{}{"question_state": "time_up"}); err != nil {
		return nil, err
	}
	return &dtos.WSEventDto{
		Type:    "SYNC_TRIGGER",
		Payload: nil,
	}, nil
}

func (u *gameUseCase) SubmitAnswer(ctx context.Context, pin string, userID uint, optionID uint) (*dtos.StudentAnswerResponseDto, *dtos.WSEventDto, error) {
	state, err := u.gameRepo.GetGameState(ctx, pin)
	if err != nil {
		return nil, nil, errors.New("game session not found")
	}

	if state.QuestionState != "answering" {
		return nil, nil, errors.New("question is not open for answers")
	}

	isNew, err := u.gameRepo.MarkUserAnswered(ctx, pin, state.CurrentQuestion, userID)
	if err != nil || !isNew {
		return nil, nil, errors.New("already answered this question")
	}

	now := time.Now()
	timeTaken := now.Sub(state.QuestionStartedAt).Seconds()

	quizDataJSON, _ := u.gameRepo.GetQuizData(ctx, pin)
	var quizSnapshot dtos.GetQuizResponseDto
	json.Unmarshal([]byte(quizDataJSON), &quizSnapshot)
	question := quizSnapshot.Questions[state.CurrentQuestion-1]

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

	points := 0
	if selectedOption.IsCorrect {
		ratio := timeTaken / float64(question.TimeLimitSeconds)
		if ratio > 1 {
			ratio = 1
		}
		baseScore := (1 - (ratio / 2)) * 1000.0
		points = int(math.Round(baseScore)) * question.PointsMultiplier
	}

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

	u.gameRepo.SavePlayerData(ctx, pin, userID, pData)
	u.gameRepo.UpdateScore(ctx, pin, userID, float64(pData.Score))

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

	logEntry := &entities.AnswerLog{
		OptionID:  optionID,
		TimeMs:    int(timeTaken * 1000),
		Points:    points,
		IsCorrect: selectedOption.IsCorrect,
	}
	u.gameRepo.SaveAnswerDetail(ctx, pin, state.CurrentQuestion, userID, logEntry)

	wsEvent := &dtos.WSEventDto{
		Type:    "SYNC_TRIGGER",
		Payload: nil,
	}

	return &dtos.StudentAnswerResponseDto{Message: "Answer submitted successfully"}, wsEvent, nil
}

// =========================================================================
// DATA PERSISTENCE (SAVE TO DATABASE)
// =========================================================================

func (u *gameUseCase) StoreGameData(ctx context.Context, pin string) error {
	state, err := u.gameRepo.GetGameState(ctx, pin)
	if err != nil {
		log.Println("Failed to get game state for logging:", err)
		return err
	}

	quizDataJSON, err := u.gameRepo.GetQuizData(ctx, pin)
	if err != nil {
		return err
	}
	var quizSnapshot dtos.GetQuizResponseDto
	json.Unmarshal([]byte(quizDataJSON), &quizSnapshot)

	material, err := u.classMaterialRepo.GetByID(state.ClassMaterialID)
	if err != nil {
		return err
	}

	logs, err := u.quizGameLogRepo.GetByClassMaterialID(material.ID)
	if err != nil || len(logs) == 0 {
		return errors.New("game log not found")
	}
	gameLog := logs[0]

	now := time.Now()
	gameLog.Status = "finished"
	gameLog.FinishedAt = &now
	if err := u.quizGameLogRepo.Update(gameLog); err != nil {
		return err
	}

	var maxGameScore float64 = 0
	for _, q := range quizSnapshot.Questions {
		maxGameScore += float64(1000 * q.PointsMultiplier)
	}
	if maxGameScore == 0 {
		maxGameScore = 1
	}

	playerIDs, err := u.gameRepo.GetLobbyPlayers(ctx, pin)
	if err != nil {
		return err
	}

	for _, userID := range playerIDs {
		pData, err := u.gameRepo.GetPlayerData(ctx, pin, userID)
		if err != nil || userID == state.HostID {
			continue
		}

		rawScore := float64(pData.Score)
		var materialPoints float64
		if material.Points != nil {
			materialPoints = float64(*material.Points)
		}

		finalScore := (rawScore / maxGameScore) * materialPoints
		finalScoreInt := int(finalScore)

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

		for i, q := range quizSnapshot.Questions {
			qIndex := i + 1
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

// =========================================================================
// HELPER FUNCTIONS
// =========================================================================

func (u *gameUseCase) getLeaderboardDtos(ctx context.Context, pin string, limit int) []dtos.PlayerDto {
	leaderboard, _ := u.gameRepo.GetLeaderboard(ctx, pin, limit)
	var result []dtos.PlayerDto
	for _, p := range leaderboard {
		pData, _ := u.gameRepo.GetPlayerData(ctx, pin, p.UserID)
		name := "Unknown"
		universityID := ""
		avatar := ""
		streak := 0
		if pData != nil {
			name = pData.Name
			universityID = pData.UniversityID
			avatar = pData.AvatarURL
			streak = pData.Streak
		}
		result = append(result, dtos.PlayerDto{
			UserID:       p.UserID,
			Name:         name,
			UniversityID: universityID,
			AvatarURL:    avatar,
			Score:        p.Score,
			Rank:         p.Rank,
			Streak:       streak,
		})
	}
	return result
}

func getOptionLabel(index int) string {
	switch index {
	case 1:
		return "A"
	case 2:
		return "B"
	case 3:
		return "C"
	case 4:
		return "D"
	default:
		return ""
	}
}
