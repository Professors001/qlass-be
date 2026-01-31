package usecases

import (
	"context"
	"errors"
	"log"
	"qlass-be/domain/entities"
	"qlass-be/domain/repositories"
	"qlass-be/dtos"
	"qlass-be/middleware"
	"qlass-be/transform"
	"time"

	"golang.org/x/crypto/bcrypt"
)

type UserUseCase interface {
	RegisterFirstStep(ctx context.Context, req *dtos.RegisterRequestStepOneDto) (*dtos.ResponseRegisterStepOneDto, error)
	RegisterSecondStep(ctx context.Context, req *dtos.RegisterRequestStepTwoDto) (*dtos.ResponseRegisterStepTwoDto, error)
	Login(ctx context.Context, req *dtos.LoginRequestDto) (*dtos.LoginResponseDto, error)
	CreateTeacher(ctx context.Context, req *dtos.CreateTeacherRequestDto) (*dtos.CreateTeacherResponseDto, error)
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

	token, err := u.jwtService.GenerateToken(user.ID, user.Role)
	if err != nil {
		return nil, errors.New("failed to generate token")
	}

	return &dtos.LoginResponseDto{
		Message: "Login successful",
		Token:   token,
		User:    userDisplay,
	}, nil
}

func (u *userUseCase) CreateTeacher(ctx context.Context, req *dtos.CreateTeacherRequestDto) (*dtos.CreateTeacherResponseDto, error) {
	// Check if user already exists
	if _, err := u.userRepo.GetByEmail(req.Email); err == nil {
		return nil, errors.New("email already exists")
	}
	if _, err := u.userRepo.GetByUniID(req.UniversityID); err == nil {
		return nil, errors.New("university ID already exists")
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	teacher := &entities.User{
		UniversityID:  req.UniversityID,
		Email:         req.Email,
		PasswordHash:  string(hashedPassword),
		FirstName:     req.FirstName,
		LastName:      req.LastName,
		Role:          "teacher",
		IsVerified:    false,
		IsActive:      true,
		ProfileImgURL: "https://ui-avatars.com/api/?name=" + req.FirstName + "+" + req.LastName,
	}

	if err := u.userRepo.Create(teacher); err != nil {
		log.Println("Error creating teacher in DB:", err)
		return nil, err
	}

	return &dtos.CreateTeacherResponseDto{
		Message: "Teacher created successfully",
		UserID:  teacher.ID,
	}, nil
}
