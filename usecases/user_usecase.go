package usecases

import (
	"context"
	"errors"
	"log"
	"qlass-be/domain/repositories"
	"qlass-be/dtos"
	"qlass-be/infrastructure/middleware"
	"qlass-be/transform"
	"time"

	"golang.org/x/crypto/bcrypt"
)

type UserUseCase interface {
	RegisterFirstStep(ctx context.Context, req *dtos.RegisterRequestStepOneDto) (*dtos.ResponseRegisterStepOneDto, error)
	RegisterSecondStep(ctx context.Context, req *dtos.RegisterRequestStepTwoDto) (*dtos.ResponseRegisterStepTwoDto, error)
	Login(ctx context.Context, req *dtos.LoginRequestDto) (*dtos.LoginResponseDto, error)
}

type userUseCase struct {
	userRepo      repositories.UserRepository
	userCacheRepo repositories.UserCacheRepository
	jwtService    middleware.JwtService
}

func NewUserUseCase(repo repositories.UserRepository, cacheRepo repositories.UserCacheRepository, jwtService middleware.JwtService) UserUseCase {
	return &userUseCase{
		userRepo:      repo,
		userCacheRepo: cacheRepo,
		jwtService:    jwtService,
	}
}

func (u *userUseCase) RegisterFirstStep(ctx context.Context, req *dtos.RegisterRequestStepOneDto) (*dtos.ResponseRegisterStepOneDto, error) {
	// Check if registration is already pending in Redis
	if _, err := u.userCacheRepo.GetRegistrationData(ctx, req.Email); err == nil {
		return nil, errors.New("registration pending, please check your email for OTP")
	}

	// Check if user already exists
	existingUserByEmail, _ := u.userRepo.GetByEmail(req.Email)
	if existingUserByEmail != nil {
		return nil, errors.New("Email already exists")
	}

	existingUserByUniID, _ := u.userRepo.GetByUniID(req.UniversityID)
	if existingUserByUniID != nil {
		return nil, errors.New("University ID already exists")
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

func (u *userUseCase) RegisterSecondStep(ctx context.Context, req *dtos.RegisterRequestStepTwoDto) (*dtos.ResponseRegisterStepTwoDto, error) {
	// Retrieve temp data from Redis
	tempData, err := u.userCacheRepo.GetRegistrationData(ctx, req.Email)
	if err != nil {
		return nil, errors.New("no pending registration found or OTP expired")
	}

	// Validate OTP
	if tempData.OTP != req.OTP {
		return nil, errors.New("invalid OTP")
	}

	// Create user in DB
	newUser := transform.TempRegisterDataDtoToUserEntity(tempData)
	err = u.userRepo.Create(newUser)
	if err != nil {
		log.Println("Error creating user in DB:", err)
		return nil, err
	}

	return &dtos.ResponseRegisterStepTwoDto{
		Message: "Registration successful",
	}, nil

}

func (u *userUseCase) Login(ctx context.Context, req *dtos.LoginRequestDto) (*dtos.LoginResponseDto, error) {
	// Retrieve user by email
	user, err := u.userRepo.GetByUniID(req.UniversityID)
	if err != nil {
		return nil, errors.New("This Account is not registered")
	}

	// Compare password
	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password))
	if err != nil {
		return nil, errors.New("Incorrect password")
	}

	userDisplay := dtos.UserDisplayData{
		UniversityID: user.UniversityID,
		Email:        user.Email,
		FirstName:    user.FirstName,
		LastName:     user.LastName,
		Role:         user.Role,
	}

	token, err := u.jwtService.GenerateToken(user.UniversityID, user.Role)
	if err != nil {
		return nil, errors.New("failed to generate token")
	}

	return &dtos.LoginResponseDto{
		Message: "Login successful",
		Token:   token,
		User:    userDisplay,
	}, nil
}

