package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"qlass-be/domain/entities"
	"qlass-be/domain/repositories"

	"github.com/redis/go-redis/v9"
)

type gameRedisRepository struct {
	client *redis.Client
}

func NewGameRedisRepository(helper *CacheHelper) repositories.GameRepository {
	return &gameRedisRepository{client: helper.cache.client}
}

const (
	TTL_GAME    = 24 * time.Hour
	TTL_SESSION = 4 * time.Hour
)

// --- 1. Teacher Session ---

func (r *gameRedisRepository) SetTeacherSession(ctx context.Context, teacherID uint, pin string) error {
	key := fmt.Sprintf("teacher:session:%d", teacherID)
	return r.client.Set(ctx, key, pin, TTL_SESSION).Err()
}

func (r *gameRedisRepository) GetTeacherSession(ctx context.Context, teacherID uint) (string, error) {
	key := fmt.Sprintf("teacher:session:%d", teacherID)
	return r.client.Get(ctx, key).Result()
}

// --- 2. Game State & Quiz Data ---

func (r *gameRedisRepository) CreateGameState(ctx context.Context, pin string, state *entities.GameStateRedis) error {
	key := fmt.Sprintf("game:%s:state", pin)
	err := r.client.HSet(ctx, key, state).Err()
	if err != nil {
		return err
	}
	return r.client.Expire(ctx, key, TTL_GAME).Err()
}

func (r *gameRedisRepository) GetGameState(ctx context.Context, pin string) (*entities.GameStateRedis, error) {
	key := fmt.Sprintf("game:%s:state", pin)
	var state entities.GameStateRedis
	err := r.client.HGetAll(ctx, key).Scan(&state)
	if err != nil {
		return nil, err
	}
	if state.Pin == "" {
		return nil, fmt.Errorf("game session not found")
	}
	return &state, nil
}

func (r *gameRedisRepository) UpdateGameState(ctx context.Context, pin string, fields map[string]interface{}) error {
	key := fmt.Sprintf("game:%s:state", pin)
	return r.client.HSet(ctx, key, fields).Err()
}

func (r *gameRedisRepository) IncrementField(ctx context.Context, pin string, field string, amount int) error {
	key := fmt.Sprintf("game:%s:state", pin)
	// Useful for incrementing option_a_count, total_players, etc.
	return r.client.HIncrBy(ctx, key, field, int64(amount)).Err()
}

func (r *gameRedisRepository) SetQuizData(ctx context.Context, pin string, quizJSON string) error {
	key := fmt.Sprintf("game:%s:quiz", pin)
	return r.client.Set(ctx, key, quizJSON, TTL_GAME).Err()
}

func (r *gameRedisRepository) GetQuizData(ctx context.Context, pin string) (string, error) {
	key := fmt.Sprintf("game:%s:quiz", pin)
	return r.client.Get(ctx, key).Result()
}

// --- 3. Players ---

func (r *gameRedisRepository) AddPlayerToLobby(ctx context.Context, pin string, userID uint) error {
	key := fmt.Sprintf("game:%s:players", pin)
	err := r.client.SAdd(ctx, key, userID).Err()
	if err != nil {
		return err
	}
	return r.client.Expire(ctx, key, TTL_GAME).Err()
}

func (r *gameRedisRepository) RemovePlayerFromLobby(ctx context.Context, pin string, userID uint) error {
	key := fmt.Sprintf("game:%s:players", pin)
	return r.client.SRem(ctx, key, userID).Err()
}

func (r *gameRedisRepository) AddAllowedPlayer(ctx context.Context, pin string, userID uint) error {
	key := fmt.Sprintf("game:%s:allowed", pin)
	err := r.client.SAdd(ctx, key, userID).Err()
	if err != nil {
		return err
	}
	return r.client.Expire(ctx, key, TTL_GAME).Err()
}

func (r *gameRedisRepository) IsPlayerAllowed(ctx context.Context, pin string, userID uint) (bool, error) {
	key := fmt.Sprintf("game:%s:allowed", pin)
	return r.client.SIsMember(ctx, key, userID).Result()
}

func (r *gameRedisRepository) GetLobbyPlayers(ctx context.Context, pin string) ([]uint, error) {
	key := fmt.Sprintf("game:%s:players", pin)
	strIDs, err := r.client.SMembers(ctx, key).Result()
	if err != nil {
		return nil, err
	}

	ids := make([]uint, len(strIDs))
	for i, s := range strIDs {
		val, _ := strconv.ParseUint(s, 10, 64)
		ids[i] = uint(val)
	}
	return ids, nil
}

func (r *gameRedisRepository) SavePlayerData(ctx context.Context, pin string, userID uint, data *entities.PlayerDataRedis) error {
	key := fmt.Sprintf("game:%s:player:%d", pin, userID)
	err := r.client.HSet(ctx, key, data).Err()
	if err != nil {
		return err
	}
	return r.client.Expire(ctx, key, TTL_GAME).Err()
}

func (r *gameRedisRepository) GetPlayerData(ctx context.Context, pin string, userID uint) (*entities.PlayerDataRedis, error) {
	key := fmt.Sprintf("game:%s:player:%d", pin, userID)
	var data entities.PlayerDataRedis
	err := r.client.HGetAll(ctx, key).Scan(&data)
	if err != nil {
		return nil, err
	}
	return &data, nil
}

// --- 4. Answers ---

func (r *gameRedisRepository) MarkUserAnswered(ctx context.Context, pin string, questionIndex int, userID uint) (bool, error) {
	key := fmt.Sprintf("game:%s:answered:%d", pin, questionIndex)
	// SAdd returns 1 if new, 0 if already exists
	added, err := r.client.SAdd(ctx, key, userID).Result()
	if err != nil {
		return false, err
	}
	r.client.Expire(ctx, key, TTL_GAME)
	return added > 0, nil
}

func (r *gameRedisRepository) SaveAnswerDetail(ctx context.Context, pin string, questionIndex int, userID uint, answerLog *entities.AnswerLog) error {
	key := fmt.Sprintf("game:%s:answers:%d", pin, questionIndex)

	jsonBytes, err := json.Marshal(answerLog)
	if err != nil {
		return err
	}

	err = r.client.HSet(ctx, key, fmt.Sprintf("%d", userID), jsonBytes).Err()
	if err != nil {
		return err
	}
	return r.client.Expire(ctx, key, TTL_GAME).Err()
}

// --- 5. Leaderboard (ZSET) ---

func (r *gameRedisRepository) UpdateScore(ctx context.Context, pin string, userID uint, totalScore float64) error {
	key := fmt.Sprintf("game:%s:scores", pin)
	// ZAdd adds or updates the score for the member
	err := r.client.ZAdd(ctx, key, redis.Z{
		Score:  totalScore,
		Member: userID,
	}).Err()
	if err != nil {
		return err
	}
	return r.client.Expire(ctx, key, TTL_GAME).Err()
}

// GetLeaderboard returns Top N users (Highest score first)
func (r *gameRedisRepository) GetLeaderboard(ctx context.Context, pin string, limit int) ([]entities.PlayerScore, error) {
	key := fmt.Sprintf("game:%s:scores", pin)

	// ZRevRangeWithScores gets highest scores first
	results, err := r.client.ZRevRangeWithScores(ctx, key, 0, int64(limit-1)).Result()
	if err != nil {
		return nil, err
	}

	leaderboard := make([]entities.PlayerScore, len(results))
	for i, z := range results {
		userIDStr := z.Member.(string)
		uid, _ := strconv.ParseUint(userIDStr, 10, 64)

		leaderboard[i] = entities.PlayerScore{
			UserID: uint(uid),
			Score:  int(z.Score),
			Rank:   i + 1,
			// Note: Username/Avatar must be fetched from PlayerDataRedis or DB
		}
	}
	return leaderboard, nil
}

func (r *gameRedisRepository) KeyExists(ctx context.Context, key string) (bool, error) {
	val, err := r.client.Exists(ctx, key).Result()
	return val > 0, err
}
