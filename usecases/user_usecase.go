package usecases

import (
	"context"
	"errors"
	"log"
	"qlass-be/domain/repositories"
	"qlass-be/dtos"
	"qlass-be/transform"
	"time"

	"golang.org/x/crypto/bcrypt"
)

type UserUseCase interface {
	RegisterFirstStep(ctx context.Context, req *dtos.RegisterRequestStepOneDto) (*dtos.ResponseRegisterStepOneDto, error)
	// GetUserByID(id uint) (*entities.User, error)
	// GetUserByUID(uuid string) (*entities.User, error)
}

type userUseCase struct {
	userRepo      repositories.UserRepository
	userCacheRepo repositories.UserCacheRepository
}

func NewUserUseCase(repo repositories.UserRepository, cacheRepo repositories.UserCacheRepository) UserUseCase {
	return &userUseCase{
		userRepo:      repo,
		userCacheRepo: cacheRepo,
	}
}

func (u *userUseCase) RegisterFirstStep(ctx context.Context, req *dtos.RegisterRequestStepOneDto) (*dtos.ResponseRegisterStepOneDto, error) {
	// Check if registration is already pending in Redis
	if _, err := u.userCacheRepo.GetRegistrationData(ctx, req.Email); err == nil {
		return nil, errors.New("registration pending, please check your email for OTP")
	}

	// Check if user already exists
	existingUser, _ := u.userRepo.GetByEmail(req.Email)
	if existingUser != nil {
		return nil, errors.New("user with this email already exists")
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	tempData := transform.RequestToTempRegisterDataDto(req, string(hashedPassword), "123456")

	// Store into Redis with Email as key (TTL 5 minutes)
	err = u.userCacheRepo.SetRegistrationData(ctx, req.Email, tempData, 5*time.Minute)
	if err != nil {
		log.Println("Error storing registration data in cache:", err)
		return nil, err
	}

	return &dtos.ResponseRegisterStepOneDto{
		Message:          "Please proceed to step 2",
		RefEmail:         req.Email,
		ExpiresInSeconds: 300,
	}, nil
}

// func (u *userUseCase) GetUserByID(id uint) (*entities.User, error) {
// 	return u.userRepo.GetByID(id)
// }

// func (u *userUseCase) GetUserByUID(uuid string) (*entities.User, error) {
// 	return u.userRepo.GetByUID(uuid)
// }
