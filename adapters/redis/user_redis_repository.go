package redis

import (
	"context"
	"qlass-be/domain/repositories"
	"qlass-be/dtos"
	"qlass-be/infrastructure/cache"
	"time"
)

type UserRedisRepository struct {
	helper *cache.CacheHelper
}

func NewUserRedisRepository(helper *cache.CacheHelper) repositories.UserCacheRepository {
	return &UserRedisRepository{
		helper: helper,
	}
}

func (r *UserRedisRepository) SetRegistrationData(ctx context.Context, key string, data dtos.TempRegisterDataDto, duration time.Duration) error {
	return r.helper.SetJSON(ctx, key, data, duration)
}

func (r *UserRedisRepository) GetRegistrationData(ctx context.Context, key string) (dtos.TempRegisterDataDto, error) {
	var data dtos.TempRegisterDataDto
	err := r.helper.GetJSON(ctx, key, &data)
	return data, err
}
