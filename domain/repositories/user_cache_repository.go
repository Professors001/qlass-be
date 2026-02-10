package repositories

import (
	"context"
	"qlass-be/dtos"
	"time"
)

type UserCacheRepository interface {
	SetRegistrationData(ctx context.Context, key string, data dtos.TempRegisterDataDto, duration time.Duration) error
	GetRegistrationData(ctx context.Context, key string) (dtos.TempRegisterDataDto, error)
}
