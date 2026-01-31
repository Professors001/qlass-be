package redis

import (
	"context"
	"log"
	"qlass-be/domain/repositories"
	"qlass-be/dtos"
	"time"
)

type UserRedisRepository struct {
	helper *CacheHelper
}

func NewUserRedisRepository(helper *CacheHelper) repositories.UserCacheRepository {
	return &UserRedisRepository{
		helper: helper,
	}
}

var keyPrefix = "reg:"

func (r *UserRedisRepository) SetRegistrationData(ctx context.Context, key string, data dtos.TempRegisterDataDto, duration time.Duration) error {
	// Log data for debugging
	key = keyPrefix + key
	log.Printf("DEBUG SetRegistrationData - Key: %s, Data: %+v", key, data)
	return r.helper.SetJSON(ctx, key, data, duration)
}

func (r *UserRedisRepository) GetRegistrationData(ctx context.Context, key string) (dtos.TempRegisterDataDto, error) {
	key = keyPrefix + key
	var data dtos.TempRegisterDataDto
	err := r.helper.GetJSON(ctx, key, &data)
	return data, err
}
