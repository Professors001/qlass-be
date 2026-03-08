package usecases

import (
	"context"
	"errors"
	"log"
	"qlass-be/domain/entities"
	"qlass-be/domain/repositories"
	"qlass-be/dtos"
	"qlass-be/middleware"
	"qlass-be/transforms"
	"qlass-be/utils"
	"strings"
	"time"

	"golang.org/x/crypto/bcrypt"
)

type UserUseCase interface {
	RegisterFirstStep(ctx context.Context, req *dtos.RegisterRequestStepOneDto) (*dtos.ResponseRegisterStepOneDto, error)
	RegisterSecondStep(ctx context.Context, req *dtos.RegisterRequestStepTwoDto) (*dtos.ResponseRegisterStepTwoDto, error)
	Login(ctx context.Context, req *dtos.LoginRequestDto) (*dtos.LoginResponseDto, error)
	CreateTeacher(ctx context.Context, req *dtos.CreateTeacherRequestDto) (*dtos.CreateTeacherResponseDto, error)
	UpdateUser(req *dtos.UpdateUserRequestDto, userID uint) (*dtos.UserDisplayData, error)
	ChangePassword(req *dtos.ChangePasswordRequestDto, userID uint) (*dtos.ChangePasswordResponseDto, error)
	ForgetPasswordStep1(ctx context.Context, req *dtos.ForgetPasswordStep1RequestDto) (*dtos.ForgetPasswordStep1ResponseDto, error)
	ForgetPasswordStep2(ctx context.Context, req *dtos.ForgetPasswordStep2RequestDto) (*dtos.ForgetPasswordStep2ResponseDto, error)
	AdminUpdateUser(req *dtos.AdminUpdateUserRequestDto) error
	GetProfileImgUrlByUserID(userID uint) string
}

type userUseCase struct {
	userRepo          repositories.UserRepository
	userCacheRepo     repositories.UserCacheRepository
	jwtService        middleware.JwtService
	emailService      EmailService
	attachmentUseCase AttachmentUseCase
}

func NewUserUseCase(repo repositories.UserRepository, cacheRepo repositories.UserCacheRepository, jwtService middleware.JwtService, email EmailService, attachmentUC AttachmentUseCase) UserUseCase {
	return &userUseCase{
		userRepo:          repo,
		userCacheRepo:     cacheRepo,
		jwtService:        jwtService,
		emailService:      email,
		attachmentUseCase: attachmentUC,
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

	otp := utils.GenerateRandomString(6)

	tempData := transforms.RequestToTempRegisterDataDto(req, string(hashedPassword), otp)

	err = u.emailService.SendOTP(req.Email, otp)
	if err != nil {
		log.Println("Error sending OTP email:", err)
		return nil, err
	}

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

	if tempData.OTP != req.OTP {
		return nil, errors.New("invalid OTP")
	}

	newUser := transforms.TempRegisterDataDtoToUserEntity(tempData)

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
	var user *entities.User
	var err error

	if strings.Contains(req.Identifier, "@") {
		log.Println("Login attempt with Email:", req.Identifier)
		user, err = u.userRepo.GetByEmail(req.Identifier)
	} else {
		log.Println("Login attempt with University ID:", req.Identifier)
		user, err = u.userRepo.GetByUniID(req.Identifier)
	}

	if err != nil {
		return nil, errors.New("This Account is not registered")
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password))
	if err != nil {
		return nil, errors.New("Incorrect password")
	}

	var profileImgURL string
	if user.ProfileImgAttachment != nil && user.ProfileImgAttachment.ObjectKey != "" {
		url, err := u.attachmentUseCase.GenerateFileURL(user.ProfileImgAttachment.ObjectKey)
		if err == nil {
			profileImgURL = url
		} else {
			log.Println("Error generating profile image URL:", err)
		}
	}

	if profileImgURL == "" {
		profileImgURL = "https://ui-avatars.com/api/?name=" + user.FirstName + "+" + user.LastName
	}

	userDisplay := transforms.UserEntityToUserDisplayResponse(user, profileImgURL)

	token, err := u.jwtService.GenerateToken(user.ID, user.Role)
	if err != nil {
		return nil, errors.New("failed to generate token")
	}

	return &dtos.LoginResponseDto{
		Message: "Login successful",
		Token:   token,
		User:    *userDisplay,
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
		UniversityID: req.UniversityID,
		Email:        req.Email,
		PasswordHash: string(hashedPassword),
		FirstName:    req.FirstName,
		LastName:     req.LastName,
		Role:         "teacher",
		IsVerified:   false,
		IsActive:     true,
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

func (u *userUseCase) UpdateUser(req *dtos.UpdateUserRequestDto, userID uint) (*dtos.UserDisplayData, error) {
	user, err := u.userRepo.GetByID(userID)
	if err != nil {
		return nil, errors.New("user not found")
	}

	if req.FirstName != "" {
		user.FirstName = req.FirstName
	}
	if req.LastName != "" {
		user.LastName = req.LastName
	}

	if req.ProfileImgAttachmentID != *user.ProfileImgAttachmentID {
		attachmentToDelete := *user.ProfileImgAttachmentID

		user.ProfileImgAttachmentID = &req.ProfileImgAttachmentID

		if err := u.attachmentUseCase.DeleteAttachment(attachmentToDelete); err != nil {
			log.Println("Error deleting attachment:", err)
			return nil, err
		}
	}

	if err := u.userRepo.Update(user); err != nil {
		log.Println("Error updating user:", err)
		return nil, err
	}

	var profileImgURL string
	if user.ProfileImgAttachmentID != nil {
		if attachment, err := u.attachmentUseCase.GetAttachmentByID(*user.ProfileImgAttachmentID); err == nil {
			profileImgURL = attachment.FileURL
		}
	}

	return transforms.UserEntityToUserDisplayResponse(user, profileImgURL), nil
}

func (u *userUseCase) ChangePassword(req *dtos.ChangePasswordRequestDto, userID uint) (*dtos.ChangePasswordResponseDto, error) {
	user, err := u.userRepo.GetByID(userID)
	if err != nil {
		return nil, errors.New("This Account is not registered")
	}

	// Compare password
	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.CurrentPassword))
	if err != nil {
		return nil, errors.New("Incorrect password")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}
	user.PasswordHash = string(hashedPassword)

	if err := u.userRepo.Update(user); err != nil {
		log.Println("Error updating user:", err)
		return nil, err
	}

	return &dtos.ChangePasswordResponseDto{
		Message: "Password changed successfully",
	}, nil
}

func (u *userUseCase) ForgetPasswordStep1(ctx context.Context, req *dtos.ForgetPasswordStep1RequestDto) (*dtos.ForgetPasswordStep1ResponseDto, error) {
	user, err := u.userRepo.GetByUniID(req.UniversityID)
	if err != nil {
		return nil, errors.New("This Account is not registered")
	}

	otp := utils.GenerateRandomString(6)

	err = u.emailService.SendOTP(user.Email, otp)

	tempData := &dtos.TempForgetPasswordData{
		UniversityID: user.UniversityID,
		OTP:          otp,
		User:         *user,
	}

	err = u.userCacheRepo.SetForgetPasswordData(ctx, user.UniversityID, tempData, 5*time.Minute)
	if err != nil {
		log.Println("Error storing registration data in cache:", err)
		return nil, err
	}

	return &dtos.ForgetPasswordStep1ResponseDto{
		Message: "Please proceed to step 2",
		Email:   user.Email,
	}, nil
}

func (u *userUseCase) ForgetPasswordStep2(ctx context.Context, req *dtos.ForgetPasswordStep2RequestDto) (*dtos.ForgetPasswordStep2ResponseDto, error) {
	tempData, err := u.userCacheRepo.GetForgetPasswordData(ctx, req.UniversityID)
	if err != nil {
		return nil, errors.New("no pending registration found or OTP expired")
	}

	// Validate OTP
	if tempData.OTP != req.OTP {
		return nil, errors.New("invalid OTP")
	}

	user := tempData.User

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}
	user.PasswordHash = string(hashedPassword)

	if err := u.userRepo.Update(&user); err != nil {
		log.Println("Error updating user:", err)
		return nil, err
	}

	return &dtos.ForgetPasswordStep2ResponseDto{
		Message: "Password changed successfully",
	}, nil
}

func (u *userUseCase) AdminUpdateUser(req *dtos.AdminUpdateUserRequestDto) error {
	user, err := u.userRepo.GetByID(req.UserID)
	if err != nil {
		return errors.New("user not found")
	}

	if req.UniversityID != "" && req.UniversityID != user.UniversityID {
		duplicatedUser, _ := u.userRepo.GetByUniID(req.UniversityID)
		if duplicatedUser != nil {
			return errors.New("university ID already exists")
		}

		user.UniversityID = req.UniversityID
	}
	if req.Email != "" && req.Email != user.Email {
		duplicatedUser, _ := u.userRepo.GetByEmail(req.Email)
		if duplicatedUser != nil {
			return errors.New("email already exists")
		}

		user.Email = req.Email
		user.IsVerified = false
	}
	if req.NewPassword != "" {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.NewPassword), bcrypt.DefaultCost)
		if err != nil {
			return err
		}
		user.PasswordHash = string(hashedPassword)
	}
	if req.FirstName != "" && req.FirstName != user.FirstName {
		user.FirstName = req.FirstName
	}
	if req.LastName != "" && req.LastName != user.LastName {
		user.LastName = req.LastName
	}
	if req.Role != "" && req.Role != user.Role {
		user.Role = req.Role
	}

	err = u.userRepo.Update(user)

	if err != nil {
		log.Println("Error updating user:", err)
		return err
	}

	return nil
}

func (u *userUseCase) GetProfileImgUrlByUserID(userID uint) string {
	owner, err := u.userRepo.GetByID(userID)
	if err != nil || owner == nil {
		return ""
	}

	var profileImgURL string
	if owner.ProfileImgAttachment != nil && owner.ProfileImgAttachment.ObjectKey != "" {
		url, err := u.attachmentUseCase.GenerateFileURL(owner.ProfileImgAttachment.ObjectKey)
		if err == nil {
			profileImgURL = url
		}
	}

	if profileImgURL == "" {
		profileImgURL = "https://ui-avatars.com/api/?name=" + owner.FirstName + "+" + owner.LastName
	}

	return profileImgURL
}
