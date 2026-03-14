package repositories

import (
	"context"
	"qlass-be/dtos"
	"time"
)

type UserCacheRepository interface {
	SetRegistrationData(ctx context.Context, key string, data dtos.TempRegisterDataDto, duration time.Duration) error
	GetRegistrationData(ctx context.Context, key string) (dtos.TempRegisterDataDto, error)
	SetForgetPasswordData(ctx context.Context, key string, data *dtos.TempForgetPasswordData, duration time.Duration) error
	GetForgetPasswordData(ctx context.Context, key string) (*dtos.TempForgetPasswordData, error)
}
